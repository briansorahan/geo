package geo

import "testing"

func TestPolygonCompare(t *testing.T) {
	cases{
		G: &Polygon{
			{
				{1.2, 3.4},
				{5.6, 7.8},
				{1.4, 9.3},
				{-1.7, 7.3},
			},
		},
		Different: []Geometry{
			&Polygon{
				{
					{1.2, 3.4},
					{5.6, 7.8},
					{1.4, 9.3},
				},
			},
			&Polygon{
				{
					{1.2, 3.4},
					{5.6, 7.8},
				},
				{
					{1.4, 9.3},
				},
			},
			&Polygon{
				{
					{1.2, 3.4},
					{5.6, 7.8},
					{1.4, 9.3},
				},
			},
			&Polygon{
				{
					{1.2, 3.4},
					{5.6, 7.8},
					{1.4, 9.3},
					{-1.4, 7.3},
				},
			},
			&Polygon{
				{
					{1.2, 3.4},
					{5.6, 7.8},
					{1.4, 9.3},
					{-1.7, 7.5},
				},
			},
			&Line{
				{1.2, 3.4},
				{5.6, 7.8},
				{1.4, 9.3},
				{-1.7, 7.5},
			},
			&Point{1.2, 3.4},
		},
	}.test(t)
}

func TestPolygonContains(t *testing.T) {
	// Contains (square)
	cases{
		// Square
		G: &Polygon{
			{
				{0, 0},
				{2, 0},
				{2, 2},
				{0, 2},
				{0, 0},
			},
		},
		Inside: []Point{
			{1, 1},
		},
		Outside: []Point{
			{4, 1},
		},
	}.test(t)

	// Contains (hexagon)
	cases{
		G: &Polygon{
			{
				{0, 1},
				{1, 2},
				{2, 1},
				{2, 0},
				{1, -1},
				{0, 0},
				{0, 1},
			},
		},
		Inside: []Point{
			{1, 0},
		},
	}.test(t)

	// Contains (quadrilateral)
	cases{
		G: &Polygon{
			{
				{-1, 10},
				{10, 1},
				{1, -10},
				{-10, -1},
				{-1, 10},
			},
		},
		Inside: []Point{
			{2, 2},
			{2, -2},
		},
	}.test(t)

	// Contains (horizontal ray intersects two vertices)
	cases{
		G: &Polygon{
			{
				{0, 0},
				{0, 4},
				{2, 4},
				{3, 2},
				{4, 4},
				{6, 4},
				{8, 2},
				{6, 0},
			},
		},
		Inside: []Point{
			{1, 2},
		},
		Outside: []Point{
			{-1, 2},
		},
	}.test(t)
}

func TestPolygonEmpty(t *testing.T) {
	var (
		p        = Polygon{}
		expected = "POLYGON EMPTY"
	)
	value, err := p.Value()
	if err != nil {
		t.Fatal(err)
	}
	got, ok := value.(string)
	if !ok {
		t.Fatalf("expected string, got %T", value)
	}
	if expected != got {
		t.Fatalf("expected %s, got %s", expected, got)
	}
}

func TestPolygonMarshal(t *testing.T) {
	// Pass
	marshalTestcases{
		{
			Input: &Polygon{
				{
					{1.2, 3.4},
					{5.6, 7.8},
				},
			},
			Expected: `{"type":"Polygon","coordinates":[[[1.2,3.4],[5.6,7.8]]]}`,
		},
		{
			Input: &Polygon{
				{
					{1.2, 3.4},
					{5.6, 7.8},
					{2, 6},
				},
				{
					{0, 0},
					{0, 1},
					{1, 0},
				},
			},
			Expected: `{"type":"Polygon","coordinates":[[[1.2,3.4],[5.6,7.8],[2,6]],[[0,0],[0,1],[1,0]]]}`,
		},
	}.pass(t)
}

