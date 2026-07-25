// Package config owns BND-005: the CTR-3 configuration schema, its embedded
// presets, inheritance, overrides and diagnostics. It never mutates geometry —
// it resolves what to ask for, and hands that to internal/geometry through one
// adapter.
package config

// SchemaVersion is the only accepted `schema` value for v1. It is checked
// exactly rather than by prefix: a document written for a future major must be
// refused by a v1 binary, not silently reinterpreted.
const SchemaVersion = "country-map/v1"

// ReservedKeys are v1 key names that no field may claim.
//
// `selection` is reserved by DEC-015, which moved the crop and
// component-selection seam to a successor spec. Reserving the name costs one
// line here and buys an additive arrival later instead of a schema major and a
// migration for every config already written. A document using it today fails as
// an unknown field like any other, which is exactly the intended behaviour.
var ReservedKeys = []string{"selection"}

// Settings is everything that can be set at any layer of the precedence chain.
// Every field is a pointer or a map so that "unset" is distinguishable from
// "set to the zero value" — without that distinction a later layer could not
// tell whether an earlier one had an opinion, and `explain` could not name the
// origin of a resolved value.
type Settings struct {
	// Profile selects the detail preset: the values geometry.Presets() exposes,
	// `card` and `hero` today. Note the vocabulary crossing, which is deliberate
	// and documented in P3-CHECKPOINT.md: this maps to geometry.Input.Preset,
	// while Boundary below maps to geometry.Input.Profile.
	Profile *string `json:"profile,omitempty" yaml:"profile,omitempty"`
	// Boundary selects the boundary posture governed by DEC-002: `un` (default)
	// or `de_facto`. Spelled with an underscore, matching what geometry accepts.
	Boundary  *string    `json:"boundary,omitempty" yaml:"boundary,omitempty"`
	Delivery  *string    `json:"delivery,omitempty" yaml:"delivery,omitempty"`
	Style     *string    `json:"style,omitempty" yaml:"style,omitempty"`
	Layout    *Layout    `json:"layout,omitempty" yaml:"layout,omitempty"`
	Tokens    *Tokens    `json:"tokens,omitempty" yaml:"tokens,omitempty"`
	Marker    *Marker    `json:"marker,omitempty" yaml:"marker,omitempty"`
	Animation *Animation `json:"animation,omitempty" yaml:"animation,omitempty"`
	// Accessibility is AC-5's "accessible/decorative metadata modes serialize
	// predictably". It is a configuration concern rather than a rendering one
	// because whether a map carries meaning is something only the author knows.
	Accessibility *Accessibility `json:"accessibility,omitempty" yaml:"accessibility,omitempty"`
	Output        *Output        `json:"output,omitempty" yaml:"output,omitempty"`
}

// Accessibility decides what assistive technology is told about the map.
//
// `decorative` is the default because the accepted site composition is card
// decorations beside text that already says which country it is, and ARCH-001's
// accessibility concern asks for no semantic reliance on a decorative map.
// Announcing every card would make the page worse, not better.
type Accessibility struct {
	Mode *string `json:"mode,omitempty" yaml:"mode,omitempty"`
	// Label overrides the entity name announced in `labelled` mode.
	Label *string `json:"label,omitempty" yaml:"label,omitempty"`
}

// Document is one configuration file: the settings it carries plus the two keys
// that only make sense at file scope.
//
// Settings is embedded anonymously and untagged so encoding/json flattens it
// into the document's own keys, which is what makes `profile:` a top-level key
// rather than nested under `settings:`.
type Document struct {
	Schema  string  `json:"schema" yaml:"schema"`
	Extends *string `json:"extends,omitempty" yaml:"extends,omitempty"`
	Settings
	// Profiles carries per-detail-profile overrides, applied when the resolved
	// profile matches the key.
	Profiles map[string]Settings `json:"profiles,omitempty" yaml:"profiles,omitempty"`
	// Countries carries per-entity overrides keyed by ISO alpha-2.
	Countries map[string]Settings `json:"countries,omitempty" yaml:"countries,omitempty"`
}

// Layout expresses REQ-13's contract. The two modes accept disjoint field sets,
// and that disjointness is enforced rather than documented: a `tight` layout
// carrying a width is a mistake about which mode the caller wanted, and
// silently ignoring the width would produce a correct-looking card at the wrong
// size.
type Layout struct {
	// Mode is `tight` (natural aspect derived from the geometry) or `contain`
	// (an explicit frame, one uniform scale, unused space centred).
	Mode *string `json:"mode,omitempty" yaml:"mode,omitempty"`
	// LongSide sizes a tight layout by its rendered long edge.
	LongSide *float64 `json:"longSide,omitempty" yaml:"longSide,omitempty"`
	// MaxWidth and MaxHeight size a tight layout by a bounding box it must fit
	// inside, preserving natural proportions.
	MaxWidth  *float64 `json:"maxWidth,omitempty" yaml:"maxWidth,omitempty"`
	MaxHeight *float64 `json:"maxHeight,omitempty" yaml:"maxHeight,omitempty"`
	// Width and Height are the contain frame.
	Width  *float64 `json:"width,omitempty" yaml:"width,omitempty"`
	Height *float64 `json:"height,omitempty" yaml:"height,omitempty"`
	// Padding is uniform when a number, per-side when an object.
	Padding *Padding `json:"padding,omitempty" yaml:"padding,omitempty"`
}

