package catalog

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

const (
	isoPath               = "sources/iso/2026-07-23/iso-3166-1.json"
	isoReconciliationPath = "sources/iso/2026-07-23/reconciliation.json"
	unPath                = "sources/natural-earth/5.1.1/ne_10m_admin_0_countries_iso.geojson"
	deFactoPath           = "sources/natural-earth/5.1.1/ne_10m_admin_0_countries.geojson"
	capitalsPath          = "sources/natural-earth/5.1.2/ne_10m_populated_places_simple.geojson"
)

type boundaryFallbackPolicy struct {
	DeFactoIdenticalToUN []string `json:"de_facto_identical_to_un"`
	Reason               string   `json:"reason"`
}

type capitalOverridePolicy struct {
	Overrides []capitalOverride `json:"overrides"`
}

type capitalOverride struct {
	Alpha2      string    `json:"alpha2"`
	Replace     bool      `json:"replace,omitempty"`
	Capitals    []Capital `json:"capitals,omitempty"`
	PrimaryName string    `json:"primary_name,omitempty"`
	PrimaryID   string    `json:"primary_id,omitempty"`
}

type profileOraclePolicy struct {
	Expectations []ProfileExpectation `json:"expectations"`
}

type protectedPolicy struct {
	Features []protectedPolicyFeature `json:"features"`
}

type protectedPolicyFeature struct {
	Alpha2       string `json:"alpha2"`
	Name         string `json:"name"`
	Anchor       Point  `json:"anchor"`
	MinimumParts int    `json:"minimum_parts"`
}

