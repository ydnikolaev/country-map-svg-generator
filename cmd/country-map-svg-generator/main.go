// Command country-map-svg-generator turns a versioned, embedded country geometry
// corpus into optimized SVG assets. It owns BND-004: the executable and the
// command surface, never geodata parsing, which belongs to internal/catalog and
// internal/geometry.
package main

import (
	"errors"
	"io"
	"os"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run is the whole program with its edges injected, so the exit class and both
// streams are assertable from a unit test as well as through the built binary.
// The e2e suite still drives the binary — this only makes the wiring itself
// testable, it does not replace that.
func run(args []string, stdout, stderr io.Writer) int {
	root, flags := newRootCommand(stdout, stderr)
	root.SetArgs(args)

	cmd, err := root.ExecuteC()
	if err == nil {
		return int(ExitOK)
	}

	command := CLIName
	if cmd != nil {
		command = cmd.CommandPath()
	}
	// flags.json is authoritative once parsing succeeded; before that the raw
	// arguments are the only evidence of what the caller asked for.
	jsonMode := flags.json || requestedJSON(args)
	// One classification serves both the reported envelope and the exit code:
	// they are the same fact, and deriving it twice invites them to disagree.
	cliErr := topLevelError(err)
	reportError(stdout, stderr, command, jsonMode, cliErr)
	return int(cliErr.Class)
}

// topLevelError maps whatever reached the top onto the taxonomy. Command bodies
// always return a *CLIError, so anything else at this layer came from cobra's
// own parsing and routing — an unknown command, a malformed flag — which is a
// caller mistake. Mapping the fallback to usage rather than to an internal class
// keeps a mistyped command from being reported as a rendering failure.
func topLevelError(err error) *CLIError {
	var cliErr *CLIError
	if errors.As(err, &cliErr) {
		return cliErr
	}
	return &CLIError{Class: ExitUsage, Code: "usage", Message: err.Error()}
}
