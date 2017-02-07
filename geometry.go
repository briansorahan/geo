package geo

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

// Geometry types.
const (
	CircleType             = "Circle"
	FeatureCollectionType  = "FeatureCollection"
	FeatureType            = "Feature"
	GeometryCollectionType = "GeometryCollection"
	MultiLineType          = "MultiLine"
	LineType               = "LineString"
	MultiPointType         = "MultiPoint"
	PointType              = "Point"
	PolygonType            = "Polygon"
)

// Geometry defines the interface of every geometry type.
type Geometry interface {
	json.Marshaler
	json.Unmarshaler
	sql.Scanner
	driver.Valuer

	Compare(g Geometry) bool
	Contains(p Point) bool
	String() string
}

// ScanGeometry scans a geometry from well known text.
func ScanGeometry(s string) (Geometry, error) {
	if i := strings.Index(s, pointWKTPrefix); i == 0 {
		pt := &Point{}
		if err := pt.Scan(s); err != nil {
			return nil, err
		}
		return pt, nil
	}
	if i := strings.Index(s, lineWKTPrefix); i == 0 {
		l := &Line{}
		if err := l.Scan(s); err != nil {
			return nil, err
		}
		return l, nil
	}
	if i := strings.Index(s, polygonWKTPrefix); i == 0 {
		p := &Polygon{}
		if err := p.Scan(s); err != nil {
			return nil, err
		}
		return p, nil
	}
	if i := strings.Index(s, circleWKTPrefix); i == 0 {
		c := &Circle{}
		if err := c.Scan(s); err != nil {
			return nil, err
		}
		return c, nil
	}
	return nil, fmt.Errorf("unrecognized geometry: %s", s)
}

// UnmarshalGeometry unmarshals a geometry from geojson data.
func UnmarshalGeometry(data []byte) (Geometry, error) {
	g := &geometry{}
	if err := json.Unmarshal(data, g); err != nil {
		return nil, err
	}
	return g.unmarshalCoordinates()
}

// geometry is a utility type used to unmarshal geometries from JSON.
type geometry struct {
	Type        string          `json:"type"`
	Coordinates json.RawMessage `json:"coordinates"`
	Radius      float64         `json:"radius"` // For circles!
}

// Geometry returns a Geometry, or an error if Type is invalid.
func (g geometry) unmarshalCoordinates() (Geometry, error) {
	switch g.Type {
	default:
		return nil, fmt.Errorf("unrecognized geometry type: %s", g.Type)
	case PointType:
		pt := [2]float64{}
		if err := json.Unmarshal(g.Coordinates, &pt); err != nil {
			return nil, err
		}
		p := Point(pt)
		return &p, nil
	case MultiPointType:
		mp := MultiPoint{}
		if err := json.Unmarshal(g.Coordinates, &mp); err != nil {
			return nil, err
		}
		return &mp, nil
	case LineType:
		ln := [][2]float64{}
		if err := json.Unmarshal(g.Coordinates, &ln); err != nil {
			return nil, err
		}
		l := Line(ln)
		return &l, nil
	case MultiLineType:
		mln := [][][2]float64{}
		if err := json.Unmarshal(g.Coordinates, &mln); err != nil {
			return nil, err
		}
		ml := MultiLine(mln)
		return &ml, nil
	case PolygonType:
		poly := [][][2]float64{}
		if err := json.Unmarshal(g.Coordinates, &poly); err != nil {
			return nil, err
		}
		p := Polygon(poly)
		return &p, nil
	case CircleType:
		center := [2]float64{}
		if err := json.Unmarshal(g.Coordinates, &center); err != nil {
			return nil, err
		}
		return &Circle{
			Coordinates: center,
			Radius:      g.Radius,
		}, nil
	}
}