func Compile(dataRoot string) (*Corpus, error) {
	receipts, err := loadAndVerifyReceipts(dataRoot)
	if err != nil {
		return nil, err
	}
	for _, path := range []string{isoPath, isoReconciliationPath, unPath, deFactoPath, capitalsPath} {
		if err := requireReceipt(receipts, path); err != nil {
			return nil, err
		}
	}
	isoEntities, err := loadISO(filepath.Join(dataRoot, filepath.FromSlash(isoPath)))
	if err != nil {
		return nil, err
	}
	if err := loadISOReconciliation(dataRoot, filepath.Join(dataRoot, filepath.FromSlash(isoReconciliationPath)), isoEntities); err != nil {
		return nil, err
	}
	unFeatures, err := loadFeatures(filepath.Join(dataRoot, filepath.FromSlash(unPath)))
	if err != nil {
		return nil, err
	}
	deFactoFeatures, err := loadFeatures(filepath.Join(dataRoot, filepath.FromSlash(deFactoPath)))
	if err != nil {
		return nil, err
	}
	capitalFeatures, err := loadFeatures(filepath.Join(dataRoot, filepath.FromSlash(capitalsPath)))
	if err != nil {
		return nil, err
	}
	var fallback boundaryFallbackPolicy
	if err := decodeStrict(filepath.Join(dataRoot, "policy", "boundary-fallbacks.json"), &fallback); err != nil {
		return nil, err
	}
	var overrides capitalOverridePolicy
	if err := decodeStrict(filepath.Join(dataRoot, "policy", "capital-overrides.json"), &overrides); err != nil {
		return nil, err
	}
	var protections protectedPolicy
	if err := decodeStrict(filepath.Join(dataRoot, "policy", "protected-features.json"), &protections); err != nil {
		return nil, err
	}
	var oracle profileOraclePolicy
	if err := decodeStrict(filepath.Join(dataRoot, "policy", "boundary-profile-oracle.json"), &oracle); err != nil {
		return nil, err
	}

	unByAlpha3, err := countryGeometryByAlpha3(filepath.Join(dataRoot, filepath.FromSlash(unPath)), unFeatures, true)
	if err != nil {
		return nil, err
	}
	deFactoByAlpha3, err := countryGeometryByAlpha3(filepath.Join(dataRoot, filepath.FromSlash(deFactoPath)), deFactoFeatures, false)
	if err != nil {
		return nil, err
	}
	fallbackSet := map[string]bool{}
	for _, code := range fallback.DeFactoIdenticalToUN {
		if !alpha3Pattern.MatchString(code) || fallbackSet[code] {
			return nil, diagnostic("boundary-fallbacks.json", code, "de_facto_identical_to_un", "fallback-unique", "invalid or duplicate fallback")
		}
		fallbackSet[code] = true
	}
	if fallback.Reason == "" {
		return nil, diagnostic("boundary-fallbacks.json", "-", "reason", "fallback-documented", "fallback reason is required")
	}
	capitals, err := compileCapitals(capitalFeatures, overrides)
	if err != nil {
		return nil, err
	}
	protected, err := compileProtected(protections)
	if err != nil {
		return nil, err
	}

	geometryByID := map[string]Geometry{}
	manifest := Manifest{SchemaVersion: SchemaVersion, CorpusVersion: CorpusVersion, EntityCount: len(isoEntities), Profiles: []string{"un", "de-facto"}}
	for _, sourceEntity := range isoEntities {
		unGeometry, ok := unByAlpha3[sourceEntity.Alpha3]
		if !ok {
			return nil, diagnostic(unPath, sourceEntity.Alpha2, "ISO_A3", "iso-geometry-complete", "missing UN profile feature for %s", sourceEntity.Alpha3)
		}
		unID, err := addGeometry(geometryByID, unGeometry)
		if err != nil {
			return nil, err
		}
		entity := Entity{
			Alpha2:    sourceEntity.Alpha2,
			Alpha3:    sourceEntity.Alpha3,
			Name:      sourceEntity.Name,
			Profiles:  BoundaryProfiles{UN: ProfileResolution{GeometryID: unID}},
			Capitals:  append([]Capital(nil), capitals[sourceEntity.Alpha2]...),
			Protected: append([]ProtectedFeature(nil), protected[sourceEntity.Alpha2]...),
		}
		if deFactoGeometry, ok := deFactoByAlpha3[sourceEntity.Alpha3]; ok {
			deFactoID, err := addGeometry(geometryByID, deFactoGeometry)
			if err != nil {
				return nil, err
			}
			if deFactoID == unID {
				entity.Profiles.DeFacto = ProfileResolution{IdenticalTo: "un"}
			} else {
				entity.Profiles.DeFacto = ProfileResolution{GeometryID: deFactoID}
			}
			if fallbackSet[sourceEntity.Alpha3] {
				return nil, diagnostic("boundary-fallbacks.json", sourceEntity.Alpha2, "de_facto_identical_to_un", "fallback-exact", "fallback is declared but upstream geometry exists")
			}
		} else {
			if !fallbackSet[sourceEntity.Alpha3] {
				return nil, diagnostic(deFactoPath, sourceEntity.Alpha2, "ISO_A3", "profile-complete", "missing de-facto geometry and explicit fallback for %s", sourceEntity.Alpha3)
			}
			entity.Profiles.DeFacto = ProfileResolution{IdenticalTo: "un"}
			delete(fallbackSet, sourceEntity.Alpha3)
		}
		manifest.Entities = append(manifest.Entities, entity)
	}
	if len(fallbackSet) != 0 {
		return nil, diagnostic("boundary-fallbacks.json", "-", "de_facto_identical_to_un", "fallback-exact", "unused fallbacks remain: %v", sortedKeys(fallbackSet))
	}
	if err := validateProfileOracle(&manifest, oracle); err != nil {
		return nil, err
	}
	geometries := make([]Geometry, 0, len(geometryByID))
	for _, geometry := range geometryByID {
		geometries = append(geometries, geometry)
	}
	sort.Slice(geometries, func(i, j int) bool { return geometries[i].ID < geometries[j].ID })
	corpus := &Corpus{Manifest: manifest, Geometries: geometries, Receipts: receipts}
	corpus.Coverage = calculateCoverage(corpus)
	corpus.Manifest.Identity, err = CorpusIdentity(corpus.Manifest, corpus.Geometries, corpus.Receipts, corpus.Coverage)
	if err != nil {
		return nil, err
	}
	if err := Validate(corpus); err != nil {
		return nil, err
	}
	return corpus, nil
}

