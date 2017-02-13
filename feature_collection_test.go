package geo

import "testing"

func TestFeatureCollection(t *testing.T) {
	// Equal
	cases{
		G: &FeatureCollection{
			Feature{
				Geometry: &Point{1, 1},
			},
		},
		Different: []Geometry{
			&FeatureCollection{},
			&FeatureCollection{
				Feature{
					Geometry: &Line{{0, 0}, {1, 1}},
				},
			},
			&Point{0, 0},
		},
	}.test(t)

	// Contains
	cases{
		G: &FeatureCollection{
			Feature{
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
		},
		Inside: []Point{
			{1, 1},
		},
		Outside: []Point{
			{12, 12},
		},
	}.test(t)
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
	scanTestcases{
		{
			Input:    `GEOMETRYCOLLECTION(POINT(0 0), LINESTRING(0 0, 1 1))`,
			Instance: &FeatureCollection{},
		},
	}.pass(t)
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

	// Fail
	unmarshalTestcases{
		{
			Input:    []byte(`$#$//&(*$/&#`),
			Instance: &FeatureCollection{},
		},
		{
			Input:    []byte(`{"type":"FeatureColleccion","features":[]}`),
			Instance: &FeatureCollection{},
		},
	}.fail(t)
}

func TestFeatureCollectionValue(t *testing.T) {
	// Pass
	valueTestcases{
		{
			Input: &FeatureCollection{
				Feature{
					Geometry: &Point{0, 0},
				},
				Feature{
					Geometry: &Line{{-1, 1}, {1, -1}},
				},
			},
			Expected: `GEOMETRYCOLLECTION(POINT(0 0), LINESTRING(-1 1, 1 -1))`,
		},
	}.pass(t)
}