func TestPolygonScan(t *testing.T) {
	// Pass
	for _, testcase := range []struct {
		WKT      string
		Expected Polygon
	}{
		{
			WKT: "POLYGON((1.2 3.4, 5.6 7.8, 6.2 1.5, 1.2 3.4))",
			Expected: Polygon{
				{
					{1.2, 3.4},
					{5.6, 7.8},
					{6.2, 1.5},
					{1.2, 3.4},
				},
			},
		},
		{
			WKT: `POLYGON((-113.14448537305 33.4192544895836,-113.140408415347 33.4192634445438,-113.140419144183 33.4209917345629,-113.144506830722 33.4210096441239,-113.14448537305 33.4192544895836))`,
			Expected: Polygon{
				{
					{-113.14448537305, 33.4192544895836},
					{-113.140408415347, 33.4192634445438},
					{-113.140419144183, 33.4209917345629},
					{-113.144506830722, 33.4210096441239},
					{-113.14448537305, 33.4192544895836},
				},
			},
		},
	} {
		p := &Polygon{}
		if err := p.Scan(testcase.WKT); err != nil {
			t.Fatal(err)
		}
		for i, coord := range testcase.Expected {
			if expected, got := coord[0], (*p)[i][0]; expected != got {
				t.Fatalf("expected %f, got %f", expected, got)
			}
			if expected, got := coord[1], (*p)[i][1]; expected != got {
				t.Fatalf("expected %f, got %f", expected, got)
			}
		}
	}
	// Fail
	for _, testcase := range []interface{}{
		"POLYGON((1.2, 3.4, 5.6, 7.8))",
		[]byte("POLYGON((1.2, 3.4, 5.6, 7.8))"),
		7,
		"POLYGON(1.2 3.4 5.6 7.8)",
		"POLYGON((1.2 3.4 5.6 7.8)}",
		"PIKACHU",
	} {
		p := &Polygon{}
		if err := p.Scan(testcase); err == nil {
			t.Fatalf("expected err, got nil for %s", testcase.(string))
		}
	}
}

func TestPolygonString(t *testing.T) {
	for _, c := range []struct {
		Input    Polygon
		Expected string
	}{
		{
			Input: Polygon{
				{
					{-4, -4},
					{-4, 4},
					{4, 4},
					{4, -4},
				},
				{
					{0, 0},
					{1, 0},
					{0, 1},
				},
			},
			Expected: `POLYGON((-4 -4, -4 4, 4 4, 4 -4),(0 0, 1 0, 0 1))`,
		},
	} {
		if expected, got := c.Expected, c.Input.String(); expected != got {
			t.Fatalf("expected %s, got %s", expected, got)
		}
	}
}

func TestPolygonUnmarshalJSON(t *testing.T) {
	unmarshalTestcases{
		{
			Input:    []byte(`{"type":"Polygon","coordinates":[[[0,0],[0,1],[1,1],[1,0]]]}`),
			Instance: &Polygon{},
			Expected: &Polygon{
				{
					{0, 0},
					{0, 1},
					{1, 1},
					{1, 0},
				},
			},
		},
	}.pass(t)

	unmarshalTestcases{
		{
			Input:    []byte(`{"type":"Porygon","coordinates":[[[0,0],[0,1],[1,1],[1,0]]]}`),
			Instance: &Polygon{},
		},
		{
			Input:    []byte(`{"type":"Polygon","coordinates":"bjork"}`),
			Instance: &Polygon{},
		},
	}.fail(t)
}

func TestPolygonValue(t *testing.T) {
	var (
		p = Polygon{
			{
				{1.2, 3.4},
				{5.6, 7.8},
				{8.7, 6.5},
				{4.3, 2.1},
			},
		}
		expected = `POLYGON((1.2 3.4, 5.6 7.8, 8.7 6.5, 4.3 2.1))`
	)
	value, err := p.Value()
	if err != nil {
		t.Fatal(err)
	}
	got, ok := value.(string)
	if !ok {
		t.Fatalf("expected string, got %T", value)
	}
	if expected != got {
		t.Fatalf("expected %s, got %s", expected, got)
	}
}
