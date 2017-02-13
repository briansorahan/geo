package geo

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

const (
	featureCollectionJSONPrefix = `{"type":"FeatureCollection","features":[`
	featureCollectionWKTPrefix  = `GEOMETRYCOLLECTION`
)

// FeatureCollection represents a feature collection.
type FeatureCollection []Feature

// Compare compares one feature collection to another.
func (coll FeatureCollection) Compare(g Geometry) bool {
	other, ok := g.(*FeatureCollection)
	if !ok {
		return false
	}
	if len(coll) != len(*other) {
		return false
	}
	for i, feat := range coll {
		if !feat.Compare(&(*other)[i]) {
			return false
		}
	}
	return true
}

// Contains returns true if every feature in the collection
// contains the p, and false otherwise.
func (coll FeatureCollection) Contains(p Point) bool {
	for _, feat := range coll {
		if !feat.Contains(p) {
			return false
		}
	}
	return true
}

// MarshalJSON marshals the feature collection to geojson.
func (coll FeatureCollection) MarshalJSON() ([]byte, error) {
	buf := []byte(featureCollectionJSONPrefix)
	for i, feature := range coll {
		feat, err := json.Marshal(feature)
		if err != nil {
			return nil, err
		}
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = append(buf, feat...)
	}
	return append(buf, ']', '}'), nil
}

// Scan scans the feature collection from WKT.
// This method expects a GEOMETRYCOLLECTION.
// TODO: implement this.
func (coll *FeatureCollection) Scan(src interface{}) error {
	return scan(coll, src)
}

// scan scans the feature collection from WKT.
// This method expects a GEOMETRYCOLLECTION.
// TODO: implement this.
func (coll *FeatureCollection) scan(s string) error {
	return nil
}

// String converts the feature collection to a GEOMETRYCOLLECTION.
func (coll FeatureCollection) String() string {
	s := featureCollectionWKTPrefix + "("
	for i, feat := range coll {
		if i == 0 {
			s += feat.Geometry.String()
		} else {
			s += ", " + feat.String()
		}
	}
	return s + ")"
}

// featureCollection is a utility type used to unmarshal a geojson FeatureCollection.
type featureCollection struct {
	Type     string     `json:"type"`
	Features []*feature `json:"features"`
	BBox     []float64  `json:"bbox"`
}

func (fc *featureCollection) ToFeatureCollection() (*FeatureCollection, error) {
	coll := make([]Feature, len(fc.Features))
	for i, feat := range fc.Features {
		f, err := feat.ToFeature()
		if err != nil {
			return nil, err
		}
		coll[i] = *f
	}
	featcoll := FeatureCollection(coll)
	return &featcoll, nil
}

// UnmarshalJSON unmarshals the feature collection from geojson.
func (coll *FeatureCollection) UnmarshalJSON(data []byte) error {
	fc, err := unmarshalFeatureCollection(data)
	if err != nil {
		return err
	}
	featcoll, err := fc.ToFeatureCollection()
	if err != nil {
		return err
	}
	*coll = *featcoll
	return nil
}

// Value returns WKT for the feature collection.
// Note that this returns a GEOMETRYCOLLECTION.
func (coll FeatureCollection) Value() (driver.Value, error) {
	return coll.String(), nil
}

func unmarshalFeatureCollection(data []byte) (*featureCollection, error) {
	fc := &featureCollection{}

	// Never fails because data is always valid JSON.
	_ = json.Unmarshal(data, fc)

	// Check the type.
	if expected, got := FeatureCollectionType, fc.Type; expected != got {
		return nil, fmt.Errorf("expected %s type, got %s", expected, got)
	}
	return fc, nil
}

func unmarshalFeatureCollectionBBox(data []byte) (Geometry, error) {
	fc, err := unmarshalFeatureCollection(data)
	if err != nil {
		return nil, err
	}
	coll, err := fc.ToFeatureCollection()
	if err != nil {
		return nil, err
	}
	if len(fc.BBox) > 0 {
		return WithBBox(fc.BBox, coll), nil
	}
	return coll, nil
}
