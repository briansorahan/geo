package geo

import "encoding/json"

// UnmarshalJSON unmarshals any GeoJSON type from a JSON blob.
func UnmarshalJSON(data []byte) (Geometry, error) {
	geom := &geometry{}
	if err := json.Unmarshal(data, geom); err != nil {
		return nil, err
	}
	switch geom.Type {
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
