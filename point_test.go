package geo

import "testing"

func TestPointCompare(t *testing.T) {
	// Different
	compareTestcases{
		{
			G1: &Point{1.2, 3.4},
			G2: &Point{1.2, 3.7},
		},
		{
			G1: &Point{1.2, 3.4},
			G2: &Point{9.2, 3.7},
		},
		{
			G1: &Point{1.2, 3.4},
			G2: &Line{{9.2, 3.7}},
		},
		{
			G1: &Point{1.2, 3.4},
			G2: &Polygon{{{9.2, 3.7}}},
		},
	}.fail(t)
}

func TestPointMarshal(t *testing.T) {
	marshalTestcases{
		{
			Input:    &Point{1.2, 3.4},
			Expected: `{"type":"Point","coordinates":[1.2,3.4]}`,
		},
	}.pass(t)
}

func TestPointScan(t *testing.T) {
	// Good
	for _, testcase := range []struct {
		WKT      interface{}
		Expected Point
	}{
		{
			WKT:      "POINT(1.2 3.4)",
			Expected: Point{1.2, 3.4},
		},
	} {
		p := &Point{}
		if err := p.Scan(testcase.WKT); err != nil {
			t.Fatal(err)
		}
		if expected, got := testcase.Expected[0], p[0]; expected != got {
			t.Fatalf("expected %f, got %f", expected, got)
		}
		if expected, got := testcase.Expected[1], p[1]; expected != got {
			t.Fatalf("expected %f, got %f", expected, got)
		}
	}

	// Bad
	for _, testcase := range []interface{}{
		"POINT(1.2, 3.4)",        // bad comma
		[]byte("PIONT(1.4 2.3)"), // typo
		7, // bad type
	} {
		p := &Point{}
		if err := p.Scan(testcase); err == nil {
			t.Fatal("expected error, got nil")
		}
	}
}

func TestPointValue(t *testing.T) {
	var (
		p        = &Point{1.2, 3.4}
		expected = `POINT(1.2 3.4)`
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
