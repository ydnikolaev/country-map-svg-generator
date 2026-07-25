package render

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/yuranikolaev/country-map-svg-generator/internal/config"
	"github.com/yuranikolaev/country-map-svg-generator/internal/geometry"
)

// HookPrefix versions the CSS surface CTR-005 defines. Class names and custom
// properties are a public compatibility surface: removing one needs a schema
// major, so the prefix is spelled once and every hook derives from it.
const HookPrefix = "country-map"

// CSS hooks, spelled here and nowhere else.
const (
	RootClass   = HookPrefix
	ShapeClass  = HookPrefix + "__shape"
	MarkerClass = HookPrefix + "__marker"
	// AnimatedClass is added only when animation is enabled, so a host
	// stylesheet can scope motion without guessing.
	AnimatedClass = HookPrefix + "--animated"
)

// customProperty renders the versioned custom-property name for a token.
func customProperty(token string) string { return "--" + HookPrefix + "-" + token }

// SVG serializes one geometry result into a complete SVG document.
//
// The serializer is deliberately small and knows nothing about style names: a
// style has already become tokens by the time it gets here, which is what REQ-4's
// "resolve through tokens rather than renderer forks" means in practice.
func SVG(result geometry.Result, iso, name string, settings config.Settings) ([]byte, error) {
	delivery := valueOr(settings.Delivery, "standalone")
	paint, err := resolvePaint(settings, delivery)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", iso, err)
	}

	var out strings.Builder
	root := element{name: "svg"}
	root.set("xmlns", "http://www.w3.org/2000/svg")
	root.set("viewBox", viewBox(result.ViewBox))

	if delivery == deliveryThemedInline {
		// The hooks exist only in themed-inline. A standalone asset is complete
		// on its own, and emitting classes nothing can target would be exactly
		// the redundant markup REQ-9 forbids.
		classes := []string{RootClass}
		if animationEnabled(settings) {
			classes = append(classes, AnimatedClass)
			if hook := settings.Animation.Hook; hook != nil {
				classes = append(classes, *hook)
			}
		}
		root.set("class", strings.Join(classes, " "))
		root.set("data-country", iso)
		if settings.Profile != nil {
			root.set("data-profile", *settings.Profile)
		}
		if settings.Boundary != nil {
			root.set("data-boundary", *settings.Boundary)
		}
	}
	applyAccessibility(&root, settings, name)

	shape := element{name: "path"}
	if delivery == deliveryThemedInline {
		shape.set("class", ShapeClass)
	}
	shape.set("d", result.Path)
	paint.applyShape(&shape)

	root.open(&out)
	shape.selfClose(&out)
	for _, marker := range markersFor(result, settings) {
		circle := element{name: "circle"}
		if delivery == deliveryThemedInline {
			circle.set("class", MarkerClass)
			circle.set("data-marker", marker.ID)
		}
		circle.set("cx", number(marker.X))
		circle.set("cy", number(marker.Y))
		circle.set("r", number(paint.markerRadius))
		paint.applyMarker(&circle)
		circle.selfClose(&out)
	}
	root.close(&out)
	return []byte(out.String()), nil
}

const (
	deliveryStandalone   = "standalone"
	deliveryThemedInline = "themed-inline"
)

// paint is the resolved presentation, already reduced to attribute values.
type paint struct {
	themed        bool
	fill          string
	fillOpacity   string
	stroke        string
	strokeOpacity string
	strokeWidth   string
	lineCap       string
	lineJoin      string
	markerFill    string
	markerStroke  string
	markerRadius  float64
}

