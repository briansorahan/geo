package geo

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
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
type MultiPolygon [][][][3]float64

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
	for _, poly := range multiPolygon {
		if !Polygon(poly).Contains(point) {
			return false
		}
	}
	return true
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
		return errors.Errorf("malformed multi polygon %s", s)
	}
	l := len(s)

	if s[l-len(multiPolygonWKTSuffix):] != multiPolygonWKTSuffix {
		return errors.Errorf("malformed multi polygon %s", s)
	}
	s = s[len(multiPolygonWKTPrefix) : l-len(multiPolygonWKTSuffix)]

	// empty the polygon
	*multiPolygon = MultiPolygon{}

	// Split the string into polygons.
	// We have to split on double parens because single parens
	// would split the polygons themselves.
	// The first polygon will have a leading double parens, i.e. "(("
	// and the last polygon will have a trailing double parens, i.e. "))"
	polygons := strings.Split(s, ")),((")

	// Get the coordinates.
	for _, polys := range polygons {
		var (
			poly   = [][][3]float64{}
			polyss = strings.Split(polys, "),(")
		)
		for _, ss := range polyss {
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
	p := [][][][3]float64{}
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
