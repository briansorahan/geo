package geo

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

const (
	mpWKTEmpty   = `MULTIPOINT EMPTY`
	mpWKTPrefix  = `MULTIPOINT`
	mpJSONPrefix = `{"type":"MultiPoint","coordinates":[`
	mpJSONSuffix = `]}`
)

// MultiPoint is a collection of points.
type MultiPoint [][2]float64

// Compare compares one MultiPoint to another.
func (mp MultiPoint) Compare(g Geometry) bool {
	other, ok := g.(*MultiPoint)
	if !ok {
		return false
	}
	return pointsCompare(mp, *other)
}

// Contains determines if the MultiPoint contains a point.
func (mp MultiPoint) Contains(p Point) bool {
	if len(mp) < 2 {
		return false
	}
	for i, vertex := range mp {
		if vertex[0] == p[0] && vertex[1] == p[1] {
			return true
		}
		if i == 0 {
			continue
		}
		if segmentContains(vertex, mp[i-1], p) {
			return true
		}
	}
	return false
}

// MarshalJSON marshals the MultiPoint to JSON.
func (mp MultiPoint) MarshalJSON() ([]byte, error) {
	return pointsMarshalJSON(mp, mpJSONPrefix, mpJSONSuffix), nil
}

// Scan scans a MultiPoint from Well Known Text.
func (mp *MultiPoint) Scan(src interface{}) error {
	return scan(mp, src)
}

// scan scans a MultiPoint from a Well Known Text string.
func (mp *MultiPoint) scan(s string) error {
	if len(s) <= len(mpWKTPrefix) {
		return fmt.Errorf("could not scan MultiPoint from %s", s)
	}
	points, err := pointsScan(s[len(mpWKTPrefix):])
	if err != nil {
		return err
	}
	*mp = points
	return nil
}

// String converts the MultiPoint to a string.
func (mp MultiPoint) String() string {
	if len(mp) == 0 {
		return mpWKTEmpty
	}
	return mpWKTPrefix + pointsString(mp)
}

// UnmarshalJSON unmarshals a MultiPoint from GeoJSON.
func (mp *MultiPoint) UnmarshalJSON(data []byte) error {
	g := geometry{}

	// Never fails because data is always valid JSON.
	_ = json.Unmarshal(data, &g)

	if expected, got := MultiPointType, g.Type; expected != got {
		return fmt.Errorf("expected type %s, got %s", expected, got)
	}

	ln := [][2]float64{}
	if err := json.Unmarshal(g.Coordinates, &ln); err != nil {
		return err
	}
	*mp = ln
	return nil
}

// Value returns a driver Value.
func (mp MultiPoint) Value() (driver.Value, error) {
	return mp.String(), nil
}
