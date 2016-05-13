package geo

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	linestringJSONPrefix1 = `{"type":"LineString","coordinates":[`
)

// Linestring is a line.
type Linestring [][2]float64

// MarshalJSON marshals the linestring to JSON.
func (linestring Linestring) MarshalJSON() ([]byte, error) {
	s := linestringJSONPrefix1
	for i, point := range linestring {
		if i == 0 {
			s += "[" + strconv.FormatFloat(point[0], 'f', -1, 64) + ","
			s += strconv.FormatFloat(point[1], 'f', -1, 64) + "]"
		} else {
			s += ",[" + strconv.FormatFloat(point[0], 'f', -1, 64) + ","
			s += strconv.FormatFloat(point[1], 'f', -1, 64) + "]"
		}
	}
	return []byte(s + "]}"), nil
}

// UnmarshalJSON unmarshals the linestring from JSON.
func (linestring *Linestring) UnmarshalJSON(data []byte) error {
	var (
		s            = string(data)
		idx          = strings.Index(s, linestringJSONPrefix1)
		lastBrackets = strings.LastIndex(s, "]}")
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
		*linestring = append(*linestring, [2]float64{c1, c2})
	}
	return nil
}
