package geo

import (
	"database/sql"
	"log"
	"testing"
)

func TestPointMarshal(t *testing.T) {
	for _, testcase := range []struct {
		P        Point
		Expected string
	}{
		{
			P:        Point{1.2, 3.4},
			Expected: `{"type":"Point","coordinates":[1.2,3.4]}`,
		},
	} {
		got, err := testcase.P.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != testcase.Expected {
			t.Fatalf("expected %s, got %s", testcase.Expected, string(got))
		}
	}
}

func TestPointScan(t *testing.T) {
	// Good
	for _, testcase := range []struct {
		WKT      interface{}
		Expected Point
	}{
		{
			WKT:      "POINT(1.2 3.4)",
			Expected: Point{1.2, 3.4},
		},
	} {
		p := &Point{}
		if err := p.Scan(testcase.WKT); err != nil {
			t.Fatal(err)
		}
		if expected, got := testcase.Expected[0], p[0]; expected != got {
			t.Fatalf("expected %f, got %f", expected, got)
		}
		if expected, got := testcase.Expected[1], p[1]; expected != got {
			t.Fatalf("expected %f, got %f", expected, got)
		}
	}

	// Bad
	for _, testcase := range []interface{}{
		"POINT(1.2, 3.4)",        // bad comma
		[]byte("PIONT(1.4 2.3)"), // typo
		7, // bad type
	} {
		p := &Point{}
		if err := p.Scan(testcase); err == nil {
			t.Fatal("expected error, got nil")
		}
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

func ExamplePoint() {
	// Insert a point into a table called "locations".
	// The location column is of type GEOMETRY(POINT).
	const InsertLocation = `
INSERT INTO locations (location)
VALUES                (ST_GeomFromText($1))
`
	db, err := sql.Open("postgres", "datasource")
	if err != nil {
		log.Fatal(err)
	}
	result, err := db.Exec(InsertLocation, &Point{1.2, 3.4})
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	if rowsAffected != 1 {
		log.Fatalf("Expected rowsAffected to be 1, got %d", rowsAffected)
	}

	// Fetch the point we just created.
	const GetHeartbeats = `
SELECT ST_AsText(location) AS location
FROM   locations
`
	rows, err := db.Query(GetHeartbeats)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		p := &Point{}
		if err := rows.Scan(p); err != nil {
			log.Fatal(err)
		}
		log.Printf("got location %s\n", p)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
