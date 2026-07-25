package main

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/yuranikolaev/country-map-svg-generator/internal/catalog"
	"github.com/yuranikolaev/country-map-svg-generator/internal/config"
	"github.com/yuranikolaev/country-map-svg-generator/internal/geometry"
)

// applyGeometryPreset resolves the preset-supplied layout and byte budget so a
// report shows what geometry will actually use, not the sparse request the
// configuration produced. Without it, `explain` would print an empty long side
// for every entity that relies on its preset — which is most of them.
func applyGeometryPreset(request geometry.Input) (geometry.Input, error) {
	applied, err := geometry.ApplyPreset(request)
	if err != nil {
		return geometry.Input{}, failf(ExitRender, "preset_unresolvable", "%v", err)
	}
	return applied, nil
}

func entityName(corpus *catalog.Corpus, iso string) string {
	for _, entity := range corpus.Manifest.Entities {
		if entity.Alpha2 == iso {
			return entity.Name
		}
	}
	return iso
}

// flattenSettings renders resolved settings as dotted path to displayable value,
// using the same JSON shape the schema is defined in. Deriving it from the wire
// form rather than from a hand-written switch is what keeps `explain` from
// silently omitting a field someone added to the schema.
func flattenSettings(settings config.Settings) map[string]string {
	encoded, err := json.Marshal(settings)
	if err != nil {
		return map[string]string{}
	}
	var generic map[string]any
	if err := json.Unmarshal(encoded, &generic); err != nil {
		return map[string]string{}
	}
	out := map[string]string{}
	var walk func(any, string)
	walk = func(value any, prefix string) {
		switch typed := value.(type) {
		case map[string]any:
			for key, child := range typed {
				walk(child, joinPath(prefix, key))
			}
		case []any:
			parts := make([]string, len(typed))
			for i, item := range typed {
				parts[i] = display(item)
			}
			out[prefix] = strings.Join(parts, ", ")
		default:
			if prefix != "" {
				out[prefix] = display(value)
			}
		}
	}
	walk(generic, "")
	return out
}

func display(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case float64:
		return trimNumber(typed)
	case bool:
		return strconv.FormatBool(typed)
	case nil:
		return ""
	}
	encoded, err := json.Marshal(value)
	if err != nil {
		return ""
	}
	return string(encoded)
}

// trimNumber renders a float without a trailing ".0", so a long side reads as
// 160 rather than 160.0 in output a human scans.
func trimNumber(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func joinPath(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + "." + key
}
