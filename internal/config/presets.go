package config

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
)

//go:embed presets/v1.json
var presetJSON []byte

// Preset is one named, inheritable settings bundle. A preset may name a single
// parent; the chain is walked oldest-ancestor-first so the nearest definition
// wins, which is the same direction the document layers resolve in.
type Preset struct {
	Name    string  `json:"name"`
	Extends *string `json:"extends,omitempty"`
	// Note explains what the preset is for. It travels with the preset rather
	// than living in a separate document so that a preset and its reason cannot
	// drift apart.
	Note     string   `json:"note,omitempty"`
	Settings Settings `json:"settings"`
}

// defaultsBlock is the weakest layer. It carries a `note` alongside the
// settings because the values are product choices that need their reasons
// travelling with them, and because keeping them in data rather than in Go is
// what stops "card" and "un" from becoming literals in the code.
type defaultsBlock struct {
	Note string `json:"note"`
	Settings
}

type presetFile struct {
	Version  string        `json:"version"`
	Note     string        `json:"note"`
	Defaults defaultsBlock `json:"defaults"`
	Presets  []Preset      `json:"presets"`
}

var (
	presetsOnce   sync.Once
	presetsByName map[string]Preset
	presetDefault Settings
	presetsErr    error
)

// Defaults is the weakest layer of the precedence chain: what a document that
// says nothing at all resolves to.
func Defaults() (Settings, error) {
	if _, err := Presets(); err != nil {
		return Settings{}, err
	}
	return presetDefault, nil
}

// Presets returns the embedded named presets. Parsing is lazy and memoized: the
// file is small, but a package-level panic on init would take down `--help`
// along with everything else.
func Presets() (map[string]Preset, error) {
	presetsOnce.Do(func() { presetsByName, presetDefault, presetsErr = parsePresets(presetJSON) })
	return presetsByName, presetsErr
}

// PresetNames lists the embedded presets in stable order, for diagnostics and
// for `init` to offer.
func PresetNames() ([]string, error) {
	all, err := Presets()
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(all))
	for name := range all {
		names = append(names, name)
	}
	sort.Strings(names)
	return names, nil
}

func parsePresets(raw []byte) (map[string]Preset, Settings, error) {
	// The embedded file goes through the same key checking documents do. A typo
	// in a shipped preset would otherwise be invisible until it silently failed
	// to apply, and it would look like a bug in inheritance rather than a typo.
	var generic struct {
		Defaults map[string]any `json:"defaults"`
		Presets  []struct {
			Name     string `json:"name"`
			Extends  string `json:"extends"`
			Settings any    `json:"settings"`
		} `json:"presets"`
	}
	if err := json.Unmarshal(raw, &generic); err != nil {
		return nil, Settings{}, fmt.Errorf("embedded presets are not valid JSON: %w", err)
	}

	var file presetFile
	decoder := json.NewDecoder(strings.NewReader(string(raw)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&file); err != nil {
		return nil, Settings{}, fmt.Errorf("embedded presets do not match the preset schema: %w", err)
	}
	if file.Version != "v1" {
		return nil, Settings{}, fmt.Errorf("embedded presets declare version %q, want v1", file.Version)
	}

	// `note` documents the block rather than configuring anything, so it is
	// stripped before the settings shape is checked.
	defaults := make(map[string]any, len(generic.Defaults))
	for key, value := range generic.Defaults {
		if key != "note" {
			defaults[key] = value
		}
	}
	if err := walkPresetSettings(defaults); err != nil {
		return nil, Settings{}, fmt.Errorf("embedded defaults: %w", err)
	}

	out := make(map[string]Preset, len(file.Presets))
	for i, preset := range file.Presets {
		if preset.Name == "" {
			return nil, Settings{}, fmt.Errorf("embedded preset %d has no name", i)
		}
		if _, clash := out[preset.Name]; clash {
			return nil, Settings{}, fmt.Errorf("embedded preset %q is defined twice", preset.Name)
		}
		if generic.Presets[i].Settings != nil {
			if err := walkPresetSettings(generic.Presets[i].Settings); err != nil {
				return nil, Settings{}, fmt.Errorf("embedded preset %q: %w", preset.Name, err)
			}
		}
		out[preset.Name] = preset
	}
	for name, preset := range out {
		if preset.Extends == nil {
			continue
		}
		if _, ok := out[*preset.Extends]; !ok {
			return nil, Settings{}, fmt.Errorf("embedded preset %q extends unknown preset %q", name, *preset.Extends)
		}
	}
	return out, file.Defaults.Settings, nil
}

func walkPresetSettings(settings any) error {
	var unknown []KeyPath
	walkKeys(settings, settingsType(), "", &unknown)
	if len(unknown) == 0 {
		return nil
	}
	sort.Strings(unknown)
	return &UnknownKeyError{Keys: unknown}
}

// PresetChain resolves a preset name into its ancestry, oldest first, so the
// caller can lay the layers down in precedence order.
//
// Inheritance is single-parent, which makes a cycle the only structural failure
// possible; it is detected rather than bounded by a depth limit, because a limit
// would turn a cycle into a confusing "too deep" message and would cap a
// legitimate deep chain at the same time.
func PresetChain(name string) ([]Preset, error) {
	all, err := Presets()
	if err != nil {
		return nil, err
	}
	var chain []Preset
	seen := map[string]bool{}
	var path []string

	for current := name; current != ""; {
		preset, ok := all[current]
		if !ok {
			known, _ := PresetNames()
			return nil, fmt.Errorf("unknown preset %q; embedded presets are %s", current, strings.Join(known, ", "))
		}
		if seen[current] {
			return nil, fmt.Errorf("preset inheritance cycle: %s", strings.Join(append(path, current), " -> "))
		}
		seen[current] = true
		path = append(path, current)
		chain = append(chain, preset)
		if preset.Extends == nil {
			break
		}
		current = *preset.Extends
	}

	// Reverse: the chain was walked child-first, and layers are applied
	// oldest-ancestor-first so the nearest definition wins.
	for i, j := 0, len(chain)-1; i < j; i, j = i+1, j-1 {
		chain[i], chain[j] = chain[j], chain[i]
	}
	return chain, nil
}
