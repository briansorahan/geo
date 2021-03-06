package geo

import (
	"testing"
)

func TestCircleEqual(t *testing.T) {
	// Different
	cases{
		G: &Circle{Radius: 1, Coordinates: Point{0, 0}},
		Different: []Geometry{
			&Circle{Radius: 1, Coordinates: Point{0, 2}},
			&Circle{Radius: 3, Coordinates: Point{0, 0}},
			&Point{1, 1},
			&Line{{0, 0}, {1, 1}},
			&Polygon{{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}}},
		},
	}.test(t)
}

func TestCircleContains(t *testing.T) {
	// equirectangular
	CircleContainsMethod = ContainsMethodEquirectangular
	cases{
		G: &Circle{Radius: 100000, Coordinates: Point{0, 4}},
		Inside: []Point{
			{0.0001, 4},
		},
		Outside: []Point{
			{4, 4},
		},
	}.test(t)

	// spherical law of cosines
	CircleContainsMethod = ContainsMethodSphericalCosines
	cases{
		G: &Circle{Radius: 100000, Coordinates: Point{0, 4}},
		Inside: []Point{
			{0.0001, 4},
		},
		Outside: []Point{
			{4, 4},
		},
	}.test(t)

	// haversine
	CircleContainsMethod = ContainsMethodHaversine
	cases{
		G: &Circle{Radius: 100000, Coordinates: Point{0, 4}},
		Inside: []Point{
			{0.0001, 4},
		},
		Outside: []Point{
			{4, 4},
		},
	}.test(t)
}

func TestCircleContainsPanic(t *testing.T) {
	defer func() {
		if val := recover(); val == nil {
			t.Fatal("Circle Contains did not panic for unknown ContainsMethod")
		}
	}()

	var (
		c = &Circle{Radius: 100000, Coordinates: Point{0, 4}}
		p = Point{0.0001, 4}
	)
	CircleContainsMethod = "foo"
	c.Contains(p)
}

func TestCircleMarshalJSON(t *testing.T) {
	marshalTestcases{
		{
			Input:    &Circle{Radius: 1.23, Coordinates: Point{0, 0}},
			Expected: `{"type":"Circle","radius":1.23,"coordinates":[0,0]}`,
		},
	}.pass(t)
}

func TestCircleScan(t *testing.T) {
	// Pass
	for _, testcase := range []struct {
		WKT      interface{}
		Expected *Circle
	}{
		{
			WKT:      "CIRCULARSTRING(1 0, 0 1, -1 0, 0 -1, 1 0)",
			Expected: &Circle{Radius: 1, Coordinates: Point{0, 0}},
		},
	} {
		c := &Circle{}
		if err := c.Scan(testcase.WKT); err != nil {
			t.Fatal(err)
		}
		if !c.Equal(testcase.Expected) {
			t.Fatalf("expected %v, got %v", c, testcase.Expected)
		}
	}

	// Fail
	for _, testcase := range []interface{}{
		"CIRCULARSTRING(1, 1, 2, 3, 4, 5)",              // bad comma
		[]byte("CIRFCRESOT(1.4 2.3, 3.6 8.2, 4.6 0.2)"), // typo
		[]byte("CIRCULARSTRING(1 0, 0 1)"),              // < 3 points
		7, // bad type
	} {
		c := &Circle{}
		if err := c.Scan(testcase); err == nil {
			t.Fatal("expected error, got nil")
		}
	}
}

func TestCircleUnmarshalJSON(t *testing.T) {
	// Pass
	unmarshalTestcases{
		{
			Instance: &Circle{},
			Input:    []byte(`{"type":"Circle","coordinates":[2,2],"radius":1.8}`),
			Expected: &Circle{Radius: 1.8, Coordinates: Point{2, 2}},
		},
	}.pass(t)

	// Fail
	unmarshalTestcases{
		{
			// Bad JSON
			Instance: &Circle{},
			Input:    []byte(`ype":"Circle","coordinates":[2,2],"radiu1."8}`),
		},
		{
			// Bad type
			Instance: &Circle{},
			Input:    []byte(`{"type":"Circulo","coordinates":[2,2],"radius":1.8}`),
		},
		{
			// Bad coordinates
			Instance: &Circle{},
			Input:    []byte(`{"type":"Circle","coordinates":[[2,2]],"radius":1.8}`),
		},
	}.fail(t)
}

func TestCircleValue(t *testing.T) {
	var (
		c        = &Circle{Radius: 2, Coordinates: Point{0, 4}}
		expected = `CIRCULARSTRING(2 4, 0 6, -2 4, 0 2, 2 4)`
	)
	value, err := c.Value()
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