func resolvePaint(settings config.Settings, delivery string) (paint, error) {
	tokens := config.Tokens{}
	if settings.Tokens != nil {
		tokens = *settings.Tokens
	}
	p := paint{themed: delivery == deliveryThemedInline, markerRadius: valueOrFloat(tokens.MarkerRadius, 2)}

	for _, binding := range []struct {
		token  string
		value  *string
		target *string
	}{
		{"fill", tokens.Fill, &p.fill},
		{"stroke", tokens.Stroke, &p.stroke},
		{"marker-fill", tokens.MarkerFill, &p.markerFill},
		{"marker-stroke", tokens.MarkerStroke, &p.markerStroke},
	} {
		resolved, err := paintValue(binding.token, binding.value, delivery)
		if err != nil {
			return paint{}, err
		}
		*binding.target = resolved
	}
	// A marker with no colour of its own inherits the shape's, so a config that
	// only turns markers on still produces something visible.
	if p.markerFill == "" {
		p.markerFill = p.fill
	}

	for _, binding := range []struct {
		token    string
		value    *float64
		fallback float64
		target   *string
	}{
		{"fill-opacity", tokens.FillOpacity, 1, &p.fillOpacity},
		{"stroke-opacity", tokens.StrokeOpacity, 1, &p.strokeOpacity},
		{"stroke-width", tokens.StrokeWidth, 1, &p.strokeWidth},
	} {
		*binding.target = numericValue(binding.token, binding.value, binding.fallback, delivery)
	}
	p.lineCap = keywordValue("line-cap", tokens.LineCap, "butt", delivery)
	p.lineJoin = keywordValue("line-join", tokens.LineJoin, "miter", delivery)

	if advanced := tokens.Advanced; advanced != nil && advanced.Gradient != nil {
		// An advanced reference replaces the fill outright; it has already been
		// bounded to a local fragment by validation.
		p.fill = *advanced.Gradient
	}
	return p, nil
}

// paintValue turns one colour token into an attribute value.
//
// In themed-inline it becomes `var(--country-map-<token>, <resolved>)`: the host
// page can set the property, and the resolved value is the fallback, which is
// REQ-6's "usable fallbacks, including currentColor".
//
// In standalone it must not depend on host CSS at all — ARCH-INV-6 says a
// standalone asset is visually complete without it. A `var()` the author wrote
// therefore collapses to its own fallback, and a `var()` with no fallback is
// refused rather than emitted as something that renders as nothing.
func paintValue(token string, value *string, delivery string) (string, error) {
	resolved := valueOr(value, "currentColor")
	if delivery == deliveryStandalone {
		return bakeVar(token, resolved)
	}
	return fmt.Sprintf("var(%s, %s)", customProperty(token), resolved), nil
}

func numericValue(token string, value *float64, fallback float64, delivery string) string {
	resolved := number(valueOrFloat(value, fallback))
	if delivery == deliveryStandalone {
		return resolved
	}
	return fmt.Sprintf("var(%s, %s)", customProperty(token), resolved)
}

func keywordValue(token string, value *string, fallback, delivery string) string {
	resolved := valueOr(value, fallback)
	if delivery == deliveryStandalone {
		return resolved
	}
	return fmt.Sprintf("var(%s, %s)", customProperty(token), resolved)
}

// bakeVar collapses a var() reference to its fallback for a standalone asset.
func bakeVar(token, value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if !strings.HasPrefix(trimmed, "var(") || !strings.HasSuffix(trimmed, ")") {
		return trimmed, nil
	}
	inner := strings.TrimSuffix(strings.TrimPrefix(trimmed, "var("), ")")
	_, fallback, hasFallback := strings.Cut(inner, ",")
	if !hasFallback || strings.TrimSpace(fallback) == "" {
		return "", fmt.Errorf(
			"tokens.%s is %q, which needs host CSS to mean anything; a standalone asset must be complete on its own, "+
				"so give the variable a fallback (var(--x, currentColor)) or choose delivery: themed-inline",
			token, trimmed)
	}
	// A nested var() in the fallback has the same problem one level down.
	return bakeVar(token, strings.TrimSpace(fallback))
}

func (p paint) applyShape(target *element) {
	target.set("fill", p.fill)
	// A default-valued attribute is redundant markup (REQ-9), so anything that
	// matches the SVG default is left out. In themed-inline every value is a
	// var() reference and none of them can be omitted, because omitting one
	// would remove the hook the host page targets.
	//
	// An opacity on a fill that is `none` is the same redundancy one step
	// further: it describes something that is not drawn.
	if !p.paintIsNone(p.fill) {
		p.setUnlessDefault(target, "fill-opacity", p.fillOpacity, "1")
	}
	target.set("stroke", p.stroke)
	if !p.paintIsNone(p.stroke) {
		p.setUnlessDefault(target, "stroke-opacity", p.strokeOpacity, "1")
		p.setUnlessDefault(target, "stroke-width", p.strokeWidth, "1")
		p.setUnlessDefault(target, "stroke-linecap", p.lineCap, "butt")
		p.setUnlessDefault(target, "stroke-linejoin", p.lineJoin, "miter")
	}
}

