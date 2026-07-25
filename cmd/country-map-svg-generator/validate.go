package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yuranikolaev/country-map-svg-generator/internal/config"
)

// ValidateReport is CTR-004's validate payload. It reports what was checked as
// well as what passed: "valid" over an unstated selection tells a caller
// nothing, and the entity count is what distinguishes a real full-catalog check
// from a config that resolved for one country.
type ValidateReport struct {
	Config    string `json:"config,omitempty"`
	Schema    string `json:"schema"`
	Entities  int    `json:"entities_checked"`
	Selection string `json:"selection"`
	Valid     bool   `json:"valid"`
}

func newValidateCommand(flags *globalFlags) *cobra.Command {
	local := &configFlags{}
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Check a configuration before it writes anything",
		Long: strings.TrimSpace(`
Check a configuration and report every problem at once.

Validation runs twice over different things. Each document layer — the file, its
preset ancestry, each per-country block — is checked on its own, so a mistake is
reported against the layer that contains it rather than against whatever
inherited it. Then the fully resolved configuration is checked for combinations
that only exist after merging: a contain layout whose frame an override removed
is legal in every layer and contradictory in the result.

With no --config the embedded defaults are checked, which is what the generator
would use.`),
		Args: exactArgs(0, CLIName+" validate [--config <path>] [--iso <code>] [--json]"),
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

			overrides := local.overrides()
			paths := make(map[string]string, len(isos))
			for _, iso := range isos {
				resolved, resolveErr := config.Resolve(doc, iso, overrides)
				if resolveErr != nil {
					return configError(resolveErr)
				}
				if validateErr := config.ValidateResolved(iso, resolved, vocab); validateErr != nil {
					return configError(validateErr)
				}
				path, pathErr := config.OutputPath(iso, resolved.Settings)
				if pathErr != nil {
					return configError(pathErr)
				}
				paths[iso] = path
			}
			// Collisions are a property of the whole selection, so they are the
			// last check and the only one that needs every entity resolved first.
			if err := config.CheckOutputCollisions(paths); err != nil {
				return configError(err)
			}

			report := ValidateReport{
				Config: local.path, Schema: config.SchemaVersion,
				Entities: len(isos), Selection: describeSelection(local), Valid: true,
			}
			return emit(cmd, flags, report, humanValidate(report))
		},
	}
	local.bind(cmd)
	return cmd
}

func describeSelection(flags *configFlags) string {
	if flags.iso == "" {
		return "full catalog"
	}
	return strings.ToUpper(flags.iso)
}

func humanValidate(report ValidateReport) string {
	source := "embedded defaults"
	if report.Config != "" {
		source = report.Config
	}
	return fmt.Sprintf("valid: %s resolves for %d entities (%s)", source, report.Entities, report.Selection)
}
