package geo

import (
	"errors"
	"testing"
)

// BadGeom is a type that always returns an error from MarshalJSON and UnmarshalJSON.
type BadGeom struct{}

// MarshalJSON always returns an error.
func (badgeom BadGeom) MarshalJSON() ([]byte, error) {
	return nil, errors.New("bad geom")
}

// UnmarshalJSON always returns an error.
func (badgeom BadGeom) UnmarshalJSON(data []byte) error {
	return errors.New("bad geom")
}

func TestFeatureMarshal(t *testing.T) {
	// Pass
	for _, testcase := range []struct {
		Feature  Feature
		Expected string
	}{
		{
			Feature: Feature{
				Geometry: &Point{1, 2},
			},
			Expected: `{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2]},"properties":null}`,
		},
	} {
		got, err := testcase.Feature.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != testcase.Expected {
			t.Fatalf("expected %s, got %s", testcase.Expected, string(got))
		}
	}
	// Fail
	for _, testcase := range []struct {
		Feature  Feature
		Expected string
	}{
		{Feature: Feature{Geometry: &BadGeom{}}},
		{Feature: Feature{Geometry: &Point{1, 2}, Properties: &BadGeom{}}},
	} {

		if _, err := testcase.Feature.MarshalJSON(); err == nil {
			t.Fatal("expected error, got nil")
		}
	}
}
