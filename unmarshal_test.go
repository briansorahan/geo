package geo

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalFail(t *testing.T) {
	for i, data := range [][]byte{
		// Invalid JSON
		[]byte(`($*&%*)(`),

		// Invalid feature
		[]byte(`{"type":"Feature","geometry":{"type":"Point","coordinates":{"foo":"bar"}}}`),

		// Invalid feature collection
		[]byte(`{"type":"FeatureCollection","features":[{"type":"Feature","geometry":{"type":"Point","coordinates":{"foo":"bar"}}}]}`),
	} {
		if _, err := UnmarshalJSON(data); err == nil {
			t.Fatalf("(case %d) expected error, got nil", i)
		}
	}
}

func TestUnmarshalFeature(t *testing.T) {
	// Pass
	for i, testcase := range []struct {
		Input    []byte
		Expected Geometry
	}{
		{
			Input:    []byte(`{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2]}}`),
			Expected: &Feature{Geometry: &Point{1, 2}},
		},
	} {
		geom, err := UnmarshalJSON(testcase.Input)
		if err != nil {
			t.Fatalf("(case %d) %s", i, err)
		}
		if expected, got := testcase.Expected, geom; !expected.Equal(got) {
			t.Fatalf("(case %d) expected %#v, got %#v", i, expected, got)
		}
	}

	// Fail
	for i, input := range [][]byte{
		[]byte(`{"type":"Feature","geometry":{"type":"Point","}}`),
	} {

		if _, err := UnmarshalJSON(input); err == nil {
			t.Fatalf("(case %d) expected error, got nil", i)
		}
	}

	// Test the private feature type used for unmarshalling.
	f := &feature{
		Geometry: json.RawMessage(`{"type":"Point}`),
	}
	if _, err := f.ToFeature(); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUnmarshalFeatureCollection(t *testing.T) {
	// Pass
	for i, testcase := range []struct {
		Input    []byte
		Expected Geometry
	}{
		{
			Input: []byte(`{"type":"FeatureCollection","features":[{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2]}}]}`),
			Expected: &FeatureCollection{
				{Geometry: &Point{1, 2}},
			},
		},
	} {
		geom, err := UnmarshalJSON(testcase.Input)
		if err != nil {
			t.Fatalf("(case %d) %s", i, err)
		}
		if expected, got := testcase.Expected, geom; !expected.Equal(got) {
			t.Fatalf("(case %d) expected %#v, got %#v", i, expected, got)
		}
	}

	// Fail
	for i, input := range [][]byte{
		[]byte(`{"type":"FeatureCollection","features":3,"bbox":[]}`),
		[]byte(``),
	} {

		if _, err := UnmarshalJSON(input); err == nil {
			t.Fatalf("(case %d) expected error, got nil", i)
		}
	}
}

func TestUnmarshalGeometryCollection(t *testing.T) {
	for i, testcase := range []struct {
		Input    []byte
		Expected Geometry
	}{
		{
			Input: []byte(`{"type":"GeometryCollection","geometries":[{"type":"Point","coordinates":[1,2]}]}`),
			Expected: &GeometryCollection{
				&Point{1, 2},
			},
		},
	} {
		geom, err := UnmarshalJSON(testcase.Input)
		if err != nil {
			t.Fatalf("(case %d) %s", i, err)
		}
		if expected, got := testcase.Expected, geom; !expected.Equal(got) {
			t.Fatalf("(case %d) expected %#v, got %#v", i, expected, got)
		}
	}

	// Fail
	for i, input := range [][]byte{
		[]byte(`{"type":"GeometryCollection","geometries":3,"bbox":[]}`),
		[]byte(`{"type":"GeometryCollection","geometries":[{"type":"Point","coordinates":42}],"bbox":[]}`),
	} {

		if _, err := UnmarshalJSON(input); err == nil {
			t.Fatalf("(case %d) expected error, got nil", i)
		}
	}
}
