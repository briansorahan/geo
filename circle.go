package geo

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

const (
	circleWKTPrefix = `CIRCULARSTRING`
)

// Circle is a circle in the XY plane.
type Circle struct {
	Center Point
	Radius float64
}

// Compare compares the circle to another geometry.
func (circle Circle) Compare(g Geometry) bool {
	c, ok := g.(*Circle)
	if !ok {
		return false
	}
	if !c.Center.Compare(&circle.Center) {
		return false
	}
	return c.Radius == circle.Radius
}

// MarshalJSON marshals a circle to GeoJSON.
// See https://github.com/geojson/geojson-spec/wiki/Proposal---Circles-and-Ellipses-Geoms
func (circle Circle) MarshalJSON() ([]byte, error) {
	return []byte(`{"type":"Circle","radius":` +
		strconv.FormatFloat(circle.Radius, 'f', -1, 64) +
		`,"coordinates":[` +
		strconv.FormatFloat(circle.Center[0], 'f', -1, 64) + `,` +
		strconv.FormatFloat(circle.Center[1], 'f', -1, 64) + `]}`), nil
}

// Scan scans a circle from well known text.
func (circle *Circle) Scan(src interface{}) error {
	return scan(circle, src)
}

// Scan scans a circle from well known text.
func (circle *Circle) scan(s string) error {
	idx := strings.Index(s, circleWKTPrefix)
	if idx != 0 {
		return fmt.Errorf("malformed circle: %s", s)
	}
	points, err := pointsScan(s[len(circleWKTPrefix):])
	if err != nil {
		return err
	}
	if len(points) < 3 {
		return fmt.Errorf("malformed circle: %s", s)
	}
	// points 0 and 2 should be on opposite sides of the circle,
	// so we can calculate the radius as 1/2 the distance between them
	// and the center as the midpoint.
	circle.Radius = Point(points[0]).DistanceFrom(points[2]) / 2
	var (
		dx = points[2][0] - points[0][0]
		dy = points[2][1] - points[0][1]
	)
	circle.Center = Point{points[0][0] + (dx / 2), points[0][1] + (dy / 2)}
	return nil
}

// String returns a string representation of the circle.
func (circle Circle) String() string {
	return "CIRCULARSTRING" + pointsString([][2]float64{
		{circle.Center[0] + circle.Radius, circle.Center[1]},
		{circle.Center[0], circle.Center[1] + circle.Radius},
		{circle.Center[0] - circle.Radius, circle.Center[1]},
		{circle.Center[0], circle.Center[1] - circle.Radius},
		{circle.Center[0] + circle.Radius, circle.Center[1]},
	})
}

// Value
func (circle Circle) Value() (driver.Value, error) {
	return circle.String(), nil
}
