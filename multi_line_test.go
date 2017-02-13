package geo

import "testing"

func TestMultiLineEqual(t *testing.T) {
	cases{
		G: &MultiLine{
			{
				{1.2, 3.4},
				{5.6, 7.8},
				{1.4, 9.3},
				{-1.7, 7.3},
			},
		},
		Different: []Geometry{
			&MultiLine{
				{
					{1.2, 3.4},
					{5.6, 7.8},
					{1.4, 9.3},
				},
			},
			&MultiLine{
				{
					{1.2, 3.4},
					{5.6, 7.8},
				},
				{
					{1.4, 9.3},
				},
			},
			&MultiLine{
				{
					{1.2, 3.4},
					{5.6, 7.8},
					{1.4, 9.3},
				},
			},
			&MultiLine{
				{
					{1.2, 3.4},
					{5.6, 7.8},
					{1.4, 9.3},
					{-1.4, 7.3},
				},
			},
			&MultiLine{
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

func TestMultiLineContains(t *testing.T) {
	// Contains (square)
	cases{
		// Square
		G: &MultiLine{
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
		G: &MultiLine{
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
		G: &MultiLine{
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
		G: &MultiLine{
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

func TestMultiLineEmpty(t *testing.T) {
	var (
		p        = MultiLine{}
		expected = "MULTILINESTRING EMPTY"
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

func TestMultiLineMarshal(t *testing.T) {
	// Pass
	marshalTestcases{
		{
			Input: &MultiLine{
				{
					{1.2, 3.4},
					{5.6, 7.8},
				},
			},
			Expected: `{"type":"MultiLineString","coordinates":[[[1.2,3.4],[5.6,7.8]]]}`,
		},
		{
			Input: &MultiLine{
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
			Expected: `{"type":"MultiLineString","coordinates":[[[1.2,3.4],[5.6,7.8],[2,6]],[[0,0],[0,1],[1,0]]]}`,
		},
	}.pass(t)
}

func TestMultiLineScan(t *testing.T) {
	// Pass
	for _, testcase := range []struct {
		WKT      string
		Expected MultiLine
	}{
		{
			WKT: "MULTILINESTRING((1.2 3.4, 5.6 7.8, 6.2 1.5, 1.2 3.4))",
			Expected: MultiLine{
				{
					{1.2, 3.4},
					{5.6, 7.8},
					{6.2, 1.5},
					{1.2, 3.4},
				},
			},
		},
		{
			WKT: `MULTILINESTRING((-113.14448537305 33.4192544895836,-113.140408415347 33.4192634445438,-113.140419144183 33.4209917345629,-113.144506830722 33.4210096441239,-113.14448537305 33.4192544895836))`,
			Expected: MultiLine{
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
		p := &MultiLine{}
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
		"MULTILINESTRING((1.2, 3.4, 5.6, 7.8))",
		[]byte("MULTILINESTRING((1.2, 3.4, 5.6, 7.8))"),
		7,
		"MULTILINESTRING(1.2 3.4 5.6 7.8)",
		"MULTILINESTRING((1.2 3.4 5.6 7.8)}",
		"PIKACHU",
	} {
		p := &MultiLine{}
		if err := p.Scan(testcase); err == nil {
			t.Fatalf("expected err, got nil for %s", testcase.(string))
		}
	}
}

func TestMultiLineString(t *testing.T) {
	for _, c := range []struct {
		Input    MultiLine
		Expected string
	}{
		{
			Input: MultiLine{
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
			Expected: `MULTILINESTRING((-4 -4, -4 4, 4 4, 4 -4),(0 0, 1 0, 0 1))`,
		},
	} {
		if expected, got := c.Expected, c.Input.String(); expected != got {
			t.Fatalf("expected %s, got %s", expected, got)
		}
	}
}

func TestMultiLineUnmarshalJSON(t *testing.T) {
	unmarshalTestcases{
		{
			Input:    []byte(`{"type":"MultiLineString","coordinates":[[[0,0],[0,1],[1,1],[1,0]]]}`),
			Instance: &MultiLine{},
			Expected: &MultiLine{
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
			Instance: &MultiLine{},
		},
		{
			Input:    []byte(`{"type":"MultiLineString","coordinates":"bjork"}`),
			Instance: &MultiLine{},
		},
	}.fail(t)
}

func TestMultiLineValue(t *testing.T) {
	var (
		p = MultiLine{
			{
				{1.2, 3.4},
				{5.6, 7.8},
				{8.7, 6.5},
				{4.3, 2.1},
			},
		}
		expected = `MULTILINESTRING((1.2 3.4, 5.6 7.8, 8.7 6.5, 4.3 2.1))`
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