func countryGeometryByAlpha3(path string, features []feature, requireUnique bool) (map[string]MultiPolygon, error) {
	out := map[string]MultiPolygon{}
	for _, f := range features {
		code := propertyString(f, "ISO_A3_EH")
		if !alpha3Pattern.MatchString(code) {
			continue
		}
		geometry, err := geometryFromFeature(path, code, f)
		if err != nil {
			return nil, err
		}
		if requireUnique {
			if _, exists := out[code]; exists {
				return nil, diagnostic(path, code, "ISO_A3", "geometry-mapping-unique", "duplicate ISO mapping")
			}
			out[code] = geometry
		} else {
			out[code] = append(out[code], geometry...)
		}
	}
	return out, nil
}

func addGeometry(all map[string]Geometry, coordinates MultiPolygon) (string, error) {
	id, err := geometryIdentity(coordinates)
	if err != nil {
		return "", err
	}
	if existing, ok := all[id]; ok {
		a, _ := CanonicalJSON(existing.Coordinates)
		b, _ := CanonicalJSON(coordinates)
		if !bytes.Equal(a, b) {
			return "", fmt.Errorf("geometry digest collision %s", id)
		}
	} else {
		all[id] = Geometry{ID: id, Coordinates: coordinates}
	}
	return id, nil
}

func compileCapitals(features []feature, policy capitalOverridePolicy) (map[string][]Capital, error) {
	out := map[string][]Capital{}
	for _, f := range features {
		adm0 := propertyInt64(f, "adm0cap") == 1
		alternate := propertyInt64(f, "capalt") == 1
		if !adm0 && !alternate {
			continue
		}
		code := propertyString(f, "iso_a2")
		if !alpha2Pattern.MatchString(code) {
			continue
		}
		point := Point{propertyFloat(f, "longitude"), propertyFloat(f, "latitude")}
		if !validCoordinate(point) {
			return nil, diagnostic(capitalsPath, code, "coordinates", "capital-coordinate", "invalid capital coordinate")
		}
		roles := []string{}
		if adm0 {
			roles = append(roles, "national")
		}
		if alternate {
			roles = append(roles, "alternate")
		}
		out[code] = append(out[code], Capital{
			ID:    fmt.Sprintf("%s-ne-%d", code, propertyInt64(f, "ne_id")),
			Name:  propertyString(f, "name"),
			Point: Point{round6(point[0]), round6(point[1])},
			Roles: roles,
		})
	}
	overrideByCode := map[string]capitalOverride{}
	for _, override := range policy.Overrides {
		if !alpha2Pattern.MatchString(override.Alpha2) {
			return nil, diagnostic("capital-overrides.json", override.Alpha2, "overrides", "capital-override-unique", "invalid or duplicate override")
		}
		if _, exists := overrideByCode[override.Alpha2]; exists {
			return nil, diagnostic("capital-overrides.json", override.Alpha2, "overrides", "capital-override-unique", "invalid or duplicate override")
		}
		if override.PrimaryName != "" && override.PrimaryID != "" {
			return nil, diagnostic("capital-overrides.json", override.Alpha2, "primary", "capital-primary-reference", "set only one of primary_name or primary_id")
		}
		overrideByCode[override.Alpha2] = override
	}
	for code, override := range overrideByCode {
		if override.Replace {
			out[code] = nil
		}
		ids := map[string]bool{}
		for _, capital := range out[code] {
			ids[capital.ID] = true
		}
		for _, capital := range override.Capitals {
			if capital.ID == "" || capital.Name == "" || !validCapitalRoles(capital.Roles) || !validCoordinate(capital.Point) || ids[capital.ID] {
				return nil, diagnostic("capital-overrides.json", code, "capitals", "capital-override-record", "invalid or duplicate capital %q", capital.ID)
			}
			capital.Primary = false
			capital.Point = Point{round6(capital.Point[0]), round6(capital.Point[1])}
			out[code] = append(out[code], capital)
			ids[capital.ID] = true
		}
	}
	for code := range out {
		sort.Slice(out[code], func(i, j int) bool {
			if out[code][i].Name != out[code][j].Name {
				return out[code][i].Name < out[code][j].Name
			}
			return out[code][i].ID < out[code][j].ID
		})
		override, hasOverride := overrideByCode[code]
		if len(out[code]) == 1 && (!hasOverride || (override.PrimaryName == "" && override.PrimaryID == "")) {
			out[code][0].Primary = true
		} else {
			matches := 0
			for i := range out[code] {
				out[code][i].Primary = false
				if (override.PrimaryName != "" && out[code][i].Name == override.PrimaryName) ||
					(override.PrimaryID != "" && out[code][i].ID == override.PrimaryID) {
					out[code][i].Primary = true
					matches++
				}
			}
			if matches != 1 {
				return nil, diagnostic("capital-overrides.json", code, "primary", "capital-primary-reference", "missing selection for %d capitals", len(out[code]))
			}
		}
		delete(overrideByCode, code)
	}
	if len(overrideByCode) != 0 {
		return nil, diagnostic("capital-overrides.json", "-", "overrides", "capital-override-reference", "unused overrides: %v", sortedCapitalOverrideKeys(overrideByCode))
	}
	return out, nil
}

