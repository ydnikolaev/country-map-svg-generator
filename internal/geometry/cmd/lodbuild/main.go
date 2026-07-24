package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"time"

	"github.com/yuranikolaev/country-map-svg-generator/internal/catalog"
	geometry "github.com/yuranikolaev/country-map-svg-generator/internal/geometry"
)

const expectedGeometryCount = 283

type recipe struct {
	SchemaVersion int    `json:"schema_version"`
	LODVersion    string `json:"lod_version"`
	SourceCorpus  string `json:"source_corpus"`
	Mapshaper     struct {
		Version              string  `json:"version"`
		Integrity            string  `json:"integrity"`
		Shasum               string  `json:"shasum"`
		Algorithm            string  `json:"algorithm"`
		Weighting            float64 `json:"weighting"`
		IntersectionRollback bool    `json:"intersection_rollback"`
		KeepShapes           bool    `json:"keep_shapes"`
		Clean                bool    `json:"clean"`
		Planar               bool    `json:"planar"`
	} `json:"mapshaper"`
	Projection struct {
		Version             string  `json:"version"`
		CoordinateSpace     string  `json:"coordinate_space"`
		BuildRotation       float64 `json:"build_rotation"`
		YAxis               string  `json:"y_axis"`
		Flatness            float64 `json:"flatness"`
		CoordinatePrecision float64 `json:"coordinate_precision"`
	} `json:"projection"`
	Tiers struct {
		Compact struct {
			Resolution         int     `json:"resolution"`
			MaximumFittedScale float64 `json:"maximum_fitted_scale"`
		} `json:"compact"`
		Standard struct {
			Resolution         int     `json:"resolution"`
			MaximumFittedScale float64 `json:"maximum_fitted_scale"`
		} `json:"standard"`
		Source struct {
			MaximumFittedScale *float64 `json:"maximum_fitted_scale"`
		} `json:"source"`
	} `json:"tiers"`
	Budgets map[string]struct {
		HardPathBytes     int `json:"hard_path_bytes"`
		AdvisoryPathBytes int `json:"advisory_path_bytes"`
	} `json:"budgets"`
	CanonicalQuantization  float64 `json:"canonical_quantization"`
	MaximumProjectedPoints int     `json:"maximum_projected_points"`
}

type omission struct {
	GeometryID    string `json:"geometry_id"`
	SourcePolygon int    `json:"source_polygon"`
	SourceRing    int    `json:"source_ring"`
	Reason        string `json:"reason"`
}

type tierArtifact struct {
	SchemaVersion int                             `json:"schema_version"`
	LODVersion    string                          `json:"lod_version"`
	Tier          string                          `json:"tier"`
	SourceCorpus  string                          `json:"source_corpus"`
	RecipeSHA256  string                          `json:"recipe_sha256"`
	ToolVersion   string                          `json:"tool_version"`
	Resolution    int                             `json:"resolution"`
	Geometries    []geometry.ProjectedLODGeometry `json:"geometries"`
	Omissions     []omission                      `json:"omissions"`
}

type buildResult struct {
	Compact, Standard                 []byte
	CompactArtifact, StandardArtifact tierArtifact
	Duration                          time.Duration
}

