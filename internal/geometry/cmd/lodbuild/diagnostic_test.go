package main

import (
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/yuranikolaev/country-map-svg-generator/internal/catalog"
	geometry "github.com/yuranikolaev/country-map-svg-generator/internal/geometry"
)

func TestDiagnosticParameterCapPredicateRequiresAProof(t *testing.T) {
	rows := []catalogRow{{
		Entity: "CA", Profile: "un", Preset: "hero",
		RequestedTier: "source", PathBytes: 675606, OverP2Cap: true,
	}}
	fits, proof, err := proveParameterCandidatesCannotFitAll(rows)
	if err != nil {
		t.Fatal(err)
	}
	if fits || proof.Method != "unchanged_noncompact_witness" ||
		proof.CandidateScope != "compact_resolution_and_weighting" ||
		proof.Entity != "CA" || proof.RequestedTier != "source" ||
		proof.PathBytes != 675606 || proof.P2Cap != 7500 {
		t.Fatalf("fits=%t proof=%+v", fits, proof)
	}
	aggregates := diagnosticAggregates{AnyParameterFitsAll: true}
	if err := applyParameterCapProof(&aggregates, rows); err != nil {
		t.Fatal(err)
	}
	if aggregates.AnyParameterFitsAll || aggregates.ParameterCapProof != proof {
		t.Fatalf("aggregate predicate was not assigned from proof: %+v", aggregates)
	}
	if _, _, err := proveParameterCandidatesCannotFitAll([]catalogRow{{
		Entity: "AQ", Profile: "un", Preset: "card",
		RequestedTier: "compact", PathBytes: 3000, OverP2Cap: true,
	}}); err == nil {
		t.Fatal("compact-only rows left the predicate at its false zero value instead of failing closed")
	}
}

func TestDiagnosticIdentityDriftChecksEveryRecordedLoadBearingIdentity(t *testing.T) {
	c, err := catalog.Embedded()
	if err != nil {
		t.Fatal(err)
	}
	base, _, err := loadDiagnosticIdentities(c)
	if err != nil {
		t.Fatal(err)
	}
	if diagnosticIdentityDrift(base) {
		t.Fatalf("current identities unexpectedly drifted: %+v", base)
	}
	wantRuntime := runtime.Version() + " " + runtime.GOOS + "/" + runtime.GOARCH
	if base.GoVersion != wantRuntime {
		t.Fatalf("recorded Go identity %q was not derived from executing runtime %q", base.GoVersion, wantRuntime)
	}
	mutations := map[string]func(*diagnosticIdentities){
		"p1":                  func(v *diagnosticIdentities) { v.P1Corpus += "-drift" },
		"baseline_commit":     func(v *diagnosticIdentities) { v.BaselineCommit += "-drift" },
		"baseline_tree":       func(v *diagnosticIdentities) { v.BaselineTree += "-drift" },
		"scratch":             func(v *diagnosticIdentities) { v.ScratchBefore += "-drift" },
		"projection":          func(v *diagnosticIdentities) { v.Projection.Flatness += .01 },
		"recipe":              func(v *diagnosticIdentities) { v.RecipeSHA256 += "-drift" },
		"mapshaper_version":   func(v *diagnosticIdentities) { v.MapshaperVersion += "-drift" },
		"mapshaper_integrity": func(v *diagnosticIdentities) { v.MapshaperIntegrity += "-drift" },
		"mapshaper_shasum":    func(v *diagnosticIdentities) { v.MapshaperShasum += "-drift" },
		"lockfile":            func(v *diagnosticIdentities) { v.LockfileSHA256 += "-drift" },
		"tool":                func(v *diagnosticIdentities) { v.ToolSHA256 += "-drift" },
		"go_version":          func(v *diagnosticIdentities) { v.GoVersion += "-drift" },
		"go_mod":              func(v *diagnosticIdentities) { v.GoModSHA256 += "-drift" },
		"go_sum":              func(v *diagnosticIdentities) { v.GoSumSHA256 += "-drift" },
		"serializer":          func(v *diagnosticIdentities) { v.SerializerSHA256 += "-drift" },
		"serializer_test":     func(v *diagnosticIdentities) { v.SerializerTestSHA += "-drift" },
		"command":             func(v *diagnosticIdentities) { v.Command += "-drift" },
	}
	for name, mutate := range mutations {
		t.Run(name, func(t *testing.T) {
			got := base
			mutate(&got)
			if !diagnosticIdentityDrift(got) {
				t.Fatalf("%s drift was not detected", name)
			}
		})
	}
}

