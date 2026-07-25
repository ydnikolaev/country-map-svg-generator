package config

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

// Vocabulary carries the value sets that this package must not invent. Profiles
// come from the geometry presets, boundaries from the geometry adapter's
// accepted set, and the ISO and capital registries from the corpus. Injecting
// them keeps two promises at once: adding a detail preset needs no edit here,
// and no country literal or magic enum lives in the config layer.
type Vocabulary struct {
	Profiles   []string
	Boundaries []string
	// ISOCodes is the set of entities the corpus actually carries. Empty means
	// "not supplied", and entity-level checks are then skipped rather than
	// failing everything closed — `validate` on a config without a corpus is a
	// legitimate operation.
	ISOCodes []string
	// CapitalIDs maps an ISO code to the capital IDs the corpus knows, used to
	// check a custom marker selection against the registry instead of against a
	// coordinate the author typed.
	CapitalIDs map[string][]string
}

// SafeDimensionCeiling bounds a layout dimension. It is a sanity bound, not a
// product limit: the largest documented use is a 700 px hero (QAB-1), so this
// sits three orders above anything real. Its job is to turn a runaway or
// mistyped value into a diagnostic instead of a viewBox whose numbers have
// stopped meaning CSS pixels.
const SafeDimensionCeiling = 1e6

// ValidationError is one problem with one key path. Diagnostics are per-key so
// an agent can fix them all in one pass, and each carries the path rather than a
// line number, because a resolved value can come from a preset, the document or
// a flag and only the path is meaningful across all three.
type ValidationError struct {
	Path    KeyPath
	Message string
}

func (e ValidationError) Error() string { return e.Path + ": " + e.Message }

// ValidationErrors is every problem found, in stable path order.
type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	parts := make([]string, len(e))
	for i, item := range e {
		parts[i] = item.Error()
	}
	return strings.Join(parts, "; ")
}

// ValidateDocument checks one document's shape and values before any layering.
// Per-layer validation matters: a typo in a preset must be reported against the
// preset, not against whichever document happened to inherit it.
func ValidateDocument(doc *Document, vocab Vocabulary) error {
	var problems []ValidationError
	add := func(path, format string, args ...any) {
		problems = append(problems, ValidationError{Path: path, Message: fmt.Sprintf(format, args...)})
	}

	if doc.Schema != SchemaVersion {
		add("schema", "must be %q, got %q", SchemaVersion, doc.Schema)
	}
	if doc.Extends != nil && strings.TrimSpace(*doc.Extends) == "" {
		add("extends", "must name a preset, not an empty string")
	}

	validateSettings(&doc.Settings, "", vocab, add)

	for name := range doc.Profiles {
		settings := doc.Profiles[name]
		path := "profiles." + name
		if !contains(vocab.Profiles, name) && len(vocab.Profiles) != 0 {
			add(path, "unknown profile; known profiles are %s", strings.Join(sorted(vocab.Profiles), ", "))
		}
		// A profile block that re-selects a profile is a contradiction rather
		// than an override: it would make `profiles.card.profile: hero` mean
		// something no precedence rule can express.
		if settings.Profile != nil {
			add(path+".profile", "a profile block cannot select a different profile")
		}
		validateSettings(&settings, path, vocab, add)
	}

	for iso := range doc.Countries {
		settings := doc.Countries[iso]
		path := "countries." + iso
		if len(vocab.ISOCodes) != 0 && !contains(vocab.ISOCodes, iso) {
			add(path, "unknown ISO alpha-2 code for this corpus")
		}
		validateSettings(&settings, path, vocab, add)
		validateMarkerSelection(&settings, iso, path, vocab, add)
	}

	if len(problems) == 0 {
		return nil
	}
	sort.Slice(problems, func(i, j int) bool {
		if problems[i].Path == problems[j].Path {
			return problems[i].Message < problems[j].Message
		}
		return problems[i].Path < problems[j].Path
	})
	return ValidationErrors(problems)
}

type addFunc func(path, format string, args ...any)

