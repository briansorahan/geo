package geo

import (
	"database/sql/driver"
	"testing"
)

func TestLineCompare(t *testing.T) {
	cases{
		G: &Line{{1.2, 3.4}, {5.6, 7.8}, {1.4, 9.3}, {-1.7, 7.3}},
		Same: []Geometry{
			&Line{{1.2, 3.4}, {5.6, 7.8}, {1.4, 9.3}, {-1.7, 7.3}},
		},
		Different: []Geometry{
			&Line{{1.2, 3.4}, {5.6, 7.8}, {1.4, 9.3}},
			&Line{{1.2, 3.4}, {5.6, 7.8}, {1.4, 9.3}, {-1.4, 7.3}},
			&Line{{1.2, 3.4}, {5.6, 7.8}, {1.4, 9.3}, {-1.7, 7.5}},
			&Polygon{{{1.2, 3.4}, {5.6, 7.8}, {1.4, 9.3}, {-1.7, 7.5}}},
			&Point{1.2, 3.4},
		},
	}.test(t)
}

func TestLineContains(t *testing.T) {
	cases{
		G: &Line{{0, 0}, {4, 4}},
		Inside: []Point{
			Point{0, 0},
			Point{2, 2},
			Point{4, 4},
		},
		Outside: []Point{
			Point{-1, 0},
			Point{1, 3},
			Point{4, 5},
		},
	}.test(t)

	cases{
		G:       &Line{{4, 4}},
		Outside: []Point{{-1, 0}},
	}.test(t)
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

func TestLineUnmarshal(t *testing.T) {
	unmarshalTestcases{
		{
			Input:    []byte(`{"type":"LineString","coordinates":[[0,0],[1,1]]}`),
			Instance: &Line{},
			Expected: &Line{{0, 0}, {1, 1}},
		},
	}.pass(t)

	unmarshalTestcases{
		{
			// Bad type
			Input:    []byte(`{"type":"Line","coordinates":[[0,0],[1,1]]}`),
			Instance: &Line{},
		},
		{
			// Bad coordinates
			Input:    []byte(`{"type":"LineString","coordinates":"foo"}`),
			Instance: &Line{},
		},
	}.fail(t)
}

func TestLineValue(t *testing.T) {
	for _, c := range []struct {
		Input    Line
		Expected driver.Value
	}{
		{
			Input:    Line{{0, 0}, {1, 1}},
			Expected: driver.Value(`LINESTRING(0 0, 1 1)`),
		},
	} {
		val, err := c.Input.Value()
		if err != nil {
			t.Fatal(err)
		}
		if val != c.Expected {
			t.Fatalf("expected %v, got %v", c.Expected, val)
		}
	}
}
