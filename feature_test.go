package geo

import (
	"encoding/json"
	"testing"
)

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

func TestFeatureUnmarshal(t *testing.T) {
	// Pass
	for _, testcase := range []struct {
		Input    []byte
		Expected Feature
		Instance *Feature
	}{
		{
			Input: []byte(`{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2]},"properties":null}`),
			Expected: Feature{
				Geometry: &Point{1, 2},
			},
			Instance: &Feature{
				Geometry: &Point{},
			},
		},
		{
			Input: []byte(`{"type":"Feature","geometry":{"type":"Point","coordinates":[3,2.5]}}`),
			Expected: Feature{
				Geometry: &Point{3, 2.5},
			},
			Instance: &Feature{
				Geometry: &Point{},
			},
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
			Instance: &Feature{
				Geometry: &Polygon{},
			},
		},
	} {
		if err := json.Unmarshal(testcase.Input, testcase.Instance); err != nil {
			t.Fatal(err)
		}
		if expected, got := testcase.Expected.Geometry, testcase.Instance.Geometry; !expected.Compare(got) {
			t.Fatalf("expected %v, got %v", expected, got)
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
	} {

		if err := testcase.Instance.UnmarshalJSON(testcase.Input); err == nil {
			t.Fatalf("expected error, got nil for %s", string(testcase.Input))
		}
	}
}
