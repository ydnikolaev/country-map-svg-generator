package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// CLIName is the binary's invoked name. It is the one place the name is spelled,
// so a rename does not have to be chased through help text and diagnostics.
//
// It follows the name the binary was actually invoked under, because the name is
// only useful if the reader can retype it. Installed through a shim as `svgmap`,
// every usage line and every hint has to say `svgmap generate` — a diagnostic
// that sends an operator to a command they do not have is worse than none.
var CLIName = defaultCLIName()

const canonicalCLIName = "country-map-svg-generator"

// defaultCLIName reads argv[0] and falls back to the canonical name for the two
// cases where argv[0] is not a name anyone can type: an empty or path-shaped
// argv[0], and the `go test` binary, which would otherwise put `*.test` into the
// help text the txtar scripts assert against.
func defaultCLIName() string {
	if len(os.Args) == 0 {
		return canonicalCLIName
	}
	base := filepath.Base(os.Args[0])
	switch {
	case base == "", base == ".", base == "/", base == string(filepath.Separator):
		return canonicalCLIName
	case strings.HasSuffix(base, ".test"), strings.HasPrefix(base, "___"):
		return canonicalCLIName
	}
	return base
}

// globalFlags are the options every command honours. The CLI is agent-first
// (US-1), so the machine surface is a root-level concern rather than a per
// command afterthought.
type globalFlags struct {
	json bool
}

func newRootCommand(stdout, stderr io.Writer) (*cobra.Command, *globalFlags) {
	flags := &globalFlags{}
	root := &cobra.Command{
		Use:   CLIName,
		Short: "Generate optimized country map SVGs from a versioned geometry corpus",
		Long: strings.TrimSpace(`
Generate optimized country map SVGs from a versioned, embedded geometry corpus.

The primary operator is an AI agent: every command is non-interactive, every
option is a flag, and --json emits a stable envelope with a typed error class.
Normal generation is offline and needs no Node, Python, GDAL or auxiliary files.`),
		SilenceUsage:  true,
		SilenceErrors: true,
		// A bare invocation is a usage error rather than a silent no-op: an agent
		// that forgot the subcommand should get a typed refusal, not exit 0.
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return failf(ExitUsage, "unknown_command", "unknown command %q", args[0]).
					withContext("hint", "run "+CLIName+" --help for the command list")
			}
			return failf(ExitUsage, "missing_command", "no command given").
				withContext("hint", "run "+CLIName+" --help for the command list")
		},
	}
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.PersistentFlags().BoolVar(&flags.json, "json", false, "emit the machine-readable envelope instead of human output")

	// Flag parsing failures are usage failures, not unclassified ones. Without
	// this they would reach the top level as a bare error and be mapped by
	// fallback rather than by intent.
	root.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		return failf(ExitUsage, "flag", "%s", err.Error()).
			withContext("command", cmd.CommandPath())
	})

	// Commands are registered by the task that implements them, never stubbed:
	// registration is what makes VAL-1's gate demand a command's scripts, and a
	// stub would satisfy it with scripts asserting that nothing happens.
	root.AddCommand(
		newDemoCommand(flags),
		newInitCommand(flags),
		newValidateCommand(flags),
		newExplainCommand(flags),
		newGenerateCommand(flags),
		newInspectCommand(flags),
		newPreviewCommand(flags),
		newSchemaCommand(flags),
		newVersionCommand(flags),
	)
	return root, flags
}

// emit renders one successful result in whichever mode the caller asked for.
// Commands call this instead of printing directly, so the envelope shape cannot
// drift command by command.
func emit(cmd *cobra.Command, flags *globalFlags, data any, human string) error {
	if flags.json {
		return writeEnvelope(cmd.OutOrStdout(), successEnvelope(cmd.CommandPath(), data))
	}
	if human == "" {
		return nil
	}
	_, err := fmt.Fprintln(cmd.OutOrStdout(), strings.TrimRight(human, "\n"))
	return err
}

// reportError renders a failure. It takes the flag value from the parsed struct
// when parsing got that far and from the raw arguments when it did not — an
// agent that passes a bad flag together with --json still needs the envelope,
// which is exactly the case where the parsed value is unavailable.
func reportError(stdout, stderr io.Writer, command string, jsonMode bool, err *CLIError) {
	if jsonMode {
		_ = writeEnvelope(stdout, errorEnvelope(command, err))
		return
	}
	fmt.Fprintf(stderr, "%s: %s: %s\n", CLIName, err.Class, err.Message)
	keys := make([]string, 0, len(err.Context))
	for key := range err.Context {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Fprintf(stderr, "  %s: %s\n", key, err.Context[key])
	}
}

// valueTakingFlags are the flags whose next argument is their value. The error
// path has to know them: without that, `--config --json` (a missing value) and
// any flag whose value happens to be the string `--json` would both flip the
// output mode, which is a confusing thing to happen while reporting a different
// mistake.
var valueTakingFlags = map[string]bool{
	"--config": true, "--out": true, "--preset": true,
	"--profile": true, "--boundary": true, "--style": true, "--delivery": true,
	"--iso": true,
}

// requestedJSON scans the raw arguments. Used only on the error path, when flag
// parsing may not have completed and the parsed value therefore does not exist.
func requestedJSON(args []string) bool {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			return false
		}
		if arg == "--json" || arg == "--json=true" {
			return true
		}
		// Skip the value position, so a value is never mistaken for a flag.
		if valueTakingFlags[arg] {
			i++
		}
	}
	return false
}

// exactArgs wraps cobra's arity check so an arity mistake carries the usage exit
// class like every other caller mistake.
func exactArgs(n int, usage string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != n {
			return failf(ExitUsage, "arity", "expected %d argument(s), got %d", n, len(args)).
				withContext("command", cmd.CommandPath(), "usage", usage)
		}
		return nil
	}
}
