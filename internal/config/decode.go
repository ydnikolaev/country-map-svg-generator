package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
)

// Format names how a document's bytes are spelled. Both formats reach the same
// typed document through the same checks — the only thing that differs is how
// bytes become a generic value.
type Format string

const (
	FormatYAML Format = "yaml"
	FormatJSON Format = "json"
)

// Extensions maps the accepted file extensions to their format. A file with any
// other extension is refused rather than sniffed: guessing would make the
// diagnostic for a genuinely malformed file read as a format mistake.
var Extensions = map[string]Format{
	".yaml": FormatYAML,
	".yml":  FormatYAML,
	".json": FormatJSON,
}

// DocumentError is any failure to turn bytes into a document the schema
// recognises: bad syntax, an unknown or wrongly-cased key, a value the schema
// cannot represent, an unsupported extension.
//
// It exists as one type so a caller classifies by type rather than by matching
// substrings of an error message. The CLI's exit classes are a frozen public
// contract that a build pipeline branches on, and deriving one of them from
// another library's wording means a dependency upgrade can silently reclassify
// a failure.
type DocumentError struct {
	// Path is the configuration file, when the failure came from one.
	Path string
	Err  error
}

func (e *DocumentError) Error() string {
	if e.Path == "" {
		return e.Err.Error()
	}
	return e.Path + ": " + e.Err.Error()
}

func (e *DocumentError) Unwrap() error { return e.Err }

func documentError(path string, err error) error {
	if err == nil {
		return nil
	}
	var already *DocumentError
	if errors.As(err, &already) {
		if already.Path == "" {
			already.Path = path
		}
		return already
	}
	return &DocumentError{Path: path, Err: err}
}

// DecodeFile reads and fully validates a configuration document.
//
// The pipeline is deliberately uniform across formats: bytes become a generic
// value using the format's own parser, that value is normalized to something
// JSON can represent, its keys are checked exactly, and only then is it decoded
// into the typed document. Routing both formats through one typed decode is what
// keeps YAML and JSON from disagreeing about which keys exist — a difference
// that would otherwise show up as a config that validates in one spelling and
// not in the other.
func DecodeFile(path string) (*Document, error) {
	extension := strings.ToLower(filepath.Ext(path))
	format, known := Extensions[extension]
	if !known {
		return nil, documentError(path, fmt.Errorf("unsupported config extension %q; use one of %s", extension, strings.Join(sortedExtensions(), ", ")))
	}
	// A read failure is deliberately not a DocumentError: the document is not
	// the problem, the path or the filesystem is, and those get a different exit
	// class because the caller fixes them differently.
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	doc, err := Decode(raw, format)
	if err != nil {
		return nil, documentError(path, err)
	}
	return doc, nil
}

// Decode turns bytes into a validated document. It does not apply presets or
// layering — those need a vocabulary and a selection, and keeping them separate
// is what lets `validate` check a document that names entities the current
// corpus does not carry.
func Decode(raw []byte, format Format) (*Document, error) {
	generic, err := parseGeneric(raw, format)
	if err != nil {
		return nil, documentError("", err)
	}
	normalized, err := normalize(generic, "")
	if err != nil {
		return nil, documentError("", err)
	}
	if err := CheckKeys(normalized); err != nil {
		return nil, documentError("", err)
	}

	// Re-marshal through JSON so exactly one typed decoder exists. The
	// alternative — a YAML decoder for YAML and a JSON decoder for JSON — gives
	// two implementations of "what does this document mean", and they differ in
	// the details that matter here: case sensitivity, embedded struct
	// flattening, and how a union type like padding is handled.
	encoded, err := json.Marshal(normalized)
	if err != nil {
		return nil, documentError("", fmt.Errorf("document could not be normalized: %w", err))
	}
	var doc Document
	decoder := json.NewDecoder(bytes.NewReader(encoded))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&doc); err != nil {
		return nil, documentError("", err)
	}
	return &doc, nil
}

func parseGeneric(raw []byte, format Format) (any, error) {
	switch format {
	case FormatJSON:
		var value any
		if err := json.Unmarshal(raw, &value); err != nil {
			return nil, fmt.Errorf("invalid JSON: %w", err)
		}
		return value, nil
	case FormatYAML:
		var value any
		if err := yaml.Unmarshal(raw, &value); err != nil {
			// goccy's formatted error carries the line, column and an excerpt of
			// the offending source, which is the whole reason it is the parser
			// here rather than a minimal one.
			return nil, fmt.Errorf("invalid YAML:\n%s", yaml.FormatError(err, false, true))
		}
		return value, nil
	}
	return nil, fmt.Errorf("unknown format %q", format)
}

// normalize converts a parsed document into a value JSON can represent, failing
// with the exact path when it cannot. YAML admits things JSON does not — a
// non-string mapping key, a NaN or an infinity — and each would otherwise
// surface far from its cause: a non-string key as an opaque marshalling error, a
// non-finite number as a viewBox full of NaN.
func normalize(value any, path KeyPath) (any, error) {
	switch typed := value.(type) {
	case map[string]any:
		out := make(map[string]any, len(typed))
		for key, child := range typed {
			normalized, err := normalize(child, join(path, key))
			if err != nil {
				return nil, err
			}
			out[key] = normalized
		}
		return out, nil

	case map[any]any:
		out := make(map[string]any, len(typed))
		for key, child := range typed {
			text, ok := key.(string)
			if !ok {
				return nil, fmt.Errorf("%s: keys must be strings, got %T", pathOrRoot(path), key)
			}
			normalized, err := normalize(child, join(path, text))
			if err != nil {
				return nil, err
			}
			out[text] = normalized
		}
		return out, nil

	case []any:
		out := make([]any, len(typed))
		for i, child := range typed {
			normalized, err := normalize(child, fmt.Sprintf("%s[%d]", path, i))
			if err != nil {
				return nil, err
			}
			out[i] = normalized
		}
		return out, nil

	case float64:
		if math.IsNaN(typed) {
			return nil, fmt.Errorf("%s: must be a number, got a not-a-number value", pathOrRoot(path))
		}
		if math.IsInf(typed, 0) {
			return nil, fmt.Errorf("%s: must be finite, got an infinity", pathOrRoot(path))
		}
		return typed, nil

	case uint64:
		return float64(typed), nil
	case int64:
		return float64(typed), nil
	case int:
		return float64(typed), nil

	default:
		return value, nil
	}
}

func pathOrRoot(path KeyPath) string {
	if path == "" {
		return "document"
	}
	return path
}

func sortedExtensions() []string {
	out := make([]string, 0, len(Extensions))
	for extension := range Extensions {
		out = append(out, extension)
	}
	return sorted(out)
}
