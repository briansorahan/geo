package geo

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	polygonPrefix      = `POLYGON((`
	polygonJSONPrefix1 = `{"type":"Polygon","coordinates":[[`
	// BUG(briansorahan): unmarshal json regardless of field order
)

// Polygon is the GeoJSON Polygon geometry.
type Polygon [][2]float64

// Compare compares one polygon to another.
func (polygon Polygon) Compare(other Polygon) bool {
	if len(polygon) != len(other) {
		return false
	}
	for i, vertex := range polygon {
		if vertex[0] != other[i][0] {
			return false
		}
		if vertex[1] != other[i][1] {
			return false
		}
	}
	return true
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
	s := polygonJSONPrefix1
	for i, point := range polygon {
		if i == 0 {
			s += "[" + strconv.FormatFloat(point[0], 'f', -1, 64) + ","
			s += strconv.FormatFloat(point[1], 'f', -1, 64) + "]"
		} else {
			s += ",[" + strconv.FormatFloat(point[0], 'f', -1, 64) + ","
			s += strconv.FormatFloat(point[1], 'f', -1, 64) + "]"
		}
	}
	return []byte(s + "]]}"), nil
}

// UnmarshalJSON unmarshals a polygon from GeoJSON.
func (polygon *Polygon) UnmarshalJSON(data []byte) error {
	var (
		s            = string(data)
		idx          = strings.Index(s, polygonJSONPrefix1)
		lastBrackets = strings.LastIndex(s, "]]}")
	)
	if idx != 0 || lastBrackets == -1 {
		return fmt.Errorf("could not unmarshal polygon from %q", s)
	}
	points := strings.Split(s[len(polygonJSONPrefix1):lastBrackets], "],[")
	if len(points) < 3 {
		return errors.New("polygon must contain at least 3 points")
	}
	for i, p := range points {
		var coords []string
		if i == 0 {
			coords = strings.Split(p[1:], ",")
		} else if i == len(points)-1 {
			coords = strings.Split(p[:len(p)-1], ",")
		} else {
			coords = strings.Split(p, ",")
		}
		if len(coords) != 2 {
			return fmt.Errorf("could not unmarshal polygon from %q", s)
		}
		c1, err := strconv.ParseFloat(coords[0], 64)
		if err != nil {
			return fmt.Errorf("could not unmarshal polygon from %q", s)
		}
		c2, err := strconv.ParseFloat(coords[1], 64)
		if err != nil {
			return fmt.Errorf("could not unmarshal polygon from %q", s)
		}
		*polygon = append(*polygon, [2]float64{c1, c2})
	}
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
	if i := strings.Index(s, polygonPrefix); i != 0 {
		return fmt.Errorf("malformed polygon %s", s)
	}
	l := len(s)
	if s[l-1] != ')' {
		return fmt.Errorf("malformed polygon %s", s)
	}
	s = s[len(polygonPrefix) : l-1]
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
	s := polygonPrefix
	s += strconv.FormatFloat(polygon[0][0], 'f', -1, 64)
	s += " " + strconv.FormatFloat(polygon[0][1], 'f', -1, 64)
	for _, coord := range polygon[1:] {
		s += ", " + strconv.FormatFloat(coord[0], 'f', -1, 64)
		s += " " + strconv.FormatFloat(coord[1], 'f', -1, 64)
	}
	return s + "))"
}
