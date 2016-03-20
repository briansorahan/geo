package geo

import "testing"

func TestPoint(t *testing.T) {
	expected := `{"type":"Point","coordinates":[1.2,3.4]}`
	got, err := MarshalJSON(Point{1.2, 3.4})
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != expected {
		t.Fatalf("expected %s, got %s", expected, string(got))
	}
}
