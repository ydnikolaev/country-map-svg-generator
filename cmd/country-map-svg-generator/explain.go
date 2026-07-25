package main

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/ydnikolaev/country-map-svg-generator/internal/config"
	"github.com/ydnikolaev/country-map-svg-generator/internal/render"
)

// ExplainReport answers "what will this produce, and why". Every resolved value
// carries the layer that set it, because the question an author actually has is
// not "what is the fill" but "which of my six layers set the fill".
type ExplainReport struct {
	Config string `json:"config,omitempty"`
	ISO    string `json:"iso"`
	Name   string `json:"name"`
	// Layers is the precedence chain in order, so the report explains the model
	// as well as the outcome.
	Layers []string `json:"layers"`
	Values []Value  `json:"values"`
	// Geometry is what the resolved configuration actually asks geometry for.
	// It is reported because the two axes cross on the way in — the
	// configuration's profile is geometry's preset and its boundary is
	// geometry's profile — and an author debugging a wrong-looking map needs to
	// see the request that was really made.
	Geometry GeometryRequestView `json:"geometry_request"`
	Output   string              `json:"output_path"`
}

type Value struct {
	Path  string `json:"path"`
	Value string `json:"value"`
	Layer string `json:"layer"`
}

type GeometryRequestView struct {
	Preset       string `json:"preset"`
	Profile      string `json:"profile"`
	LayoutMode   string `json:"layout_mode"`
	LongSide     string `json:"long_side,omitempty"`
	Width        string `json:"width,omitempty"`
	Height       string `json:"height,omitempty"`
	MaxPathBytes int    `json:"max_path_bytes"`
}

func newExplainCommand(flags *globalFlags) *cobra.Command {
	local := &configFlags{}
	cmd := &cobra.Command{
		Use:   "explain",
		Short: "Show the resolved configuration and where each value came from",
		Long: strings.TrimSpace(`
Resolve one entity's configuration and report every value with its origin layer.

The precedence chain is: embedded defaults, the preset ancestry named by
extends, the document's own settings, the block for the resolved profile, the
per-country override, and finally the command-line flags. Each preset in a chain
is named separately, so "which preset set this" has an answer.

Nothing is written. This is the command to run before generate, and the one to
run when generate produced something unexpected.`),
		Args: exactArgs(0, CLIName+" explain --iso <code> [--config <path>] [--json]"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if local.iso == "" {
				return failf(ExitUsage, "missing_iso", "explain resolves one entity; pass --iso").
					withContext("hint", "run "+CLIName+" validate to check the whole selection instead")
			}
			doc, err := local.loadDocument()
			if err != nil {
				return err
			}
			vocab, corpus, err := vocabulary()
			if err != nil {
				return err
			}
			if doc != nil {
				if err := config.ValidateDocument(doc, vocab); err != nil {
					return configError(err)
				}
			}
			isos, err := selection(local, corpus, vocab)
			if err != nil {
				return err
			}
			iso := isos[0]

			resolved, err := config.Resolve(doc, iso, local.overrides())
			if err != nil {
				return configError(err)
			}
			if err := config.ValidateResolved(iso, resolved, vocab); err != nil {
				return configError(err)
			}
			request, err := render.GeometryRequest(corpus, iso, resolved.Settings)
			if err != nil {
				return failf(ExitRender, "request_unbuildable", "%v", err)
			}
			applied, err := applyGeometryPreset(request)
			if err != nil {
				return err
			}
			outputPath, err := config.OutputPath(iso, resolved.Settings)
			if err != nil {
				return configError(err)
			}

			report := ExplainReport{
				Config: local.path, ISO: iso, Name: entityName(corpus, iso),
				Layers: config.LayerNames,
				Values: valuesOf(resolved),
				Geometry: GeometryRequestView{
					Preset: applied.Preset, Profile: applied.Profile,
					LayoutMode:   string(applied.Layout.Mode),
					LongSide:     optionalNumber(applied.Layout.LongSide),
					Width:        optionalNumber(applied.Layout.Width),
					Height:       optionalNumber(applied.Layout.Height),
					MaxPathBytes: applied.MaxPathBytes,
				},
				Output: outputPath,
			}
			return emit(cmd, flags, report, humanExplain(report))
		},
	}
	local.bind(cmd)
	return cmd
}

func valuesOf(resolved config.Resolved) []Value {
	origins := resolved.Provenance.Origins()
	flat := flattenSettings(resolved.Settings)
	out := make([]Value, 0, len(origins))
	for _, origin := range origins {
		value, ok := flat[origin.Path]
		if !ok {
			continue
		}
		out = append(out, Value{Path: origin.Path, Value: value, Layer: origin.Layer})
	}
	return out
}

func humanExplain(report ExplainReport) string {
	var out strings.Builder
	fmt.Fprintf(&out, "%s — %s\n\n", report.ISO, report.Name)

	writer := tabwriter.NewWriter(&out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(writer, "KEY\tVALUE\tFROM")
	for _, value := range report.Values {
		fmt.Fprintf(writer, "%s\t%s\t%s\n", value.Path, value.Value, value.Layer)
	}
	writer.Flush()

	fmt.Fprintf(&out, "\ngeometry request: preset %s, boundary profile %s, layout %s",
		report.Geometry.Preset, report.Geometry.Profile, report.Geometry.LayoutMode)
	if report.Geometry.LongSide != "" {
		fmt.Fprintf(&out, " long side %s", report.Geometry.LongSide)
	}
	if report.Geometry.Width != "" {
		fmt.Fprintf(&out, " %sx%s", report.Geometry.Width, report.Geometry.Height)
	}
	fmt.Fprintf(&out, ", budget %d bytes\n", report.Geometry.MaxPathBytes)
	fmt.Fprintf(&out, "output: %s", report.Output)
	return out.String()
}

func optionalNumber(value float64) string {
	if value == 0 {
		return ""
	}
	return trimNumber(value)
}
