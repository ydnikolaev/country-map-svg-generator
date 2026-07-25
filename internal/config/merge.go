package config

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
)

// Layer is one step of the precedence chain, carrying the settings it
// contributes and the name `explain` reports as the origin of anything it wins.
// The name is descriptive rather than symbolic — "preset site-default",
// "country US" — because it is read by a human trying to find where a value
// came from.
type Layer struct {
	Name     string
	Settings Settings
}

// LayerNames are the fixed precedence steps, weakest first. The order is the
// contract; the fact that it lives here as data rather than as the shape of
// some merge function is what lets `explain` describe it without restating it.
var LayerNames = []string{
	"embedded defaults",
	"preset ancestry",
	"document globals",
	"profile block",
	"country override",
	"command-line flags",
}

// Provenance maps each resolved key path to the layer that last set it. A path
// absent from the map was never set by anyone.
type Provenance map[KeyPath]string

// Merge folds the layers in order, later winning, and records where each value
// came from.
//
// It merges typed settings rather than generic maps so that every layer has
// already been through key checking and validation by the time it gets here — a
// generic merge would defer both to the end and report a preset's typo against
// whichever document inherited it.
func Merge(layers []Layer) (Settings, Provenance) {
	var out Settings
	provenance := Provenance{}
	target := reflect.ValueOf(&out).Elem()
	for _, layer := range layers {
		mergeStruct(target, reflect.ValueOf(layer.Settings), "", layer.Name, provenance)
	}
	return out, provenance
}

// mergeStruct copies every field the overlay has an opinion about onto the
// target, recursing into nested settings blocks so that a layer setting one
// token does not erase the rest.
func mergeStruct(target, overlay reflect.Value, prefix, layer string, provenance Provenance) {
	structType := target.Type()
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if !field.IsExported() {
			continue
		}
		name, onWire := wireName(field)
		path := prefix
		switch {
		case name == "" && field.Anonymous:
			// An embedded settings block contributes its fields at the parent's
			// path, matching the flattened wire form.
		case onWire:
			path = join(prefix, name)
		default:
			// An off-wire field (Padding.Uniform) still merges — it carries a
			// value the author wrote — but it has no path of its own, so it
			// reports under its parent.
		}
		mergeValue(target.Field(i), overlay.Field(i), path, layer, provenance)
	}
}

func mergeValue(target, overlay reflect.Value, path, layer string, provenance Provenance) {
	switch overlay.Kind() {
	case reflect.Struct:
		mergeStruct(target, overlay, path, layer, provenance)

	case reflect.Ptr:
		if overlay.IsNil() {
			// The layer said nothing about this field. Saying nothing is not the
			// same as saying zero, which is the entire reason every field is a
			// pointer.
			return
		}
		if overlay.Type().Elem().Kind() == reflect.Struct {
			// A nested block merges field by field, so setting layout.longSide
			// does not discard an inherited layout.mode.
			if target.IsNil() {
				target.Set(reflect.New(target.Type().Elem()))
			}
			mergeStruct(target.Elem(), overlay.Elem(), path, layer, provenance)
			return
		}
		target.Set(overlay)
		provenance[path] = layer

	case reflect.Slice:
		if overlay.IsNil() {
			return
		}
		// A list replaces wholesale rather than appending. Appending would make
		// a country override unable to shorten an inherited list, and would make
		// the resolved value depend on how many ancestors happened to mention it.
		target.Set(overlay)
		provenance[path] = layer

	case reflect.Map:
		if overlay.IsNil() {
			return
		}
		if target.IsNil() {
			target.Set(reflect.MakeMap(overlay.Type()))
		}
		iter := overlay.MapRange()
		for iter.Next() {
			target.SetMapIndex(iter.Key(), iter.Value())
			provenance[join(path, fmt.Sprint(iter.Key().Interface()))] = layer
		}
	}
}

// Origins renders the provenance as a stable, sorted list for `explain`.
func (p Provenance) Origins() []Origin {
	out := make([]Origin, 0, len(p))
	for path, layer := range p {
		out = append(out, Origin{Path: path, Layer: layer})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out
}

// Origin is one resolved key and the layer it came from.
type Origin struct {
	Path  KeyPath `json:"path"`
	Layer string  `json:"layer"`
}

// UnmarshalJSON lets padding be written as a single number or as four sides.
// A uniform padding is overwhelmingly the common case, and forcing four keys for
// it would make every config noisier to read for no gain.
func (p *Padding) UnmarshalJSON(raw []byte) error {
	var uniform float64
	if err := json.Unmarshal(raw, &uniform); err == nil {
		p.Uniform = &uniform
		return nil
	}
	// The alias breaks the recursion into this method.
	type sides Padding
	var value sides
	if err := json.Unmarshal(raw, &value); err != nil {
		return fmt.Errorf("padding must be a number or an object with top, right, bottom and left: %w", err)
	}
	*p = Padding(value)
	return nil
}

// MarshalJSON round-trips the uniform spelling so a resolved config re-serializes
// to something an author would recognise as what they wrote.
func (p Padding) MarshalJSON() ([]byte, error) {
	if p.Uniform != nil {
		return json.Marshal(*p.Uniform)
	}
	type sides Padding
	return json.Marshal(sides(p))
}

// Resolved is the outcome of the precedence chain for one entity: the settings
// plus where each one came from.
type Resolved struct {
	Settings   Settings   `json:"settings"`
	Provenance Provenance `json:"-"`
}
