package geo

import "testing"

func TestPolygonMarshal(t *testing.T) {
	p := &Polygon{
		{1.2, 3.4},
		{5.6, 7.8},
	}
	expected := `{"type":"Polygon","coordinates":[[1.2,3.4],[5.6,7.8]]}`
	got, err := p.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != expected {
		t.Fatalf("expected %s, got %s", expected, string(got))
	}
}

func TestPolygonScan(t *testing.T) {
	var (
		p    = &Polygon{}
		good = "POLYGON(1.2 3.4, 5.6 7.8, 6.2 1.5, 1.2 3.4)"
		bad1 = "POLYGON(1.2, 3.4, 5.6, 7.8)"
		bad2 = "POLYGON(1.2, 3.4, 5.6, 7.8"
		bad3 = "PIKACHU"
	)
	// good scan
	if err := p.Scan(good); err != nil {
		t.Fatal(err)
	}
	for i, coord := range [][2]float64{
		{1.2, 3.4},
		{5.6, 7.8},
		{6.2, 1.5},
		{1.2, 3.4},
	} {
		if expected, got := coord[0], (*p)[i][0]; expected != got {
			t.Fatalf("expected %f, got %f", expected, got)
		}
		if expected, got := coord[1], (*p)[i][1]; expected != got {
			t.Fatalf("expected %f, got %f", expected, got)
		}
	}
	// bad scan
	if err := p.Scan(bad1); err == nil {
		t.Fatalf("expected err, got nil")
	}
	// bad scan with bytes
	if err := p.Scan([]byte(bad1)); err == nil {
		t.Fatal("expected err, got nil")
	}
	// scan with bad type
	if err := p.Scan(7); err == nil {
		t.Fatal("expected error, got nil")
	}
	// scan with bad strings
	if err := p.Scan(bad2); err == nil {
		t.Fatal("expected err, got nil")
	}
	if err := p.Scan(bad3); err == nil {
		t.Fatal("expected err, got nil")
	}
}

func TestPolygonValue(t *testing.T) {
	var (
		p = &Polygon{
			{1.2, 3.4},
			{5.6, 7.8},
			{8.7, 6.5},
			{4.3, 2.1},
		}
		expected = `POLYGON(1.2 3.4, 5.6 7.8, 8.7 6.5, 4.3 2.1)`
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

func TestPolygonEmpty(t *testing.T) {
	var (
		p        = &Polygon{}
		expected = "POLYGON EMPTY"
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
