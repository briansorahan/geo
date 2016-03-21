package geo

import "testing"

func TestPointMarshal(t *testing.T) {
	p := &Point{1.2, 3.4}
	expected := `{"type":"Point","coordinates":[1.2,3.4]}`
	got, err := p.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != expected {
		t.Fatalf("expected %s, got %s", expected, string(got))
	}
}

func TestPointScan(t *testing.T) {
	var (
		p    = &Point{}
		good = "POINT(1.2 3.4)"
		bad  = "POINT(1.2, 3.4)"
	)
	// good scan
	if err := p.Scan(good); err != nil {
		t.Fatal(err)
	}
	if expected, got := 1.2, p[0]; expected != got {
		t.Fatalf("expected %f, got %f", expected, got)
	}
	if expected, got := 3.4, p[1]; expected != got {
		t.Fatalf("expected %f, got %f", expected, got)
	}
	// bad scan
	if err := p.Scan(bad); err == nil {
		t.Fatalf("expected err, got nil")
	}
	// bad scan with bytes
	if err := p.Scan([]byte(bad)); err == nil {
		t.Fatal("expected err, got nil")
	}
	// scan with bad type
	if err := p.Scan(7); err == nil {
		t.Fatal("expected error, got nil")
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
