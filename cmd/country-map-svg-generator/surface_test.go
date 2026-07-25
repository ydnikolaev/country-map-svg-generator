package main

import (
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// scriptDir is where the testscript scenarios live. The gate below reads it as
// a directory listing rather than as a hard-coded list, so a script added by
// hand counts and a script deleted by hand is noticed.
const scriptDir = "testdata/script"

// requiredScenarios is VAL-1's obligation restated as file names. Every runnable
// command leaf needs all four: the help text an agent discovers the command
// with, the happy path, an argument mistake, and the same mistake under --json
// so the machine surface is proven on the error path too — the path an agent is
// most likely to hit and least able to recover from if the shape drifts.
var requiredScenarios = []string{"help", "happy", "invalid", "json_error"}

// TestCommandSurfaceIsFullyCovered enumerates the command tree from cobra itself
// rather than from a maintained list. That is the point: a new command is
// registered in one place, and this gate reddens until its scripts exist. A gate
// that read its own list of commands would prove only that the list matches
// itself.
func TestCommandSurfaceIsFullyCovered(t *testing.T) {
	leaves := runnableLeaves(t)
	if len(leaves) == 0 {
		t.Fatal("command tree exposes no runnable leaves; the enumeration is broken, not the CLI")
	}

	present := scriptsPresent(t)
	for _, leaf := range leaves {
		for _, scenario := range requiredScenarios {
			name := leaf + "_" + scenario + ".txtar"
			if !present[name] {
				t.Errorf("command %q has no %s scenario: expected %s", leaf, scenario, filepath.Join(scriptDir, name))
			}
		}
	}

	// The root is not a leaf — it has children — but its own refusals are part
	// of CTR-004 and an agent hits them first.
	if !present["root_usage.txtar"] {
		t.Errorf("root usage refusals are uncovered: expected %s", filepath.Join(scriptDir, "root_usage.txtar"))
	}
}

// TestEveryScriptBelongsToACommand is the other direction. Without it a renamed
// command leaves its old scripts behind, still green, testing a surface that no
// longer exists.
func TestEveryScriptBelongsToACommand(t *testing.T) {
	known := map[string]bool{"root_usage": true}
	for _, leaf := range runnableLeaves(t) {
		for _, scenario := range requiredScenarios {
			known[leaf+"_"+scenario] = true
		}
	}
	for name := range scriptsPresent(t) {
		base := strings.TrimSuffix(name, ".txtar")
		if !known[base] {
			t.Errorf("script %s matches no command scenario; rename or delete it", filepath.Join(scriptDir, name))
		}
	}
}

// TestEveryExitClassIsNamed keeps the taxonomy and its rendering in step. An
// unnamed class would serialize as "unknown" in the envelope, which tells an
// agent nothing and is exactly the diagnostic quality REQ-12 forbids.
func TestEveryExitClassIsNamed(t *testing.T) {
	for class := ExitOK; class <= ExitFilesystem; class++ {
		if class.String() == "unknown" {
			t.Errorf("exit class %d has no name", int(class))
		}
	}
	if ExitClass(99).String() != "unknown" {
		t.Error("an unregistered class must render as unknown rather than as a neighbour")
	}
}

func runnableLeaves(t *testing.T) []string {
	t.Helper()
	root, _ := newRootCommand(io.Discard, io.Discard)
	return leavesOf(root)
}

// leavesOf is separated from the gate so the enumeration itself can be tested
// against a tree the gate never sees. A gate whose enumeration is only ever run
// on the real tree cannot distinguish "every command is covered" from "no
// command was found".
func leavesOf(root *cobra.Command) []string {
	var leaves []string
	var walk func(*cobra.Command)
	walk = func(cmd *cobra.Command) {
		children := cmd.Commands()
		userFacing := 0
		for _, child := range children {
			if child.Hidden || child.Name() == "help" || child.Name() == "completion" {
				continue
			}
			userFacing++
			walk(child)
		}
		if userFacing == 0 && cmd.Runnable() && cmd != root {
			leaves = append(leaves, strings.ReplaceAll(strings.TrimPrefix(cmd.CommandPath(), root.Name()+" "), " ", "_"))
		}
	}
	walk(root)
	sort.Strings(leaves)
	return leaves
}

// TestLeafEnumerationCatchesANewCommand is the coverage gate's tooth. Both
// TestCommandSurfaceIsFullyCovered and TestEveryScriptBelongsToACommand are
// green when the enumeration returns nothing useful, so the enumeration is
// exercised against a tree with the shapes it has to tell apart: a leaf, a
// parent that is not a leaf, a nested leaf, and the three kinds of command that
// must not be demanded of a script author.
func TestLeafEnumerationCatchesANewCommand(t *testing.T) {
	root := &cobra.Command{Use: "tool", Run: func(*cobra.Command, []string) {}}
	root.AddCommand(&cobra.Command{Use: "alpha", Run: func(*cobra.Command, []string) {}})

	parent := &cobra.Command{Use: "beta", Run: func(*cobra.Command, []string) {}}
	parent.AddCommand(&cobra.Command{Use: "inner", Run: func(*cobra.Command, []string) {}})
	root.AddCommand(parent)

	// Not user-facing surface: no script may be demanded for these.
	root.AddCommand(&cobra.Command{Use: "secret", Hidden: true, Run: func(*cobra.Command, []string) {}})
	root.AddCommand(&cobra.Command{Use: "completion", Run: func(*cobra.Command, []string) {}})
	// A grouping command with no Run of its own is not runnable and not a leaf.
	root.AddCommand(&cobra.Command{Use: "group"})

	got := leavesOf(root)
	want := []string{"alpha", "beta_inner"}
	if len(got) != len(want) {
		t.Fatalf("leaves = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("leaves = %v, want %v", got, want)
		}
	}
}

func scriptsPresent(t *testing.T) map[string]bool {
	t.Helper()
	entries, err := os.ReadDir(scriptDir)
	if err != nil {
		t.Fatalf("read %s: %v", scriptDir, err)
	}
	present := map[string]bool{}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".txtar") {
			present[entry.Name()] = true
		}
	}
	return present
}
