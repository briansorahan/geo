package geo

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Polygon is the GeoJSON Polygon geometry.
type Polygon [][2]float64

// MarshalJSON returns the GeoJSON representation of the polygon.
func (polygon *Polygon) MarshalJSON() ([]byte, error) {
	s := []byte(`{"type":"Polygon","coordinates":`)
	coords, err := json.Marshal([][2]float64(*polygon))
	if err != nil {
		return nil, err
	}
	return bytes.Join([][]byte{s, coords, []byte("}")}, []byte{}), nil
}

// Scan scans a polygon from Well Known Text.
func (polygon *Polygon) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		return polygon.scan(string(v))
	case string:
		return polygon.scan(v)
	default:
		return ErrScan
	}
}

// scan scans a polygon from a Well Known Text string.
func (polygon *Polygon) scan(s string) error {
	if i := strings.Index(s, "POLYGON("); i != 0 {
		return fmt.Errorf("malformed polygon %s", s)
	}
	l := len(s)
	if s[l-1] != ')' {
		return fmt.Errorf("malformed polygon %s", s)
	}
	s = s[8 : l-1]
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
func (polygon *Polygon) Value() (driver.Value, error) {
	return polygon.String(), nil
}

// String converts the polygon to a string.
func (polygon *Polygon) String() string {
	if len(*polygon) == 0 {
		return "POLYGON EMPTY"
	}
	s := "POLYGON("
	s += strconv.FormatFloat((*polygon)[0][0], 'f', -1, 64)
	s += " " + strconv.FormatFloat((*polygon)[0][1], 'f', -1, 64)
	for _, coord := range (*polygon)[1:] {
		s += ", " + strconv.FormatFloat(coord[0], 'f', -1, 64)
		s += " " + strconv.FormatFloat(coord[1], 'f', -1, 64)
	}
	return s + ")"
}