func main() {
	out := flag.String("out", "", "output directory for temporary spike artifacts")
	diagnostic := flag.Bool("diagnostic", false, "run the RUN-010 T0A.1 semantic diagnostic")
	representative := flag.Bool("representative", false, "build RUN-011 T0A.2 representative oracle artifacts")
	indonesiaBoundary := flag.String("indonesia-boundary-out", "", "write RUN-012 Indonesia 120/121 regression JSON")
	semanticOut := flag.String("semantic-out", "", "canonical semantic diagnostic JSON")
	timingOut := flag.String("timing-out", "", "diagnostic timing receipt JSON")
	ladder := flag.Bool("ladder", false, "build the AM-005/DEC-009 committed silhouette ladder artifact")
	ladderArtifactOut := flag.String("ladder-artifact-out", "", "output path for the ladder artifact JSON (defaults to internal/geometry/lod/ladder.artifact.json)")
	ladderExplain := flag.String("ladder-explain", "", "replay the ladder search for one alpha2 and print every rejected rung")
	flag.Parse()
	if *ladderExplain != "" {
		if *out == "" {
			fmt.Fprintln(os.Stderr, "-out is required with -ladder-explain")
			os.Exit(2)
		}
		c, err := catalog.Embedded()
		if err != nil {
			fatal(err)
		}
		if err := explainLadder(c, *ladderExplain, *out); err != nil {
			fatal(err)
		}
		return
	}
	if *ladder {
		if *out == "" {
			fmt.Fprintln(os.Stderr, "-out is required with -ladder")
			os.Exit(2)
		}
		artifactPath := *ladderArtifactOut
		if artifactPath == "" {
			artifactPath = filepath.Join(sourceRoot(), "internal/geometry/lod/ladder.artifact.json")
		}
		if err := runLadderBuild(*out, artifactPath); err != nil {
			fatal(err)
		}
		return
	}
	if *indonesiaBoundary != "" {
		c, err := catalog.Embedded()
		if err != nil {
			fatal(err)
		}
		if err := buildIndonesiaBoundary(c, *indonesiaBoundary); err != nil {
			fatal(err)
		}
		return
	}
	if *representative {
		if *out == "" {
			fmt.Fprintln(os.Stderr, "-out is required with -representative")
			os.Exit(2)
		}
		c, err := catalog.Embedded()
		if err != nil {
			fatal(err)
		}
		if err := buildRepresentative(c, *out); err != nil {
			fatal(err)
		}
		return
	}
	if *diagnostic {
		if *semanticOut == "" || *timingOut == "" {
			fmt.Fprintln(os.Stderr, "-semantic-out and -timing-out are required with -diagnostic")
			os.Exit(2)
		}
		c, err := catalog.Embedded()
		if err != nil {
			fatal(err)
		}
		if err := runDiagnostic(c, *semanticOut, *timingOut); err != nil {
			fatal(err)
		}
		return
	}
	if *out == "" {
		fmt.Fprintln(os.Stderr, "-out is required")
		os.Exit(2)
	}
	c, err := catalog.Embedded()
	if err != nil {
		fatal(err)
	}
	result, err := build(c, *out)
	if err != nil {
		fatal(err)
	}
	fmt.Printf("compact=%d standard=%d geometries=%d duration=%s\n", len(result.Compact), len(result.Standard), len(result.CompactArtifact.Geometries), result.Duration)
}

func fatal(err error) { fmt.Fprintln(os.Stderr, err); os.Exit(1) }

