package geo

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"testing"
)

type cases struct {
	G Geometry

	// Equal tests.
	Same      []Geometry
	Different []Geometry

	// Contains tests.
	Inside  []Point
	Outside []Point
}

// pass runs the test cases that should pass.
func (c cases) test(t *testing.T) {
	for i, same := range c.Same {
		if ok := c.G.Equal(same); !ok {
			t.Fatalf("(case %d) expected %s and %s to be the same", i, c.G.String(), same.String())
		}
	}
	for i, diff := range c.Different {
		if ok := c.G.Equal(diff); ok {
			t.Fatalf("(case %d) expected %s to not equal %s", i, c.G.String(), diff.String())
		}
	}
	for i, p := range c.Inside {
		if ok := c.G.Contains(p); !ok {
			t.Fatalf("(case %d) expected %s to contain %s", i, c.G.String(), p.String())
		}
	}
	for i, p := range c.Outside {
		if ok := c.G.Contains(p); ok {
			t.Fatalf("(case %d) expected %s to not contain %s", i, c.G.String(), p.String())
		}
	}
}

// marshalTestcases is a helper type for MarshalJSON tests.
type marshalTestcases []struct {
	Input    Geometry
	Expected string
}

// pass runs the test cases that should pass.
func (tests marshalTestcases) pass(t *testing.T) {
	for i, testcase := range tests {
		got, err := testcase.Input.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != testcase.Expected {
			t.Fatalf("(case %d) expected %s, got %s", i, testcase.Expected, string(got))
		}
	}
}

// scanTestcases is a helper type for Scan tests.
type scanTestcases []struct {
	Input    interface{}
	Instance Geometry
	Expected Geometry
}

// fail runs test cases that should fail.
func (tests scanTestcases) fail(t *testing.T) {
	for i, c := range tests {
		if err := c.Instance.Scan(c.Input); err == nil {
			t.Fatalf("(case %d) expected error, got nil", i)
		}
	}
}

// pass runs test cases that should pass.
func (tests scanTestcases) pass(t *testing.T) {
	for i, c := range tests {
		if err := c.Instance.Scan(c.Input); err != nil {
			t.Fatalf("(case %d) failed to scan: %s", i, err)
		}
	}
}

// unmarshalTestcases is a helper type for UnmarshalJSON tests.
type unmarshalTestcases []struct {
	Input    []byte
	Expected Geometry
	Instance Geometry
}

// pass runs the test cases that should pass.
func (tests unmarshalTestcases) pass(t *testing.T) {
	for i, c := range tests {
		if err := json.Unmarshal(c.Input, c.Instance); err != nil {
			t.Fatal(err)
		}
		if !c.Instance.Equal(c.Expected) {
			t.Fatalf("(unmarshal %T pass case %d) expected %s to equal %s", c.Instance, i, c.Instance.String(), c.Expected.String())
		}
	}
}

// fail runs the test cases that should fail.
func (tests unmarshalTestcases) fail(t *testing.T) {
	for i, c := range tests {
		if err := json.Unmarshal(c.Input, c.Instance); err == nil {
			t.Fatalf("(unmarshal %T fail case %d) expected error, got nil", c.Instance, i)
		}
	}
}

// valueTestcases is a helper type for MarshalJSON tests.
type valueTestcases []struct {
	Input    Geometry
	Expected interface{}
}

// pass runs the test cases that should pass.
func (tests valueTestcases) pass(t *testing.T) {
	for i, testcase := range tests {
		got, err := testcase.Input.Value()
		if err != nil {
			t.Fatal(err)
		}
		if got != testcase.Expected {
			t.Fatalf("(case %d) expected %v, got %v", i, testcase.Expected, got)
		}
	}
}

// badGeom is a type that always returns an error from MarshalJSON and UnmarshalJSON.
type badGeom struct{}

// Equal always returns false
func (badgeom badGeom) Equal(g Geometry) bool {
	return false
}

// Contains always returns false.
func (badgeom badGeom) Contains(p Point) bool {
	return false
}

// MarshalJSON always returns an error.
func (badgeom badGeom) MarshalJSON() ([]byte, error) {
	return nil, errors.New("bad geom")
}

// Scan always returns an error.
func (badgeom badGeom) Scan(src interface{}) error {
	return errors.New("bad geom")
}

// String always returns "badgeom"
func (badgeom badGeom) String() string {
	return "badgeom"
}

// UnmarshalJSON always returns an error.
func (badgeom badGeom) UnmarshalJSON(data []byte) error {
	return errors.New("bad geom")
}

// Value always returns an error.
func (badgeom badGeom) Value() (driver.Value, error) {
	return nil, errors.New("bad geom")
}

// Transform transforms the geometry point by point.
func (badgeom badGeom) Transform(t Transformer) {
}

// VisitCoordinates visits each point in the geometry.
func (badgeom badGeom) VisitCoordinates(v Visitor) {
}
