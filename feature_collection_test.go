package geo

import "testing"

func TestFeatureCollectionCompare(t *testing.T) {
	// Fail
	compareTestcases{
		{
			G1: &FeatureCollection{},
			G2: &Point{},
		},
		{
			G1: &FeatureCollection{
				Feature{
					Geometry: &Point{1, 1},
				},
			},
			G2: &FeatureCollection{},
		},
		{
			G1: &FeatureCollection{
				Feature{
					Geometry: &Point{1, 1},
				},
			},
			G2: &FeatureCollection{
				Feature{
					Geometry: &Line{{0, 0}, {1, 1}},
				},
			},
		},
	}.fail(t)
}

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
						{
							{1, 1},
							{1, -1},
							{-1, -1},
							{-1, 1},
						},
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
					Geometry: badGeom{},
				},
				{
					Geometry: &Polygon{
						{
							{1, 1},
							{1, -1},
							{-1, -1},
							{-1, 1},
						},
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

func TestFeatureCollectionScan(t *testing.T) {
	// TODO
}

func TestFeatureCollectionUnmarshal(t *testing.T) {
	// Pass
	unmarshalTestcases{
		{
			Input: []byte(`{"type": "FeatureCollection", "features": [{"geometry": {"type": "Polygon", "coordinates": [[[-113.131642956287, 33.4246272997084], [-113.133949656039, 33.4246272997084], [-113.133960384876, 33.4210006893439], [-113.131696600467, 33.4210006893439], [-113.131642956287, 33.4246272997084]]]}, "type": "Feature", "properties": {"id": "3f166f14-e8f3-454b-b7aa-401c8ee81c8d"}}]}`),
			Expected: &FeatureCollection{
				Feature{
					Geometry: &Polygon{
						{
							{-113.131642956287, 33.4246272997084},
							{-113.133949656039, 33.4246272997084},
							{-113.133960384876, 33.4210006893439},
							{-113.131696600467, 33.4210006893439},
							{-113.131642956287, 33.4246272997084},
						},
					},
				},
			},
			Instance: &FeatureCollection{},
		},
	}.pass(t)
}

func TestFeatureCollectionValue(t *testing.T) {
	// TODO
}
