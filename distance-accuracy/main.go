package main

import (
	"encoding/json"
	"os"

	"github.com/briansorahan/geo"
)

func main() {
	var (
		inside = &geo.Line{
			{-111.953309476376, 40.0652636206423},
		}
		outside = &geo.Line{}
		pivot   = &geo.Circle{
			Coordinates: {-111.953309476376, 40.0652636206423},
			Radius:      1336.23605143716,
			RadiusUnits: "feet",
		}
	)
	for i := 0; i < 1000; i++ {
		next := [2]float64{
			(*line)[i][0] + 0.0000075,
			(*line)[0][1],
		}
		*inside = append(*inside, next)
	}
	feat := geo.Feature{
		Geometry: line,
		Properties: map[string]interface{}{
			"stroke": "#ee6677",
		},
	}
	if err := json.NewEncoder(os.Stdout).Encode(feat); err != nil {
		panic(err)
	}
}
