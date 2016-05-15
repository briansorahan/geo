package geo

const (
	linestringWKTEmpty   = `LINESTRING EMPTY`
	linestringWKTPrefix  = `LINESTRING(`
	linestringWKTSuffix  = `)`
	linestringJSONPrefix = `{"type":"LineString","coordinates":[`
	linestringJSONSuffix = `]}`
)

// Line is a line.
type Line [][2]float64

// Compare compares one linestring to another.
func (linestring Line) Compare(other Geometry) bool {
	ls, ok := other.(*Line)
	if !ok {
		return false
	}
	return pointsCompare(linestring, *ls)
}

// MarshalJSON marshals the linestring to JSON.
func (linestring Line) MarshalJSON() ([]byte, error) {
	return pointsMarshalJSON(linestring, linestringJSONPrefix, linestringJSONSuffix), nil
}

// UnmarshalJSON unmarshals the linestring from JSON.
func (linestring *Line) UnmarshalJSON(data []byte) error {
	points, err := pointsUnmarshalJSON(data, linestringJSONPrefix, linestringJSONSuffix)
	if err != nil {
		return err
	}
	*linestring = points
	return nil
}

// String converts the linestring to a string.
func (linestring Line) String() string {
	if len(linestring) == 0 {
		return linestringWKTEmpty
	}
	return pointsString(linestring, linestringWKTPrefix, linestringWKTSuffix)
}
