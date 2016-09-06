package geo

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
)

const (
	earthRadiusMeters = 6371e3 // earth's radius in meters
	feetToMeters      = 0.3048 // convert feet to meters
	circleWKTPrefix   = `CIRCULARSTRING`
)

// Circle is a circle in the XY plane.
type Circle struct {
	Coordinates Point   `json:"coordinates"`
	Radius      float64 `json:"radius"`
}

// Compare compares the circle to another geometry.
func (c Circle) Compare(g Geometry) bool {
	c2, ok := g.(*Circle)
	if !ok {
		return false
	}
	if !c.Coordinates.Compare(&c2.Coordinates) {
		return false
	}
	return c.Radius == c2.Radius
}

// Contains determines if the circle contains the point.
// This assumes radius is specified in feet.
func (c Circle) Contains(p Point) bool {
	return c.ContainsHaversine(p)
}

// ContainsHaversine uses the haversine formula to determine if the
// point is contained in the circle.
func (c Circle) ContainsHaversine(p Point) bool {
	var (
		lat1 = toRadians(c.Coordinates[1])
		lat2 = toRadians(p[1])
		dLat = toRadians(p[1] - c.Coordinates[1])
		dLng = toRadians(p[0] - c.Coordinates[0])
		a    = (math.Sin(dLat/2) * math.Sin(dLat/2)) +
			(math.Cos(lat1) * math.Cos(lat2) * math.Sin(dLng/2) * math.Sin(dLng/2))
		d = earthRadiusMeters * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	)
	return d < (feetToMeters * c.Radius)
}

// ContainsSLC uses the spherical law of cosines to determine if
// the point is contained in the circle.
func (c Circle) ContainsSLC(p Point) bool {
	var (
		lat1 = toRadians(c.Coordinates[1])
		lat2 = toRadians(p[1])
		dLng = toRadians(p[0] - c.Coordinates[0])
		a    = (math.Sin(lat1) * math.Sin(lat2)) +
			(math.Cos(lat1) * math.Cos(lat2) * math.Cos(dLng))
		d = earthRadiusMeters * math.Acos(a)
	)
	return d < (feetToMeters * c.Radius)
}

// ContainsEquirectangular uses equirectangular projection to
// determine if the point is contained in the circle.
func (c Circle) ContainsEquirectangular(p Point) bool {
	var (
		dLng = toRadians(p[0] - c.Coordinates[0])
		mLat = toRadians(p[1]+c.Coordinates[1]) / float64(2)
		y    = toRadians(p[1] - c.Coordinates[1])
		x    = dLng * math.Cos(mLat)
		d    = earthRadiusMeters * math.Sqrt((x*x)+(y*y))
	)
	return d < (feetToMeters * c.Radius)
}

// toRadians converts from degrees to radians.
func toRadians(degrees float64) float64 {
	return (math.Pi * degrees) / 180
}

// MarshalJSON marshals a circle to GeoJSON.
// See https://github.com/geojson/geojson-spec/wiki/Proposal---Circles-and-Ellipses-Geoms
func (c Circle) MarshalJSON() ([]byte, error) {
	return []byte(`{"type":"Circle","radius":` +
		strconv.FormatFloat(c.Radius, 'f', -1, 64) +
		`,"coordinates":[` +
		strconv.FormatFloat(c.Coordinates[0], 'f', -1, 64) + `,` +
		strconv.FormatFloat(c.Coordinates[1], 'f', -1, 64) + `]}`), nil
}

// Scan scans a circle from well known text.
func (c *Circle) Scan(src interface{}) error {
	return scan(c, src)
}

// Scan scans a circle from well known text.
func (c *Circle) scan(s string) error {
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
	c.Radius = Point(points[0]).DistanceFrom(points[2]) / 2
	var (
		dx = points[2][0] - points[0][0]
		dy = points[2][1] - points[0][1]
	)
	c.Coordinates = Point{points[0][0] + (dx / 2), points[0][1] + (dy / 2)}
	return nil
}

// String returns a string representation of the circle.
func (c Circle) String() string {
	return "CIRCULARSTRING" + pointsString([][2]float64{
		{c.Coordinates[0] + c.Radius, c.Coordinates[1]},
		{c.Coordinates[0], c.Coordinates[1] + c.Radius},
		{c.Coordinates[0] - c.Radius, c.Coordinates[1]},
		{c.Coordinates[0], c.Coordinates[1] - c.Radius},
		{c.Coordinates[0] + c.Radius, c.Coordinates[1]},
	})
}

// UnmarshalJSON unmarshals the circle from GeoJSON.
func (c *Circle) UnmarshalJSON(data []byte) error {
	g := &geometry{}

	if err := json.Unmarshal(data, g); err != nil {
		return err
	}

	// TODO: check type

	if err := json.Unmarshal(g.Coordinates, &c.Coordinates); err != nil {
		return err
	}

	c.Radius = g.Radius

	return nil
}

// Value returns a sql driver value.
func (c Circle) Value() (driver.Value, error) {
	return c.String(), nil
}
