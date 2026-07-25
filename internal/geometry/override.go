package geometry

import (
	"bytes"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"io"
	"math"
	"sort"
	"strings"

	"github.com/ydnikolaev/country-map-svg-generator/internal/catalog"
)

//go:embed overrides/v1.json
var overrideJSON []byte

//go:embed overrides/group-anchors.v1.json
var groupAnchorJSON []byte

type Override struct {
	ISO, Version, Corpus, Profile, Reason string
	Rotation                              *float64
	Padding                               *Insets
	Quality                               *Quality
	RetentionAnchors                      []Point
	MinimumParts                          int
	MarkerOffsets                         map[string]Point
}

type overrideWire struct {
	ISO          string   `json:"iso"`
	Version      string   `json:"version"`
	Corpus       string   `json:"corpus"`
	Profile      string   `json:"profile"`
	Reason       string   `json:"reason"`
	Rotation     *float64 `json:"rotation,omitempty"`
	MinimumParts int      `json:"minimum_parts,omitempty"`
}

type GroupAnchorSet struct {
	SchemaVersion int           `json:"schema_version"`
	Version       string        `json:"version"`
	Groups        []GroupAnchor `json:"groups"`
}

type GroupAnchor struct {
	Name    string            `json:"name"`
	ISO     string            `json:"iso"`
	Corpus  string            `json:"corpus"`
	Profile string            `json:"profile"`
	Anchor  catalog.Point     `json:"anchor"`
	Review  GroupAnchorReview `json:"review"`
}

type GroupAnchorReview struct {
	ID                string `json:"id"`
	Reviewer          string `json:"reviewer"`
	Status            string `json:"status"`
	SHA256            string `json:"sha256"`
	Corpus            string `json:"corpus"`
	ProjectionVersion string `json:"projection_version"`
	RecipeVersion     string `json:"recipe_version"`
	VisibilityPolicy  string `json:"visibility_policy"`
}

type GroupAnchorResolution struct {
	Name, ReviewID, Reviewer, ReviewSHA256 string
	SourceOrder                            int
	SourceDigest, IdentityReason           string
}

type GroupAnchorApplication struct {
	BaseSourceOrders []int
	Groups           []GroupAnchorResolution
}

func Overrides() ([]Override, error) {
	var rows []overrideWire
	d := json.NewDecoder(bytes.NewReader(overrideJSON))
	d.DisallowUnknownFields()
	if err := d.Decode(&rows); err != nil {
		return nil, err
	}
	out := make([]Override, len(rows))
	for i, r := range rows {
		out[i] = Override{ISO: r.ISO, Version: r.Version, Corpus: r.Corpus, Profile: r.Profile, Reason: r.Reason, Rotation: r.Rotation, MinimumParts: r.MinimumParts}
		if err := ValidateOverride(out[i], "", ""); err != nil {
			return nil, err
		}
	}
	return out, nil
}

func EmbeddedGroupAnchors() (GroupAnchorSet, error) {
	return ParseGroupAnchors(groupAnchorJSON)
}

func ParseGroupAnchors(raw []byte) (GroupAnchorSet, error) {
	var set GroupAnchorSet
	d := json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	if err := d.Decode(&set); err != nil {
		return set, err
	}
	if err := d.Decode(&struct{}{}); err != io.EOF {
		return set, fail(ErrOverride, "", "group_anchor_schema", "exactly one group-anchor document is required")
	}
	if set.SchemaVersion != 1 || set.Version != "v1" {
		return set, fail(ErrOverride, "", "group_anchor_schema", "schema_version=1 and version=v1 are required")
	}
	names := make(map[string]bool, len(set.Groups))
	for _, group := range set.Groups {
		if group.Name == "" || names[group.Name] {
			return set, fail(ErrOverride, group.ISO, "group_anchor_name", "group names must be non-empty and unique")
		}
		names[group.Name] = true
		if len(group.ISO) != 2 || group.Corpus == "" || group.Profile == "" ||
			!finite(group.Anchor[0]) || !finite(group.Anchor[1]) ||
			group.Anchor[0] < -180 || group.Anchor[0] > 180 ||
			group.Anchor[1] < -90 || group.Anchor[1] > 90 {
			return set, fail(ErrOverride, group.ISO, "group_anchor_metadata", "bounded iso, corpus, profile, and anchor are required")
		}
		review := group.Review
		sum, sumErr := hex.DecodeString(review.SHA256)
		if review.ID == "" || review.Reviewer == "" || review.Status != "approved" ||
			len(sum) != sha256.Size || sumErr != nil || hex.EncodeToString(sum) != review.SHA256 ||
			strings.Trim(review.SHA256, "0") == "" ||
			review.Corpus == "" || review.Corpus != group.Corpus || review.ProjectionVersion != "v1" ||
			review.RecipeVersion != "v2" || review.VisibilityPolicy != "protected-visibility/v1" {
			return set, fail(ErrOverride, group.ISO, "group_anchor_review", "current approved review provenance is required")
		}
	}
	return set, nil
}

