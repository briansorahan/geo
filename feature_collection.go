package geo

import (
	"database/sql/driver"
	"encoding/json"
)

const featureCollectionJSONPrefix = `{"type":"FeatureCollection","features":[`

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
func (coll *FeatureCollection) Scan(src interface{}) error {
	return nil
}

// scan scans the feature collection from WKT.
// This method expects a GEOMETRYCOLLECTION.
func (coll *FeatureCollection) scan(s string) error {
	return nil
}

// String converts the feature collection to a GEOMETRYCOLLECTION.
func (coll FeatureCollection) String() string {
	return ""
}

// featureCollection is a utility type used to unmarshal geojson FeatureCollection's.
type featureCollection struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

// UnmarshalJSON unmarshals the feature collection from geojson.
func (coll *FeatureCollection) UnmarshalJSON(data []byte) error {
	fc := featureCollection{}
	if err := json.Unmarshal(data, &fc); err != nil {
		return err
	}
	*coll = make([]Feature, len(fc.Features))
	for i, feat := range fc.Features {
		(*coll)[i] = feat
	}
	return nil
}

// Value returns WKT for the feature collection.
// Note that this returns a GEOMETRYCOLLECTION.
func (coll FeatureCollection) Value() (driver.Value, error) {
	return coll.String(), nil
}
