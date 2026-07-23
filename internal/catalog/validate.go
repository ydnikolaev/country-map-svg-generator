package catalog

import (
	"fmt"
	"sort"
)

var requiredProtected = map[string]string{
	"ID": "Indonesian archipelago",
	"KI": "Kiribati island groups",
	"MC": "Monaco microstate",
	"NR": "Nauru island",
	"SM": "San Marino microstate",
	"VA": "Vatican City microstate",
}

var requiredDistinctProfiles = map[string]bool{"CN": true, "CY": true, "IL": true, "IN": true, "RU": true}

func Validate(corpus *Corpus) error {
	if corpus.Manifest.SchemaVersion != SchemaVersion || corpus.Manifest.CorpusVersion != CorpusVersion {
		return diagnostic("manifest", "-", "version", "corpus-version", "unsupported schema/corpus version")
	}
	if len(corpus.Manifest.Entities) != 249 || corpus.Manifest.EntityCount != 249 {
		return diagnostic("manifest", "-", "entities", "iso-count-249", "manifest has %d entities", len(corpus.Manifest.Entities))
	}
	if len(corpus.Manifest.Profiles) != 2 || corpus.Manifest.Profiles[0] != "un" || corpus.Manifest.Profiles[1] != "de-facto" {
		return diagnostic("manifest", "-", "profiles", "boundary-profiles", "expected ordered un and de-facto profiles")
	}
	geometryIDs := map[string]bool{}
	for _, geometry := range corpus.Geometries {
		if geometry.ID == "" || geometryIDs[geometry.ID] || len(geometry.Coordinates) == 0 {
			return diagnostic("geometries", "-", "id", "geometry-unique-nonempty", "duplicate or empty geometry %q", geometry.ID)
		}
		want, err := geometryIdentity(geometry.Coordinates)
		if err != nil || want != geometry.ID {
			return diagnostic("geometries", "-", "id", "geometry-identity", "geometry id mismatch for %q", geometry.ID)
		}
		geometryIDs[geometry.ID] = true
	}
	if !sort.SliceIsSorted(corpus.Geometries, func(i, j int) bool { return corpus.Geometries[i].ID < corpus.Geometries[j].ID }) {
		return diagnostic("geometries", "-", "records", "geometry-sorted", "geometry records are not sorted")
	}
	seen, protected := map[string]bool{}, map[string]bool{}
	for i, entity := range corpus.Manifest.Entities {
		if !alpha2Pattern.MatchString(entity.Alpha2) || !alpha3Pattern.MatchString(entity.Alpha3) || entity.Name == "" {
			return diagnostic("manifest", entity.Alpha2, "identity", "iso-record", "invalid entity identity")
		}
		if seen[entity.Alpha2] || (i > 0 && corpus.Manifest.Entities[i-1].Alpha2 >= entity.Alpha2) {
			return diagnostic("manifest", entity.Alpha2, "alpha2", "iso-unique-sorted", "duplicate or unsorted entity")
		}
		seen[entity.Alpha2] = true
		if err := validateResolution(entity.Alpha2, "un", entity.Profiles.UN, geometryIDs); err != nil {
			return err
		}
		if entity.Profiles.UN.IdenticalTo != "" {
			return diagnostic("manifest", entity.Alpha2, "profiles.un", "profile-explicit", "un must own explicit geometry")
		}
		if err := validateResolution(entity.Alpha2, "de-facto", entity.Profiles.DeFacto, geometryIDs); err != nil {
			return err
		}
		if requiredDistinctProfiles[entity.Alpha2] && entity.Profiles.DeFacto.GeometryID == "" {
			return diagnostic("manifest", entity.Alpha2, "profiles.de_facto", "dispute-profile-expectation", "reviewed profile must remain distinct from un")
		}
		if entity.Profiles.DeFacto.IdenticalTo != "" && entity.Profiles.DeFacto.IdenticalTo != "un" {
			return diagnostic("manifest", entity.Alpha2, "profiles.de_facto.identical_to", "profile-reference", "only un is a valid identical reference")
		}
		primary := 0
		capitalIDs := map[string]bool{}
		for _, capital := range entity.Capitals {
			if capital.ID == "" || capital.Name == "" || capitalIDs[capital.ID] || len(capital.Roles) == 0 || !validCoordinate(capital.Point) {
				return diagnostic("manifest", entity.Alpha2, "capitals", "capital-record", "invalid or duplicate capital %q", capital.ID)
			}
			capitalIDs[capital.ID] = true
			if capital.Primary {
				primary++
			}
		}
		if len(entity.Capitals) > 0 && primary != 1 {
			return diagnostic("manifest", entity.Alpha2, "capitals.primary", "capital-primary", "expected one primary, got %d", primary)
		}
		for _, feature := range entity.Protected {
			if feature.Name == "" || feature.MinimumParts < 1 || !validCoordinate(feature.Anchor) {
				return diagnostic("manifest", entity.Alpha2, "protected_features", "protected-record", "invalid protected feature %q", feature.Name)
			}
			if feature.Name == requiredProtected[entity.Alpha2] {
				protected[entity.Alpha2] = true
			}
		}
	}
	for code, name := range requiredProtected {
		if !protected[code] {
			return diagnostic("manifest", code, "protected_features", "protected-coverage", "required feature %q is uncovered", name)
		}
	}
	want, err := CorpusIdentity(corpus.Manifest, corpus.Geometries)
	if err != nil || want != corpus.Manifest.Identity {
		return diagnostic("manifest", "-", "identity", "corpus-identity", "content identity mismatch")
	}
	return nil
}

func validateResolution(entity, profile string, resolution ProfileResolution, geometries map[string]bool) error {
	if (resolution.GeometryID == "") == (resolution.IdenticalTo == "") {
		return diagnostic("manifest", entity, "profiles."+profile, "profile-resolution", "must set exactly one of geometry_id or identical_to")
	}
	if resolution.GeometryID != "" && !geometries[resolution.GeometryID] {
		return diagnostic("manifest", entity, "profiles."+profile+".geometry_id", "geometry-reference", "dangling geometry %q", resolution.GeometryID)
	}
	return nil
}

func ResolveGeometry(corpus *Corpus, alpha2, profile string) (Geometry, error) {
	var entity *Entity
	for i := range corpus.Manifest.Entities {
		if corpus.Manifest.Entities[i].Alpha2 == alpha2 {
			entity = &corpus.Manifest.Entities[i]
			break
		}
	}
	if entity == nil {
		return Geometry{}, fmt.Errorf("unknown ISO alpha-2 %q", alpha2)
	}
	var resolution ProfileResolution
	switch profile {
	case "un":
		resolution = entity.Profiles.UN
	case "de-facto":
		resolution = entity.Profiles.DeFacto
		if resolution.IdenticalTo == "un" {
			resolution = entity.Profiles.UN
		}
	default:
		return Geometry{}, fmt.Errorf("unknown boundary profile %q", profile)
	}
	for _, geometry := range corpus.Geometries {
		if geometry.ID == resolution.GeometryID {
			return cloneGeometry(geometry), nil
		}
	}
	return Geometry{}, fmt.Errorf("entity %s has dangling geometry %q", alpha2, resolution.GeometryID)
}

func cloneGeometry(geometry Geometry) Geometry {
	out := Geometry{ID: geometry.ID, Coordinates: make(MultiPolygon, len(geometry.Coordinates))}
	for i, polygon := range geometry.Coordinates {
		out.Coordinates[i] = make(Polygon, len(polygon))
		for j, ring := range polygon {
			out.Coordinates[i][j] = append(Ring(nil), ring...)
		}
	}
	return out
}
