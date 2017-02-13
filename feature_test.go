package geo

import (
	"encoding/json"
	"testing"
)

func TestFeatureEqual(t *testing.T) {
	// Point
	cases{
		G: &Feature{
			Geometry: &Point{1, 1},
		},
		Different: []Geometry{
			&Line{{0, 0}, {1, 1}},
		},
	}.test(t)

	// Polygon
	cases{
		G: &Feature{
			Geometry: &Polygon{
				{
					{0, 0},
					{2, 0},
					{2, 2},
					{0, 2},
					{0, 0},
				},
			},
		},
		Inside: []Point{
			{1, 1},
		},
		Outside: []Point{
			{12, 12},
		},
	}.test(t)
}

func TestFeatureMarshal(t *testing.T) {
	// Pass
	for _, testcase := range []struct {
		Feature  Feature
		Expected string
	}{
		{
			Feature: Feature{
				Geometry: &Point{1, 2},
			},
			Expected: `{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2]},"properties":null}`,
		},
		{
			Feature: Feature{
				Geometry: &Line{{0, 0}, {1, 2}},
			},
			Expected: `{"type":"Feature","geometry":{"type":"LineString","coordinates":[[0,0],[1,2]]},"properties":null}`,
		},
	} {
		got, err := testcase.Feature.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != testcase.Expected {
			t.Fatalf("expected %s, got %s", testcase.Expected, string(got))
		}
	}
	// Fail
	for _, testcase := range []struct {
		Feature  Feature
		Expected string
	}{
		{Feature: Feature{Geometry: badGeom{}}},
		{Feature: Feature{Geometry: &Point{1, 2}, Properties: badGeom{}}},
	} {

		if _, err := testcase.Feature.MarshalJSON(); err == nil {
			t.Fatal("expected error, got nil")
		}
	}
}

func TestFeatureScan(t *testing.T) {
	// Pass
	scanTestcases{
		{
			Input: `POINT(0 0)`,
			Instance: &Feature{
				Geometry: &Point{0, 0},
			},
		},
		{
			Input: `LINESTRING(0 0, 1 1)`,
			Instance: &Feature{
				Geometry: &Line{{0, 0}, {1, 1}},
			},
		},
		{
			Input: `POLYGON((-1 1, 1 1, 1 -1, -1 -1, -1 1))`,
			Instance: &Feature{
				Geometry: &Polygon{{{-1, 1}, {1, 1}, {1, -1}, {-1, -1}, {-1, 1}}},
			},
		},
		{
			Input: "CIRCULARSTRING(1 0, 0 1, -1 0, 0 -1, 1 0)",
			Instance: &Feature{
				Geometry: &Circle{Radius: 1, Coordinates: Point{0, 0}},
			},
		},
	}.pass(t)

	// Fail
	scanTestcases{
		{
			Input:    `POINT(0,0)`,
			Instance: &Feature{},
		},
		{
			Input:    `LINESTRING(#$*@#$%%)`,
			Instance: &Feature{},
		},
		{
			Input:    `POLYGON(#$*@#$%%)`,
			Instance: &Feature{},
		},
		{
			Input:    `CIRCULARSTRING(#$*@#$%%)`,
			Instance: &Feature{},
		},
		{
			Input:    `GARBAGE(#$*@#$%%)`,
			Instance: &Feature{},
		},
	}.fail(t)
}

