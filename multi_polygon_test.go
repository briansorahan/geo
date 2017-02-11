package geo

import "testing"

func TestMultiPolygonCompare(t *testing.T) {
	cases{
		G: &MultiPolygon{
			{
				{
					{1.2, 3.4},
					{5.6, 7.8},
					{1.4, 9.3},
					{-1.7, 7.3},
				},
			},
		},
		Different: []Geometry{
			&MultiPolygon{
				{
					{
						{1.2, 3.4},
						{5.6, 7.8},
						{1.4, 9.3},
					},
				},
			},
			&MultiPolygon{
				{
					{
						{1.2, 3.4},
						{5.6, 7.8},
					},
					{
						{1.4, 9.3},
					},
				},
			},
			&MultiPolygon{
				{
					{
						{1.2, 3.4},
						{5.6, 7.8},
						{1.4, 9.3},
					},
				},
			},
			&MultiPolygon{
				{
					{
						{1.2, 3.4},
						{5.6, 7.8},
						{1.4, 9.3},
						{-1.4, 7.3},
					},
				},
			},
			&MultiPolygon{
				{
					{
						{1.2, 3.4},
						{5.6, 7.8},
						{1.4, 9.3},
						{-1.7, 7.5},
					},
				},
			},
			&MultiPolygon{
				{
					{
						{1.2, 3.4},
						{5.6, 7.8},
					},
				},
				{
					{
						{1.4, 9.3},
						{-1.7, 7.5},
					},
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

func TestMultiPolygonContains(t *testing.T) {
	// Contains (square)
	cases{
		// Square
		G: &MultiPolygon{
			{
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
			{4, 1},
		},
	}.test(t)

	// Contains (hexagon)
	cases{
		G: &MultiPolygon{
			{
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
		},
		Inside: []Point{
			{1, 0},
		},
	}.test(t)

	// Contains (quadrilateral)
	cases{
		G: &MultiPolygon{
			{
				{
					{-1, 10},
					{10, 1},
					{1, -10},
					{-10, -1},
					{-1, 10},
				},
			},
		},
		Inside: []Point{
			{2, 2},
			{2, -2},
		},
	}.test(t)

	// Contains (horizontal ray intersects two vertices)
	cases{
		G: &MultiPolygon{
			{
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
		},
		Inside: []Point{
			{1, 2},
		},
		Outside: []Point{
			{-1, 2},
		},
	}.test(t)
}

func TestMultiPolygonEmpty(t *testing.T) {
	var (
		p        = MultiPolygon{}
		expected = "MULTIPOLYGON EMPTY"
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

func TestMultiPolygonMarshal(t *testing.T) {
	// Pass
	marshalTestcases{
		{
			Input: &MultiPolygon{
				{
					{
						{1.2, 3.4},
						{5.6, 7.8},
					},
				},
			},
			Expected: `{"type":"MultiPolygon","coordinates":[[[[1.2,3.4],[5.6,7.8]]]]}`,
		},
		{
			Input: &MultiPolygon{
				{
					{
						{1.2, 3.4},
						{5.6, 7.8},
						{2, 6},
					},
				},
				{
					{
						{0, 0},
						{0, 1},
						{1, 0},
					},
					{
						{3, 3},
						{3, 4},
						{4, 3},
					},
				},
			},
			Expected: `{"type":"MultiPolygon","coordinates":[[[[1.2,3.4],[5.6,7.8],[2,6]]],[[[0,0],[0,1],[1,0]],[[3,3],[3,4],[4,3]]]]}`,
		},
	}.pass(t)
}

func TestMultiPolygonScan(t *testing.T) {
	// Pass
	for i, testcase := range []struct {
		WKT      string
		Expected Geometry
	}{
		{
			WKT: "MULTIPOLYGON(((1.2 3.4, 5.6 7.8, 6.2 1.5, 1.2 3.4)),((-4 -4, -4 4, 4 4, 4 -4)))",
			Expected: &MultiPolygon{
				{
					{
						{1.2, 3.4},
						{5.6, 7.8},
						{6.2, 1.5},
						{1.2, 3.4},
					},
				},
				{
					{
						{-4, -4},
						{-4, 4},
						{4, 4},
						{4, -4},
					},
				},
			},
		},
		{
			WKT: `MULTIPOLYGON(((-113.14448537305 33.4192544895836,-113.140408415347 33.4192634445438,-113.140419144183 33.4209917345629,-113.144506830722 33.4210096441239,-113.14448537305 33.4192544895836)))`,
			Expected: &MultiPolygon{
				{
					{
						{-113.14448537305, 33.4192544895836},
						{-113.140408415347, 33.4192634445438},
						{-113.140419144183, 33.4209917345629},
						{-113.144506830722, 33.4210096441239},
						{-113.14448537305, 33.4192544895836},
					},
				},
			},
		},
	} {
		p := &MultiPolygon{}
		if err := p.Scan(testcase.WKT); err != nil {
			t.Fatalf("(case %d) %s", i, err)
		}
		if expected, got := testcase.Expected, p; !got.Compare(expected) {
			t.Fatalf("(case %d) expected %f, got %f", i, expected, got)
		}
	}
	// Fail
	for _, testcase := range []interface{}{
		"MULTIPOLYGON(((1.2, 3.4, 5.6, 7.8)))",
		"MULTIPOLYGON(((1.2 3.4, 5.6 7.8)),((1.2, 3.4),(a)",
		[]byte("MULTIPOLYGON(((1.2, 3.4, 5.6, 7.8)))"),
		7,
		"MULTIPOLYGON(1.2 3.4 5.6 7.8)",
		"MULTIPOLYGON((1.2 3.4 5.6 7.8)}",
		"PIKACHU",
	} {
		p := &MultiPolygon{}
		if err := p.Scan(testcase); err == nil {
			t.Fatalf("expected err, got nil for %s", testcase.(string))
		}
	}
}

func TestMultiPolygonString(t *testing.T) {
	for _, c := range []struct {
		Input    MultiPolygon
		Expected string
	}{
		{
			Input: MultiPolygon{
				{
					{
						{-4, -4},
						{-4, 4},
						{4, 4},
						{4, -4},
					},
					{
						{5, 5},
						{5, 6},
						{6, 5},
					},
				},
				{
					{
						{0, 0},
						{1, 0},
						{0, 1},
					},
				},
			},
			Expected: `MULTIPOLYGON(((-4 -4, -4 4, 4 4, 4 -4),(5 5, 5 6, 6 5)),((0 0, 1 0, 0 1)))`,
		},
	} {
		if expected, got := c.Expected, c.Input.String(); expected != got {
			t.Fatalf("expected %s, got %s", expected, got)
		}
	}
}

func TestMultiPolygonUnmarshalJSON(t *testing.T) {
	unmarshalTestcases{
		{
			Input:    []byte(`{"type":"MultiPolygon","coordinates":[[[[0,0],[0,1],[1,1],[1,0]]]]}`),
			Instance: &MultiPolygon{},
			Expected: &MultiPolygon{
				{
					{
						{0, 0},
						{0, 1},
						{1, 1},
						{1, 0},
					},
				},
			},
		},
	}.pass(t)

	unmarshalTestcases{
		{
			Input:    []byte(`{"type":"Porygon","coordinates":[[[0,0],[0,1],[1,1],[1,0]]]}`),
			Instance: &MultiPolygon{},
		},
		{
			Input:    []byte(`{"type":"MultiPolygon","coordinates":"bjork"}`),
			Instance: &MultiPolygon{},
		},
	}.fail(t)
}

func TestMultiPolygonValue(t *testing.T) {
	var (
		p = MultiPolygon{
			{
				{
					{1.2, 3.4},
					{5.6, 7.8},
					{8.7, 6.5},
					{4.3, 2.1},
				},
			},
		}
		expected = `MULTIPOLYGON(((1.2 3.4, 5.6 7.8, 8.7 6.5, 4.3 2.1)))`
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
