package geo

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

const (
	pointWKTPrefix  = `POINT`
	pointWKT        = `(%f %f)`
	pointJSONPrefix = `{"type":"Point","coordinates":[`
)

// Point defines a point.
type Point [2]float64

// Compare compares one point to another.
func (point Point) Compare(g Geometry) bool {
	pt, ok := g.(*Point)
	if !ok {
		return false
	}
	if point[0] != (*pt)[0] {
		return false
	}
	if point[1] != (*pt)[1] {
		return false
	}
	return true
}

// Contains is the exact same as Compare.
func (point Point) Contains(other Point) bool {
	return point.Compare(&other)
}

// DistanceFrom computes the distance from one point to another.
func (point Point) DistanceFrom(other Point) float64 {
	return math.Sqrt(math.Pow(point[1]-other[1], 2) + math.Pow(point[0]-other[0], 2))
}

// MarshalJSON returns the GeoJSON representation of the point.
func (point Point) MarshalJSON() ([]byte, error) {
	s := pointJSONPrefix
	s += strconv.FormatFloat(point[0], 'f', -1, 64) + ","
	s += strconv.FormatFloat(point[1], 'f', -1, 64) + "]}"
	return []byte(s), nil
}

// RayhIntersects returns true if the horizontal ray going from
// point to positive infinity intersects the line that connects a and b.
func (point Point) RayhIntersects(a, b Point) bool {
	var (
		left   = math.Min(a[0], b[0])
		right  = math.Max(a[0], b[0])
		bottom = math.Min(a[1], b[1])
		top    = math.Max(a[1], b[1])
	)
	if point[0] > right {
		return false
	}
	if point[1] >= top || point[1] < bottom {
		return false
	}
	if point[0] <= left {
		return true
	}
	slope := (b[1] - a[1]) / (b[0] - a[0])
	if slope >= 0 {
		return ((point[1] - bottom) / (point[0] - left)) >= slope
	}
	return ((point[1] - top) / (point[0] - left)) <= slope
}

// Scan scans a point from Well Known Text.
func (point *Point) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		if _, err := fmt.Sscanf(string(v), pointWKTPrefix+pointWKT, &point[0], &point[1]); err != nil {
			return err
		}
	case string:
		if _, err := fmt.Sscanf(v, pointWKTPrefix+pointWKT, &point[0], &point[1]); err != nil {
			return err
		}
	default:
		return fmt.Errorf("could not scan point from %T", src)
	}
	return nil
}

// String convert the point to a string.
func (point Point) String() string {
	s := "POINT("
	s += strconv.FormatFloat(point[0], 'f', -1, 64)
	s += " " + strconv.FormatFloat(point[1], 'f', -1, 64) + ")"
	return s
}

// UnmarshalJSON unmarshals a point from GeoJSON.
func (point *Point) UnmarshalJSON(data []byte) error {
	g := geometry{}

	// Never fails because data is always valid JSON.
	_ = json.Unmarshal(data, &g)

	if expected, got := PointType, g.Type; expected != got {
		return fmt.Errorf("expected %s type, got %s", expected, got)
	}

	pt := [2]float64{}
	if err := json.Unmarshal(g.Coordinates, &pt); err != nil {
		return err
	}
	*point = pt
	return nil
}

// Value converts a point to Well Known Text.
func (point Point) Value() (driver.Value, error) {
	return point.String(), nil
}
