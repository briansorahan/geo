package geo

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

const (
	polygonWKTPrefix = `POLYGON(`
	polygonWKTSuffix = `)`
)

var (
	polygonJSONPrefix = []byte(`{"type":"Polygon","coordinates":[`)
	polygonJSONSuffix = []byte(`]}`)
)

// Polygon is a GeoJSON Polygon.
type Polygon [][][2]float64

// Compare compares one polygon to another.
func (p1 Polygon) Compare(g Geometry) bool {
	p2, ok := g.(*Polygon)
	if !ok {
		return false
	}
	if len(p1) != len(*p2) {
		return false
	}
	for i, pp1 := range p1 {
		pp2 := (*p2)[i]
		if len(pp1) != len(pp2) {
			return false
		}
		if !pointsCompare(pp1, pp2) {
			return false
		}
	}
	return true
}

// Contains uses the ray casting algorithm to decide
// if the point is contained in the polygon.
func (polygon Polygon) Contains(point Point) bool {
	intersections := 0
	for _, poly := range polygon {
		for j, vertex := range poly {
			if point.RayhIntersects(poly[(j+1)%len(poly)], vertex) {
				intersections++
			}
		}
	}
	return (intersections % 2) == 1
}

// MarshalJSON returns the GeoJSON representation of the polygon.
func (polygon Polygon) MarshalJSON() ([]byte, error) {
	s := polygonJSONPrefix
	for i, poly := range polygon {
		if i == 0 {
			s = append(s, '[')
		} else {
			s = append(s, ',', '[')
		}
		s = append(s, pointsMarshalJSON(poly, "", "")...)
		s = append(s, ']')
	}
	return append(s, polygonJSONSuffix...), nil
}

// Scan scans a polygon from Well Known Text.
func (polygon *Polygon) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		return polygon.scan(string(v))
	case string:
		return polygon.scan(v)
	default:
		return fmt.Errorf("could not scan polygon from %T", src)
	}
}

// scan scans a polygon from a Well Known Text string.
func (polygon *Polygon) scan(s string) error {
	if i := strings.Index(s, polygonWKTPrefix); i != 0 {
		return fmt.Errorf("malformed polygon %s", s)
	}
	l := len(s)

	if s[l-len(polygonWKTSuffix):] != polygonWKTSuffix {
		return fmt.Errorf("malformed polygon %s", s)
	}
	s = s[len(polygonWKTPrefix) : l-len(polygonWKTSuffix)]

	// empty the polygon
	*polygon = Polygon{}

	// get the coordinates
	polygons := strings.Split(s, "),(")
	for _, ss := range polygons {
		points, err := pointsScan(ss)
		if err != nil {
			return err
		}
		*polygon = append(*polygon, points)
	}
	return nil
}

// Value converts a point to Well Known Text.
func (polygon Polygon) Value() (driver.Value, error) {
	return polygon.String(), nil
}

// String converts the polygon to a string.
func (polygon Polygon) String() string {
	if len(polygon) == 0 {
		return "POLYGON EMPTY"
	}
	s := polygonWKTPrefix
	for i, points := range polygon {
		if i == 0 {
			s += pointsString(points)
		} else {
			s += "," + pointsString(points)
		}
	}
	return s + polygonWKTSuffix
}