func validateSettings(settings *Settings, prefix string, vocab Vocabulary, add addFunc) {
	at := func(key string) string { return join(prefix, key) }

	if settings.Profile != nil && len(vocab.Profiles) != 0 && !contains(vocab.Profiles, *settings.Profile) {
		add(at("profile"), "unknown profile %q; known profiles are %s", *settings.Profile, strings.Join(sorted(vocab.Profiles), ", "))
	}
	if settings.Boundary != nil && len(vocab.Boundaries) != 0 && !contains(vocab.Boundaries, *settings.Boundary) {
		// The likeliest mistake here is a hyphen: the corpus manifest declares
		// its profile list as "de-facto" while the runtime accepts "de_facto",
		// so an author reading the manifest writes the spelling geometry does
		// not take. Naming the accepted set is what makes that recoverable.
		add(at("boundary"), "unknown boundary %q; accepted values are %s", *settings.Boundary, strings.Join(sorted(vocab.Boundaries), ", "))
	}
	if settings.Delivery != nil && !contains(Deliveries, *settings.Delivery) {
		add(at("delivery"), "unknown delivery %q; accepted values are %s", *settings.Delivery, strings.Join(Deliveries, ", "))
	}
	if settings.Style != nil && !contains(StyleNames, *settings.Style) {
		add(at("style"), "unknown style %q; accepted values are %s", *settings.Style, strings.Join(StyleNames, ", "))
	}
	if settings.Layout != nil {
		validateLayout(settings.Layout, at("layout"), add)
	}
	if settings.Tokens != nil {
		validateTokens(settings.Tokens, at("tokens"), add)
	}
	if settings.Marker != nil && settings.Marker.Mode != nil && !contains(MarkerModes, *settings.Marker.Mode) {
		add(at("marker.mode"), "unknown marker mode %q; accepted values are %s", *settings.Marker.Mode, strings.Join(MarkerModes, ", "))
	}
	if settings.Animation != nil {
		validateAnimation(settings.Animation, at("animation"), add)
	}
	if settings.Accessibility != nil {
		if mode := settings.Accessibility.Mode; mode != nil && !contains(AccessibilityModes, *mode) {
			add(at("accessibility.mode"), "unknown accessibility mode %q; accepted values are %s", *mode, strings.Join(AccessibilityModes, ", "))
		}
		if label := settings.Accessibility.Label; label != nil {
			if strings.TrimSpace(*label) == "" {
				add(at("accessibility.label"), "must not be empty; omit it to use the entity name")
			} else if strings.ContainsAny(*label, "<>&\"") {
				add(at("accessibility.label"), "must not contain markup characters")
			}
		}
	}
	if settings.Output != nil && settings.Output.Filename != nil {
		validateFilename(*settings.Output.Filename, at("output.filename"), add)
	}
}

// validateLayout enforces REQ-13. The two modes take disjoint field sets, and
// the disjointness is enforced rather than tolerated: a `tight` layout carrying
// a width is a mistake about which mode the author wanted, and ignoring the
// width silently would produce a plausible card at the wrong size — the failure
// shape that costs the most to notice.
func validateLayout(layout *Layout, prefix string, add addFunc) {
	at := func(key string) string { return join(prefix, key) }

	for key, value := range map[string]*float64{
		"longSide": layout.LongSide, "maxWidth": layout.MaxWidth, "maxHeight": layout.MaxHeight,
		"width": layout.Width, "height": layout.Height,
	} {
		if value != nil {
			checkDimension(*value, at(key), add)
		}
	}
	if layout.Padding != nil {
		validatePadding(layout.Padding, at("padding"), add)
	}

	if layout.Mode == nil {
		// Absent here is legal: a later layer or a preset supplies it, and the
		// resolved value is checked again by ValidateResolved.
		return
	}
	if !contains(LayoutModes, *layout.Mode) {
		add(at("mode"), "unknown layout mode %q; accepted values are %s", *layout.Mode, strings.Join(LayoutModes, ", "))
		return
	}

	switch *layout.Mode {
	case "tight":
		for _, key := range []string{"width", "height"} {
			if fieldOf(layout, key) != nil {
				add(at(key), "not accepted in tight mode; tight derives natural proportions, use longSide or maxWidth/maxHeight")
			}
		}
		if layout.LongSide != nil && (layout.MaxWidth != nil || layout.MaxHeight != nil) {
			add(at("longSide"), "cannot be combined with maxWidth or maxHeight; they are two ways to size the same layout")
		}
	case "contain":
		for _, key := range []string{"longSide", "maxWidth", "maxHeight"} {
			if fieldOf(layout, key) != nil {
				add(at(key), "not accepted in contain mode; contain takes an explicit width and height")
			}
		}
	}
}

