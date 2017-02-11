package geo

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	multiPolygonWKTPrefix = `MULTIPOLYGON(`
	multiPolygonWKTSuffix = `)`
)

var (
	multiPolygonJSONPrefix = []byte(`{"type":"MultiPolygon","coordinates":[`)
	multiPolygonJSONSuffix = []byte(`]}`)
)

// MultiPolygon is a GeoJSON MultiPolygon.
type MultiPolygon [][][][2]float64

// Compare compares one polygon to another.
func (multiPolygon MultiPolygon) Compare(g Geometry) bool {
	p, ok := g.(*MultiPolygon)
	if !ok {
		return false
	}
	if len(multiPolygon) != len(*p) {
		return false
	}
	for i, p1 := range multiPolygon {
		for j, p2 := range p1 {
			p3 := (*p)[i][j]
			if len(p2) != len(p3) {
				return false
			}
			if !pointsCompare(p2, p3) {
				return false
			}
		}
	}
	return true
}

// Contains uses the ray casting algorithm to decide
// if the point is contained in the polygon.
func (multiPolygon MultiPolygon) Contains(point Point) bool {
	intersections := 0
	for _, poly := range multiPolygon {
		for _, edge := range poly {
			for j, vertex := range edge {
				if point.RayhIntersects(edge[(j+1)%len(poly)], vertex) {
					intersections++
				}
			}
		}
	}
	return (intersections % 2) == 1
}

// MarshalJSON returns the GeoJSON representation of the polygon.
func (multiPolygon MultiPolygon) MarshalJSON() ([]byte, error) {
	s := multiPolygonJSONPrefix
	for i, poly := range multiPolygon {
		if i == 0 {
			s = append(s, '[')
		} else {
			s = append(s, ',', '[')
		}
		for j, line := range poly {
			if j == 0 {
				s = append(s, '[')
			} else {
				s = append(s, ',', '[')
			}
			s = append(s, pointsMarshalJSON(line, "", "")...)
			s = append(s, ']')
		}
		s = append(s, ']')
	}
	return append(s, multiPolygonJSONSuffix...), nil
}

// Scan scans a polygon from Well Known Text.
func (multiPolygon *MultiPolygon) Scan(src interface{}) error {
	return scan(multiPolygon, src)
}

// scan scans a polygon from a Well Known Text string.
func (multiPolygon *MultiPolygon) scan(s string) error {
	if i := strings.Index(s, multiPolygonWKTPrefix); i != 0 {
		return fmt.Errorf("malformed polygon %s", s)
	}
	l := len(s)

	if s[l-len(multiPolygonWKTSuffix):] != multiPolygonWKTSuffix {
		return fmt.Errorf("malformed polygon %s", s)
	}
	s = s[len(multiPolygonWKTPrefix) : l-len(multiPolygonWKTSuffix)]

	// empty the polygon
	*multiPolygon = MultiPolygon{}

	// get the coordinates
	multiPolygons := strings.Split(s, "),(")
	for _, polyss := range multiPolygons {
		var (
			polys = strings.Split(polyss, "),(")
			poly  = [][][2]float64{}
		)
		for _, ss := range polys {
			points, err := pointsScan(ss)
			if err != nil {
				return err
			}
			poly = append(poly, points)
		}
		*multiPolygon = append(*multiPolygon, poly)
	}
	return nil
}

// String converts the polygon to a string.
func (multiPolygon MultiPolygon) String() string {
	if len(multiPolygon) == 0 {
		return "MULTIPOLYGON EMPTY"
	}
	s := multiPolygonWKTPrefix
	for i, polys := range multiPolygon {
		if i == 0 {
			s += "("
		} else {
			s += ",("
		}
		for j, points := range polys {
			if j == 0 {
				s += pointsString(points)
			} else {
				s += "," + pointsString(points)
			}
		}
		s += ")"
	}
	return s + multiPolygonWKTSuffix
}

// UnmarshalJSON unmarshals the polygon from GeoJSON.
func (multiPolygon *MultiPolygon) UnmarshalJSON(data []byte) error {
	g := &geometry{}

	// Never fails since data is always valid JSON.
	_ = json.Unmarshal(data, g)

	if expected, got := MultiPolygonType, g.Type; expected != got {
		return fmt.Errorf("expected type %s, got %s", expected, got)
	}

	p := [][][][2]float64{}

	if err := json.Unmarshal(g.Coordinates, &p); err != nil {
		return err
	}

	*multiPolygon = MultiPolygon(p)

	return nil
}

// Value converts a point to Well Known Text.
func (multiPolygon MultiPolygon) Value() (driver.Value, error) {
	return multiPolygon.String(), nil
}
