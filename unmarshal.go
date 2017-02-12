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
		var (
			fc  = &FeatureCollection{}
			err = fc.UnmarshalJSON(data)
		)
		return fc, err
	case GeometryCollectionType:
		var (
			gc  = &GeometryCollection{}
			err = gc.UnmarshalJSON(data)
		)
		return gc, err
	}
}
