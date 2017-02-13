package geo

import "encoding/json"

type blob struct {
	*geometry

	Features   []feature  `json:"features"`
	Geometries []geometry `json:"geometries"`
}

// UnmarshalJSON unmarshals any GeoJSON type from a JSON blob.
// TODO: unmarshal bbox'es for Feature, FeatureCollection and GeometryCollection
func UnmarshalJSON(data []byte) (Geometry, error) {
	var (
		geom = &geometry{}
		b    = blob{geometry: geom}
	)
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, err
	}
	switch b.Type {
	default:
		return geom.unmarshalCoordinates()
	case FeatureType:
		return unmarshalFeatureBBox(data)
	case FeatureCollectionType:
		return unmarshalFeatureCollectionBBox(data)
	case GeometryCollectionType:
		return unmarshalGeometryCollectionBBox(data)
	}
}
