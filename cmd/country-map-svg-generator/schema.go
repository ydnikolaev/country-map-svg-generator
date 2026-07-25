package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ydnikolaev/country-map-svg-generator/internal/config"
)

// The schema command answers "what may I set, and to what?" — the question every
// other surface answers only by refusing something first.
//
// `explain` reports what a document resolved TO; it cannot report what a key
// ACCEPTS, because a resolved value carries no memory of its vocabulary. An
// operator who wanted the accepted set had to provoke a validation error per
// key and read it out of the diagnostic. That works, but it means the only way
// to learn the surface is to fail against it, which is a poor contract for the
// agent this CLI is built for.
//
// Nothing here is written down twice. Every field below is read from the same
// value the validator checks against: the key paths come from the struct tags
// by reflection, the closed vocabularies from the vars validation compares to,
// and profile and boundary from geometry and the corpus at runtime — which is
// why adding a detail preset shows up here with no edit to this file.
type SchemaInfo struct {
	Schema string `json:"schema"`
	// Keys is every settable dotted path the document schema defines.
	Keys []string `json:"keys"`
	// Vocabularies maps a key to the closed set of values it accepts. A key that
	// is absent takes free-form input bounded by a rule rather than by a set —
	// the paint and reference tokens, whose rules are described in Notes.
	Vocabularies map[string][]string `json:"vocabularies"`
	// FilenamePlaceholders is the closed set output.filename may interpolate.
	FilenamePlaceholders []string `json:"filename_placeholders"`
	// Styles carries each style with the sentence that says what it is for, so a
	// caller choosing between them does not have to render all five.
	Styles []SchemaStyle `json:"styles"`
	// Notes are the semantics a vocabulary cannot carry: what a value means, and
	// which bound applies to a key with no closed set.
	Notes map[string]string `json:"notes"`
}

type SchemaStyle struct {
	Name string `json:"name"`
	Note string `json:"note"`
}

func newSchemaCommand(flags *globalFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "schema",
		Short: "Report every settable key and the values it accepts",
		Long: strings.TrimSpace(`
Report the configuration surface: every key a document may set, the closed set
of values each key accepts, the placeholders a filename template may use, and
what each style is for.

This is the discovery surface. Every other command tells you a value is wrong
after you have written it; this one tells you what the accepted values are
before you do. The output is derived from the same definitions validation checks
against, so it cannot drift from what the binary will actually accept.`),
		Args: exactArgs(0, CLIName+" schema [--json]"),
		RunE: func(cmd *cobra.Command, args []string) error {
			info, err := schemaInfo()
			if err != nil {
				return err
			}
			return emit(cmd, flags, info, humanSchema(info))
		},
	}
}