func (p paint) applyMarker(target *element) {
	target.set("fill", p.markerFill)
	if p.markerStroke != "" && !p.paintIsNone(p.markerStroke) {
		target.set("stroke", p.markerStroke)
	}
}

// paintIsNone reports whether a paint attribute resolves to nothing drawn, so
// the attributes that only describe a stroke can be omitted with it.
func (p paint) paintIsNone(value string) bool {
	if p.themed {
		// The host can always set the property to something visible, so a themed
		// asset keeps its stroke attributes regardless of the fallback.
		return false
	}
	return value == "none"
}

func (p paint) setUnlessDefault(target *element, name, value, defaultValue string) {
	if !p.themed && value == defaultValue {
		return
	}
	target.set(name, value)
}

func applyAccessibility(root *element, settings config.Settings, name string) {
	mode := "decorative"
	label := name
	if settings.Accessibility != nil {
		mode = valueOr(settings.Accessibility.Mode, mode)
		label = valueOr(settings.Accessibility.Label, label)
	}
	if mode == "labelled" {
		root.set("role", "img")
		root.set("aria-label", label)
		return
	}
	// Decorative is the default: a card beside text that already names the
	// country should not be announced again. `focusable` is set because some
	// browsers put an inline SVG in the tab order otherwise, which is a keyboard
	// trap on a decoration.
	root.set("aria-hidden", "true")
	root.set("focusable", "false")
}

func markersFor(result geometry.Result, settings config.Settings) []geometry.Marker {
	if settings.Marker == nil || settings.Marker.Mode == nil {
		return nil
	}
	switch *settings.Marker.Mode {
	case "none":
		return nil
	case "all-capitals":
		return result.Markers
	case "capital":
		if len(result.Markers) == 0 {
			return nil
		}
		// The corpus orders capitals by role, so the first is the primary one.
		return result.Markers[:1]
	case "custom":
		wanted := map[string]bool{}
		for _, id := range settings.Marker.Custom {
			wanted[id] = true
		}
		var out []geometry.Marker
		for _, marker := range result.Markers {
			if wanted[marker.ID] {
				out = append(out, marker)
			}
		}
		return out
	}
	return nil
}

func animationEnabled(settings config.Settings) bool {
	return settings.Animation != nil && settings.Animation.Enabled != nil && *settings.Animation.Enabled
}

func viewBox(bounds geometry.Bounds) string {
	return strings.Join([]string{
		number(bounds.MinX), number(bounds.MinY), number(bounds.Width()), number(bounds.Height()),
	}, " ")
}

// number renders a value at the output grid's own precision.
//
// Determinism (REQ-9, ARCH-INV-4) comes from this being the only float
// formatter in the package. The rounding is a byte-budget decision found by
// looking at real output: a tight layout derives its viewBox height from a
// scaled bound, which produced `viewBox="0 0 144 122.9026241596183"` — thirteen
// decimals of precision on a document whose path coordinates are quantized to
// q=0.01. That is bytes in every asset for a distinction no renderer can draw
// and no reader can see.
//
// The precision matches the quantization the geometry pipeline already applies,
// so the frame is stated no more precisely than the shape inside it.
const outputDecimals = 2

// Number is the package's one float formatter, exported so a caller describing
// an asset reports the same value the asset carries. The manifest quoting a
// dimension the viewBox does not have is worse than either being wrong: a
// consumer laying out a grid from the manifest would reserve a box the file does
// not fill.
func Number(value float64) string { return number(value) }

func number(value float64) string {
	rounded := math.Round(value*100) / 100
	// Round(-0) is -0, which formats as "-0" and would make two identical
	// documents differ by a sign.
	if rounded == 0 {
		rounded = 0
	}
	return strconv.FormatFloat(rounded, 'f', -1, 64)
}

func valueOr(value *string, fallback string) string {
	if value == nil {
		return fallback
	}
	return *value
}

func valueOrFloat(value *float64, fallback float64) float64 {
	if value == nil {
		return fallback
	}
	return *value
}
