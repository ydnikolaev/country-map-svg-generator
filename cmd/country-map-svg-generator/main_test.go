package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

// binaryDir holds the directory the built CLI lives in for this package's test
// run. Scripts get it on PATH.
var binaryDir string

func TestMain(m *testing.M) {
	os.Exit(runPackageTests(m))
}

// runPackageTests builds the real binary once, the way a stranger's checkout
// would build it, and hands it to the scripts.
//
// It deliberately does not use testscript.RunMain. RunMain re-executes the test
// binary as the command, which means the thing under test is compiled by `go
// test` with whatever workspace is ambient. The Go CLI validation profile asks
// for the opposite on both counts: black-box the built artifact, and build it
// with GOWORK=off so a green suite cannot rest on a go.work a consumer will not
// have. The cost is one `go build` per package run.
func runPackageTests(m *testing.M) int {
	dir, err := os.MkdirTemp("", "country-map-cli-e2e-")
	if err != nil {
		fmt.Fprintf(os.Stderr, "e2e: temp dir: %v\n", err)
		return 1
	}
	defer os.RemoveAll(dir)

	name := CLIName
	if runtimeIsWindows() {
		name += ".exe"
	}
	binary := filepath.Join(dir, name)

	build := exec.Command("go", "build", "-o", binary, ".")
	build.Env = append(os.Environ(), "GOWORK=off")
	if output, err := build.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "e2e: GOWORK=off go build failed: %v\n%s\n", err, output)
		return 1
	}
	binaryDir = dir
	return m.Run()
}

func TestScript(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: filepath.Join("testdata", "script"),
		// Every command a script runs must be spelled out with exec, so a script
		// can never accidentally exercise a shell builtin or a host tool.
		RequireExplicitExec: true,
		Setup: func(env *testscript.Env) error {
			// The built binary is the only thing added to the environment. No
			// network, no host state: the CLI reads its corpus from its own
			// embedded data, which is what makes ARCH-INV-2 provable here.
			env.Vars = append(env.Vars,
				"PATH="+binaryDir+string(os.PathListSeparator)+os.Getenv("PATH"),
				"HOME="+env.WorkDir,
			)
			return nil
		},
	})
}

func runtimeIsWindows() bool { return os.PathListSeparator == ';' }
