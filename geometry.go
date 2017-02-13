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
	MultiLineType          = "MultiLineString"
	LineType               = "LineString"
	MultiPointType         = "MultiPoint"
	PointType              = "Point"
	PolygonType            = "Polygon"
	MultiPolygonType       = "MultiPolygon"
)

// Geometry defines the interface of every geometry type.
type Geometry interface {
	json.Marshaler
	json.Unmarshaler
	sql.Scanner
	driver.Valuer

	Equal(g Geometry) bool
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

// geometry is a utility type used to unmarshal geometries from JSON.
type geometry struct {
	Type        string          `json:"type"`
	Coordinates json.RawMessage `json:"coordinates"`
	Radius      float64         `json:"radius"` // For circles!
	BBox        []float64       `json:"bbox"`
}

// Geometry returns a Geometry, or an error if Type is invalid.
func (g geometry) unmarshalCoordinates() (geom Geometry, err error) {
	switch g.Type {
	default:
		return nil, fmt.Errorf("unrecognized geometry type: %s", g.Type)
	case PointType:
		pt := [3]float64{}
		err = json.Unmarshal(g.Coordinates, &pt)
		p := Point(pt)
		geom = &p
	case MultiPointType:
		mpt := [][3]float64{}
		err = json.Unmarshal(g.Coordinates, &mpt)
		mp := MultiPoint(mpt)
		geom = &mp
	case LineType:
		ln := [][3]float64{}
		err = json.Unmarshal(g.Coordinates, &ln)
		l := Line(ln)
		geom = &l
	case MultiLineType:
		mln := [][][3]float64{}
		err = json.Unmarshal(g.Coordinates, &mln)
		ml := MultiLine(mln)
		geom = &ml
	case PolygonType:
		poly := [][][3]float64{}
		err = json.Unmarshal(g.Coordinates, &poly)
		p := Polygon(poly)
		geom = &p
	case MultiPolygonType:
		mpoly := [][][][3]float64{}
		err = json.Unmarshal(g.Coordinates, &mpoly)
		mp := MultiPolygon(mpoly)
		geom = &mp
	case CircleType:
		center := [3]float64{}
		err = json.Unmarshal(g.Coordinates, &center)
		geom = &Circle{
			Coordinates: center,
			Radius:      g.Radius,
		}
	}
	if len(g.BBox) > 0 {
		return WithBBox(g.BBox, geom), err
	}
	return geom, err
}
