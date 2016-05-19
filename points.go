package geo

import (
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

// pointsScan scans a list of points from Well Known Text.
// The points should look like (X0 Y0,X1 Y1,X2 Y2)
func pointsScan(s string) ([][2]float64, error) {
	if s[0] != '(' || s[len(s)-1] != ')' {
		return nil, fmt.Errorf("could not scan points from %s", s)
	}
	points := [][2]float64{}
	for _, coords := range strings.Split(s[1:len(s)-1], ",") {
		var (
			pair = [2]float64{}
			xy   = strings.Split(strings.TrimSpace(coords), " ")
		)
		if len(xy) != 2 {
			return nil, fmt.Errorf("could not scan points from %s", s)
		}
		for i, val := range xy {
			f, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return nil, err
			}
			pair[i] = f
		}
		points = append(points, pair)
	}
	return points, nil
}

// pointsString converts a slice of points to Well Known Text.
func pointsString(points [][2]float64) string {
	s := "("
	s += strconv.FormatFloat(points[0][0], 'f', -1, 64)
	s += " " + strconv.FormatFloat(points[0][1], 'f', -1, 64)
	for _, coord := range points[1:] {
		s += ", " + strconv.FormatFloat(coord[0], 'f', -1, 64)
		s += " " + strconv.FormatFloat(coord[1], 'f', -1, 64)
	}
	return s + ")"
}
