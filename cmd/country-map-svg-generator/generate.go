package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ydnikolaev/country-map-svg-generator/internal/catalog"
	"github.com/ydnikolaev/country-map-svg-generator/internal/config"
	"github.com/ydnikolaev/country-map-svg-generator/internal/geometry"
	"github.com/ydnikolaev/country-map-svg-generator/internal/render"
)

// GenerateReport is what a caller reads after a batch. It reports the skipped
// entities as prominently as the written ones: a catalog that is quietly short
// by three is the failure mode this whole command is shaped to avoid.
type GenerateReport struct {
	Output     string           `json:"output"`
	Manifest   string           `json:"manifest"`
	Written    int              `json:"written"`
	Skipped    []render.Skipped `json:"skipped"`
	TotalBytes int              `json:"total_bytes"`
	DryRun     bool             `json:"dry_run"`
}

func newGenerateCommand(flags *globalFlags) *cobra.Command {
	local := &configFlags{}
	var (
		out    string
		dryRun bool
	)
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate SVG assets and a manifest",
		Long: strings.TrimSpace(`
Generate assets for one entity, an explicit set, or the whole catalog.

The batch is built and structurally validated in memory before anything touches
the output directory, and it is then published atomically. A failure at any
point leaves the previous output exactly as it was — a half-regenerated catalog
whose manifest agrees with it is indistinguishable from a complete one, and that
is the failure this ordering exists to prevent.

Generation is offline: the corpus, the presets and the schemas are embedded, and
nothing is fetched.`),
		Args: exactArgs(0, CLIName+" generate [--config <path>] [--iso <code>] [--out <dir>] [--json]"),
		RunE: func(cmd *cobra.Command, args []string) error {
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

			batch, outputDir, err := buildBatch(corpus, doc, local, vocab, isos, out)
			if err != nil {
				return err
			}

			report := GenerateReport{
				Output: outputDir, Manifest: outputDir + "/" + render.ManifestName,
				Written: len(batch.Assets), Skipped: batch.Manifest.Skipped,
				TotalBytes: batch.Manifest.TotalBytes, DryRun: dryRun,
			}
			if dryRun {
				return emit(cmd, flags, report, humanGenerate(report))
			}
			if err := render.Publish(batch, outputDir); err != nil {
				var publishErr *render.PublishError
				if errors.As(err, &publishErr) {
					return failf(ExitFilesystem, "publish_failed", "%s", publishErr.Error()).
						withContext("output", outputDir,
							"previous_output", boolWord(publishErr.PreviousIntact, "unchanged", "possibly inconsistent"))
				}
				return failf(ExitFilesystem, "publish_failed", "%v", err)
			}
			return emit(cmd, flags, report, humanGenerate(report))
		},
	}
	local.bind(cmd)
	cmd.Flags().StringVar(&out, "out", "", "output directory (default: output.dir from the configuration)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "build and validate the batch without writing anything")
	return cmd
}

// buildBatch resolves, renders and validates every selected entity. Nothing is
// written here: the whole batch has to succeed before the output directory is
// touched at all.
func buildBatch(
	corpus *catalog.Corpus, doc *config.Document, local *configFlags,
	vocab config.Vocabulary, isos []string, outFlag string,
) (*render.Batch, string, error) {
	overrides := local.overrides()
	batch := &render.Batch{}
	paths := make(map[string]string, len(isos))
	outputDir := outFlag

	var firstSettings config.Settings
	for _, iso := range isos {
		resolved, err := config.Resolve(doc, iso, overrides)
		if err != nil {
			return nil, "", configError(err)
		}
		if err := config.ValidateResolved(iso, resolved, vocab); err != nil {
			return nil, "", configError(err)
		}
		settings := resolved.Settings
		if firstSettings.Profile == nil {
			firstSettings = settings
		}

		path, err := config.OutputPath(iso, settings)
		if err != nil {
			return nil, "", configError(err)
		}
		paths[iso] = path
		if outputDir == "" && settings.Output != nil && settings.Output.Dir != nil {
			outputDir = *settings.Output.Dir
		}

		request, err := render.GeometryRequest(corpus, iso, settings)
		if err != nil {
			return nil, "", failf(ExitRender, "request_unbuildable", "%v", err)
		}
		result, err := geometry.Generate(request)
		if err != nil {
			// DEC-009's typed absence is a recorded decision, not a failure. Falling
			// back to a source rendering here would emit the exact silhouette the
			// oracle refused, which is what makes the typed outcome worth having.
			if geometry.IsNoArtifact(err) {
				batch.Skipped = append(batch.Skipped, render.Skipped{
					ISO: iso, Name: entityName(corpus, iso), Reason: err.Error(),
				})
				continue
			}
			// Geometry's own budget refusal reaches here before the CLI's check
			// can run, and its message ("linear path bytes 59892 exceed hard
			// maximum 2500") is true and says nothing about what to change.
			// Reclassifying it to the budget class and explaining the cliff is
			// the difference between a diagnostic and a number.
			if pipelineErr := (*geometry.PipelineError)(nil); errors.As(err, &pipelineErr) && pipelineErr.Code == geometry.ErrBudget {
				applied, applyErr := geometry.ApplyPreset(request)
				if applyErr != nil {
					applied = request
				}
				return nil, "", failf(ExitBudget, "path_over_budget", "%s: %v", iso, err).
					withContext("file", path, "hint", overBudgetHint(applied))
			}
			return nil, "", failf(ExitRender, "generation_failed", "%s: %v", iso, err)
		}

		document, err := render.SVG(result, iso, entityName(corpus, iso), settings)
		if err != nil {
			return nil, "", failf(ExitRender, "serialization_failed", "%v", err)
		}
		// The structural gate runs on every asset before any of them is written,
		// so a forbidden construct fails the batch rather than reaching a page.
		if err := render.Structure(document); err != nil {
			return nil, "", failf(ExitValidation, "unsafe_output", "%s: %v", iso, err).
				withContext("file", path)
		}
		if err := checkBudget(iso, path, request, result, len(document)); err != nil {
			return nil, "", err
		}

		batch.Assets = append(batch.Assets, render.StagedAsset{
			Asset: render.Asset{
				ISO: iso, Name: entityName(corpus, iso), File: path,
				Profile:  valueOrEmpty(settings.Profile),
				Boundary: valueOrEmpty(settings.Boundary),
				Style:    valueOrEmpty(settings.Style),
				Delivery: valueOrEmpty(settings.Delivery),
				Bytes:    len(document),
				// Rendered by the serializer's own formatter, so the manifest can
				// never quote a dimension the viewBox does not carry.
				Width:     render.Number(result.ViewBox.Width()),
				Height:    render.Number(result.ViewBox.Height()),
				Markers:   strings.Count(string(document), "<circle"),
				SHA256:    render.Digest(document),
				PathBytes: result.Metrics.PathBytes,
			},
			Content: document,
		})
	}

	// Collisions are a property of the whole selection, so they are checked once
	// the set is known and always before anything is written.
	if err := config.CheckOutputCollisions(paths); err != nil {
		return nil, "", configError(err)
	}
	if outputDir == "" {
		outputDir = "."
	}

	configDigest, err := render.ConfigDigest(firstSettings)
	if err != nil {
		return nil, "", failf(ExitConfig, "config_undigestable", "%v", err)
	}
	batch.Finalize(buildVersion, corpus.Manifest.CorpusVersion, corpus.Manifest.Identity, configDigest)
	return batch, outputDir, nil
}

