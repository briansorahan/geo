package geo

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

const (
	polygonWKTPrefix  = `POLYGON((`
	polygonWKTSuffix  = `))`
	polygonJSONPrefix = `{"type":"Polygon","coordinates":[[`
	polygonJSONSuffix = `]]}`
	// BUG(briansorahan): unmarshal json regardless of field order
)

// Polygon is the GeoJSON Polygon geometry.
type Polygon [][2]float64

// Compare compares one polygon to another.
func (polygon Polygon) Compare(other Geometry) bool {
	poly, ok := other.(*Polygon)
	if !ok {
		return false
	}
	return pointsCompare(polygon, *poly)
}

// Contains uses the ray casting algorithm to decide
// if the point is contained in the polygon.
func (polygon Polygon) Contains(point Point) bool {
	intersections := 0
	for i, vertex := range polygon {
		if point.RayhIntersects(polygon[(i+1)%len(polygon)], vertex) {
			intersections++
		}
	}
	return (intersections % 2) == 1
}

// MarshalJSON returns the GeoJSON representation of the polygon.
func (polygon Polygon) MarshalJSON() ([]byte, error) {
	return pointsMarshalJSON(polygon, polygonJSONPrefix, polygonJSONSuffix), nil
}

// UnmarshalJSON unmarshals a polygon from GeoJSON.
func (polygon *Polygon) UnmarshalJSON(data []byte) error {
	points, err := pointsUnmarshalJSON(data, polygonJSONPrefix, polygonJSONSuffix)
	if err != nil {
		return err
	}
	*polygon = points
	return nil
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
	if s[l-2:] != polygonWKTSuffix {
		return fmt.Errorf("malformed polygon %s", s)
	}
	s = s[len(polygonWKTPrefix) : l-1]
	// empty the polygon
	*polygon = Polygon{}
	// get the coordinates
	coords := strings.Split(s, ",")
	for _, coord := range coords {
		points := [2]float64{}
		if _, err := fmt.Sscanf(strings.TrimSpace(coord), "%f %f", &points[0], &points[1]); err != nil {
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
	return pointsString(polygon, polygonWKTPrefix, polygonWKTSuffix)
}
