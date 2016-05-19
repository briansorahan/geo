package geo

import "fmt"

const (
	lineWKTEmpty   = `LINESTRING EMPTY`
	lineWKTPrefix  = `LINESTRING`
	lineJSONPrefix = `{"type":"LineString","coordinates":[`
	lineJSONSuffix = `]}`
)

// Line is a line.
type Line [][2]float64

// Compare compares one linestring to another.
func (line Line) Compare(other Geometry) bool {
	ls, ok := other.(*Line)
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
	switch v := src.(type) {
	case []byte:
		return line.scan(string(v))
	case string:
		return line.scan(v)
	default:
		return fmt.Errorf("could not scan line from %T", src)
	}
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
