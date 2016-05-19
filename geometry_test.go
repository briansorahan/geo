package geo

import (
	"encoding/json"
	"testing"
)

func TestGeometry(t *testing.T) {
	// Pass
	for _, testcase := range []struct {
		Input    geometry
		Expected Geometry
	}{
		{
			Input: geometry{
				Type:        LineType,
				Coordinates: json.RawMessage(`[[0,0],[1,1]]`),
			},
			Expected: &Line{
				{0, 0},
				{1, 1},
			},
		},
	} {
		g, err := testcase.Input.Geometry()
		if err != nil {
			t.Fatalf("failed to get geometry for %v: %s", testcase.Input, err)
		}
		if !testcase.Expected.Compare(g) {
			t.Fatalf("expected %v, got %v", testcase.Expected, g)
		}
	}
	// Fail
	for _, g := range []geometry{
		geometry{
			Type:        PointType,
			Coordinates: json.RawMessage(`{/}`),
		},
		geometry{
			Type:        LineType,
			Coordinates: json.RawMessage(`{/}`),
		},
		geometry{
			Type:        PolygonType,
			Coordinates: json.RawMessage(`{/}`),
		},
	} {
		if _, err := g.Geometry(); err == nil {
			t.Fatalf("expected error, got nil for %v", g)
		}
	}
}
