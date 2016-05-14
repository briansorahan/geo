package geo

import "testing"

func TestFeatureCollectionMarshal(t *testing.T) {
	// Pass
	for _, testcase := range []struct {
		Coll     FeatureCollection
		Expected string
	}{
		{
			Coll: FeatureCollection{
				{
					Geometry: &Point{1, 2},
				},
				{
					Geometry: &Polygon{
						{1, 1},
						{1, -1},
						{-1, -1},
						{-1, 1},
					},
				},
			},
			Expected: `{"type":"FeatureCollection","features":[{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2]},"properties":null},{"type":"Feature","geometry":{"type":"Polygon","coordinates":[[[1,1],[1,-1],[-1,-1],[-1,1]]]},"properties":null}]}`,
		},
	} {
		got, err := testcase.Coll.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != testcase.Expected {
			t.Fatalf("expected %s, got %s", testcase.Expected, string(got))
		}
	}
	// Fail
	for _, testcase := range []struct {
		Coll     FeatureCollection
		Expected string
	}{
		{
			Coll: FeatureCollection{
				{
					Geometry: &BadGeom{},
				},
				{
					Geometry: &Polygon{
						{1, 1},
						{1, -1},
						{-1, -1},
						{-1, 1},
					},
				},
			},
			Expected: `{"type":"FeatureCollection","features":[{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2]},"properties":null},{"type":"Feature","geometry":{"type":"Polygon","coordinates":[[[1,1],[1,-1],[-1,-1],[-1,1]]]},"properties":null}]}`,
		},
	} {
		if _, err := testcase.Coll.MarshalJSON(); err == nil {
			t.Fatal("expected error, got nil")
		}
	}
}
