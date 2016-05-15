package geo

import "testing"

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
