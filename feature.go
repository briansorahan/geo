package geo

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

var (
	featureJSONPrefix = []byte(`{"type":"Feature","geometry":`)
	propertiesJSONKey = []byte(`,"properties":`)
)

// Feature is a GeoJSON feature.
type Feature struct {
	Geometry   Geometry    `json:"geometry"`
	Properties interface{} `json:"properties,omitempty"`
}

// Compare compares one feature to another.
// Note that this method does not compare properties.
func (f Feature) Compare(g Geometry) bool {
	other, ok := g.(*Feature)
	if !ok {
		return false
	}
	return f.Geometry.Compare(other.Geometry)
}

// Contains determines if the feature's geometry contains the point.
func (f Feature) Contains(p Point) bool {
	return f.Geometry.Contains(p)
}

// MarshalJSON marshals the feature to GeoJSON.
func (f Feature) MarshalJSON() ([]byte, error) {
	geom, err := f.Geometry.MarshalJSON()
	if err != nil {
		return nil, err
	}
	props, err := json.Marshal(f.Properties)
	if err != nil {
		return nil, err
	}
	buf := append(featureJSONPrefix, geom...)
	buf = append(buf, propertiesJSONKey...)
	buf = append(buf, props...)
	return append(buf, '}'), nil
}

// Scan scans a feature from well known text.
func (f *Feature) Scan(src interface{}) error {
	return scan(f, src)
}

// scan scans a feature from well known text.
func (f *Feature) scan(s string) error {
	geom, err := ScanGeometry(s)
	if err != nil {
		return err
	}
	f.Geometry = geom
	return nil
}

// String converts the feature to a WKT string.
func (f Feature) String() string {
	return f.Geometry.String()
}

// feature is a utility type used to unmarshal geojson Feature's.
type feature struct {
	Type       string          `json:"type"`
	Geometry   json.RawMessage `json:"geometry"`
	Properties interface{}     `json:"properties"`
}

// UnmarshalJSON unmarshals a feature from JSON.
func (f *Feature) UnmarshalJSON(data []byte) error {
	feat := feature{}
	if err := json.Unmarshal(data, &feat); err != nil {
		return err
	}
	if feat.Type != FeatureType {
		return fmt.Errorf("could not unmarshal feature from %s", string(data))
	}
	g := geometry{}
	if err := json.Unmarshal(feat.Geometry, &g); err != nil {
		return err
	}
	// Unmarshal the coordinates into one of our Geometry types.
	geom, err := g.unmarshalCoordinates()
	if err != nil {
		return err
	}

	f.Geometry = geom
	f.Properties = feat.Properties

	return nil
}

// Value returns well known text for the feature.
func (f Feature) Value() (driver.Value, error) {
	return f.Geometry.Value()
}