func fieldOf(layout *Layout, key string) *float64 {
	switch key {
	case "longSide":
		return layout.LongSide
	case "maxWidth":
		return layout.MaxWidth
	case "maxHeight":
		return layout.MaxHeight
	case "width":
		return layout.Width
	case "height":
		return layout.Height
	}
	return nil
}

func checkDimension(value float64, path string, add addFunc) {
	switch {
	case math.IsNaN(value):
		add(path, "must be a number")
	case math.IsInf(value, 0):
		add(path, "must be finite")
	case value <= 0:
		add(path, "must be greater than zero, got %g", value)
	case value > SafeDimensionCeiling:
		add(path, "must not exceed %g, got %g", SafeDimensionCeiling, value)
	}
}

func validatePadding(padding *Padding, prefix string, add addFunc) {
	sides := map[string]*float64{
		"top": padding.Top, "right": padding.Right, "bottom": padding.Bottom, "left": padding.Left,
	}
	anySide := false
	for key, value := range sides {
		if value == nil {
			continue
		}
		anySide = true
		checkPadding(*value, join(prefix, key), add)
	}
	if padding.Uniform != nil {
		checkPadding(*padding.Uniform, prefix, add)
		if anySide {
			add(prefix, "is either a single number or per-side values, not both")
		}
	}
}

// checkPadding differs from checkDimension in one place: zero padding is
// legitimate and common, so only negatives are refused.
func checkPadding(value float64, path string, add addFunc) {
	switch {
	case math.IsNaN(value):
		add(path, "must be a number")
	case math.IsInf(value, 0):
		add(path, "must be finite")
	case value < 0:
		add(path, "must not be negative, got %g", value)
	case value > SafeDimensionCeiling:
		add(path, "must not exceed %g, got %g", SafeDimensionCeiling, value)
	}
}

func validateTokens(tokens *Tokens, prefix string, add addFunc) {
	at := func(key string) string { return join(prefix, key) }

	for key, value := range map[string]*float64{
		"fillOpacity": tokens.FillOpacity, "strokeOpacity": tokens.StrokeOpacity,
	} {
		if value == nil {
			continue
		}
		if math.IsNaN(*value) || *value < 0 || *value > 1 {
			add(at(key), "must be between 0 and 1, got %g", *value)
		}
	}
	for key, value := range map[string]*float64{
		"strokeWidth": tokens.StrokeWidth, "markerRadius": tokens.MarkerRadius,
	} {
		if value == nil {
			continue
		}
		switch {
		case math.IsNaN(*value):
			add(at(key), "must be a number")
		case math.IsInf(*value, 0):
			add(at(key), "must be finite")
		case *value < 0:
			add(at(key), "must not be negative, got %g", *value)
		case *value > SafeDimensionCeiling:
			add(at(key), "must not exceed %g, got %g", SafeDimensionCeiling, *value)
		}
	}
	if tokens.LineCap != nil && !contains(LineCaps, *tokens.LineCap) {
		add(at("lineCap"), "unknown line cap %q; accepted values are %s", *tokens.LineCap, strings.Join(LineCaps, ", "))
	}
	if tokens.LineJoin != nil && !contains(LineJoins, *tokens.LineJoin) {
		add(at("lineJoin"), "unknown line join %q; accepted values are %s", *tokens.LineJoin, strings.Join(LineJoins, ", "))
	}
	for key, value := range map[string]*string{
		"fill": tokens.Fill, "stroke": tokens.Stroke,
		"markerFill": tokens.MarkerFill, "markerStroke": tokens.MarkerStroke,
	} {
		if value != nil {
			checkPaintValue(*value, at(key), add)
		}
	}
	if tokens.Advanced != nil {
		for key, value := range map[string]*string{
			"gradient": tokens.Advanced.Gradient, "pattern": tokens.Advanced.Pattern,
			"filter": tokens.Advanced.Filter,
		} {
			if value != nil {
				checkLocalReference(*value, at("advanced."+key), add)
			}
		}
	}
}

// checkPaintValue is where REQ-9's "no external URL" reaches the config layer.
// A paint token is author-supplied text that lands in an SVG attribute, so it is
// validated at this boundary rather than trusted and filtered later: by the time
// a renderer sees it, the context that would make a good diagnostic is gone.
func checkPaintValue(value, path string, add addFunc) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		add(path, "must not be empty")
		return
	}
	if strings.ContainsAny(trimmed, "<>\"'") {
		add(path, "must not contain markup characters")
		return
	}
	if hasExternalReference(trimmed) {
		add(path, "must not reference an external resource; only local fragment references are allowed")
	}
}