// Padding accepts either a single number or four sides. Both spellings are kept
// because a uniform padding is overwhelmingly the common case and forcing four
// keys for it would make every config noisier to read.
type Padding struct {
	Top    *float64 `json:"top,omitempty" yaml:"top,omitempty"`
	Right  *float64 `json:"right,omitempty" yaml:"right,omitempty"`
	Bottom *float64 `json:"bottom,omitempty" yaml:"bottom,omitempty"`
	Left   *float64 `json:"left,omitempty" yaml:"left,omitempty"`
	// Uniform is set when the document wrote a bare number. It is mutually
	// exclusive with the four sides.
	Uniform *float64 `json:"-" yaml:"-"`
}

// Tokens is REQ-5's presentation vocabulary. Every style resolves through these
// rather than through a renderer fork, which is what keeps ARCH-INV-3's
// geometry/presentation separation real: adding a style adds token defaults, not
// a code path.
type Tokens struct {
	Fill          *string  `json:"fill,omitempty" yaml:"fill,omitempty"`
	FillOpacity   *float64 `json:"fillOpacity,omitempty" yaml:"fillOpacity,omitempty"`
	Stroke        *string  `json:"stroke,omitempty" yaml:"stroke,omitempty"`
	StrokeOpacity *float64 `json:"strokeOpacity,omitempty" yaml:"strokeOpacity,omitempty"`
	StrokeWidth   *float64 `json:"strokeWidth,omitempty" yaml:"strokeWidth,omitempty"`
	LineCap       *string  `json:"lineCap,omitempty" yaml:"lineCap,omitempty"`
	LineJoin      *string  `json:"lineJoin,omitempty" yaml:"lineJoin,omitempty"`
	MarkerFill    *string  `json:"markerFill,omitempty" yaml:"markerFill,omitempty"`
	MarkerStroke  *string  `json:"markerStroke,omitempty" yaml:"markerStroke,omitempty"`
	MarkerRadius  *float64 `json:"markerRadius,omitempty" yaml:"markerRadius,omitempty"`
	// Advanced carries the optional gradient, pattern and filter references
	// REQ-5 allows "under explicit safe bounds". The bound is that each must be a
	// local fragment reference into the document's own defs — never an external
	// URL, which REQ-9 forbids outright.
	Advanced *AdvancedTokens `json:"advanced,omitempty" yaml:"advanced,omitempty"`
}

type AdvancedTokens struct {
	Gradient *string `json:"gradient,omitempty" yaml:"gradient,omitempty"`
	Pattern  *string `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	Filter   *string `json:"filter,omitempty" yaml:"filter,omitempty"`
}

// Marker is REQ-7. The default is `none`: a decorative card should not sprout a
// capital dot because the corpus happens to know one.
type Marker struct {
	Mode *string `json:"mode,omitempty" yaml:"mode,omitempty"`
	// Custom names capital IDs from the corpus registry when Mode is `custom`.
	// It selects from the registry rather than accepting coordinates, so a
	// marker can never disagree with the corpus about where a place is.
	Custom []string `json:"custom,omitempty" yaml:"custom,omitempty"`
}

// Animation is REQ-8: opt-in, subtle, transform/opacity only, hook-based.
type Animation struct {
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	// Hook names the CSS class the animation attaches to, so motion lives in the
	// host stylesheet rather than in generated markup.
	Hook *string `json:"hook,omitempty" yaml:"hook,omitempty"`
}

type Output struct {
	Dir *string `json:"dir,omitempty" yaml:"dir,omitempty"`
	// Filename is a template over a fixed, validated placeholder set. It is not
	// an arbitrary format string: an unknown placeholder fails validation rather
	// than appearing literally in a file name.
	Filename *string `json:"filename,omitempty" yaml:"filename,omitempty"`
}

// The closed vocabularies. Style and delivery are enumerated here because they
// are this package's own contract; profile and boundary are deliberately absent
// — those are enumerated from geometry and the corpus at validation time, so a
// new detail preset appears in the CLI without a code edit here.
var (
	Deliveries  = []string{"standalone", "themed-inline"}
	StyleNames  = []string{"outline", "filled", "bold-soft", "silhouette", "ghost"}
	LayoutModes = []string{"tight", "contain"}
	MarkerModes = []string{"none", "capital", "all-capitals", "custom"}
	// AccessibilityModes: `decorative` hides the map from assistive technology,
	// `labelled` announces it as an image with a name.
	AccessibilityModes = []string{"decorative", "labelled"}
	LineCaps           = []string{"butt", "round", "square"}
	LineJoins          = []string{"miter", "round", "bevel"}
)

// FilenamePlaceholders is the closed set a filename template may use.
var FilenamePlaceholders = []string{"iso", "profile", "boundary", "style", "delivery"}