// checkBudget enforces the frozen per-profile byte ceiling. It is its own exit
// class because it is the one failure a caller fixes by asking for less detail
// rather than by fixing a defect.
func checkBudget(iso, path string, request geometry.Input, result geometry.Result, fileBytes int) error {
	applied, err := geometry.ApplyPreset(request)
	if err != nil {
		return failf(ExitRender, "preset_unresolvable", "%s: %v", iso, err)
	}
	if applied.MaxPathBytes > 0 && result.Metrics.PathBytes > applied.MaxPathBytes {
		return failf(ExitBudget, "path_over_budget",
			"%s: the silhouette is %d bytes against a %d byte ceiling for profile %q",
			iso, result.Metrics.PathBytes, applied.MaxPathBytes, applied.Preset).
			withContext("file", path, "hint", overBudgetHint(applied))
	}
	return nil
}

// overBudgetHint explains the cliff rather than restating the number.
//
// A request whose long side is not the profile's own leaves the committed
// ladder and takes the source path, and for a large entity the source geometry
// is orders of magnitude over the frozen ceiling — Greenland is 59892 bytes
// against the card profile's 2500. Verified: it renders at the card profile's
// own 128 and fails at 160, which is one of that profile's documented reference
// sizes and the value the specification's example config uses. Tracked as
// WKI-C35A01E965DC.
//
// Without this, the diagnostic is geometry's own "linear path bytes 59892 exceed
// hard maximum 2500", which is true and tells an author nothing about what to
// change.
func overBudgetHint(applied geometry.Input) string {
	presets, err := geometry.Presets()
	if err != nil {
		return "choose a smaller profile, or a layout whose long side asks for less detail"
	}
	preset, known := presets[applied.Preset]
	if !known || applied.Layout.LongSide <= 0 || applied.Layout.LongSide == preset.Layout.LongSide {
		return "choose a larger profile, or an entity-specific override"
	}
	return fmt.Sprintf(
		"this asked for a long side of %s, and only the profile's own %s is served from the committed ladder; "+
			"any other size falls back to the full-detail source geometry, which is far over the ceiling for a large entity. "+
			"Use the profile's long side, choose the larger profile, or override this entity",
		trimNumber(applied.Layout.LongSide), trimNumber(preset.Layout.LongSide))
}

func humanGenerate(report GenerateReport) string {
	var out strings.Builder
	verb := "wrote"
	if report.DryRun {
		verb = "would write"
	}
	fmt.Fprintf(&out, "%s %d assets to %s (%d bytes)\n", verb, report.Written, report.Output, report.TotalBytes)
	if len(report.Skipped) != 0 {
		// Skipped entities are named, not counted. A number is easy to read past,
		// and a catalog quietly short by three is exactly what must not happen.
		fmt.Fprintf(&out, "skipped %d with a recorded reason:\n", len(report.Skipped))
		for _, skipped := range report.Skipped {
			fmt.Fprintf(&out, "  %s (%s): %s\n", skipped.ISO, skipped.Name, skipped.Reason)
		}
	}
	fmt.Fprintf(&out, "manifest: %s", report.Manifest)
	return out.String()
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func boolWord(value bool, whenTrue, whenFalse string) string {
	if value {
		return whenTrue
	}
	return whenFalse
}
