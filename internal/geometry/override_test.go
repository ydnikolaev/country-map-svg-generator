package geometry

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/yuranikolaev/country-map-svg-generator/internal/catalog"
)

func TestOverrideFailsClosed(t *testing.T) {
	bad := 200.0
	for name, o := range map[string]Override{
		"missing metadata": {},
		"rotation":         {ISO: "ZZ", Version: "v1", Reason: "test", Rotation: &bad},
		"parts":            {ISO: "ZZ", Version: "v1", Reason: "test", MinimumParts: -1},
	} {
		t.Run(name, func(t *testing.T) {
			if ValidateOverride(o, "", "") == nil {
				t.Fatal("invalid override passed")
			}
		})
	}
	if err := ValidateOverride(Override{ISO: "ZZ", Version: "v1", Reason: "test", Corpus: "a"}, "b", ""); err == nil {
		t.Fatal("wrong corpus passed")
	}
}

func TestGroupAnchorSchemaMutationTeeth(t *testing.T) {
	if set, err := EmbeddedGroupAnchors(); err != nil || len(set.Groups) != 0 {
		t.Fatalf("empty embedded group anchors set=%+v err=%v", set, err)
	}
	for _, field := range []string{"threshold", "path_budget", "candidate_order", "entity_branch", "mode"} {
		t.Run(field, func(t *testing.T) {
			raw := []byte(`{"schema_version":1,"version":"v1","groups":[],"` + field + `":"forbidden"}`)
			if _, err := ParseGroupAnchors(raw); err == nil {
				t.Fatalf("%s executable/policy mutation passed", field)
			}
		})
	}
	valid := reviewedGroup("archipelago", catalog.Point{21, 0})
	for name, mutate := range map[string]func(*GroupAnchor){
		"stale projection": func(group *GroupAnchor) { group.Review.ProjectionVersion = "v0" },
		"invalid digest":   func(group *GroupAnchor) { group.Review.SHA256 = "not-a-digest" },
		"unapproved":       func(group *GroupAnchor) { group.Review.Status = "pending" },
	} {
		t.Run(name, func(t *testing.T) {
			group := valid
			mutate(&group)
			raw, err := json.Marshal(GroupAnchorSet{SchemaVersion: 1, Version: "v1", Groups: []GroupAnchor{group}})
			if err != nil {
				t.Fatal(err)
			}
			if _, err := ParseGroupAnchors(raw); err == nil {
				t.Fatalf("%s review mutation passed", name)
			}
		})
	}
	conflict := valid
	conflict.Anchor = catalog.Point{31, 0}
	raw, err := json.Marshal(GroupAnchorSet{SchemaVersion: 1, Version: "v1", Groups: []GroupAnchor{valid, conflict}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParseGroupAnchors(raw); err == nil {
		t.Fatal("duplicate/conflicting group name passed")
	}
}

func TestGroupAnchorIsAdditiveAndCannotEvictBase(t *testing.T) {
	in := groupAnchorInput()
	band := SilhouetteBand{GridLongSide: 140, ContributionThreshold: 10_000}
	set := GroupAnchorSet{SchemaVersion: 1, Version: "v1", Groups: []GroupAnchor{
		reviewedGroup("archipelago", catalog.Point{21, 0}),
	}}
	resolved, application, err := ApplyGroupAnchors(in, set, band)
	if err != nil {
		t.Fatal(err)
	}
	if len(application.BaseSourceOrders) != 1 || application.BaseSourceOrders[0] != 0 ||
		len(application.Groups) != 1 ||
		application.Groups[0].IdentityReason != "group_anchor_additive" {
		t.Fatalf("application=%+v", application)
	}
	full, err := normalize(fromCatalog(in.Geometry.Coordinates))
	if err != nil {
		t.Fatal(err)
	}
	prj := centeredProjector(full, 0)
	projected, _, err := projectGeometry(full, prj, LODProjectionContractV1().Flatness)
	if err != nil {
		t.Fatal(err)
	}
	groupSource := application.Groups[0].SourceOrder
	visibility, err := buildProtectedVisibility(resolved, projected, MultiPolygon{projected[0], projected[groupSource]}, prj, band)
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateGroupAnchorResult(application, visibility); err != nil {
		t.Fatal(err)
	}
	if visibility[0].IdentityRank == 0 || visibility[groupSource].IdentityRank == 0 {
		t.Fatalf("non-additive identity=%+v", visibility)
	}
	for sourceOrder, component := range visibility {
		if sourceOrder != 0 && sourceOrder != groupSource && component.IdentityRank != 0 {
			t.Fatalf("undeclared identity member=%+v", component)
		}
	}
	evicted := append([]ProtectedVisibilityComponent(nil), visibility...)
	evicted[0].IdentityRank = 0
	if ValidateGroupAnchorResult(application, evicted) == nil {
		t.Fatal("base-member eviction passed")
	}
	nonAdditive := append([]ProtectedVisibilityComponent(nil), visibility...)
	nonAdditive[groupSource].Retained = false
	if ValidateGroupAnchorResult(application, nonAdditive) == nil {
		t.Fatal("non-additive group retention passed")
	}
	if _, err := buildProtectedVisibility(resolved, projected, MultiPolygon{projected[0]}, prj, band); err == nil {
		t.Fatal("group-anchor component loss passed")
	}
}

func TestGroupAnchorRejectsStaleReviewAndMissingLineage(t *testing.T) {
	in := groupAnchorInput()
	band := SilhouetteBand{GridLongSide: 140, ContributionThreshold: 10_000}
	stale := reviewedGroup("archipelago", catalog.Point{21, 0})
	stale.Review.Corpus = "stale-corpus"
	if _, _, err := ApplyGroupAnchors(in, GroupAnchorSet{Groups: []GroupAnchor{stale}}, band); err == nil {
		t.Fatal("stale review corpus passed")
	}
	outside := reviewedGroup("archipelago", catalog.Point{-100, 80})
	if _, _, err := ApplyGroupAnchors(in, GroupAnchorSet{Groups: []GroupAnchor{outside}}, band); err == nil {
		t.Fatal("anchor outside source passed")
	}
}

func reviewedGroup(name string, anchor catalog.Point) GroupAnchor {
	return GroupAnchor{
		Name: name, ISO: "ZZ", Corpus: "corpus-v1", Profile: "un", Anchor: anchor,
		Review: GroupAnchorReview{
			ID: "OWNER-001", Reviewer: "geometry-owner", Status: "approved",
			SHA256: strings.Repeat("a", 64), Corpus: "corpus-v1",
			ProjectionVersion: "v1", RecipeVersion: "v2", VisibilityPolicy: "protected-visibility/v1",
		},
	}
}

func groupAnchorInput() Input {
	square := func(minX, minY, maxX, maxY float64) catalog.Polygon {
		return catalog.Polygon{catalog.Ring{
			{minX, minY}, {maxX, minY}, {maxX, maxY}, {minX, maxY}, {minX, minY},
		}}
	}
	return Input{
		Entity: catalog.Entity{Alpha2: "ZZ"},
		Geometry: catalog.Geometry{ID: "group-anchor-fixture", Coordinates: catalog.MultiPolygon{
			square(-5, -5, 5, 5),
			square(20, -1, 22, 1),
			square(30, -2, 34, 2),
		}},
		CorpusID: "corpus-v1", Profile: "un",
	}
}