// checkLocalReference bounds REQ-5's advanced references. A gradient, pattern or
// filter must point into the document's own defs — url(#id) — so generated
// output stays self-contained and cannot fetch anything at render time.
func checkLocalReference(value, path string, add addFunc) {
	trimmed := strings.TrimSpace(value)
	if !strings.HasPrefix(trimmed, "url(#") || !strings.HasSuffix(trimmed, ")") {
		add(path, "must be a local fragment reference of the form url(#id), got %q", value)
		return
	}
	id := strings.TrimSuffix(strings.TrimPrefix(trimmed, "url(#"), ")")
	if id == "" {
		add(path, "names an empty fragment id")
		return
	}
	for _, r := range id {
		isSafe := r == '-' || r == '_' ||
			(r >= '0' && r <= '9') || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
		if !isSafe {
			add(path, "fragment id %q may contain only letters, digits, hyphen and underscore", id)
			return
		}
	}
}

func hasExternalReference(value string) bool {
	lowered := strings.ToLower(value)
	for _, scheme := range []string{"http://", "https://", "//", "data:", "file:", "javascript:"} {
		if strings.Contains(lowered, scheme) {
			return true
		}
	}
	// url(...) is allowed only in its local fragment form.
	if index := strings.Index(lowered, "url("); index >= 0 {
		return !strings.HasPrefix(lowered[index:], "url(#")
	}
	return false
}

func validateAnimation(animation *Animation, prefix string, add addFunc) {
	if animation.Hook == nil {
		return
	}
	hook := strings.TrimSpace(*animation.Hook)
	if hook == "" {
		add(join(prefix, "hook"), "must name a CSS class, not an empty string")
		return
	}
	// The hook becomes a class attribute value, so it is constrained to a CSS
	// identifier rather than accepted as free text.
	for i, r := range hook {
		isSafe := r == '-' || r == '_' ||
			(r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9' && i > 0)
		if !isSafe {
			add(join(prefix, "hook"), "must be a CSS identifier: letters, digits, hyphen and underscore, not starting with a digit")
			return
		}
	}
}

// validateFilename bounds the output template to a closed placeholder set. It is
// not a format string: an unknown placeholder is a mistake that would otherwise
// appear literally in a file name and be discovered only by looking at the
// output directory.
func validateFilename(template, path string, add addFunc) {
	if strings.TrimSpace(template) == "" {
		add(path, "must not be empty")
		return
	}
	if strings.ContainsAny(template, `/\`) {
		add(path, "must be a file name, not a path; use output.dir for the directory")
	}
	rest := template
	for {
		open := strings.Index(rest, "{")
		if open < 0 {
			break
		}
		close := strings.Index(rest[open:], "}")
		if close < 0 {
			add(path, "has an unclosed placeholder")
			return
		}
		name := rest[open+1 : open+close]
		if !contains(FilenamePlaceholders, name) {
			add(path, "unknown placeholder {%s}; accepted placeholders are %s", name, strings.Join(FilenamePlaceholders, ", "))
		}
		rest = rest[open+close+1:]
	}
	if strings.Contains(rest, "}") {
		add(path, "has a closing brace with no placeholder")
	}
	// A template without {iso} is deliberately not refused here. It is a
	// legitimate way to write a single-entity configuration, and for a batch the
	// output-collision check refuses it with a better diagnostic — one that names
	// which entities actually collided. Requiring {iso} at this layer would make
	// that check, which REQ-3 asks for by name, unreachable.
}

func validateMarkerSelection(settings *Settings, iso, prefix string, vocab Vocabulary, add addFunc) {
	if settings.Marker == nil || len(settings.Marker.Custom) == 0 {
		return
	}
	path := join(prefix, "marker.custom")
	if settings.Marker.Mode == nil || *settings.Marker.Mode != "custom" {
		add(path, "is only meaningful when marker.mode is custom")
	}
	known, ok := vocab.CapitalIDs[iso]
	if !ok {
		return
	}
	for _, id := range settings.Marker.Custom {
		if !contains(known, id) {
			add(path, "unknown capital id %q for %s; known ids are %s", id, iso, strings.Join(sorted(known), ", "))
		}
	}
}

func contains(set []string, value string) bool {
	for _, item := range set {
		if item == value {
			return true
		}
	}
	return false
}

func sorted(values []string) []string {
	out := append([]string(nil), values...)
	sort.Strings(out)
	return out
}
