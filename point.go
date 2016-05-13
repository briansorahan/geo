package geo

import (
	"database/sql/driver"
	"fmt"
	"math"
	"strconv"
	"strings"
)

const (
	pointWKT        = `POINT(%f %f)`
	pointJSONPrefix = `{"type":"Point","coordinates":[`
)

// Point defines a point.
type Point [2]float64

// Compare compares one point to another.
func (point Point) Compare(other Point) bool {
	if point[0] != other[0] {
		return false
	}
	if point[1] != other[1] {
		return false
	}
	return true
}

// MarshalJSON returns the GeoJSON representation of the point.
func (point Point) MarshalJSON() ([]byte, error) {
	s := pointJSONPrefix
	s += strconv.FormatFloat(point[0], 'f', -1, 64) + ","
	s += strconv.FormatFloat(point[1], 'f', -1, 64) + "]}"
	return []byte(s), nil
}

// UnmarshalJSON unmarshals a point from GeoJSON.
func (point *Point) UnmarshalJSON(data []byte) error {
	var (
		s           = string(data)
		idx         = strings.Index(s, pointJSONPrefix)
		lastBracket = strings.LastIndex(s, "]")
	)
	if idx != 0 || lastBracket == -1 {
		return fmt.Errorf("could not unmarshal point from %q", s)
	}
	coords := strings.Split(s[len(pointJSONPrefix):lastBracket], ",")
	if len(coords) != 2 {
		return fmt.Errorf("could not unmarshal point from %q", s)
	}
	c1, err := strconv.ParseFloat(coords[0], 64)
	if err != nil {
		return fmt.Errorf("could not unmarshal point from %q", s)
	}
	c2, err := strconv.ParseFloat(coords[1], 64)
	if err != nil {
		return fmt.Errorf("could not unmarshal point from %q", s)
	}
	point[0], point[1] = c1, c2
	return nil
}

// Scan scans a point from Well Known Text.
func (point *Point) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		if _, err := fmt.Sscanf(string(v), pointWKT, &point[0], &point[1]); err != nil {
			return err
		}
	case string:
		if _, err := fmt.Sscanf(v, pointWKT, &point[0], &point[1]); err != nil {
			return err
		}
	default:
		return fmt.Errorf("could not scan point from %T", src)
	}
	return nil
}

// Value converts a point to Well Known Text.
func (point Point) Value() (driver.Value, error) {
	return point.String(), nil
}

// String convert the point to a string.
func (point Point) String() string {
	s := "POINT("
	s += strconv.FormatFloat(point[0], 'f', -1, 64)
	s += " " + strconv.FormatFloat(point[1], 'f', -1, 64) + ")"
	return s
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
