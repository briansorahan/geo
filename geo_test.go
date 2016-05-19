package geo

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"testing"
)

// compareTestcases is a helper type for Compare tests.
type compareTestcases []struct {
	G1 Geometry
	G2 Geometry
}

// pass runs the test cases that should pass.
func (tests compareTestcases) pass(t *testing.T) {
	for _, c := range tests {
		if same := c.G1.Compare(c.G2); !same {
			t.Fatalf("expected %s and %s to be the same", c.G1.String(), c.G2.String())
		}
	}
}

// fail runs the test cases that should fail.
func (tests compareTestcases) fail(t *testing.T) {
	for _, c := range tests {
		if same := c.G1.Compare(c.G2); same {
			t.Fatalf("expected %s to not equal %s", c.G1.String(), c.G2.String())
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
	for _, testcase := range tests {
		got, err := testcase.Input.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != testcase.Expected {
			t.Fatalf("expected %s, got %s", testcase.Expected, string(got))
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
	for _, c := range tests {
		if err := json.Unmarshal(c.Input, c.Instance); err != nil {
			t.Fatal(err)
		}
		if !c.Instance.Compare(c.Expected) {
			t.Fatalf("expected %s to equal %s", c.Instance.String(), c.Expected.String())
		}
	}
}

// fail runs the test cases that should fail.
func (tests unmarshalTestcases) fail(t *testing.T) {
	for _, c := range tests {
		if err := json.Unmarshal(c.Input, c.Instance); err == nil {
			t.Fatal("expected error, got nil")
		}
	}
}

// badGeom is a type that always returns an error from MarshalJSON and UnmarshalJSON.
type badGeom struct{}

// Compare always returns false
func (badgeom badGeom) Compare(g Geometry) bool {
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
