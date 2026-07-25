package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ydnikolaev/country-map-svg-generator/internal/config"
)

// InitReport tells a caller what was written and what to do next, so an agent
// does not have to guess the follow-up command.
type InitReport struct {
	Path    string   `json:"path"`
	Preset  string   `json:"preset"`
	Schema  string   `json:"schema"`
	Presets []string `json:"available_presets"`
	// Notes carries the guidance a YAML starter holds as comments, for the JSON
	// path where comments do not exist.
	Notes      []string `json:"notes,omitempty"`
	NextAction string   `json:"next_action"`
}

func newInitCommand(flags *globalFlags) *cobra.Command {
	var (
		out    string
		preset string
		force  bool
	)
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Write a starter configuration",
		Long: strings.TrimSpace(`
Write a starter configuration that is valid as generated.

The file is commented, extends one of the embedded presets, and carries only the
keys worth changing first — a configuration that restated every default would be
noise, and would attribute those values to the file rather than to where they
actually come from.

Nothing is prompted for: the CLI is non-interactive by construction, so every
choice is a flag. An existing file is never overwritten without --force.`),
		Args: exactArgs(0, CLIName+" init [--out <path>] [--preset <name>] [--force] [--json]"),
		RunE: func(cmd *cobra.Command, args []string) error {
			known, err := config.PresetNames()
			if err != nil {
				return failf(ExitConfig, "presets_unreadable", "embedded presets could not be read: %v", err)
			}
			if !containsString(known, preset) {
				return failf(ExitUsage, "unknown_preset", "unknown preset %q", preset).
					withContext("available", strings.Join(known, ", "))
			}
			if extension := strings.ToLower(filepath.Ext(out)); config.Extensions[extension] == "" {
				return failf(ExitUsage, "unsupported_extension", "unsupported config extension %q", extension).
					withContext("accepted", ".yaml, .yml, .json")
			}

			if _, err := os.Stat(out); err == nil && !force {
				return failf(ExitFilesystem, "exists", "%s already exists", out).
					withContext("hint", "pass --force to overwrite it")
			} else if err != nil && !os.IsNotExist(err) {
				return failf(ExitFilesystem, "stat_failed", "%s: %v", out, err)
			}

			format := config.Extensions[strings.ToLower(filepath.Ext(out))]
			body := starterConfig(preset, format)
			// The generated file goes through the same validation an authored one
			// does. A starter config that does not validate is worse than none:
			// it teaches the wrong shape and blames the author for it. This is
			// what caught the JSON path writing a commented YAML body.
			if _, err := config.Decode([]byte(body), format); err != nil {
				return failf(ExitConfig, "starter_invalid", "the generated starter configuration does not validate: %v", err)
			}
			if err := os.WriteFile(out, []byte(body), 0o644); err != nil {
				return failf(ExitFilesystem, "write_failed", "%s: %v", out, err)
			}

			report := InitReport{
				Path: out, Preset: preset, Schema: config.SchemaVersion, Presets: known,
				NextAction: CLIName + " validate --config " + out,
			}
			human := fmt.Sprintf("wrote %s (extends %s)\nnext: %s", out, preset, report.NextAction)
			if format == config.FormatJSON {
				// JSON has no comments, so the guidance a YAML starter carries
				// inline is printed instead. Dropping it silently would make the
				// JSON path a worse product for no stated reason.
				report.Notes = starterNotes()
				human = fmt.Sprintf("wrote %s (extends %s)\n\n%s\n\nnext: %s",
					out, preset, strings.Join(starterNotes(), "\n"), report.NextAction)
			}
			return emit(cmd, flags, report, human)
		},
	}
	cmd.Flags().StringVar(&out, "out", "country-map.yaml", "path to write")
	cmd.Flags().StringVar(&preset, "preset", "site-default", "embedded preset to extend")
	cmd.Flags().BoolVar(&force, "force", false, "overwrite an existing file")
	return cmd
}

// starterNotes is the guidance a YAML starter carries as comments. JSON has no
// comments, so the JSON path prints these instead of dropping them.
func starterNotes() []string {
	return []string{
		"Every key is optional: the file only needs to carry what differs from the preset it extends.",
		"profile selects the detail preset (card, hero) — this is geometry's preset, not its boundary profile.",
		"boundary selects the boundary posture (un, de_facto) — the underscore is not optional.",
		"layout tight derives natural proportions; contain takes an explicit width and height. The two modes accept disjoint fields.",
		"countries.<ISO> overrides any of the above for one entity.",
		"Run " + CLIName + " explain --iso <code> to see the resolved value of anything and the layer it came from.",
	}
}

// starterConfig is written as text rather than marshalled from a struct, because
// for YAML the comments are the point: they are what make the file a starting
// point instead of a dump.
func starterConfig(preset string, format config.Format) string {
	if format == config.FormatJSON {
		return fmt.Sprintf(`{
  "schema": %q,
  "extends": %q,
  "profile": "card",
  "boundary": "un",
  "layout": { "mode": "tight", "longSide": 160 }
}
`, config.SchemaVersion, preset)
	}
	return fmt.Sprintf(`# %s configuration.
# Every key is optional — this file only needs to carry what differs from the
# preset it extends. Run "%s explain --iso <code> --config <this file>" to see
# the resolved value of anything and the layer it came from.
schema: %s

# The named preset this file inherits from. Embedded presets are versioned
# public surface; run "%s init --preset <name>" to start from a different one.
extends: %s

# Detail profile: which byte budget and natural size the map is built for.
# Note that this is geometry's *preset*, not its boundary profile.
profile: card

# Boundary posture (DEC-002). "un" is the default; "de_facto" is the other
# accepted value. The underscore is not optional.
boundary: un

# tight derives natural proportions from the geometry; contain takes an explicit
# width and height and centres the unused space. The two modes accept disjoint
# fields, so a width here would be refused rather than ignored.
#
# No longSide is set on purpose. Only the profile's own long side is served from
# the committed detail ladder; any other size falls back to the full-detail
# source geometry, which is far over the byte ceiling for a large entity. Set one
# only for a selection you have generated successfully.
layout:
  mode: tight

# Per-country overrides. The key is an ISO alpha-2 code, and everything above
# can be overridden here for one entity.
# countries:
#   US:
#     profile: hero
#     layout: { mode: contain, width: 720, height: 420 }
`, CLIName, CLIName, config.SchemaVersion, CLIName, preset)
}

func containsString(set []string, value string) bool {
	for _, item := range set {
		if item == value {
			return true
		}
	}
	return false
}