func validCapitalRoles(roles []string) bool {
	if len(roles) == 0 {
		return false
	}
	seen := map[string]bool{}
	for _, role := range roles {
		if role == "" || seen[role] {
			return false
		}
		seen[role] = true
	}
	return true
}

func validateProfileOracle(manifest *Manifest, policy profileOraclePolicy) error {
	if len(policy.Expectations) != len(reviewedProfileOracle) {
		return diagnostic("boundary-profile-oracle.json", "-", "expectations", "profile-oracle-exact", "expected %d reviewed rows", len(reviewedProfileOracle))
	}
	entities := map[string]Entity{}
	for _, entity := range manifest.Entities {
		entities[entity.Alpha2] = entity
	}
	seen := map[string]bool{}
	for _, expectation := range policy.Expectations {
		if seen[expectation.Alpha2] || reviewedProfileOracle[expectation.Alpha2] != expectation {
			return diagnostic("boundary-profile-oracle.json", expectation.Alpha2, "expectations", "profile-oracle-exact", "row is duplicate or differs from reviewed oracle")
		}
		seen[expectation.Alpha2] = true
		entity := entities[expectation.Alpha2]
		if entity.Profiles.UN.GeometryID != expectation.UNGeometryID || entity.Profiles.DeFacto.GeometryID != expectation.DeFactoGeometry {
			return diagnostic("boundary-profile-oracle.json", expectation.Alpha2, "expectations", "profile-oracle-match", "compiled profiles differ from reviewed identities")
		}
	}
	return nil
}

func compileProtected(policy protectedPolicy) (map[string][]ProtectedFeature, error) {
	out := map[string][]ProtectedFeature{}
	for _, item := range policy.Features {
		if !alpha2Pattern.MatchString(item.Alpha2) || item.Name == "" || item.MinimumParts < 1 || !validCoordinate(item.Anchor) {
			return nil, diagnostic("protected-features.json", item.Alpha2, "features", "protected-policy", "invalid protected feature %q", item.Name)
		}
		out[item.Alpha2] = append(out[item.Alpha2], ProtectedFeature{Name: item.Name, Anchor: item.Anchor, MinimumParts: item.MinimumParts})
	}
	return out, nil
}

func calculateCoverage(corpus *Corpus) Coverage {
	coverage := Coverage{EntityCount: len(corpus.Manifest.Entities)}
	for _, entity := range corpus.Manifest.Entities {
		coverage.UNExplicit++
		if entity.Profiles.DeFacto.GeometryID != "" {
			coverage.DeFactoExplicit++
		} else {
			coverage.DeFactoIdentical++
		}
		if len(entity.Capitals) > 0 {
			coverage.CapitalEntities++
			coverage.CapitalRecords += len(entity.Capitals)
		}
		if len(entity.Protected) > 0 {
			coverage.Protected = append(coverage.Protected, entity.Alpha2)
		}
	}
	return coverage
}

func sortedKeys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for key := range m {
		out = append(out, key)
	}
	sort.Strings(out)
	return out
}

func sortedCapitalOverrideKeys(m map[string]capitalOverride) []string {
	out := make([]string, 0, len(m))
	for key := range m {
		out = append(out, key)
	}
	sort.Strings(out)
	return out
}

func writeJSON(path string, value any) error {
	b, err := CanonicalJSON(value)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}
