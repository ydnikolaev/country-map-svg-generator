package main

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/yuranikolaev/country-map-svg-generator/internal/catalog"
	"github.com/yuranikolaev/country-map-svg-generator/internal/config"
	"github.com/yuranikolaev/country-map-svg-generator/internal/geometry"
	"github.com/yuranikolaev/country-map-svg-generator/internal/render"
)

// InspectReport answers "what does this entity actually resolve to" without
// writing anything. It reports what was *selected and drawn*, not what was
// requested: the difference between the two is where every framing and budget
// surprise in this epic has lived.
type InspectReport struct {
	ISO      string `json:"iso"`
	Name     string `json:"name"`
	Alpha3   string `json:"alpha3"`
	Profile  string `json:"profile"`
	Boundary string `json:"boundary"`

	Tier       string `json:"selected_tier"`
	ViewBox    string `json:"view_box"`
	PathBytes  int    `json:"path_bytes"`
	FileBytes  int    `json:"file_bytes"`
	Budget     int    `json:"path_budget"`
	Points     int    `json:"output_points"`
	Components int    `json:"components_drawn"`
	// Removals is what the silhouette does not draw. It is reported because an
	// entity that drops components is framed for more than it shows, which is
	// the shape of WKI-37F18A2AA6A5.
	Removals int `json:"components_removed"`

	Capitals []CapitalFact `json:"capitals"`
	// Diagnostics are geometry's own warnings, passed through rather than
	// summarized: they name things no metric here would.
	Diagnostics []string `json:"diagnostics,omitempty"`
	NoArtifact  string   `json:"no_artifact,omitempty"`
}

type CapitalFact struct {
	ID      string `json:"id"`
	X       string `json:"x"`
	Y       string `json:"y"`
	Anomaly string `json:"anomaly,omitempty"`
}

func newInspectCommand(flags *globalFlags) *cobra.Command {
	local := &configFlags{}
	cmd := &cobra.Command{
		Use:   "inspect",
		Short: "Report what one entity resolves to, without writing anything",
		Long: strings.TrimSpace(`
Report the geometry and marker facts for one entity under a resolved
configuration.

It reports what was selected and drawn rather than what was requested. That
distinction is where the surprises live: an entity can be framed for components
its silhouette does not draw, and a request can fall off the committed detail
ladder onto full-detail source geometry without saying so.

Nothing is written. Use explain for where a configuration value came from, and
inspect for what the pipeline did with it.`),
		Args: exactArgs(0, CLIName+" inspect --iso <code> [--config <path>] [--json]"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if local.iso == "" {
				return failf(ExitUsage, "missing_iso", "inspect reports one entity; pass --iso").
					withContext("hint", "run "+CLIName+" generate --dry-run to check a whole selection")
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
			settings := resolved.Settings

			request, err := render.GeometryRequest(corpus, iso, settings)
			if err != nil {
				return failf(ExitRender, "request_unbuildable", "%v", err)
			}
			applied, err := applyGeometryPreset(request)
			if err != nil {
				return err
			}

			report := InspectReport{
				ISO: iso, Name: entityName(corpus, iso), Alpha3: entityAlpha3(corpus, iso),
				Profile:  valueOrEmpty(settings.Profile),
				Boundary: valueOrEmpty(settings.Boundary),
				Budget:   applied.MaxPathBytes,
			}

			result, genErr := geometry.Generate(request)
			if genErr != nil {
				if geometry.IsNoArtifact(genErr) {
					// A typed absence is a fact about the entity, not a failure of
					// the command: inspect exists to report exactly this.
					report.NoArtifact = genErr.Error()
					return emit(cmd, flags, report, humanInspect(report))
				}
				return failf(ExitRender, "generation_failed", "%s: %v", iso, genErr)
			}

			document, err := render.SVG(result, iso, report.Name, settings)
			if err != nil {
				return failf(ExitRender, "serialization_failed", "%v", err)
			}
			report.Tier = result.LOD.SelectedTier
			report.ViewBox = fmt.Sprintf("%s x %s",
				render.Number(result.ViewBox.Width()), render.Number(result.ViewBox.Height()))
			report.PathBytes = result.Metrics.PathBytes
			report.FileBytes = len(document)
			report.Points = result.Metrics.OutputPoints
			report.Components = len(result.Components)
			report.Removals = len(result.Removals)
			for _, marker := range result.Markers {
				report.Capitals = append(report.Capitals, CapitalFact{
					ID: marker.ID, X: render.Number(marker.X), Y: render.Number(marker.Y),
					Anomaly: marker.Anomaly,
				})
			}
			for _, diagnostic := range result.Diagnostics {
				report.Diagnostics = append(report.Diagnostics,
					diagnostic.Severity+" "+diagnostic.Code+": "+diagnostic.Message)
			}
			return emit(cmd, flags, report, humanInspect(report))
		},
	}
	local.bind(cmd)
	return cmd
}

func humanInspect(report InspectReport) string {
	var out strings.Builder
	fmt.Fprintf(&out, "%s (%s) — %s\n", report.ISO, report.Alpha3, report.Name)
	fmt.Fprintf(&out, "profile %s, boundary %s\n\n", report.Profile, report.Boundary)

	if report.NoArtifact != "" {
		// DEC-009's outcome is stated as a decision, not as an error, because
		// that is what it is: the acceptance oracle refused every rung.
		fmt.Fprintf(&out, "no artifact: %s\n", report.NoArtifact)
		fmt.Fprint(&out, "This is a recorded outcome, not a failure — no rung satisfied the acceptance oracle.")
		return out.String()
	}

	writer := tabwriter.NewWriter(&out, 0, 0, 2, ' ', 0)
	fmt.Fprintf(writer, "tier\t%s\n", report.Tier)
	fmt.Fprintf(writer, "viewBox\t%s\n", report.ViewBox)
	fmt.Fprintf(writer, "path bytes\t%d of %d\n", report.PathBytes, report.Budget)
	fmt.Fprintf(writer, "file bytes\t%d\n", report.FileBytes)
	fmt.Fprintf(writer, "points\t%d\n", report.Points)
	fmt.Fprintf(writer, "components drawn\t%d\n", report.Components)
	fmt.Fprintf(writer, "components removed\t%d\n", report.Removals)
	writer.Flush()

	if len(report.Capitals) != 0 {
		fmt.Fprintf(&out, "\ncapitals (%d):\n", len(report.Capitals))
		for _, capital := range report.Capitals {
			fmt.Fprintf(&out, "  %s at %s,%s", capital.ID, capital.X, capital.Y)
			if capital.Anomaly != "" {
				fmt.Fprintf(&out, " (%s)", capital.Anomaly)
			}
			out.WriteByte('\n')
		}
	}
	for _, diagnostic := range report.Diagnostics {
		fmt.Fprintf(&out, "\n%s", diagnostic)
	}
	return strings.TrimRight(out.String(), "\n")
}

func entityAlpha3(corpus *catalog.Corpus, iso string) string {
	for _, entity := range corpus.Manifest.Entities {
		if entity.Alpha2 == iso {
			return entity.Alpha3
		}
	}
	return ""
}
