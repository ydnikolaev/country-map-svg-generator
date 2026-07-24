package geometry

import (
	"fmt"
	"github.com/yuranikolaev/country-map-svg-generator/internal/catalog"
)

func InputFromCatalog(c *catalog.Corpus, alpha2, profile, preset string) (Input, error) {
	var entity *catalog.Entity
	for i := range c.Manifest.Entities {
		if c.Manifest.Entities[i].Alpha2 == alpha2 {
			entity = &c.Manifest.Entities[i]
			break
		}
	}
	if entity == nil {
		return Input{}, fmt.Errorf("unknown entity %q", alpha2)
	}
	var id string
	switch profile {
	case "un":
		id = entity.Profiles.UN.GeometryID
	case "de_facto":
		id = entity.Profiles.DeFacto.GeometryID
		if id == "" && entity.Profiles.DeFacto.IdenticalTo == "un" {
			id = entity.Profiles.UN.GeometryID
		}
	default:
		return Input{}, fmt.Errorf("unknown profile %q", profile)
	}
	var source *catalog.Geometry
	for i := range c.Geometries {
		if c.Geometries[i].ID == id {
			source = &c.Geometries[i]
			break
		}
	}
	if source == nil {
		return Input{}, fmt.Errorf("geometry %q missing", id)
	}
	markers := make([]MarkerInput, 0, len(entity.Capitals))
	for _, m := range entity.Capitals {
		markers = append(markers, MarkerInput{ID: m.ID, Lon: m.Point[0], Lat: m.Point[1]})
	}
	return Input{Entity: *entity, Geometry: *source, Profile: profile, CorpusID: c.Manifest.Identity, Preset: preset, Markers: markers}, nil
}