func build(c *catalog.Corpus, out string) (buildResult, error) {
	start := time.Now()
	root := sourceRoot()
	recipePath := filepath.Join(root, "internal/geometry/lod/v1.recipe.json")
	rawRecipe, err := os.ReadFile(recipePath)
	if err != nil {
		return buildResult{}, err
	}
	var r recipe
	dec := json.NewDecoder(bytes.NewReader(rawRecipe))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&r); err != nil {
		return buildResult{}, fmt.Errorf("recipe: %w", err)
	}
	if err := validateRecipe(r, c); err != nil {
		return buildResult{}, err
	}
	if err := os.MkdirAll(out, 0o755); err != nil {
		return buildResult{}, err
	}
	recipeHash := sha256.Sum256(rawRecipe)
	projectedCatalog := make([]catalog.Geometry, len(c.Geometries))
	projectedByID := make(map[string]geometry.ProjectedLODGeometry, len(c.Geometries))
	for i, source := range c.Geometries {
		record, projectErr := geometry.ProjectLODGeometry(source, r.Projection.Flatness, r.Projection.CoordinatePrecision)
		if projectErr != nil {
			return buildResult{}, fmt.Errorf("project %s: %w", source.ID, projectErr)
		}
		record.SourceCorpus = r.SourceCorpus
		record.RecipeSHA256 = hex.EncodeToString(recipeHash[:])
		projectedCatalog[i], projectedByID[source.ID] = record.Geometry, record
	}
	input := filepath.Join(out, "source.geojson")
	if err := writeSource(input, projectedCatalog); err != nil {
		return buildResult{}, err
	}
	buildTier := func(name string, resolution int) (tierArtifact, []byte, error) {
		geoPath := filepath.Join(out, name+".geojson")
		script := filepath.Join(root, "internal/geometry/lod/tool/build.mjs")
		cmd := exec.Command("node", script, input, geoPath, fmt.Sprint(resolution), fmt.Sprint(r.Mapshaper.Weighting))
		cmd.Dir = root
		output, runErr := cmd.CombinedOutput()
		if runErr != nil {
			return tierArtifact{}, nil, fmt.Errorf("mapshaper %s: %w: %s", name, runErr, output)
		}
		outputGeometries, err := readOutput(geoPath)
		if err != nil {
			return tierArtifact{}, nil, err
		}
		if len(outputGeometries) != expectedGeometryCount {
			return tierArtifact{}, nil, fmt.Errorf("%s coverage: got %d geometries, want %d", name, len(outputGeometries), expectedGeometryCount)
		}
		geometries := make([]geometry.ProjectedLODGeometry, len(outputGeometries))
		for i, outputGeometry := range outputGeometries {
			record, ok := projectedByID[outputGeometry.ID]
			if !ok {
				return tierArtifact{}, nil, fmt.Errorf("%s output has unknown geometry %q", name, outputGeometry.ID)
			}
			record.Geometry = outputGeometry
			geometries[i] = record
		}
		a := tierArtifact{SchemaVersion: 1, LODVersion: r.LODVersion, Tier: name, SourceCorpus: r.SourceCorpus, RecipeSHA256: hex.EncodeToString(recipeHash[:]), ToolVersion: r.Mapshaper.Version, Resolution: resolution, Geometries: geometries, Omissions: collectOmissions(c.Geometries, outputGeometries)}
		raw, err := json.Marshal(a)
		if err != nil {
			return tierArtifact{}, nil, err
		}
		raw = append(raw, '\n')
		return a, raw, os.WriteFile(filepath.Join(out, name+".artifact.json"), raw, 0o644)
	}
	compact, compactRaw, err := buildTier("compact", r.Tiers.Compact.Resolution)
	if err != nil {
		return buildResult{}, err
	}
	standard, standardRaw, err := buildTier("standard", r.Tiers.Standard.Resolution)
	if err != nil {
		return buildResult{}, err
	}
	return buildResult{Compact: compactRaw, Standard: standardRaw, CompactArtifact: compact, StandardArtifact: standard, Duration: time.Since(start)}, nil
}

func sourceRoot() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Clean(filepath.Join(filepath.Dir(file), "../../../.."))
}

func validateRecipe(r recipe, c *catalog.Corpus) error {
	if r.SchemaVersion != 1 || r.LODVersion != "v1" {
		return errors.New("unsupported recipe version")
	}
	if r.SourceCorpus != c.Manifest.Identity {
		return fmt.Errorf("recipe corpus %q does not bind embedded %q", r.SourceCorpus, c.Manifest.Identity)
	}
	if r.Mapshaper.Version != "0.7.44" || r.Mapshaper.Integrity != "sha512-3Cx+IABMXt1G28Y8J7oalW5P5VYyt1vHz5FO+KkV51EHnJioo9h9maO+u+4IyCPZ9Mh1hqchgomFmc68GtFwQQ==" || r.Mapshaper.Shasum != "e08fc40d50347698ea07d2f82ead42a98b1c62f0" {
		return errors.New("mapshaper pin mismatch")
	}
	if r.Mapshaper.Algorithm != "weighted_visvalingam" || r.Mapshaper.Weighting != .7 || !r.Mapshaper.IntersectionRollback || !r.Mapshaper.KeepShapes || !r.Mapshaper.Clean || !r.Mapshaper.Planar {
		return errors.New("mapshaper policy mismatch")
	}
	contract := geometry.LODProjectionContractV1()
	if r.Projection.Version != contract.Version || r.Projection.CoordinateSpace != contract.CoordinateSpace ||
		r.Projection.BuildRotation != contract.BuildRotation || r.Projection.YAxis != contract.YAxis ||
		r.Projection.Flatness != contract.Flatness || r.Projection.CoordinatePrecision != contract.CoordinatePrecision {
		return errors.New("projection policy mismatch")
	}
	if r.Tiers.Compact.Resolution < 8 || r.Tiers.Compact.Resolution > 64 || r.Tiers.Standard.Resolution < 64 || r.Tiers.Standard.Resolution > 256 {
		return errors.New("resolution outside plan bounds")
	}
	if r.Tiers.Compact.MaximumFittedScale != 240 || r.Tiers.Standard.MaximumFittedScale != 700 {
		return errors.New("scale thresholds outside accepted calibration")
	}
	return nil
}