func TestDiagnosticSourceInventoryIsExactAndBiting(t *testing.T) {
	root := sourceRoot()
	var discovered []string
	for _, pattern := range []string{
		filepath.Join(root, "internal/geometry/*.go"),
		filepath.Join(root, "internal/geometry/cmd/lodbuild/*.go"),
	} {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			t.Fatal(err)
		}
		for _, match := range matches {
			if strings.HasSuffix(match, "_test.go") {
				continue
			}
			relative, err := filepath.Rel(root, match)
			if err != nil {
				t.Fatal(err)
			}
			discovered = append(discovered, filepath.ToSlash(relative))
		}
	}
	sort.Strings(discovered)
	// Files added to internal/geometry(/cmd/lodbuild) after the RUN-010 T0A.1
	// diagnostic froze diagnosticSourcePaths. Each new maintainer build-tool
	// or runtime file lands here deliberately, not silently.
	postDiagnosticAdditions := []string{
		"internal/geometry/cmd/lodbuild/ladder.go",
		"internal/geometry/cmd/lodbuild/representative.go",
		"internal/geometry/silhouette.go",
	}
	additionSet := make(map[string]bool, len(postDiagnosticAdditions))
	for _, path := range postDiagnosticAdditions {
		additionSet[path] = true
	}
	var preSpike, additions []string
	for _, path := range discovered {
		if additionSet[path] {
			additions = append(additions, path)
		} else {
			preSpike = append(preSpike, path)
		}
	}
	sort.Strings(postDiagnosticAdditions)
	if len(preSpike) == 0 || !reflect.DeepEqual(preSpike, diagnosticSourcePaths) ||
		!reflect.DeepEqual(additions, postDiagnosticAdditions) {
		t.Fatalf("pre-spike=%v additions=%v want pre-spike=%v additions=%v", preSpike, additions, diagnosticSourcePaths, postDiagnosticAdditions)
	}
	items, aggregate, err := loadDiagnosticSourceInventory(root)
	if err != nil {
		t.Fatal(err)
	}
	if diagnosticSourceInventoryDrift(items, aggregate) {
		t.Fatalf("current source inventory unexpectedly drifted: items=%v aggregate=%s", items, aggregate)
	}
	mutated := append([]sourceIdentity(nil), items...)
	mutated[0].SHA256 = strings.Repeat("0", 64)
	if !diagnosticSourceInventoryDrift(mutated, aggregate) {
		t.Fatal("source hash drift was not detected")
	}
	if !diagnosticSourceInventoryDrift(items[:len(items)-1], aggregate) {
		t.Fatal("missing executing source path was not detected")
	}
}

func TestDiagnosticTypedFailureIdentityIsExact(t *testing.T) {
	failure128, _ := expectedTypedFailure(128)
	failure160, _ := expectedTypedFailure(160)
	rows := []parameterRow{
		{Value: 128, Outcome: "typed_topology_failure", Failure: failure128},
		{Value: 160, Outcome: "typed_topology_failure", Failure: failure160},
	}
	if typedFailureMismatch(rows, 2) {
		t.Fatal("exact typed failures were rejected")
	}
	rows[0].Failure += "-drift"
	if !typedFailureMismatch(rows, 2) {
		t.Fatal("typed failure message drift was not detected")
	}
	rows[0].Failure = failure128
	rows[1].Outcome = "typed_other_failure"
	if !typedFailureMismatch(rows, 2) {
		t.Fatal("typed failure outcome drift was not detected")
	}
}

func TestDiagnosticRepresentativeIdentityRetentionIsMonotonic(t *testing.T) {
	compact := representativeRow{Visibility: []geometry.ProtectedVisibilityComponent{
		{SourceOrder: 3, IdentityRank: 1, Retained: true},
	}}
	standard := representativeRow{Visibility: []geometry.ProtectedVisibilityComponent{
		{SourceOrder: 3, IdentityRank: 1, Retained: true},
	}}
	if err := validateMonotonicIdentity(compact, standard); err != nil {
		t.Fatal(err)
	}
	standard.Visibility[0].Retained = false
	if err := validateMonotonicIdentity(compact, standard); err == nil {
		t.Fatal("non-monotonic larger-band retention passed")
	}
}
