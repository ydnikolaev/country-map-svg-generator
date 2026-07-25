package config

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
)

//go:embed styles/v1.json
var styleJSON []byte

// Style is one of REQ-4's base styles, expressed entirely as token defaults.
//
// A style never reaches the serializer by name. That is what "all resolve
// through tokens rather than renderer forks" means in practice: adding a style
// adds an entry to the embedded file, and no rendering code learns about it.
type Style struct {
	Name   string `json:"name"`
	Note   string `json:"note"`
	Tokens Tokens `json:"tokens"`
}

type styleFile struct {
	Version string  `json:"version"`
	Note    string  `json:"note"`
	Styles  []Style `json:"styles"`
}

var (
	stylesOnce   sync.Once
	stylesByName map[string]Style
	stylesErr    error
)

// StyleDefaults returns the token defaults for one style. They enter the
// precedence chain as their own layer, immediately after the embedded defaults
// and beneath everything a document can say, so a token an author sets always
// wins and `explain` can name the style as the origin of one they did not.
func StyleDefaults(name string) (Settings, error) {
	all, err := Styles()
	if err != nil {
		return Settings{}, err
	}
	style, ok := all[name]
	if !ok {
		return Settings{}, fmt.Errorf("unknown style %q; accepted styles are %s", name, strings.Join(StyleNames, ", "))
	}
	tokens := style.Tokens
	return Settings{Tokens: &tokens}, nil
}

// Styles returns the embedded styles by name.
//
// StyleNames exists alongside it on purpose: validation needs the vocabulary
// without paying to parse the file, and a v1 binary must refuse an unknown style
// identically whether or not the style data happens to be readable. parseStyles
// asserts the two agree in both directions.
func Styles() (map[string]Style, error) {
	stylesOnce.Do(func() { stylesByName, stylesErr = parseStyles(styleJSON) })
	return stylesByName, stylesErr
}

func parseStyles(raw []byte) (map[string]Style, error) {
	var generic struct {
		Styles []struct {
			Name   string `json:"name"`
			Tokens any    `json:"tokens"`
		} `json:"styles"`
	}
	if err := json.Unmarshal(raw, &generic); err != nil {
		return nil, fmt.Errorf("embedded styles are not valid JSON: %w", err)
	}

	var file styleFile
	decoder := json.NewDecoder(strings.NewReader(string(raw)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&file); err != nil {
		return nil, fmt.Errorf("embedded styles do not match the style schema: %w", err)
	}
	if file.Version != "v1" {
		return nil, fmt.Errorf("embedded styles declare version %q, want v1", file.Version)
	}

	out := make(map[string]Style, len(file.Styles))
	for i, style := range file.Styles {
		if style.Name == "" {
			return nil, fmt.Errorf("embedded style %d has no name", i)
		}
		if _, clash := out[style.Name]; clash {
			return nil, fmt.Errorf("embedded style %q is defined twice", style.Name)
		}
		// The token block goes through the same key checking a document does, so
		// a typo in a shipped style fails loudly instead of silently not applying.
		if generic.Styles[i].Tokens != nil {
			var unknown []KeyPath
			walkKeys(generic.Styles[i].Tokens, tokensType(), "", &unknown)
			if len(unknown) != 0 {
				sort.Strings(unknown)
				return nil, fmt.Errorf("embedded style %q: %w", style.Name, &UnknownKeyError{Keys: unknown})
			}
		}
		out[style.Name] = style
	}

	// The embedded data and the vocabulary validation uses must agree, or a
	// document could name a style that validates and then has no tokens.
	for _, name := range StyleNames {
		if _, ok := out[name]; !ok {
			return nil, fmt.Errorf("style %q is accepted by validation but has no embedded token defaults", name)
		}
	}
	for name := range out {
		if !contains(StyleNames, name) {
			return nil, fmt.Errorf("embedded style %q is not in the accepted style vocabulary", name)
		}
	}
	return out, nil
}
