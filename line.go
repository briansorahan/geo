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
type Line [][3]float64

// Equal compares one line to another.
func (line Line) Equal(g Geometry) bool {
	ls, ok := g.(*Line)
	if !ok {
		return false
	}
	return pointsEqual(line, *ls)
}

// Contains determines if the line contains a point.
func (line Line) Contains(p Point) bool {
	return pointsContain(line, p)
}

// segmentContains returns true if p lies on the line segment that connects s1 and s2.
func segmentContains(s1, s2, p [3]float64) bool {
	// Return false if p is outside of the bounding box around s1 and s2.
	if (p[0] > s1[0] && p[0] > s2[0]) || (p[0] < s1[0] && p[0] < s2[0]) {
		return false
	}
	if (p[1] > s1[1] && p[1] > s2[1]) || (p[1] < s1[1] && p[1] < s2[1]) {
		return false
	}
	// Equal the slope of the segment between s1 and p
	// to the slope of the segment between s1 and s2.
	var (
		segmentSlope = (s2[1] - s1[1]) / (s2[0] - s1[0])
		pSlope       = (p[1] - s1[1]) / (p[0] - s1[0])
	)
	return segmentSlope == pSlope
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

// UnmarshalJSON unmarshals a line from GeoJSON.
func (line *Line) UnmarshalJSON(data []byte) error {
	pts, err := pointsUnmarshal(data, LineType)
	if err != nil {
		return err
	}
	*line = pts
	return nil
}

// Value returns a driver Value.
func (line Line) Value() (driver.Value, error) {
	return line.String(), nil
}

// Transform transforms the geometry point by point.
func (line *Line) Transform(t Transformer) {
	nl := make([][3]float64, len(*line))
	for i, point := range *line {
		nl[i] = [3]float64(t.Transform(point))
	}
	*line = nl
}

// Visit visits each point in the geometry.
func (line Line) Visit(v Visitor) {
	for _, point := range line {
		v.Visit(point)
	}
}
