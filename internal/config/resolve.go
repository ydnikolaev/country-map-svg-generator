package config

import (
	"fmt"
	"sort"
	"strings"
)

// Resolve folds one entity's configuration down the precedence chain and
// reports where every value came from.
//
// The chain is: embedded defaults, the preset ancestry named by `extends`, the
// document's own globals, the block for the resolved profile, the per-country
// override, and finally the command-line flags. Each step is a Layer, so
// `explain` can name the origin of any value without restating the order.
func Resolve(doc *Document, iso string, flags Settings) (Resolved, error) {
	layers, err := layersFor(doc, iso, flags)
	if err != nil {
		return Resolved{}, err
	}
	settings, provenance := Merge(layers)
	pruneSupersededLayoutFields(&settings, provenance, layerOrder(layers))
	return Resolved{Settings: settings, Provenance: provenance}, nil
}

func layerOrder(layers []Layer) map[string]int {
	order := make(map[string]int, len(layers))
	for i, layer := range layers {
		order[layer.Name] = i
	}
	return order
}

// pruneSupersededLayoutFields drops sizing that belonged to a layout mode a
// later layer replaced.
//
// This is the spec's own example: the document sets `mode: tight, longSide:
// 160` and a country override sets `mode: contain, width: 720, height: 420`.
// Without pruning, the merged layout carries a contain mode and an inherited
// long side, and the resolved check refuses the documented shape.
//
// Only fields set *earlier* than the mode are dropped. A layer that sets a mode
// and a foreign field together is stating a contradiction about its own intent,
// and ValidateDocument already refuses that per layer — pruning it here would
// silently accept the mistake instead.
func pruneSupersededLayoutFields(settings *Settings, provenance Provenance, order map[string]int) {
	if settings.Layout == nil || settings.Layout.Mode == nil {
		return
	}
	modeLayer, ok := order[provenance["layout.mode"]]
	if !ok {
		return
	}
	foreign := map[string]**float64{
		"layout.width": &settings.Layout.Width, "layout.height": &settings.Layout.Height,
	}
	if *settings.Layout.Mode == "contain" {
		foreign = map[string]**float64{
			"layout.longSide":  &settings.Layout.LongSide,
			"layout.maxWidth":  &settings.Layout.MaxWidth,
			"layout.maxHeight": &settings.Layout.MaxHeight,
		}
	}
	for path, field := range foreign {
		if *field == nil {
			continue
		}
		if source, known := order[provenance[path]]; known && source < modeLayer {
			*field = nil
			delete(provenance, path)
		}
	}
}

func layersFor(doc *Document, iso string, flags Settings) ([]Layer, error) {
	defaults, err := Defaults()
	if err != nil {
		return nil, err
	}
	layers := []Layer{{Name: LayerNames[0], Settings: defaults}}

	if doc != nil && doc.Extends != nil {
		chain, err := PresetChain(*doc.Extends)
		if err != nil {
			return nil, err
		}
		for _, preset := range chain {
			// Each ancestor is named individually rather than collapsed into one
			// "preset ancestry" layer: with a chain of three, "which preset set
			// this" is the question an author actually has.
			layers = append(layers, Layer{Name: LayerNames[1] + " " + preset.Name, Settings: preset.Settings})
		}
	}
	if doc != nil {
		layers = append(layers, Layer{Name: LayerNames[2], Settings: doc.Settings})
	}

	// The profile block to apply depends on the resolved profile, and the
	// resolved profile can itself be set by the country override or a flag —
	// `countries: {US: {profile: hero}}` is the documented example. So the
	// profile is resolved first, from every layer that may carry it except the
	// profile blocks themselves, and only then is the matching block inserted.
	// Doing it in one pass would make the result depend on an order that has no
	// non-arbitrary answer.
	profile := resolveProfile(layers, doc, iso, flags)
	if doc != nil && profile != "" {
		if block, ok := doc.Profiles[profile]; ok {
			layers = append(layers, Layer{Name: LayerNames[3] + " " + profile, Settings: block})
		}
	}
	if doc != nil && iso != "" {
		if override, ok := doc.Countries[iso]; ok {
			layers = append(layers, Layer{Name: LayerNames[4] + " " + iso, Settings: override})
		}
	}
	layers = append(layers, Layer{Name: LayerNames[5], Settings: flags})
	return layers, nil
}

// resolveProfile answers "which profile block applies" using the same
// precedence as everything else, minus the profile blocks. A profile block may
// not select a profile — ValidateDocument refuses that — so this cannot be
// circular.
func resolveProfile(base []Layer, doc *Document, iso string, flags Settings) string {
	candidates := append([]Layer(nil), base...)
	if doc != nil && iso != "" {
		if override, ok := doc.Countries[iso]; ok {
			candidates = append(candidates, Layer{Name: LayerNames[4], Settings: override})
		}
	}
	candidates = append(candidates, Layer{Name: LayerNames[5], Settings: flags})

	resolved, _ := Merge(candidates)
	if resolved.Profile == nil {
		return ""
	}
	return *resolved.Profile
}

