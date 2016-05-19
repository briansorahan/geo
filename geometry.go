package geo

import (
	"encoding/json"
	"fmt"
)

// Geometry types.
const (
	PointType   = "Point"
	LineType    = "LineString"
	PolygonType = "Polygon"
)

// Geometry defines the interface of every geometry type.
type Geometry interface {
	json.Marshaler
	Compare(other Geometry) bool
	String() string
}

// geometry is a utility type used to unmarshal geometries from JSON.
type geometry struct {
	Type        string          `json:"type"`
	Coordinates json.RawMessage `json:"coordinates"`
}

// Geometry returns a Geometry, or an error if Type is invalid.
func (g geometry) Geometry() (Geometry, error) {
	switch g.Type {
	default:
		return nil, fmt.Errorf("unrecognized geometry type: %s", g.Type)
	case PointType:
		p := &Point{}
		if err := json.Unmarshal(g.Coordinates, p); err != nil {
			return nil, err
		}
		return p, nil
	case LineType:
		l := &Line{}
		if err := json.Unmarshal(g.Coordinates, l); err != nil {
			return nil, err
		}
		return l, nil
	case PolygonType:
		p := &Polygon{}
		if err := json.Unmarshal(g.Coordinates, p); err != nil {
			return nil, err
		}
		return p, nil
	}
}
