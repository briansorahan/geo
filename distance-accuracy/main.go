package main

import (
	"encoding/json"
	"flag"
	"math"
	"os"

	"github.com/briansorahan/geo"
)

func contains(p geo.Circle, lng, lat float64) bool {
	var (
		dlng = (p.Coordinates[0] - lng) * float64(228200)
		dlat = (p.Coordinates[1] - lat) * float64(364000)
		d    = math.Sqrt(math.Pow(dlng, 2) + math.Pow(dlat, 2))
	)
	return d <= p.Radius
}

func main() {
	var (
		algo = flag.String("algo", "", "great circle distance algorithm (haversine, slc, equirectangular)")
		pts  = geo.Line{
			{-111.990439295769, 39.8762828965173},
		}
		pivot = geo.Circle{
			Coordinates: [2]float64{-111.990439295769, 39.8762828965173},
			Radius:      1527.76981726524,
		}
		inside, outside geo.Line
	)
	flag.Parse()

	for i := 0; i < 1000; i++ {
		next := [2]float64{
			pts[i][0] + 0.0000075,
			pts[0][1],
		}
		pts = append(pts, next)
		switch *algo {
		default:
			if !contains(pivot, next[0], next[1]) {
				outside = append(outside, next)
			} else {
				inside = append(inside, next)
			}
		case "haversine":
			if !pivot.ContainsHaversine(next) {
				outside = append(outside, next)
			} else {
				inside = append(inside, next)
			}
		case "slc":
			if !pivot.ContainsSLC(next) {
				outside = append(outside, next)
			} else {
				inside = append(inside, next)
			}
		case "equirectangular":
			if !pivot.ContainsEquirectangular(next) {
				outside = append(outside, next)
			} else {
				inside = append(inside, next)
			}
		}
	}
	feats := geo.FeatureCollection{
		geo.Feature{
			Geometry: &inside,
			Properties: map[string]interface{}{
				"stroke": "#ee6677",
			},
		},
		geo.Feature{
			Geometry: &outside,
			Properties: map[string]interface{}{
				"stroke": "#3344dd",
			},
		},
	}
	if err := json.NewEncoder(os.Stdout).Encode(feats); err != nil {
		panic(err)
	}
}