// ValidateResolved re-checks a fully resolved configuration. It is not
// redundant with ValidateDocument: a per-layer check cannot see a combination
// that only exists after merging — a document setting `mode: contain` and a
// country override supplying only a `longSide` is legal in each layer and
// contradictory in the result.
func ValidateResolved(iso string, resolved Resolved, vocab Vocabulary) error {
	var problems []ValidationError
	add := func(path, format string, args ...any) {
		problems = append(problems, ValidationError{Path: path, Message: fmt.Sprintf(format, args...)})
	}
	settings := resolved.Settings

	validateSettings(&settings, "", vocab, add)

	for key, value := range map[string]*string{
		"profile": settings.Profile, "boundary": settings.Boundary,
		"delivery": settings.Delivery, "style": settings.Style,
	} {
		if value == nil {
			add(key, "is not set by any layer and has no default")
		}
	}
	if settings.Layout == nil || settings.Layout.Mode == nil {
		add("layout.mode", "is not set by any layer and has no default")
	} else {
		// The cross-layer contradiction: each layer was legal on its own.
		mode := *settings.Layout.Mode
		if mode == "contain" && (settings.Layout.Width == nil || settings.Layout.Height == nil) {
			add("layout", "contain needs both width and height; after merging, %s", missingContainFields(settings.Layout))
		}
		if mode == "tight" && settings.Layout.Width != nil {
			add("layout.width", "survived into a tight layout; a later layer set the mode without clearing the frame")
		}
	}
	if settings.Marker != nil && settings.Marker.Mode != nil && *settings.Marker.Mode == "custom" && len(settings.Marker.Custom) == 0 {
		add("marker.custom", "custom marker mode selects nothing; name at least one capital id")
	}

	if len(problems) == 0 {
		return nil
	}
	for i := range problems {
		if iso != "" {
			problems[i].Path = iso + ": " + problems[i].Path
		}
	}
	sort.Slice(problems, func(i, j int) bool { return problems[i].Path < problems[j].Path })
	return ValidationErrors(problems)
}

func missingContainFields(layout *Layout) string {
	var missing []string
	if layout.Width == nil {
		missing = append(missing, "width")
	}
	if layout.Height == nil {
		missing = append(missing, "height")
	}
	return strings.Join(missing, " and ") + " is missing"
}

// OutputPath renders the resolved filename template for one entity. It is the
// only place placeholders are substituted, so the closed set validation checks
// and the set this expands cannot drift.
func OutputPath(iso string, settings Settings) (string, error) {
	if settings.Output == nil || settings.Output.Filename == nil {
		return "", fmt.Errorf("output.filename is not set by any layer")
	}
	values := map[string]string{"iso": iso}
	for key, value := range map[string]*string{
		"profile": settings.Profile, "boundary": settings.Boundary,
		"style": settings.Style, "delivery": settings.Delivery,
	} {
		if value != nil {
			values[key] = *value
		}
	}
	name := *settings.Output.Filename
	for _, placeholder := range FilenamePlaceholders {
		name = strings.ReplaceAll(name, "{"+placeholder+"}", values[placeholder])
	}
	dir := "."
	if settings.Output.Dir != nil {
		dir = *settings.Output.Dir
	}
	if dir == "" || dir == "." {
		return name, nil
	}
	return strings.TrimRight(dir, "/") + "/" + name, nil
}

// CheckOutputCollisions refuses a batch in which two entities resolve to the
// same file, before anything is written.
//
// REQ-3 puts this before generation for a reason: discovered afterwards, the
// symptom is a catalog that is quietly short by however many entities collided,
// and the manifest would agree with the directory because both were written by
// the same losing pass.
func CheckOutputCollisions(paths map[string]string) error {
	owners := map[string][]string{}
	for iso, path := range paths {
		owners[path] = append(owners[path], iso)
	}
	var problems []ValidationError
	for path, isos := range owners {
		if len(isos) < 2 {
			continue
		}
		sort.Strings(isos)
		problems = append(problems, ValidationError{
			Path:    path,
			Message: fmt.Sprintf("is the output path for %s; add {iso}, {profile} or {boundary} to output.filename to separate them", strings.Join(isos, ", ")),
		})
	}
	if len(problems) == 0 {
		return nil
	}
	sort.Slice(problems, func(i, j int) bool { return problems[i].Path < problems[j].Path })
	return ValidationErrors(problems)
}
