package geo

import "testing"

func TestPointEqual(t *testing.T) {
	// Different
	cases{
		G: &Point{1.2, 3.4},
		Different: []Geometry{
			&Point{1.2, 3.7},
			&Point{9.2, 3.7},
			&Line{{9.2, 3.7}},
			&Polygon{{{9.2, 3.7}}},
		},
	}.test(t)
}

func TestPointContains(t *testing.T) {
	cases{
		G:       &Point{1, 1},
		Inside:  []Point{{1, 1}},
		Outside: []Point{{1, 1.2}},
	}.test(t)
}

func TestPointMarshalJSON(t *testing.T) {
	marshalTestcases{
		{
			Input:    &Point{1.2, 3.4},
			Expected: `{"type":"Point","coordinates":[1.2,3.4]}`,
		},
	}.pass(t)
}

func TestPointScan(t *testing.T) {
	// Good
	scanTestcases{
		{
			Input:    `POINT(1.2 3.4)`,
			Instance: &Point{},
			Expected: &Point{1.2, 3.4},
		},
	}.pass(t)

	// Bad
	scanTestcases{
		{
			Input:    "POINT(1.2, 3.4)", // bad comma
			Instance: &Point{},
		},
		{
			Input:    []byte("PIONT(1.4 2.3)"), // typo
			Instance: &Point{},
		},
		{
			Input:    7, // bad type
			Instance: &Point{},
		},
	}.fail(t)
}

func TestPointUnmarshal(t *testing.T) {
	unmarshalTestcases{
		{
			Input:    []byte(`{"type":"Point","coordinates":[1,1]}`),
			Instance: &Point{},
			Expected: &Point{1, 1},
		},
	}.pass(t)

	unmarshalTestcases{
		{
			Input:    []byte(`{"type":"Punto","coordinates":[1,1]}`),
			Instance: &Point{},
		},
		{
			Input:    []byte(`{"type":"Point","coordinates":"bork"}`),
			Instance: &Point{},
		},
	}.fail(t)
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
