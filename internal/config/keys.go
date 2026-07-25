package config

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// This file exists because encoding/json matches field names
// case-insensitively when no tag matches exactly. `{"layout":{"longside":160}}`
// therefore binds to LongSide and decodes cleanly, which for a public schema
// frozen at v1 would make every casing of every key part of the contract
// forever — and would make a genuine typo look like a working document.
//
// REQ-3 says validation rejects unknown fields. Case-insensitive binding is the
// same defect wearing a different hat: the caller wrote a key the schema does
// not define, and something happened anyway.
//
// The allowed key set is derived from the struct tags by reflection rather than
// written out, so the schema has exactly one source of truth. The same
// traversal answers "what keys exist", which is what `validate` and the schema
// discovery surface report.

// KeyPath is a dotted path into a document, used for diagnostics and for
// provenance. Map keys chosen by the author (an ISO code, a profile name)
// appear verbatim.
type KeyPath = string

// UnknownKeyError names every key a document carries that the schema does not
// define, with its exact path. All of them are reported at once: an author
// fixing a config wants the whole list, not one key per run.
type UnknownKeyError struct {
	Keys []KeyPath
}

func (e *UnknownKeyError) Error() string {
	if len(e.Keys) == 1 {
		return fmt.Sprintf("unknown field %s", e.Keys[0])
	}
	return fmt.Sprintf("unknown fields %s", strings.Join(e.Keys, ", "))
}

// CheckKeys walks a decoded document against the schema and reports every key
// the schema does not define, matching case exactly.
//
// It runs before the typed decode, so a document is refused for the key the
// author actually mistyped rather than for whatever type error that mistyping
// happens to cause downstream.
func CheckKeys(decoded any) error {
	var unknown []KeyPath
	walkKeys(decoded, reflect.TypeOf(Document{}), "", &unknown)
	if len(unknown) == 0 {
		return nil
	}
	sort.Strings(unknown)
	return &UnknownKeyError{Keys: unknown}
}

// settingsType is the type a settings-shaped fragment is checked against — a
// preset body, or a per-country or per-profile overlay. It exists so those
// fragments go through the same key checking a whole document does.
func settingsType() reflect.Type { return reflect.TypeOf(Settings{}) }

// tokensType is the type a bare token block is checked against — a style's
// defaults, which carry tokens and nothing else.
func tokensType() reflect.Type { return reflect.TypeOf(Tokens{}) }

// SchemaKeys reports every key path the schema defines, in stable order. It is
// the discovery surface REQ-1 asks for and the enumeration a diagnostic uses to
// suggest a near miss.
func SchemaKeys() []KeyPath {
	var keys []KeyPath
	collectKeys(reflect.TypeOf(Document{}), "", &keys, map[reflect.Type]bool{})
	sort.Strings(keys)
	return keys
}

func walkKeys(value any, target reflect.Type, prefix string, unknown *[]KeyPath) {
	object, ok := value.(map[string]any)
	if !ok {
		// A non-object where the schema wants a struct is a type error, which the
		// typed decode reports with better context than this walk could.
		return
	}
	target = deref(target)

	// A map-typed node accepts author-chosen keys; the schema constrains only
	// what lives beneath each one.
	if target.Kind() == reflect.Map {
		for key, child := range object {
			walkKeys(child, target.Elem(), join(prefix, key), unknown)
		}
		return
	}
	if target.Kind() != reflect.Struct {
		return
	}

	fields := fieldsOf(target)
	for key, child := range object {
		field, known := fields[key]
		if !known {
			*unknown = append(*unknown, join(prefix, key))
			continue
		}
		walkKeys(child, field, join(prefix, key), unknown)
	}
}

// fieldsOf maps a struct's exact JSON key names to their field types,
// flattening anonymous embedded structs the way encoding/json does. A field
// tagged "-" is not part of the wire form and is therefore not an allowed key.
func fieldsOf(target reflect.Type) map[string]reflect.Type {
	out := map[string]reflect.Type{}
	for i := 0; i < target.NumField(); i++ {
		field := target.Field(i)
		if !field.IsExported() {
			continue
		}
		name, ok := wireName(field)
		if !ok {
			continue
		}
		if name == "" && field.Anonymous {
			for key, embedded := range fieldsOf(deref(field.Type)) {
				out[key] = embedded
			}
			continue
		}
		out[name] = field.Type
	}
	return out
}

// wireName returns the exact key a field is spelled with, and whether the field
// appears on the wire at all.
func wireName(field reflect.StructField) (string, bool) {
	tag, tagged := field.Tag.Lookup("json")
	if !tagged {
		if field.Anonymous {
			return "", true
		}
		return field.Name, true
	}
	name, _, _ := strings.Cut(tag, ",")
	if name == "-" {
		return "", false
	}
	if name == "" {
		if field.Anonymous {
			return "", true
		}
		return field.Name, true
	}
	return name, true
}

func collectKeys(target reflect.Type, prefix string, keys *[]KeyPath, seen map[reflect.Type]bool) {
	target = deref(target)
	switch target.Kind() {
	case reflect.Map:
		collectKeys(target.Elem(), join(prefix, "<key>"), keys, seen)
		return
	case reflect.Struct:
	default:
		return
	}
	// Guard against a type that transitively contains itself. No v1 type does,
	// but a recursive schema added later would otherwise hang enumeration
	// rather than fail visibly.
	if seen[target] {
		return
	}
	seen[target] = true
	defer delete(seen, target)

	for name, field := range fieldsOf(target) {
		path := join(prefix, name)
		*keys = append(*keys, path)
		collectKeys(field, path, keys, seen)
	}
}

func deref(target reflect.Type) reflect.Type {
	for target.Kind() == reflect.Ptr {
		target = target.Elem()
	}
	return target
}

func join(prefix, key string) KeyPath {
	if prefix == "" {
		return key
	}
	return prefix + "." + key
}
