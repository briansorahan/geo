package geo

import (
	"encoding/json"
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

// pointsContain returns true if the
func pointsContain(pts [][2]float64, pt [2]float64) bool {
	if len(pts) < 2 {
		return false
	}
	for i, vertex := range pts {
		if vertex[0] == pt[0] && vertex[1] == pt[1] {
			return true
		}
		if i == 0 {
			continue
		}
		if segmentContains(vertex, pts[i-1], pt) {
			return true
		}
	}
	return false
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
	if s[0] == 40 {
		s = s[1 : len(s)-1]
	}
	if len(s)-2 >= 0 && s[len(s)-1] == 41 {
		s = s[:len(s)-2]
	}

	points := [][2]float64{}
	for _, coords := range strings.Split(s, ",") {
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

// pointsScanPrefix scans a string form points, and expects the given prefix.
func pointsScanPrefix(s, prefix, typeName string) ([][2]float64, error) {
	if len(s) <= len(prefix) {
		return nil, fmt.Errorf("could not scan %s from %s", typeName, s)
	}
	return pointsScan(s[len(prefix):])
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

// pointsUnmarshal unmarshals a slice of points
func pointsUnmarshal(data []byte, expectedType string) ([][2]float64, error) {
	g := geometry{}

	// Never fails because data is always valid JSON.
	_ = json.Unmarshal(data, &g)

	if expected, got := expectedType, g.Type; expected != got {
		return nil, fmt.Errorf("expected type %s, got %s", expected, got)
	}

	pts := [][2]float64{}
	if err := json.Unmarshal(g.Coordinates, &pts); err != nil {
		return nil, err
	}
	return pts, nil
}
