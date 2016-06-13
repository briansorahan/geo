package geo

import (
	"testing"

	kdgeo "github.com/kellydunn/golang-geo"
	"github.com/paulsmith/gogeos/geos"
)

func BenchmarkGeosPaulSmit(b *testing.B) {
	// Setup
	polygon, err := geos.NewPolygon([]geos.Coord{
		{0, 1, 0},
		{1, 2, 0},
		{2, 1, 0},
		{2, 0, 0},
		{1, -1, 0},
		{0, 0, 0},
		{0, 1, 0},
	})
	if err != nil {
		b.Fatal(err)
	}
	point, err := geos.NewPoint(geos.Coord{X: 1, Y: 0, Z: 0})
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	// Test
	for i := 0; i < b.N; i++ {
		contains, err := polygon.Contains(point)
		if err != nil {
			b.Fatal(err)
		}
		if !contains {
			b.Fatal("not contains")
		}
	}
}

func BenchmarkGeoKellyDunn(b *testing.B) {
	// Setup
	polygon := kdgeo.NewPolygon([]*kdgeo.Point{
		kdgeo.NewPoint(0, 1),
		kdgeo.NewPoint(1, 2),
		kdgeo.NewPoint(2, 1),
		kdgeo.NewPoint(2, 0),
		kdgeo.NewPoint(1, -1),
		kdgeo.NewPoint(0, 0),
		kdgeo.NewPoint(0, 1),
	})
	point := kdgeo.NewPoint(1, 0)

	b.ResetTimer()

	// Test
	for i := 0; i < b.N; i++ {
		if contains := polygon.Contains(point); !contains {
			b.Fatal("not contains")
		}
	}
}

// TODO: make gdal easy to use

// const WGS84 = `
// GEOGCS["WGS 84",
//     DATUM["WGS_1984",
//         SPHEROID["WGS 84",6378137,298.257223563,
//             AUTHORITY["EPSG","7030"]],
//         AUTHORITY["EPSG","6326"]],
//     PRIMEM["Greenwich",0,
//         AUTHORITY["EPSG","8901"]],
//     UNIT["degree",0.01745329251994328,
//         AUTHORITY["EPSG","9122"]],
//     AUTHORITY["EPSG","4326"]]
// `

// func BenchmarkSorahanGdal(b *testing.B) {
// 	// Setup
// 	sref := gdal.CreateSpatialReference(WGS84)
// 	polygon, err := gdal.CreateFromWKT("POLYGON(0 1, 1 2, 2 1, 2 0, 1 -1, 0 0, 0 1)", sref)
// 	if err != nil {
// 		b.Fatal(err)
// 	}
// 	point, err := gdal.CreateFromWKT("POINT(1 0)", sref)
// 	if err != nil {
// 		b.Fatal(err)
// 	}

// 	b.ResetTimer()

// 	// Test
// 	for i := 0; i < b.N; i++ {
// 		if contains := polygon.Contains(point); !contains {
// 			b.Fatal("not contains")
// 		}
// 	}
// }

func BenchmarkBrian(b *testing.B) {
	// Setup
	polygon := Polygon([][][2]float64{
		{
			{0, 1},
			{1, 2},
			{2, 1},
			{2, 0},
			{1, -1},
			{0, 0},
			{0, 1},
		},
	})
	point := Point{1, 0}

	b.ResetTimer()

	// Test
	for i := 0; i < b.N; i++ {
		if contains := polygon.Contains(point); !contains {
			b.Fatal("not contains")
		}
	}
}

// func BenchmarkYan(b *testing.B) {
// 	// Setup
// 	polygon := Coordinates{
// 		{
// 			{0, 1},
// 			{1, 2},
// 			{2, 1},
// 			{2, 0},
// 			{1, -1},
// 			{0, 0},
// 			{0, 1},
// 		},
// 	}
// 	point := Coordinate{1, 0}

// 	b.ResetTimer()

// 	// Test
// 	for i := 0; i < b.N; i++ {
// 		if contains := InsidePolygon(polygon, point); !contains {
// 			b.Fatal("not contains")
// 		}
// 	}
// }