func TestFeatureUnmarshal(t *testing.T) {
	// Pass
	for i, testcase := range []struct {
		Input    []byte
		Expected Feature
		Instance *Feature
	}{
		{
			Input: []byte(`{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2]},"properties":null}`),
			Expected: Feature{
				Geometry: &Point{1, 2},
			},
			Instance: &Feature{},
		},
		{
			Input: []byte(`{"type":"Feature","geometry":{"type":"Point","coordinates":[3,2.5]}}`),
			Expected: Feature{
				Geometry: &Point{3, 2.5},
			},
			Instance: &Feature{},
		},
		{
			Input: []byte(`{"type":"Feature","geometry":{"type":"LineString","coordinates":[[0,0],[3,2.5]]}}`),
			Expected: Feature{
				Geometry: &Line{{0, 0}, {3, 2.5}},
			},
			Instance: &Feature{},
		},
		{
			Input: []byte(`{"type":"Feature","geometry":{"type":"Polygon","coordinates":[[[-113.1454418321263,33.52932582146817],[-113.1454418321263,33.52897252424949],[-113.1454027724575,33.52897252424949],[-113.1454027724575,33.52932582146817],[-113.1454418321263,33.52932582146817]]]}}`),
			Expected: Feature{
				Geometry: &Polygon{
					{
						{-113.1454418321263, 33.52932582146817},
						{-113.1454418321263, 33.52897252424949},
						{-113.1454027724575, 33.52897252424949},
						{-113.1454027724575, 33.52932582146817},
						{-113.1454418321263, 33.52932582146817},
					},
				},
			},
			Instance: &Feature{},
		},
		{
			Input: []byte(`{"type":"Feature","geometry":{"type":"MultiPoint","coordinates":[[0,0],[3,2.5]]}}`),
			Expected: Feature{
				Geometry: &MultiPoint{{0, 0}, {3, 2.5}},
			},
			Instance: &Feature{},
		},
		{
			Input: []byte(`{"type":"Feature","geometry":{"type":"MultiLineString","coordinates":[[[-113.1454418321263,33.52932582146817],[-113.1454418321263,33.52897252424949],[-113.1454027724575,33.52897252424949],[-113.1454027724575,33.52932582146817],[-113.1454418321263,33.52932582146817]]]}}`),
			Expected: Feature{
				Geometry: &MultiLine{
					{
						{-113.1454418321263, 33.52932582146817},
						{-113.1454418321263, 33.52897252424949},
						{-113.1454027724575, 33.52897252424949},
						{-113.1454027724575, 33.52932582146817},
						{-113.1454418321263, 33.52932582146817},
					},
				},
			},
			Instance: &Feature{},
		},
		{
			Input: []byte(`{"type":"Feature","geometry":{"type":"MultiPolygon","coordinates":[[[[-113.1454418321263,33.52932582146817],[-113.1454418321263,33.52897252424949],[-113.1454027724575,33.52897252424949],[-113.1454027724575,33.52932582146817],[-113.1454418321263,33.52932582146817]]]]}}`),
			Expected: Feature{
				Geometry: &MultiPolygon{
					{
						{
							{-113.1454418321263, 33.52932582146817},
							{-113.1454418321263, 33.52897252424949},
							{-113.1454027724575, 33.52897252424949},
							{-113.1454027724575, 33.52932582146817},
							{-113.1454418321263, 33.52932582146817},
						},
					},
				},
			},
			Instance: &Feature{},
		},
		{
			Input: []byte(`{"type":"Feature","geometry":{"type":"Circle","coordinates":[3,2.5],"radius":1}}`),
			Expected: Feature{
				Geometry: &Circle{
					Coordinates: Point{3, 2.5},
					Radius:      1,
				},
			},
			Instance: &Feature{},
		},
	} {
		if err := json.Unmarshal(testcase.Input, testcase.Instance); err != nil {
			t.Fatalf("fail case %d: %s", i, err)
		}
		if expected, got := testcase.Expected.Geometry, testcase.Instance.Geometry; !expected.Equal(got) {
			t.Fatalf("(case %d) expected %v, got %v", i, expected, got)
		}
	}
	// Fail
	for _, testcase := range []struct {
		Input    []byte
		Instance *Feature
	}{
		// "Feature" is misspelled
		{
			Input:    []byte(`{"type":"Faeture","geometry":{"type":"Point","coordinates":[1,2]},"properties":null}`),
			Instance: &Feature{},
		},
		// Bad Geometry type
		{
			Input:    []byte(`{"type":"Feature","geometry":{"type":"NoBueno","coordinates":[1,2]},"properties":null}`),
			Instance: &Feature{},
		},
		// Bad Geometry JSON
		{
			Input:    []byte(`{"type":"Feature","geometry":"foo"}`),
			Instance: &Feature{},
		},
		// Bad JSON for properties
		{
			Input:    []byte(`{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2]},"properties":{[['','':}}`),
			Instance: &Feature{},
		},
		// Bad circle coordinates
		{
			Input:    []byte(`{"type":"Feature","geometry":{"type":"Circle","coordinates":"foo"}}`),
			Instance: &Feature{},
		},
	} {

		if err := testcase.Instance.UnmarshalJSON(testcase.Input); err == nil {
			t.Fatalf("expected error, got nil for %s", string(testcase.Input))
		}
	}
}

func TestFeatureValue(t *testing.T) {
	valueTestcases{
		{
			Input: &Feature{
				Geometry: &Point{0, 0},
			},
			Expected: `POINT(0 0)`,
		},
	}.pass(t)
}
