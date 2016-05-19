package geo

import (
	"database/sql/driver"
	"fmt"
)

const (
	lineWKTEmpty   = `LINESTRING EMPTY`
	lineWKTPrefix  = `LINESTRING`
	lineJSONPrefix = `{"type":"LineString","coordinates":[`
	lineJSONSuffix = `]}`
)

// Line is a line.
type Line [][2]float64

// Compare compares one line to another.
func (line Line) Compare(g Geometry) bool {
	ls, ok := g.(*Line)
	if !ok {
		return false
	}
	return pointsCompare(line, *ls)
}

// MarshalJSON marshals the line to JSON.
func (line Line) MarshalJSON() ([]byte, error) {
	return pointsMarshalJSON(line, lineJSONPrefix, lineJSONSuffix), nil
}

// Scan scans a line from Well Known Text.
func (line *Line) Scan(src interface{}) error {
	return scan(line, src)
}

// scan scans a line from a Well Known Text string.
func (line *Line) scan(s string) error {
	if len(s) <= len(lineWKTPrefix) {
		return fmt.Errorf("could not scan line from %s", s)
	}
	points, err := pointsScan(s[len(lineWKTPrefix):])
	if err != nil {
		return err
	}
	*line = points
	return nil
}

// String converts the line to a string.
func (line Line) String() string {
	if len(line) == 0 {
		return lineWKTEmpty
	}
	return lineWKTPrefix + pointsString(line)
}

// Value returns a driver Value.
func (line Line) Value() (driver.Value, error) {
	return line.String(), nil
}
