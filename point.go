package geo

// Point defines a point.
type Point [2]float64

// BoundingBox returns the bounding box for this geometry.
func (point Point) BoundingBox() [][2]float64 {
	return [][2]float64{point}
}
