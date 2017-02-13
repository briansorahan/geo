package geo

import "testing"

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
}

func TestUnmarshalFeatureCollection(t *testing.T) {
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
}
