package geo

import "testing"

func TestGeometryCollection(t *testing.T) {
	// Equal
	cases{
		G: &GeometryCollection{
			&Point{1, 1},
		},
		Different: []Geometry{
			&GeometryCollection{},
			&GeometryCollection{
				&Line{{0, 0}, {1, 1}},
			},
			&Point{0, 0},
		},
	}.test(t)

	// Contains
	cases{
		G: &GeometryCollection{
			&Polygon{
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

	// Empty GeometryCollection
	cases{
		G: &GeometryCollection{},
		Outside: []Point{
			{12, 12},
		},
	}.test(t)
}

func TestGeometryCollectionMarshal(t *testing.T) {
	// Pass
	for _, testcase := range []struct {
		Coll     GeometryCollection
		Expected string
	}{
		{
			Coll: GeometryCollection{
				&Point{1, 2},
				&Polygon{
					{
						{1, 1},
						{1, -1},
						{-1, -1},
						{-1, 1},
					},
				},
			},
			Expected: `{"type":"GeometryCollection","geometries":[{"type":"Point","coordinates":[1,2]},{"type":"Polygon","coordinates":[[[1,1],[1,-1],[-1,-1],[-1,1]]]}]}`,
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
		Coll     GeometryCollection
		Expected string
	}{
		{
			Coll: GeometryCollection{
				badGeom{},
				&Polygon{
					{
						{1, 1},
						{1, -1},
						{-1, -1},
						{-1, 1},
					},
				},
			},
			Expected: `{"type":"GeometryCollection","geometries":[{"type":"Point","coordinates":[1,2]},{"type":"Polygon","coordinates":[[[1,1],[1,-1],[-1,-1],[-1,1]]]}]}`,
		},
	} {
		if _, err := testcase.Coll.MarshalJSON(); err == nil {
			t.Fatal("expected error, got nil")
		}
	}
}

func TestGeometryCollectionScan(t *testing.T) {
	scanTestcases{
		{
			Input:    `GEOMETRYCOLLECTION(POINT(0 0), LINESTRING(0 0, 1 1))`,
			Instance: &GeometryCollection{},
		},
	}.pass(t)
}

func TestGeometryCollectionUnmarshal(t *testing.T) {
	// Pass
	unmarshalTestcases{
		{
			Input: []byte(`{"type": "GeometryCollection", "geometries": [{"type": "Polygon", "coordinates": [[[-113.131642956287, 33.4246272997084], [-113.133949656039, 33.4246272997084], [-113.133960384876, 33.4210006893439], [-113.131696600467, 33.4210006893439], [-113.131642956287, 33.4246272997084]]]}]}`),
			Expected: &GeometryCollection{
				&Polygon{
					{
						{-113.131642956287, 33.4246272997084},
						{-113.133949656039, 33.4246272997084},
						{-113.133960384876, 33.4210006893439},
						{-113.131696600467, 33.4210006893439},
						{-113.131642956287, 33.4246272997084},
					},
				},
			},
			Instance: &GeometryCollection{},
		},
	}.pass(t)

	// Fail
	unmarshalTestcases{
		// Garbage
		{
			Input:    []byte(`$#$//&(*$/&#`),
			Instance: &GeometryCollection{},
		},
		// Bad type
		{
			Input:    []byte(`{"type":"GeometryColleccion","geometries":[]}`),
			Instance: &GeometryCollection{},
		},
		// Bad geometry
		{
			Input:    []byte(`{"type":"GeometryCollection","geometries":[{"type":"Polygon","coordinates":[{"key":"value"}]}]}`),
			Instance: &GeometryCollection{},
		},
	}.fail(t)
}

func TestGeometryCollectionValue(t *testing.T) {
	// Pass
	valueTestcases{
		{
			Input: &GeometryCollection{
				&Point{0, 0},
				&Line{{-1, 1}, {1, -1}},
			},
			Expected: `GEOMETRYCOLLECTION(POINT(0 0), LINESTRING(-1 1, 1 -1))`,
		},
	}.pass(t)
}
