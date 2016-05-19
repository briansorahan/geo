package geo

import "testing"

func TestLineCompare(t *testing.T) {
	// Pass
	compareTestcases{
		{
			G1: &Line{{1.2, 3.4}, {5.6, 7.8}, {1.4, 9.3}, {-1.7, 7.3}},
			G2: &Line{{1.2, 3.4}, {5.6, 7.8}, {1.4, 9.3}, {-1.7, 7.3}},
		},
	}.pass(t)

	// Fail
	compareTestcases{
		{
			G1: &Line{{1.2, 3.4}, {5.6, 7.8}, {1.4, 9.3}, {-1.7, 7.3}},
			G2: &Line{{1.2, 3.4}, {5.6, 7.8}, {1.4, 9.3}},
		},
		{
			G1: &Line{{1.2, 3.4}, {5.6, 7.8}, {1.4, 9.3}, {-1.7, 7.3}},
			G2: &Line{{1.2, 3.4}, {5.6, 7.8}, {1.4, 9.3}, {-1.4, 7.3}},
		},
		{
			G1: &Line{{1.2, 3.4}, {5.6, 7.8}, {1.4, 9.3}, {-1.7, 7.3}},
			G2: &Line{{1.2, 3.4}, {5.6, 7.8}, {1.4, 9.3}, {-1.7, 7.5}},
		},
		{
			G1: &Line{{1.2, 3.4}, {5.6, 7.8}, {1.4, 9.3}, {-1.7, 7.3}},
			G2: &Polygon{{{1.2, 3.4}, {5.6, 7.8}, {1.4, 9.3}, {-1.7, 7.5}}},
		},
		{
			G1: &Line{{1.2, 3.4}, {5.6, 7.8}, {1.4, 9.3}, {-1.7, 7.3}},
			G2: &Point{1.2, 3.4},
		},
	}.fail(t)
}

func TestLineMarshal(t *testing.T) {
	// Pass
	marshalTestcases{
		{
			Input:    &Line{{1.2, 3.4}, {5.6, 7.8}},
			Expected: `{"type":"LineString","coordinates":[[1.2,3.4],[5.6,7.8]]}`,
		},
	}.pass(t)
}

func TestLineScan(t *testing.T) {
	// Pass
	for _, c := range []struct {
		Input    interface{}
		Expected Geometry
	}{
		{
			Input:    []byte(`LINESTRING(0 0, 1 1)`),
			Expected: &Line{{0, 0}, {1, 1}},
		},
		{
			Input:    `LINESTRING(0 2, 0 3, -1.12654 5.985)`,
			Expected: &Line{{0, 2}, {0, 3}, {-1.12654, 5.985}},
		},
	} {
		l := &Line{}
		if err := l.Scan(c.Input); err != nil {
			t.Fatalf("could not scan %v: %s", c.Input, err)
		}
		if !l.Compare(c.Expected) {
			t.Fatalf("expected %v, got %v", c.Expected, l)
		}
	}

	// Fail
	for _, c := range []interface{}{
		4,        // wrong type
		`LINE()`, // wrong prefix
		[]byte(`LINESTRING((3 3, 2 2))`), // too many parentheses
		`LINESTRING(0, 0, 1, 1)`,         // should be spaces in between coordinates
	} {
		l := &Line{}
		if err := l.Scan(c); err == nil {
			t.Fatalf("expected error, got nil for %v", c)
		}
	}
}

func TestLineString(t *testing.T) {
	// Pass
	for _, testcase := range []struct {
		Input    Line
		Expected string
	}{
		{
			Input:    Line{},
			Expected: `LINESTRING EMPTY`,
		},
		{
			Input: Line{
				{1.2, 3.4},
				{5.6, 7.8},
				{5.8, 1.6},
			},
			Expected: `LINESTRING(1.2 3.4, 5.6 7.8, 5.8 1.6)`,
		},
	} {
		if expected, got := testcase.Expected, testcase.Input.String(); expected != got {
			t.Fatalf("expected %s, got %s", expected, got)
		}
	}
}
