package geo

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestBBoxEqual(t *testing.T) {
	for i, testcase := range []struct {
		Input     Geometry
		Same      []Geometry
		Different []Geometry
	}{
		{
			Input: WithBBox([]float64{1, 2}, &Point{1, 2}),
			Same: []Geometry{
				WithBBox([]float64{1, 2}, &Point{1, 2}),
			},
			Different: []Geometry{
				&Point{1, 2},
				WithBBox([]float64{3, 4}, &Point{1, 2}),
				WithBBox([]float64{1, 2, 3, 4}, &Point{1, 2}),
			},
		},
	} {
		for _, same := range testcase.Same {
			if !testcase.Input.Equal(same) {
				t.Fatalf("(case %d) expected %#v to equal %#v", i, same, testcase.Input)
			}
		}
		for _, other := range testcase.Different {
			if testcase.Input.Equal(other) {
				t.Fatalf("(case %d) expected %#v to not equal %#v", i, other, testcase.Input)
			}
		}
	}
}

func TestBBoxMarshal(t *testing.T) {
	// Pass
	for i, testcase := range []struct {
		Input    Geometry
		Expected []byte
	}{
		{
			Input:    WithBBox([]float64{1, 2}, &Point{1, 2}),
			Expected: []byte(`{"type":"Point","coordinates":[1,2],"bbox":[1,2]}`),
		},
	} {
		got, err := json.Marshal(testcase.Input)
		if err != nil {
			t.Fatal(err)
		}
		if expected := testcase.Expected; !bytes.Equal(expected, got) {
			t.Fatalf("(case %d) expected %q, got %q", i, expected, got)
		}
	}

	// Fail
	for i, geom := range []Geometry{
		badGeom{},
	} {
		geom := WithBBox([]float64{}, geom)
		if _, err := json.Marshal(geom); err == nil {
			t.Fatalf("(case %d) expected error, got nil", i)
		}
	}
}

func TestBBoxUnmarshal(t *testing.T) {
	// Pass
	for i, testcase := range []struct {
		Input    []byte
		Expected Geometry
	}{
		// Geometry
		{
			Input:    []byte(`{"type":"Point","coordinates":[1,2],"bbox":[1,2]}`),
			Expected: WithBBox([]float64{1, 2}, &Point{1, 2}),
		},
		// Feature
		{
			Input:    []byte(`{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2]},"bbox":[1,2]}`),
			Expected: WithBBox([]float64{1, 2}, &Feature{Geometry: &Point{1, 2}}),
		},
		// FeatureCollection
		{
			Input: []byte(`{"type":"FeatureCollection","features":[{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2]}}],"bbox":[1,2]}`),
			Expected: WithBBox([]float64{1, 2}, &FeatureCollection{
				&Feature{Geometry: &Point{1, 2}},
			}),
		},
		// GeometryCollection
		{
			Input:    []byte(`{"type":"GeometryCollection","geometries":[{"type":"Point","coordinates":[1,2]}],"bbox":[1,2]}`),
			Expected: WithBBox([]float64{1, 2}, &GeometryCollection{&Point{1, 2}}),
		},
	} {
		geom, err := UnmarshalJSON(testcase.Input)
		if err != nil {
			t.Fatalf("(case %d) %s", i, err)
		}
		if expected, got := testcase.Expected, geom; !expected.Equal(got) {
			t.Fatalf("(case %d) expected %#v to equal %#v", i, got, expected)
		}
	}
}
