package geo

import (
	"encoding/json"
	"errors"
)

// Subset of GeoJSON Types.
const (
	TypePoint   = "Point"
	TypePolygon = "Polygon"
)

// Geometry is the basic geometry type.
type Geometry interface {
	BoundingBox() [][2]float64
}

// geometry provides an easy way to json-marshal Geometry types.
type geometry struct {
	Type        string   `json:"type"`
	Coordinates Geometry `json:"coordinates"`
}

// Feature is a geojson feature.
type Feature struct {
	Geom       Geometry
	Properties map[string]string
}

// MarshalJSON marshals the given geometry to JSON.
func MarshalJSON(g Geometry) ([]byte, error) {
	var typ string
	switch h := g.(type) {
	default:
		return nil, errors.New("unrecognized geometry")
	case Point:
		typ = TypePoint
	case Polygon:
		if err := h.validate(); err != nil {
			return nil, err
		}
		typ = TypePolygon
	}
	return json.Marshal(geometry{Type: typ, Coordinates: g})
}
