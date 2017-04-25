package geo

import "testing"

func TestVisit(t *testing.T) {
	for i, testcase := range []struct {
		g Geometry
		q quadrants
	}{
		{
			g: &Point{1, 1},
			q: quadrants{UR: 1},
		},
		{
			g: &Line{{-1, -1}, {1, 1}},
			q: quadrants{UR: 1, LL: 1},
		},
		{
			g: &Polygon{{{0, 0}, {0, 1}}, {{0, 1}, {1, 1}}, {{1, 1}, {1, 0}}, {{1, 0}, {0, 0}}},
			q: quadrants{UR: 8},
		},
		{
			g: &MultiPoint{{0, 0}, {1, 1}},
			q: quadrants{UR: 2},
		},
		{
			g: &MultiLine{{{0, 0}, {0, 1}}, {{0, 1}, {1, 1}}, {{1, 1}, {1, 0}}, {{1, 0}, {0, 0}}},
			q: quadrants{UR: 8},
		},
		{
			g: &MultiPolygon{
				{{{0, 0}, {0, 1}}, {{0, 1}, {1, 1}}, {{1, 1}, {1, 0}}, {{1, 0}, {0, 0}}},
				{{{3, -3}, {6, -3}}, {{6, -3}, {6, -7}}, {{6, -7}, {3, -3}}},
			},
			q: quadrants{UR: 8, LR: 6},
		},
		{
			g: &Feature{
				Geometry: &Line{{0, 0}, {1, 1}},
			},
			q: quadrants{UR: 2},
		},
		{
			g: &FeatureCollection{
				{
					Geometry: &Line{{0, 0}, {1, 1}},
				},
				{
					Geometry: &Polygon{{{0, 0}, {0, 1}}, {{0, 1}, {1, 1}}, {{1, 1}, {1, 0}}, {{1, 0}, {0, 0}}},
				},
			},
			q: quadrants{UR: 10},
		},
		{
			g: &GeometryCollection{
				&Line{{0, 0}, {1, 1}},
				&Polygon{{{0, 0}, {0, 1}}, {{0, 1}, {1, 1}}, {{1, 1}, {1, 0}}, {{1, 0}, {0, 0}}},
			},
			q: quadrants{UR: 10},
		},
	} {
		q := &quadrants{}
		testcase.g.Visit(q)
		if expected, got := testcase.q, q; !got.Equal(expected) {
			t.Fatalf("(test case %d) expected %#v to equal %#v", i, expected, *got)
		}
	}
}

// quadrants is a Visitor that determines how many of the points of a geometry
// lie in each quadrant of the cartesian plane.
type quadrants struct {
	UL int
	UR int
	LR int
	LL int
}

func (q *quadrants) Equal(other quadrants) bool {
	return q.UL == other.UL && q.UR == other.UR && q.LR == other.LR && q.LL == other.LL
}

func (q *quadrants) Visit(p Point) {
	switch x := p[0]; {
	case x < 0:
		switch y := p[1]; {
		case y < 0:
			q.LL++
		case y >= 0:
			q.UL++
		}
	case x >= 0:
		switch y := p[1]; {
		case y < 0:
			q.LR++
		case y >= 0:
			q.UR++
		}
	}
}
