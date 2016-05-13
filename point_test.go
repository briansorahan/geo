package geo

import "testing"

func TestPointCompare(t *testing.T) {
	// Different
	for _, testcase := range []struct {
		P1 Point
		P2 Point
	}{
		{
			P1: Point{1.2, 3.4},
			P2: Point{1.2, 3.7},
		},
		{
			P1: Point{1.2, 3.4},
			P2: Point{9.2, 3.7},
		},
	} {

		if same := testcase.P1.Compare(testcase.P2); same {
			t.Fatalf("expected %s, got %s", testcase.P1.String(), testcase.P2.String())
		}
	}
}

func TestPointMarshal(t *testing.T) {
	for _, testcase := range []struct {
		P        Point
		Expected string
	}{
		{
			P:        Point{1.2, 3.4},
			Expected: `{"type":"Point","coordinates":[1.2,3.4]}`,
		},
	} {
		got, err := testcase.P.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != testcase.Expected {
			t.Fatalf("expected %s, got %s", testcase.Expected, string(got))
		}
	}
}

func TestPointUnmarshal(t *testing.T) {
	// Pass
	for _, testcase := range []struct {
		Input    string
		Expected Point
	}{
		{
			Input:    `{"type":"Point","coordinates":[1.2,3.4]}`,
			Expected: Point{1.2, 3.4},
		},
	} {
		p := &Point{}
		if err := p.UnmarshalJSON([]byte(testcase.Input)); err != nil {
			t.Fatal(err)
		}
		if !p.Compare(testcase.Expected) {
			t.Fatalf("expected %s, got %s", testcase.Expected.String(), p.String())
		}
	}
	// Fail
	for _, testcase := range []struct {
		Input    string
		Expected Point
	}{
		{
			Input:    `{"type":"Pont","coordinates":[1.2,3.4]}`,
			Expected: Point{1.2, 3.4},
		},
		{
			Input:    `{"type":"Point","coordinates":[1.2]}`,
			Expected: Point{1.2, 3.4},
		},
		{
			Input:    `{"type":"Point","coordinates":[1.2,3.4}}`,
			Expected: Point{1.2, 3.4},
		},
		{
			Input:    `{"type":"Point","coordinates":[abc,3.4]}`,
			Expected: Point{1.2, 3.4},
		},
		{
			Input:    `{"type":"Point","coordinates":[1.2,abc]}`,
			Expected: Point{1.2, 3.4},
		},
	} {
		p := &Point{}
		if err := p.UnmarshalJSON([]byte(testcase.Input)); err == nil {
			t.Fatal("expected error, but got nil")
		}
	}
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
