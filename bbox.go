package geo

import "encoding/json"

// WithBBox returns a geometry that contains a "bbox" property.
// TODO: it might make sense to calculate the bounding box from the geometry's coordinates.
// TODO: validate the size of the bounding box relative to the dimensionality of the geometry.
// TODO: verify that all axes of the most southwesterly point are followed by all axes of the more northeasterly point (should also be done for geometries)
func WithBBox(bbox []float64, geom Geometry) Geometry {
	return &boundingBox{
		Geometry: geom,
		Box:      bbox,
	}
}

// boundingBox is a utility type for decorating geometries with a bounding box.
type boundingBox struct {
	Geometry

	Box []float64 `json:"bbox"`
}

// Equal compares two geometries that have bounding boxes.
func (bbox *boundingBox) Equal(g Geometry) bool {
	other, ok := g.(*boundingBox)
	if !ok {
		return false
	}
	if !boxEqual(bbox.Box, other.Box) {
		return false
	}
	return bbox.Geometry.Equal(other.Geometry)
}

// MarshalJSON marshals a geometry with a bounding box.
func (bbox *boundingBox) MarshalJSON() ([]byte, error) {
	gdata, err := bbox.Geometry.MarshalJSON()
	if err != nil {
		return nil, err
	}
	bboxData, _ := json.Marshal(bbox.Box) // Never fails.
	bboxData = append([]byte(`,"bbox":`), append(bboxData, '}')...)
	return append(gdata[:len(gdata)-1], bboxData...), nil
}

// boxEqual compares two float64 slices.
func boxEqual(b1, b2 []float64) bool {
	if len(b1) != len(b2) {
		return false
	}
	for i, x := range b1 {
		if b2[i] != x {
			return false
		}
	}
	return true
}
