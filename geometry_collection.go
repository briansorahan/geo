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

// Compare compares the GeometryCollection to a Geometry.
func (gc GeometryCollection) Compare(g Geometry) bool {
	other, ok := g.(*GeometryCollection)
	if !ok || len(gc) != len(*other) {
		return false
	}
	for i, g := range gc {
		if !g.Compare((*other)[i]) {
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
	Type       string     `json:"type"`
	Geometries []Geometry `json:"geometries"`
}

// UnmarshalJSON unmarshals the geometry collection from geojson.
func (gc *GeometryCollection) UnmarshalJSON(data []byte) error {
	fc := geometryCollection{}

	// Never fails because data is always valid JSON.
	_ = json.Unmarshal(data, &fc)

	// Check the type.
	if expected, got := GeometryCollectionType, fc.Type; expected != got {
		return fmt.Errorf("expected %s type, got %s", expected, got)
	}

	// Clear the collection and copy the unmarshalled geometries.
	*gc = make([]Geometry, len(fc.Geometries))
	copy(*gc, fc.Geometries)

	return nil
}

// Value returns WKT for the geometry collection.
// Note that this returns a GEOMETRYCOLLECTION.
func (gc GeometryCollection) Value() (driver.Value, error) {
	return gc.String(), nil
}
