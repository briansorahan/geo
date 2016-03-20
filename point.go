package geo

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
)

const pointWKT = `POINT(%f %f)`

// Point defines a point.
type Point [2]float64

// MarshalJSON returns the GeoJSON representation of the point.
func (point *Point) MarshalJSON() ([]byte, error) {
	s := []byte(`{"type":"Point","coordinates":`)
	coords, err := json.Marshal([2]float64(*point))
	if err != nil {
		return nil, err
	}
	return bytes.Join([][]byte{s, coords, []byte("}")}, []byte{}), nil
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
		return ErrScan
	}
	return nil
}

// Value converts a point to Well Known Text.
func (point *Point) Value() (driver.Value, error) {
	return point.String(), nil
}

// String convert the point to a string.
func (point *Point) String() string {
	s := "POINT("
	s += strconv.FormatFloat(point[0], 'f', -1, 64)
	s += " " + strconv.FormatFloat(point[1], 'f', -1, 64) + ")"
	return s
}
