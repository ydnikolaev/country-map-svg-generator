package catalog

import (
	"fmt"
	"reflect"
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

var reviewedProfileOracle = map[string]ProfileExpectation{
	"CN": {Alpha2: "CN", UNGeometryID: "geo-219b4dc2e9b65c0531a17bc720badedb3e9eff0265df5a4ef380dba7da2f49f2", DeFactoGeometry: "geo-e74d60cf38d89b0bdac8cd48b05ae4fcadc6d3194cb4c95e5cd0dd906ddc33e8"},
	"CY": {Alpha2: "CY", UNGeometryID: "geo-ea39f523134d18aff6536caa94268e94ecf18ed50542ede9e2ce91fec4d80f16", DeFactoGeometry: "geo-c9877083e0158735cd3d7fad1c01fe84c74a603d9229fbb2d09314e86ecccf93"},
	"IL": {Alpha2: "IL", UNGeometryID: "geo-da6a6c11619b5c92a4245714d54e9a0fe9be65529fc204e9a944da5de8b1ed2f", DeFactoGeometry: "geo-511e6b21a65c6725c89eb88e6e6bcadf430bb5d3a9e5b248b7cbda2d8e37bcc2"},
	"IN": {Alpha2: "IN", UNGeometryID: "geo-c5bf1d1cc94463dcae3c6123ac1385fefed8883258c3cc7ddb18de6fe6059143", DeFactoGeometry: "geo-490c45be3be0547c4bf2b99a8ff54262a5c59e3647fc8725ebbc4bc9177bd432"},
	"RU": {Alpha2: "RU", UNGeometryID: "geo-0ecaec385bf37de21bcc17468f32b224fbdb3cc262e009a08fc68a25dab0c377", DeFactoGeometry: "geo-08425a00146a135ad5f5344deebfd13fd87206ad90b12d983d3b8917435f247f"},
}

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
	geometryByID := map[string]Geometry{}
	for _, geometry := range corpus.Geometries {
		if geometry.ID == "" || geometryIDs[geometry.ID] || len(geometry.Coordinates) == 0 {
			return diagnostic("geometries", "-", "id", "geometry-unique-nonempty", "duplicate or empty geometry %q", geometry.ID)
		}
		want, err := geometryIdentity(geometry.Coordinates)
		if err != nil || want != geometry.ID {
			return diagnostic("geometries", "-", "id", "geometry-identity", "geometry id mismatch for %q", geometry.ID)
		}
		geometryIDs[geometry.ID] = true
		geometryByID[geometry.ID] = geometry
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
		if expected, ok := reviewedProfileOracle[entity.Alpha2]; ok &&
			(entity.Profiles.UN.GeometryID != expected.UNGeometryID || entity.Profiles.DeFacto.GeometryID != expected.DeFactoGeometry) {
			return diagnostic("manifest", entity.Alpha2, "profiles", "profile-oracle-exact", "profile identities differ from reviewed expectation")
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
			roleSeen := map[string]bool{}
			for _, role := range capital.Roles {
				if role == "" || roleSeen[role] {
					return diagnostic("manifest", entity.Alpha2, "capitals.roles", "capital-roles", "empty or duplicate role for %q", capital.ID)
				}
				roleSeen[role] = true
			}
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
			resolution := entity.Profiles.UN
			geometry := geometryByID[resolution.GeometryID]
			if len(geometry.Coordinates) < feature.MinimumParts {
				return diagnostic("manifest", entity.Alpha2, "protected_features.minimum_parts", "protected-parts", "feature %q requires %d parts, geometry has %d", feature.Name, feature.MinimumParts, len(geometry.Coordinates))
			}
			if !pointInMultiPolygon(feature.Anchor, geometry.Coordinates) {
				return diagnostic("manifest", entity.Alpha2, "protected_features.anchor", "protected-containment", "feature %q anchor is outside entity geometry", feature.Name)
			}
		}
	}
	for code, name := range requiredProtected {
		if !protected[code] {
			return diagnostic("manifest", code, "protected_features", "protected-coverage", "required feature %q is uncovered", name)
		}
	}
	if err := validatePublishedReceipts(corpus.Receipts); err != nil {
		return err
	}
	wantCoverage := calculateCoverage(corpus)
	if !reflect.DeepEqual(wantCoverage, corpus.Coverage) {
		return diagnostic("coverage", "-", "records", "coverage-exact", "published coverage does not match corpus")
	}
	want, err := CorpusIdentity(corpus.Manifest, corpus.Geometries, corpus.Receipts, corpus.Coverage)
	if err != nil || want != corpus.Manifest.Identity {
		return diagnostic("manifest", "-", "identity", "corpus-identity", "content identity mismatch")
	}
	return nil
}

func validatePublishedReceipts(receipts []Receipt) error {
	seen := map[string]bool{}
	for i, receipt := range receipts {
		if receipt.ID == "" || receipt.Path == "" || receipt.Origin == "" || receipt.UpstreamVersion == "" ||
			receipt.License == "" || receipt.Transformation == "" || !digestPattern.MatchString(receipt.SHA256) ||
			seen[receipt.ID] || (i > 0 && receipts[i-1].ID >= receipt.ID) {
			return diagnostic("receipts", receipt.ID, "records", "receipt-index-exact", "receipt is incomplete, duplicate, or unsorted")
		}
		seen[receipt.ID] = true
	}
	if len(receipts) != 5 {
		return diagnostic("receipts", "-", "records", "receipt-index-exact", "expected 5 authoritative receipts, got %d", len(receipts))
	}
	return nil
}

func pointInMultiPolygon(point Point, multi MultiPolygon) bool {
	for _, polygon := range multi {
		if len(polygon) == 0 || !pointInRing(point, polygon[0]) {
			continue
		}
		insideHole := false
		for _, hole := range polygon[1:] {
			if pointInRing(point, hole) {
				insideHole = true
				break
			}
		}
		if !insideHole {
			return true
		}
	}
	return false
}

func pointInRing(point Point, ring Ring) bool {
	inside := false
	for i, j := 0, len(ring)-1; i < len(ring); j, i = i, i+1 {
		a, b := ring[i], ring[j]
		if pointOnSegment(point, a, b) {
			return true
		}
		if (a[1] > point[1]) != (b[1] > point[1]) &&
			point[0] < (b[0]-a[0])*(point[1]-a[1])/(b[1]-a[1])+a[0] {
			inside = !inside
		}
	}
	return inside
}

func pointOnSegment(p, a, b Point) bool {
	const epsilon = 1e-9
	cross := (p[1]-a[1])*(b[0]-a[0]) - (p[0]-a[0])*(b[1]-a[1])
	if cross < -epsilon || cross > epsilon {
		return false
	}
	return p[0] >= min(a[0], b[0])-epsilon && p[0] <= max(a[0], b[0])+epsilon &&
		p[1] >= min(a[1], b[1])-epsilon && p[1] <= max(a[1], b[1])+epsilon
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
