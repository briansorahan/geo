package geo

import "testing"

func TestTransform(t *testing.T) {
	for i, testcase := range []struct {
		in  Geometry
		out Geometry
		t   Transformer
	}{
		{
			in:  &Point{0, 0},
			out: &Point{1, 1},
			t:   pointShifter(1),
		},
		{
			in:  &Line{{0, 0}, {1, 1}},
			out: &Line{{1, 1}, {2, 2}},
			t:   pointShifter(1),
		},
		{
			in:  &Polygon{{{0, 0}, {0, 1}}, {{0, 1}, {1, 1}}, {{1, 1}, {1, 0}}, {{1, 0}, {0, 0}}},
			out: &Polygon{{{1, 1}, {1, 2}}, {{1, 2}, {2, 2}}, {{2, 2}, {2, 1}}, {{2, 1}, {1, 1}}},
			t:   pointShifter(1),
		},
		{
			in:  &MultiPoint{{0, 0}, {1, 1}},
			out: &MultiPoint{{1, 1}, {2, 2}},
			t:   pointShifter(1),
		},
		{
			in:  &MultiLine{{{0, 0}, {0, 1}}, {{0, 1}, {1, 1}}, {{1, 1}, {1, 0}}, {{1, 0}, {0, 0}}},
			out: &MultiLine{{{1, 1}, {1, 2}}, {{1, 2}, {2, 2}}, {{2, 2}, {2, 1}}, {{2, 1}, {1, 1}}},
			t:   pointShifter(1),
		},
		{
			in: &MultiPolygon{
				{{{0, 0}, {0, 1}}, {{0, 1}, {1, 1}}, {{1, 1}, {1, 0}}, {{1, 0}, {0, 0}}},
				{{{3, 3}, {6, 3}}, {{6, 3}}},
			},
			out: &MultiPolygon{
				{{{1, 1}, {1, 2}}, {{1, 2}, {2, 2}}, {{2, 2}, {2, 1}}, {{2, 1}, {1, 1}}},
				{{{4, 4}, {7, 4}}, {{7, 4}}},
			},
			t: pointShifter(1),
		},
		{
			in: &Feature{
				Geometry: &Line{{0, 0}, {1, 1}},
			},
			out: &Feature{
				Geometry: &Line{{1, 1}, {2, 2}},
			},
			t: pointShifter(1),
		},
		{
			in: &FeatureCollection{
				{
					Geometry: &Line{{0, 0}, {1, 1}},
				},
				{
					Geometry: &Polygon{{{0, 0}, {0, 1}}, {{0, 1}, {1, 1}}, {{1, 1}, {1, 0}}, {{1, 0}, {0, 0}}},
				},
			},
			out: &FeatureCollection{
				{
					Geometry: &Line{{1, 1}, {2, 2}},
				},
				{
					Geometry: &Polygon{{{1, 1}, {1, 2}}, {{1, 2}, {2, 2}}, {{2, 2}, {2, 1}}, {{2, 1}, {1, 1}}},
				},
			},
			t: pointShifter(1),
		},
		{
			in: &GeometryCollection{
				&Line{{0, 0}, {1, 1}},
				&Polygon{{{0, 0}, {0, 1}}, {{0, 1}, {1, 1}}, {{1, 1}, {1, 0}}, {{1, 0}, {0, 0}}},
			},
			out: &GeometryCollection{
				&Line{{1, 1}, {2, 2}},
				&Polygon{{{1, 1}, {1, 2}}, {{1, 2}, {2, 2}}, {{2, 2}, {2, 1}}, {{2, 1}, {1, 1}}},
			},
			t: pointShifter(1),
		},
	} {
		testcase.in.Transform(testcase.t)
		if !testcase.out.Equal(testcase.in) {
			t.Fatalf("(test case %d) expected %#v to equal %#v", i, testcase.in, testcase.out)
		}
	}
}

type pointShifter float64

func (ps pointShifter) Transform(p Point) Point {
	return Point{
		p[0] + float64(ps),
		p[1] + float64(ps),
	}
}