// ApplyGroupAnchors resolves reviewed data through the same centered projection
// and source-component lineage used by protected anchors. It expands the base
// count only for group components outside the existing post-anchor base, so a
// group can add identity but cannot remove or replace an existing member.
func ApplyGroupAnchors(in Input, set GroupAnchorSet, band SilhouetteBand) (Input, GroupAnchorApplication, error) {
	raw, err := json.Marshal(set)
	if err != nil {
		return in, GroupAnchorApplication{}, err
	}
	set, err = ParseGroupAnchors(raw)
	if err != nil {
		return in, GroupAnchorApplication{}, err
	}
	var selected []GroupAnchor
	for _, group := range set.Groups {
		if group.ISO == in.Entity.Alpha2 && group.Corpus == in.CorpusID && group.Profile == in.Profile {
			selected = append(selected, group)
		}
	}
	if len(selected) == 0 {
		return in, GroupAnchorApplication{}, nil
	}
	full, err := normalize(fromCatalog(in.Geometry.Coordinates))
	if err != nil {
		return in, GroupAnchorApplication{}, err
	}
	prj := centeredProjector(full, 0)
	projected, _, err := projectGeometry(full, prj, LODProjectionContractV1().Flatness)
	if err != nil {
		return in, GroupAnchorApplication{}, err
	}
	base, err := protectedVisibilityBase(in, projected, prj, band)
	if err != nil {
		return in, GroupAnchorApplication{}, err
	}
	application := GroupAnchorApplication{BaseSourceOrders: sortedSourceOrders(base)}
	additions := map[int]bool{}
	for _, group := range selected {
		if group.Review.Corpus != in.CorpusID {
			return in, GroupAnchorApplication{}, fail(ErrOverride, group.ISO, "group_anchor_review", "review corpus is stale")
		}
		anchor, projectErr := prj.project(group.Anchor[0], group.Anchor[1])
		if projectErr != nil {
			return in, GroupAnchorApplication{}, projectErr
		}
		sourceOrder := containingSourceComponent(projected, anchor)
		if sourceOrder < 0 {
			return in, GroupAnchorApplication{}, fail(ErrProtected, group.ISO, "group_anchor_lineage", "group %q anchor has no source component", group.Name)
		}
		raw, _ := json.Marshal(projected[sourceOrder])
		sum := sha256.Sum256(raw)
		application.Groups = append(application.Groups, GroupAnchorResolution{
			Name: group.Name, ReviewID: group.Review.ID, Reviewer: group.Review.Reviewer,
			ReviewSHA256: group.Review.SHA256, SourceOrder: sourceOrder,
			SourceDigest: hex.EncodeToString(sum[:]), IdentityReason: "group_anchor_additive",
		})
		if !base[sourceOrder] {
			additions[sourceOrder] = true
		}
	}
	minimum := len(base) + len(additions)
	if minimum > len(projected) {
		return in, GroupAnchorApplication{}, fail(ErrProtected, in.Entity.Alpha2, "group_anchor_count", "additive identity exceeds source component count")
	}
	out := in
	out.Entity.Protected = append([]catalog.ProtectedFeature(nil), in.Entity.Protected...)
	for _, group := range selected {
		out.Entity.Protected = append(out.Entity.Protected, catalog.ProtectedFeature{
			Name: group.Name, Anchor: group.Anchor, MinimumParts: minimum,
		})
	}
	return out, application, nil
}

func ValidateGroupAnchorResult(application GroupAnchorApplication, visibility []ProtectedVisibilityComponent) error {
	byOrder := make(map[int]ProtectedVisibilityComponent, len(visibility))
	for _, component := range visibility {
		byOrder[component.SourceOrder] = component
	}
	for _, sourceOrder := range application.BaseSourceOrders {
		component, ok := byOrder[sourceOrder]
		if !ok || component.IdentityRank == 0 || !component.Retained {
			return fail(ErrProtected, "", "group_anchor_base", "base member %d was evicted", sourceOrder)
		}
	}
	for _, group := range application.Groups {
		component, ok := byOrder[group.SourceOrder]
		if !ok || component.SourceDigest != group.SourceDigest ||
			component.IdentityRank == 0 || !component.Retained ||
			group.IdentityReason != "group_anchor_additive" {
			return fail(ErrProtected, "", "group_anchor_additive", "group %q was not additively retained", group.Name)
		}
	}
	return nil
}

