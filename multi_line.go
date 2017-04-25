package geo

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	mlWKTPrefix = `MULTILINESTRING(`
	mlWKTSuffix = `)`
)

var (
	mlJSONPrefix = []byte(`{"type":"MultiLineString","coordinates":[`)
	mlJSONSuffix = []byte(`]}`)
)

// MultiLine is an array of Line's.
type MultiLine [][][3]float64

// Equal compares one MultiLine to another.
func (ml MultiLine) Equal(g Geometry) bool {
	p, ok := g.(*MultiLine)
	if !ok {
		return false
	}
	if len(ml) != len(*p) {
		return false
	}
	for i, p1 := range ml {
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
// if the point is contained in the MultiLine.
func (ml MultiLine) Contains(point Point) bool {
	intersections := 0
	for _, poly := range ml {
		for j, vertex := range poly {
			if point.RayhIntersects(poly[(j+1)%len(poly)], vertex) {
				intersections++
			}
		}
	}
	return (intersections % 2) == 1
}

// MarshalJSON returns the GeoJSON representation of the MultiLine.
func (ml MultiLine) MarshalJSON() ([]byte, error) {
	s := mlJSONPrefix
	for i, poly := range ml {
		if i == 0 {
			s = append(s, '[')
		} else {
			s = append(s, ',', '[')
		}
		s = append(s, pointsMarshalJSON(poly, "", "")...)
		s = append(s, ']')
	}
	return append(s, mlJSONSuffix...), nil
}

// Scan scans a MultiLine from Well Known Text.
func (ml *MultiLine) Scan(src interface{}) error {
	return scan(ml, src)
}

// scan scans a MultiLine from a Well Known Text string.
func (ml *MultiLine) scan(s string) error {
	if i := strings.Index(s, mlWKTPrefix); i != 0 {
		return fmt.Errorf("malformed MultiLine %s", s)
	}
	l := len(s)

	if s[l-len(mlWKTSuffix):] != mlWKTSuffix {
		return fmt.Errorf("malformed MultiLine %s", s)
	}
	s = s[len(mlWKTPrefix) : l-len(mlWKTSuffix)]

	// empty the MultiLine
	*ml = MultiLine{}

	// get the coordinates
	mls := strings.Split(s, "),(")
	for _, ss := range mls {
		points, err := pointsScan(ss)
		if err != nil {
			return err
		}
		*ml = append(*ml, points)
	}
	return nil
}

// String converts the MultiLine to a string.
func (ml MultiLine) String() string {
	if len(ml) == 0 {
		return "MULTILINESTRING EMPTY"
	}
	s := mlWKTPrefix
	for i, points := range ml {
		if i == 0 {
			s += pointsString(points)
		} else {
			s += "," + pointsString(points)
		}
	}
	return s + mlWKTSuffix
}

// UnmarshalJSON unmarshals the ml from GeoJSON.
func (ml *MultiLine) UnmarshalJSON(data []byte) error {
	g := &geometry{}

	// Never fails since data is always valid JSON.
	_ = json.Unmarshal(data, g)

	if expected, got := MultiLineType, g.Type; expected != got {
		return fmt.Errorf("expected type %s, got %s", expected, got)
	}

	p := [][][3]float64{}

	if err := json.Unmarshal(g.Coordinates, &p); err != nil {
		return err
	}

	*ml = MultiLine(p)

	return nil
}

// Value converts a point to Well Known Text.
func (ml MultiLine) Value() (driver.Value, error) {
	return ml.String(), nil
}

// Transform transforms the geometry point by point.
func (ml *MultiLine) Transform(t Transformer) {
	np := make([][][3]float64, len(*ml))
	for i, line := range *ml {
		nl := make([][3]float64, len(line))
		for j, point := range line {
			nl[j] = t.Transform(point)
		}
		np[i] = nl
	}
	*ml = np
}

// Visit visits each point in the geometry.
func (ml MultiLine) Visit(v Visitor) {
	for _, line := range ml {
		for _, point := range line {
			v.Visit(point)
		}
	}
}