type featureCollection struct {
	Type     string    `json:"type"`
	Features []feature `json:"features"`
}
type feature struct {
	Type       string            `json:"type"`
	Properties map[string]string `json:"properties"`
	Geometry   featureGeometry   `json:"geometry"`
}
type featureGeometry struct {
	Type        string               `json:"type"`
	Coordinates catalog.MultiPolygon `json:"coordinates"`
}

func writeSource(path string, geometries []catalog.Geometry) error {
	sorted := append([]catalog.Geometry(nil), geometries...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].ID < sorted[j].ID })
	fc := featureCollection{Type: "FeatureCollection", Features: make([]feature, len(sorted))}
	for i, g := range sorted {
		fc.Features[i] = feature{Type: "Feature", Properties: map[string]string{"geometry_id": g.ID}, Geometry: featureGeometry{Type: "MultiPolygon", Coordinates: g.Coordinates}}
	}
	raw, err := json.Marshal(fc)
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(raw, '\n'), 0o644)
}
func readOutput(path string) ([]catalog.Geometry, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var fc struct {
		Type     string `json:"type"`
		Features []struct {
			Properties map[string]string `json:"properties"`
			Geometry   struct {
				Type        string          `json:"type"`
				Coordinates json.RawMessage `json:"coordinates"`
			} `json:"geometry"`
		} `json:"features"`
	}
	if err := json.Unmarshal(raw, &fc); err != nil {
		return nil, err
	}
	out := make([]catalog.Geometry, 0, len(fc.Features))
	seen := map[string]bool{}
	for _, f := range fc.Features {
		id := f.Properties["geometry_id"]
		if id == "" || seen[id] {
			return nil, fmt.Errorf("invalid geometry id %q", id)
		}
		seen[id] = true
		var coordinates catalog.MultiPolygon
		switch f.Geometry.Type {
		case "MultiPolygon":
			if err := json.Unmarshal(f.Geometry.Coordinates, &coordinates); err != nil {
				return nil, err
			}
		case "Polygon":
			var polygon catalog.Polygon
			if err := json.Unmarshal(f.Geometry.Coordinates, &polygon); err != nil {
				return nil, err
			}
			coordinates = catalog.MultiPolygon{polygon}
		default:
			return nil, fmt.Errorf("geometry %q has unsupported type %q", id, f.Geometry.Type)
		}
		out = append(out, catalog.Geometry{ID: id, Coordinates: coordinates})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}
func collectOmissions(source, lod []catalog.Geometry) []omission {
	index := map[string]catalog.Geometry{}
	for _, g := range lod {
		index[g.ID] = g
	}
	var out []omission
	for _, g := range source {
		candidate := index[g.ID]
		for pi, p := range g.Coordinates {
			if pi >= len(candidate.Coordinates) {
				out = append(out, omission{g.ID, pi, -1, "component_omitted"})
				continue
			}
			for ri := range p {
				if ri >= len(candidate.Coordinates[pi]) {
					out = append(out, omission{g.ID, pi, ri, "ring_omitted"})
				}
			}
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].GeometryID != out[j].GeometryID {
			return out[i].GeometryID < out[j].GeometryID
		}
		if out[i].SourcePolygon != out[j].SourcePolygon {
			return out[i].SourcePolygon < out[j].SourcePolygon
		}
		return out[i].SourceRing < out[j].SourceRing
	})
	return out
}
