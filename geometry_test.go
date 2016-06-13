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
		g, err := testcase.Input.unmarshalCoordinates()
		if err != nil {
			t.Fatalf("failed to get geometry for %v: %s", testcase.Input, err)
		}
		if !testcase.Expected.Compare(g) {
			t.Fatalf("expected %v, got %v", testcase.Expected, g)
		}
	}
	// Fail
	for _, g := range []geometry{
		{
			Type:        PointType,
			Coordinates: json.RawMessage(`{/}`),
		},
		{
			Type:        LineType,
			Coordinates: json.RawMessage(`{/}`),
		},
		{
			Type:        PolygonType,
			Coordinates: json.RawMessage(`{/}`),
		},
	} {
		if _, err := g.unmarshalCoordinates(); err == nil {
			t.Fatalf("expected error, got nil for %v", g)
		}
	}
}

func TestUnmarshalGeometry(t *testing.T) {
	// Pass
	for i, c := range []struct {
		Input    []byte
		Expected Geometry
	}{
		{
			Input:    []byte(`{"type":"Point","coordinates":[0,0]}`),
			Expected: &Point{0, 0},
		},
	} {
		g, err := UnmarshalGeometry(c.Input)
		if err != nil {
			t.Fatalf("(case %d) failed to unmarshal geometry %s", i, err)
		}
		if !g.Compare(c.Expected) {
			t.Fatalf("(case %d) expected %s to be the same as %s", i, g, c.Expected)
		}
	}

	// Fail
	for i, data := range [][]byte{
		[]byte(`{&%%$$`),
	} {
		if _, err := UnmarshalGeometry(data); err == nil {
			t.Fatalf("(case %d) expected error, got nil", i)
		}
	}
}