func schemaInfo() (SchemaInfo, error) {
	// Profile and boundary are not this package's to enumerate: they come from
	// the geometry presets and the corpus, exactly as validation gets them.
	vocab, _, err := vocabulary()
	if err != nil {
		return SchemaInfo{}, err
	}

	styles, err := config.Styles()
	if err != nil {
		return SchemaInfo{}, failf(ExitData, "styles_unreadable", "embedded styles could not be decoded: %v", err)
	}
	named := make([]SchemaStyle, 0, len(styles))
	for _, name := range config.StyleNames {
		style, ok := styles[name]
		if !ok {
			// StyleNames and the embedded file are asserted to agree in both
			// directions elsewhere; if that ever breaks, report the name rather
			// than dropping it silently.
			named = append(named, SchemaStyle{Name: name})
			continue
		}
		named = append(named, SchemaStyle{Name: name, Note: style.Note})
	}

	return SchemaInfo{
		Schema: EnvelopeSchema,
		Keys:   config.SchemaKeys(),
		Vocabularies: map[string][]string{
			"profile":            sorted(vocab.Profiles),
			"boundary":           sorted(vocab.Boundaries),
			"style":              copyOf(config.StyleNames),
			"delivery":           copyOf(config.Deliveries),
			"layout.mode":        copyOf(config.LayoutModes),
			"marker.mode":        copyOf(config.MarkerModes),
			"accessibility.mode": copyOf(config.AccessibilityModes),
			"tokens.lineCap":     copyOf(config.LineCaps),
			"tokens.lineJoin":    copyOf(config.LineJoins),
		},
		FilenamePlaceholders: copyOf(config.FilenamePlaceholders),
		Styles:               named,
		Notes: map[string]string{
			"layout.mode":              "tight derives the viewBox from the projected geometry; contain takes an explicit width and height and centres the unused space. The two modes accept disjoint fields, so a width under tight is refused rather than ignored.",
			"layout.longSide":          "Only the profile's own long side is served from the committed detail ladder. Any other size falls back to full-detail source geometry, which is over the byte ceiling for a large entity.",
			"marker.mode":              "custom selects capital IDs from the corpus registry, never coordinates, so a marker cannot disagree with the corpus about where a place is.",
			"accessibility.mode":       "decorative hides the map from assistive technology; labelled announces it as an image with a name.",
			"delivery":                 "standalone is visually complete with no host CSS; themed-inline emits every value as a var() the host page can set, with the resolved value as the fallback.",
			"tokens":                   "Paint tokens accept any CSS colour, currentColor, or none. An external URL is refused.",
			"tokens.advanced.gradient": "A local fragment reference such as url(#id). It replaces the fill outright.",
			"tokens.advanced.pattern":  "A local fragment reference such as url(#id). It replaces the fill outright, and cannot be combined with gradient — both are paint servers competing for the same attribute.",
			"tokens.advanced.filter":   "A local fragment reference such as url(#id). It is its own attribute, so it composes with whatever fill applies.",
			"output.filename":          "A template over the closed placeholder set. An unknown placeholder is refused rather than emitted literally.",
			"countries":                "A per-entity block keyed by ISO alpha-2. Everything settable at the top level is settable there for one entity.",
			"profiles":                 "A per-profile block keyed by profile name, applied when that profile resolves.",
			"extends":                  "The named embedded preset this document inherits from. Presets form a chain, and explain names each link separately.",
		},
	}, nil
}

func humanSchema(info SchemaInfo) string {
	var out strings.Builder

	out.WriteString("SETTABLE KEYS\n")
	for _, key := range info.Keys {
		if values, ok := info.Vocabularies[key]; ok {
			fmt.Fprintf(&out, "  %-28s %s\n", key, strings.Join(values, " | "))
			continue
		}
		fmt.Fprintf(&out, "  %s\n", key)
	}

	// A vocabulary for a key the schema does not list as a leaf path would
	// otherwise be invisible; print the remainder rather than dropping it.
	listed := make(map[string]bool, len(info.Keys))
	for _, key := range info.Keys {
		listed[key] = true
	}
	extra := make([]string, 0, len(info.Vocabularies))
	for key := range info.Vocabularies {
		if !listed[key] {
			extra = append(extra, key)
		}
	}
	sort.Strings(extra)
	if len(extra) != 0 {
		out.WriteString("\nACCEPTED VALUES\n")
		for _, key := range extra {
			fmt.Fprintf(&out, "  %-28s %s\n", key, strings.Join(info.Vocabularies[key], " | "))
		}
	}

	out.WriteString("\nSTYLES\n")
	for _, style := range info.Styles {
		fmt.Fprintf(&out, "  %-12s %s\n", style.Name, style.Note)
	}

	fmt.Fprintf(&out, "\nFILENAME PLACEHOLDERS\n  %s\n", "{"+strings.Join(info.FilenamePlaceholders, "} {")+"}")

	out.WriteString("\nNOTES\n")
	keys := make([]string, 0, len(info.Notes))
	for key := range info.Notes {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Fprintf(&out, "  %s\n    %s\n", key, info.Notes[key])
	}

	return strings.TrimRight(out.String(), "\n")
}

func sorted(values []string) []string {
	out := copyOf(values)
	sort.Strings(out)
	return out
}

func copyOf(values []string) []string {
	return append([]string(nil), values...)
}
