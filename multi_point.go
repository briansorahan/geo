package geo

import "database/sql/driver"

const (
	mpWKTEmpty   = `MULTIPOINT EMPTY`
	mpWKTPrefix  = `MULTIPOINT`
	mpJSONPrefix = `{"type":"MultiPoint","coordinates":[`
	mpJSONSuffix = `]}`
)

// MultiPoint is a collection of points.
type MultiPoint [][3]float64

// Equal compares one MultiPoint to another.
func (mp MultiPoint) Equal(g Geometry) bool {
	other, ok := g.(*MultiPoint)
	if !ok {
		return false
	}
	return pointsEqual(mp, *other)
}

// Contains determines if the MultiPoint contains a point.
func (mp MultiPoint) Contains(p Point) bool {
	return pointsContain(mp, p)
}

// MarshalJSON marshals the MultiPoint to JSON.
func (mp MultiPoint) MarshalJSON() ([]byte, error) {
	return pointsMarshalJSON(mp, mpJSONPrefix, mpJSONSuffix), nil
}

// Scan scans a MultiPoint from Well Known Text.
func (mp *MultiPoint) Scan(src interface{}) error {
	return scan(mp, src)
}

// scan scans a MultiPoint from a Well Known Text string.
func (mp *MultiPoint) scan(s string) error {
	points, err := pointsScanPrefix(s, mpWKTPrefix, MultiPointType)
	if err != nil {
		return err
	}
	*mp = points
	return nil
}

// String converts the MultiPoint to a string.
func (mp MultiPoint) String() string {
	if len(mp) == 0 {
		return mpWKTEmpty
	}
	return mpWKTPrefix + pointsString(mp)
}

// UnmarshalJSON unmarshals a MultiPoint from GeoJSON.
func (mp *MultiPoint) UnmarshalJSON(data []byte) error {
	pts, err := pointsUnmarshal(data, MultiPointType)
	if err != nil {
		return err
	}
	*mp = pts
	return nil
}

// Value returns a driver Value.
func (mp MultiPoint) Value() (driver.Value, error) {
	return mp.String(), nil
}
