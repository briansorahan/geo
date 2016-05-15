package geo

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// pointsCompare compares two slices of points.
func pointsCompare(p1, p2 [][2]float64) bool {
	if len(p1) != len(p2) {
		return false
	}
	for i, vertex := range p1 {
		if vertex[0] != p2[i][0] {
			return false
		}
		if vertex[1] != p2[i][1] {
			return false
		}
	}
	return true
}

// pointsMarshalJSON converts a list of points to JSON.
func pointsMarshalJSON(points [][2]float64, prefix, suffix string) []byte {
	s := prefix
	for i, point := range points {
		if i == 0 {
			s += "[" + strconv.FormatFloat(point[0], 'f', -1, 64) + ","
			s += strconv.FormatFloat(point[1], 'f', -1, 64) + "]"
		} else {
			s += ",[" + strconv.FormatFloat(point[0], 'f', -1, 64) + ","
			s += strconv.FormatFloat(point[1], 'f', -1, 64) + "]"
		}
	}
	return []byte(s + suffix)
}

// pointsUnmarshalJSON unmarshals some points from JSON.
func pointsUnmarshalJSON(data []byte, prefix, suffix string) ([][2]float64, error) {
	var (
		points       = [][2]float64{}
		s            = string(data)
		idx          = strings.Index(s, prefix)
		lastBrackets = strings.LastIndex(s, suffix)
	)
	if idx != 0 || lastBrackets == -1 {
		return nil, fmt.Errorf("could not unmarshal points from %q", s)
	}
	pts := strings.Split(s[len(prefix):lastBrackets], "],[")
	if len(pts) < 3 {
		return nil, errors.New("polygon must contain at least 3 points")
	}
	for i, p := range pts {
		var coords []string
		if i == 0 {
			coords = strings.Split(p[1:], ",")
		} else if i == len(pts)-1 {
			coords = strings.Split(p[:len(p)-1], ",")
		} else {
			coords = strings.Split(p, ",")
		}
		if len(coords) != 2 {
			return nil, fmt.Errorf("could not unmarshal points from %q", s)
		}
		c1, err := strconv.ParseFloat(coords[0], 64)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal points from %q", s)
		}
		c2, err := strconv.ParseFloat(coords[1], 64)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal points from %q", s)
		}
		points = append(points, [2]float64{c1, c2})
	}
	return points, nil
}

// pointsString converts a slice of points to a string.
func pointsString(points [][2]float64, prefix, suffix string) string {
	s := prefix
	s += strconv.FormatFloat(points[0][0], 'f', -1, 64)
	s += " " + strconv.FormatFloat(points[0][1], 'f', -1, 64)
	for _, coord := range points[1:] {
		s += ", " + strconv.FormatFloat(coord[0], 'f', -1, 64)
		s += " " + strconv.FormatFloat(coord[1], 'f', -1, 64)
	}
	return s + suffix
}
