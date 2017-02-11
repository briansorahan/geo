package geo

import (
	"database/sql/driver"
	"testing"
)

func TestMultiPointCompare(t *testing.T) {
	cases{
		G: &MultiPoint{{1.2, 3.4}, {5.6, 7.8}, {1.4, 9.3}, {-1.7, 7.3}},
		Same: []Geometry{
			&MultiPoint{{1.2, 3.4}, {5.6, 7.8}, {1.4, 9.3}, {-1.7, 7.3}},
		},
		Different: []Geometry{
			&MultiPoint{{1.2, 3.4}, {5.6, 7.8}, {1.4, 9.3}},
			&MultiPoint{{1.2, 3.4}, {5.6, 7.8}, {1.4, 9.3}, {-1.4, 7.3}},
			&MultiPoint{{1.2, 3.4}, {5.6, 7.8}, {1.4, 9.3}, {-1.7, 7.5}},
			&Polygon{{{1.2, 3.4}, {5.6, 7.8}, {1.4, 9.3}, {-1.7, 7.5}}},
			&Point{1.2, 3.4},
		},
	}.test(t)
}

func TestMultiPointContains(t *testing.T) {
	cases{
		G: &MultiPoint{{0, 0}, {4, 4}},
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
		G:       &MultiPoint{{4, 4}},
		Outside: []Point{{-1, 0}},
	}.test(t)
}

func TestMultiPointMarshal(t *testing.T) {
	// Pass
	marshalTestcases{
		{
			Input:    &MultiPoint{{1.2, 3.4}, {5.6, 7.8}},
			Expected: `{"type":"MultiPoint","coordinates":[[1.2,3.4],[5.6,7.8]]}`,
		},
	}.pass(t)
}

func TestMultiPointScan(t *testing.T) {
	// Pass
	for _, c := range []struct {
		Input    interface{}
		Expected Geometry
	}{
		{
			Input:    []byte(`MULTIPOINT(0 0, 1 1)`),
			Expected: &MultiPoint{{0, 0}, {1, 1}},
		},
		{
			Input:    `MULTIPOINT(0 2, 0 3, -1.12654 5.985)`,
			Expected: &MultiPoint{{0, 2}, {0, 3}, {-1.12654, 5.985}},
		},
	} {
		l := &MultiPoint{}
		if err := l.Scan(c.Input); err != nil {
			t.Fatalf("could not scan %v: %s", c.Input, err)
		}
		if !l.Compare(c.Expected) {
			t.Fatalf("expected %v, got %v", c.Expected, l)
		}
	}

	// Fail
	for _, c := range []interface{}{
		4,                        // wrong type
		`LINE()`,                 // wrong prefix
		`MULTIPOINT(0, 0, 1, 1)`, // should be spaces in between coordinates
	} {
		l := &MultiPoint{}
		if err := l.Scan(c); err == nil {
			t.Fatalf("expected error, got nil for %v", c)
		}
	}
}

func TestMultiPointString(t *testing.T) {
	// Pass
	for _, testcase := range []struct {
		Input    MultiPoint
		Expected string
	}{
		{
			Input:    MultiPoint{},
			Expected: `MULTIPOINT EMPTY`,
		},
		{
			Input: MultiPoint{
				{1.2, 3.4},
				{5.6, 7.8},
				{5.8, 1.6},
			},
			Expected: `MULTIPOINT(1.2 3.4, 5.6 7.8, 5.8 1.6)`,
		},
	} {
		if expected, got := testcase.Expected, testcase.Input.String(); expected != got {
			t.Fatalf("expected %s, got %s", expected, got)
		}
	}
}

func TestMultiPointUnmarshal(t *testing.T) {
	unmarshalTestcases{
		{
			Input:    []byte(`{"type":"MultiPoint","coordinates":[[0,0],[1,1]]}`),
			Instance: &MultiPoint{},
			Expected: &MultiPoint{{0, 0}, {1, 1}},
		},
	}.pass(t)

	unmarshalTestcases{
		{
			// Bad type
			Input:    []byte(`{"type":"Line","coordinates":[[0,0],[1,1]]}`),
			Instance: &MultiPoint{},
		},
		{
			// Bad coordinates
			Input:    []byte(`{"type":"MultiPoint","coordinates":"foo"}`),
			Instance: &MultiPoint{},
		},
	}.fail(t)
}

func TestMultiPointValue(t *testing.T) {
	for _, c := range []struct {
		Input    MultiPoint
		Expected driver.Value
	}{
		{
			Input:    MultiPoint{{0, 0}, {1, 1}},
			Expected: driver.Value(`MULTIPOINT(0 0, 1 1)`),
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
