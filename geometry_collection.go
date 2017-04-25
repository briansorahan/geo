package geo

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

const (
	geometryCollectionJSONPrefix = `{"type":"GeometryCollection","geometries":[`
	geometryCollectionWKTPrefix  = `GEOMETRYCOLLECTION`
)

// GeometryCollection is a collection of geometries.
type GeometryCollection []Geometry

// Equal compares the GeometryCollection to a Geometry.
func (gc GeometryCollection) Equal(g Geometry) bool {
	other, ok := g.(*GeometryCollection)
	if !ok || len(gc) != len(*other) {
		return false
	}
	for i, g := range gc {
		if !g.Equal((*other)[i]) {
			return false
		}
	}
	return true
}

// Contains returns true if the GeometryCollection contains the point
// and false otherwise.
func (gc GeometryCollection) Contains(p Point) bool {
	if len(gc) == 0 {
		return false
	}
	for _, g := range gc {
		if !g.Contains(p) {
			return false
		}
	}
	return true
}

// MarshalJSON marshals the GeometryCollection to JSON.
func (gc GeometryCollection) MarshalJSON() ([]byte, error) {
	buf := []byte(geometryCollectionJSONPrefix)
	for i, geometry := range gc {
		data, err := json.Marshal(geometry)
		if err != nil {
			return nil, err
		}
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = append(buf, data...)
	}
	return append(buf, ']', '}'), nil
}

// Scan scans the feature collection from WKT.
// This method expects a GEOMETRYCOLLECTION.
// TODO: implement this.
func (gc *GeometryCollection) Scan(src interface{}) error {
	return scan(gc, src)
}

// scan scans the feature collection from WKT.
// This method expects a GEOMETRYCOLLECTION.
// TODO: implement this.
func (gc *GeometryCollection) scan(s string) error {
	return nil
}

// String converts the feature collection to a GEOMETRYCOLLECTION.
func (gc GeometryCollection) String() string {
	s := geometryCollectionWKTPrefix + "("
	for i, geometry := range gc {
		if i == 0 {
			s += geometry.String()
		} else {
			s += ", " + geometry.String()
		}
	}
	return s + ")"
}

// geometryCollection is a utility type used to unmarshal a geojson GeometryCollection.
type geometryCollection struct {
	Type       string      `json:"type"`
	Geometries []*geometry `json:"geometries"`
	BBox       []float64   `json:"bbox"`
}

func (coll geometryCollection) ToGeometryCollection() (*GeometryCollection, error) {
	geometries := make([]Geometry, len(coll.Geometries))
	for i, g := range coll.Geometries {
		gg, err := g.unmarshalCoordinates()
		if err != nil {
			return nil, err
		}
		geometries[i] = gg
	}
	gc := GeometryCollection(geometries)
	return &gc, nil
}

// UnmarshalJSON unmarshals the geometry collection from geojson.
func (gc *GeometryCollection) UnmarshalJSON(data []byte) error {
	coll, err := unmarshalGeometryCollection(data)
	if err != nil {
		return err
	}
	geometries, err := coll.ToGeometryCollection()
	if err != nil {
		return err
	}
	*gc = *geometries
	return nil
}

// Value returns WKT for the geometry collection.
// Note that this returns a GEOMETRYCOLLECTION.
func (gc GeometryCollection) Value() (driver.Value, error) {
	return gc.String(), nil
}

// Transform transforms the geometry point by point.
func (gc *GeometryCollection) Transform(t Transformer) {
}

// Visit visits each point in the geometry.
func (gc GeometryCollection) Visit(v Visitor) {
}

func unmarshalGeometryCollection(data []byte) (*geometryCollection, error) {
	coll := &geometryCollection{}

	// Never fails because data is always valid JSON.
	_ = json.Unmarshal(data, coll)

	// Check the type.
	if expected, got := GeometryCollectionType, coll.Type; expected != got {
		return nil, fmt.Errorf("expected %s type, got %s", expected, got)
	}
	return coll, nil
}

func unmarshalGeometryCollectionBBox(data []byte) (Geometry, error) {
	coll, err := unmarshalGeometryCollection(data)
	if err != nil {
		return nil, err
	}
	geometries, err := coll.ToGeometryCollection()
	if err != nil {
		return nil, err
	}
	if len(coll.BBox) > 0 {
		return WithBBox(coll.BBox, geometries), nil
	}
	return geometries, nil
}
