package geo

import "testing"

func TestLineString(t *testing.T) {
	expected := `{"type":"Polygon","coordinates":[[1.2,3.4],[5.6,7.8]]}`
	got, err := MarshalJSON(Polygon{
		{1.2, 3.4},
		{5.6, 7.8},
	})
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != expected {
		t.Fatalf("expected %s, got %s", expected, string(got))
	}
}

func TestInvalidLineString(t *testing.T) {
	if _, err := MarshalJSON(Polygon{}); err == nil {
		t.Fatal("expected error for invalid LineString")
	}
}
