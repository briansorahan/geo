package geo

import (
	"encoding/json"
	"errors"
)

// geoJSON is a utility type used to unmarshal GeoJSON blobs.
type geoJSON struct {
	Type        string          `json:"type"`
	Geometries  []geometry      `json:"geometries"`
	Geometry    geometry        `json:"geometry"`
	Coordinates json.RawMessage `json:"coordinates"`
	Properties  interface{}     `json:"properties"`
	Radius      float64         `json:"radius"`
}

// Unmarshal unmarshals a chunk of GeoJSON.
func Unmarshal(data []byte) (Geometry, error) {
	gj := geoJSON{}

	if err := json.Unmarshal(data, &gj); err != nil {
		return nil, err
	}
	switch gj.Type {
	case CircleType, LineType, MultiPointType, PointType, PolygonType:
		return UnmarshalGeometry(data)
	case FeatureType:
		var (
			f   = &Feature{}
			err = f.UnmarshalJSON(data)
		)
		return f, err
	case FeatureCollectionType:
		var (
			fc  = &FeatureCollection{}
			err = fc.UnmarshalJSON(data)
		)
		return fc, err
	case GeometryCollectionType:
		var (
			gc  = &GeometryCollection{}
			err = json.Unmarshal(data, gc)
		)
		return gc, err
	default:
		return nil, errors.New("unrecognized type: " + gj.Type)
	}
}
