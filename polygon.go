package geo

import (
	"database/sql/driver"
	"encoding/json"
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
type Polygon [][][3]float64

// Equal compares one polygon to another.
func (polygon Polygon) Equal(g Geometry) bool {
	p, ok := g.(*Polygon)
	if !ok {
		return false
	}
	if len(polygon) != len(*p) {
		return false
	}
	for i, p1 := range polygon {
		p2 := (*p)[i]
		if len(p1) != len(p2) {
			return false
		}
		if !pointsEqual(p1, p2) {
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
	return scan(polygon, src)
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

// UnmarshalJSON unmarshals the polygon from GeoJSON.
func (polygon *Polygon) UnmarshalJSON(data []byte) error {
	g := &geometry{}

	// Never fails since data is always valid JSON.
	_ = json.Unmarshal(data, g)

	if expected, got := PolygonType, g.Type; expected != got {
		return fmt.Errorf("expected type %s, got %s", expected, got)
	}

	p := [][][3]float64{}

	if err := json.Unmarshal(g.Coordinates, &p); err != nil {
		return err
	}

	*polygon = Polygon(p)

	return nil
}

// Value converts a point to Well Known Text.
func (polygon Polygon) Value() (driver.Value, error) {
	return polygon.String(), nil
}

// Transform transforms the geometry point by point.
func (polygon *Polygon) Transform(t Transformer) {
	np := make([][][3]float64, len(*polygon))
	for i, line := range *polygon {
		nl := make([][3]float64, len(line))
		for j, point := range line {
			nl[j] = t.Transform(point)
		}
		np[i] = nl
	}
	*polygon = np
}

// VisitCoordinates visits each point in the geometry.
func (polygon Polygon) VisitCoordinates(v Visitor) {
	for _, line := range polygon {
		for _, point := range line {
			v.Visit(point)
		}
	}
}
