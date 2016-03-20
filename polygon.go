package geo

import "fmt"

// Polygon is the GeoJSON Polygon geometry.
type Polygon [][2]float64

// BoundingBox returns the bounding box for this geometry.
func (polygon Polygon) BoundingBox() [][2]float64 {
	return polygon
}

// validate validates the Polygon.
func (ls Polygon) validate() error {
	if len(ls) > 1 {
		return nil
	}
	return fmt.Errorf("Polygon only contains %d positions", len(ls))
}
