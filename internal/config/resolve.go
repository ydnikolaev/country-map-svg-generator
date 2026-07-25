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
	// Two values have to be known before the layer list can be built, because
	// each of them selects a layer: the profile selects which `profiles:` block
	// applies, and the style selects which token defaults go in. Both can be set
	// by a country override or a flag — `countries: {US: {profile: hero}}` is the
	// documented example — so resolving them in a single pass would apply the
	// card block to a hero card and the wrong style's tokens along with it.
	//
	// This cannot recurse: the profile is resolved without consulting any
	// profile block, and a profile block may not select a profile, which
	// ValidateDocument refuses.
	profile := resolveScalar(defaults, doc, iso, flags, "", func(s Settings) *string { return s.Profile })
	style := resolveScalar(defaults, doc, iso, flags, profile, func(s Settings) *string { return s.Style })

	layers := []Layer{{Name: LayerNames[0], Settings: defaults}}

	// The style's token defaults sit immediately above the embedded defaults and
	// beneath everything a document can say, so a token an author sets always
	// wins and `explain` can name the style as the origin of one they did not.
	if style != "" {
		styleDefaults, err := StyleDefaults(style)
		if err != nil {
			return nil, err
		}
		layers = append(layers, Layer{Name: LayerNames[1] + " " + style, Settings: styleDefaults})
	}

	layers = append(layers, documentLayers(doc, iso, profile)...)
	layers = append(layers, Layer{Name: LayerNames[6], Settings: flags})
	return layers, nil
}

// documentLayers is everything between the style defaults and the flags: the
// preset ancestry, the document's own settings, the resolved profile's block and
// the per-country override.
func documentLayers(doc *Document, iso, profile string) []Layer {
	if doc == nil {
		return nil
	}
	var layers []Layer
	if doc.Extends != nil {
		// PresetChain's error is not returned here because layersFor has already
		// resolved it once; an unknown preset fails there.
		chain, err := PresetChain(*doc.Extends)
		if err != nil {
			return nil
		}
		for _, preset := range chain {
			// Each ancestor is named individually rather than collapsed into one
			// "preset ancestry" layer: with a chain of three, "which preset set
			// this" is the question an author actually has.
			layers = append(layers, Layer{Name: LayerNames[2] + " " + preset.Name, Settings: preset.Settings})
		}
	}
	layers = append(layers, Layer{Name: LayerNames[3], Settings: doc.Settings})
	if profile != "" {
		if block, ok := doc.Profiles[profile]; ok {
			layers = append(layers, Layer{Name: LayerNames[4] + " " + profile, Settings: block})
		}
	}
	if iso != "" {
		if override, ok := doc.Countries[iso]; ok {
			layers = append(layers, Layer{Name: LayerNames[5] + " " + iso, Settings: override})
		}
	}
	return layers
}

// resolveScalar answers "which value of this key wins" using the real
// precedence chain, so a selector cannot disagree with the resolution it
// selects a layer for. Pass an empty profile to exclude the profile blocks,
// which is what resolving the profile itself requires.
func resolveScalar(defaults Settings, doc *Document, iso string, flags Settings, profile string, pick func(Settings) *string) string {
	candidates := []Layer{{Name: LayerNames[0], Settings: defaults}}
	candidates = append(candidates, documentLayers(doc, iso, profile)...)
	candidates = append(candidates, Layer{Name: LayerNames[6], Settings: flags})

	resolved, _ := Merge(candidates)
	if value := pick(resolved); value != nil {
		return *value
	}
	return ""
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
	// Animation is hook-based (REQ-8): the asset carries a class and the host
	// stylesheet owns the motion, which is also what makes prefers-reduced-motion
	// the host's to honour. A standalone asset has no host stylesheet, so
	// enabling animation there asks for motion nothing can deliver — and worse,
	// motion nothing could switch off.
	if settings.Animation != nil && settings.Animation.Enabled != nil && *settings.Animation.Enabled &&
		settings.Delivery != nil && *settings.Delivery == "standalone" {
		add("animation.enabled", "animation is hook-based and needs a host stylesheet; use delivery: themed-inline, "+
			"or leave animation off for a standalone asset")
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