func protectedVisibilityBase(in Input, reference MultiPolygon, prj projector, band SilhouetteBand) (map[int]bool, error) {
	bounds := geometryBounds(reference)
	type rankedComponent struct {
		index, contribution int
		digest              string
	}
	ranked := make([]rankedComponent, len(reference))
	largest := 0
	for i, polygon := range reference {
		raw, _ := json.Marshal(polygon)
		sum := sha256.Sum256(raw)
		raster, err := RasterizeSilhouette(MultiPolygon{polygon}, bounds, band.GridLongSide)
		if err != nil {
			return nil, err
		}
		for _, filled := range raster.Pixels {
			if filled {
				ranked[i].contribution++
			}
		}
		ranked[i].index, ranked[i].digest = i, hex.EncodeToString(sum[:])
		if polygonArea(polygon) > polygonArea(reference[largest]) {
			largest = i
		}
	}
	sort.Slice(ranked, func(i, j int) bool {
		if ranked[i].contribution != ranked[j].contribution {
			return ranked[i].contribution > ranked[j].contribution
		}
		if ranked[i].index != ranked[j].index {
			return ranked[i].index < ranked[j].index
		}
		return ranked[i].digest < ranked[j].digest
	})
	minimum := 1
	for _, feature := range in.Entity.Protected {
		if feature.MinimumParts > minimum {
			minimum = feature.MinimumParts
		}
	}
	if minimum > len(reference) {
		minimum = len(reference)
	}
	base := make(map[int]bool, minimum)
	for _, component := range ranked[:minimum] {
		base[component.index] = true
	}
	anchors := make(map[int]bool, len(in.Entity.Protected))
	for _, feature := range in.Entity.Protected {
		anchor, err := prj.project(feature.Anchor[0], feature.Anchor[1])
		if err != nil {
			return nil, err
		}
		sourceOrder := containingSourceComponent(reference, anchor)
		if sourceOrder < 0 {
			return nil, fail(ErrProtected, in.Entity.Alpha2, "anchor_lineage", "anchor %q has no source component", feature.Name)
		}
		anchors[sourceOrder] = true
	}
	for _, sourceOrder := range sortedSourceOrders(anchors) {
		if base[sourceOrder] {
			continue
		}
		evicted := false
		for rank := minimum - 1; rank >= 0; rank-- {
			index := ranked[rank].index
			if base[index] && !anchors[index] && index != largest {
				delete(base, ranked[rank].index)
				evicted = true
				break
			}
		}
		if !evicted {
			return nil, fail(ErrProtected, in.Entity.Alpha2, "anchor_rank", "no replaceable base member")
		}
		base[sourceOrder] = true
	}
	return base, nil
}

func containingSourceComponent(reference MultiPolygon, anchor Point) int {
	for i, polygon := range reference {
		if len(polygon) > 0 && pointInRing(anchor, polygon[0]) {
			return i
		}
	}
	return -1
}

func sortedSourceOrders(set map[int]bool) []int {
	out := make([]int, 0, len(set))
	for sourceOrder := range set {
		out = append(out, sourceOrder)
	}
	sort.Ints(out)
	return out
}

func ValidateOverride(o Override, corpus, profile string) error {
	if len(o.ISO) != 2 || o.Version == "" || o.Reason == "" {
		return fail(ErrOverride, o.ISO, "metadata", "iso, version, and reason are required")
	}
	if corpus != "" && o.Corpus != "" && o.Corpus != corpus {
		return fail(ErrOverride, o.ISO, "corpus", "override does not apply")
	}
	if profile != "" && o.Profile != "" && o.Profile != profile {
		return fail(ErrOverride, o.ISO, "profile", "override does not apply")
	}
	if o.Rotation != nil && (!finite(*o.Rotation) || math.Abs(*o.Rotation) > 180) {
		return fail(ErrOverride, o.ISO, "rotation", "must be within [-180,180]")
	}
	if o.MinimumParts < 0 {
		return fail(ErrOverride, o.ISO, "minimum_parts", "must be nonnegative")
	}
	return nil
}
