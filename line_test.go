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
			G2: &Polygon{{1.2, 3.4}, {5.6, 7.8}, {1.4, 9.3}, {-1.7, 7.5}},
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

func TestLineUnmarshal(t *testing.T) {
	// Pass
	unmarshalTestcases{
		{
			Input: `{"type":"LineString","coordinates":[[1.2,3.4],[5.6,7.8],[5.8,1.6]]}`,
			Expected: &Line{
				{1.2, 3.4},
				{5.6, 7.8},
				{5.8, 1.6},
			},
			Instance: &Line{},
		},
	}.pass(t)

	// Fail
	for _, testcase := range []string{
		`{"type":"LineStirng","coordinates":[[1.2,3.4],[5.6,7.8],[5.8,1.6]]}`,
		`{"type":"LineString","coordinates":[[1.2,3.4],[5.6,7.8],[5.8,1.6}}}`,
		`{"type":"LineString","coordinates":[[1.2,3.4],[5.6,7.8>>>`,
		`{"type":"LineString","coordinates":[[1.2,3.4],[5.6,7.8],[5.8]]}`,
		`{"type":"LineString","coordinates":[[1.2,3.4],[5.6,7.8],[abc,-7.4]]}`,
		`{"type":"LineString","coordinates":[[1.2,3.4],[5.6,7.8],[-7.4,abc]]}`,
	} {
		p := &Line{}
		if err := p.UnmarshalJSON([]byte(testcase)); err == nil {
			t.Fatal("expected error, got nil")
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
