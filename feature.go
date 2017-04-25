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

// Equal compares one feature to another.
// Note that this method does not compare properties.
func (f Feature) Equal(g Geometry) bool {
	other, ok := g.(*Feature)
	if !ok {
		return false
	}
	return f.Geometry.Equal(other.Geometry)
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
	BBox       []float64       `json:"bbox"`
}

// ToFeature converts the private feature type to the public one.
func (f *feature) ToFeature() (*Feature, error) {
	g := geometry{}

	if err := json.Unmarshal(f.Geometry, &g); err != nil {
		return nil, err
	}

	// Unmarshal the coordinates into one of our Geometry types.
	geom, err := g.unmarshalCoordinates()
	if err != nil {
		return nil, err
	}
	feat := &Feature{}
	feat.Geometry = geom
	feat.Properties = f.Properties
	return feat, nil
}

// UnmarshalJSON unmarshals a feature from JSON.
func (f *Feature) UnmarshalJSON(data []byte) error {
	ff, _, err := unmarshalFeature(data)
	if err != nil {
		return err
	}
	*f = *ff
	return nil
}

// Value returns well known text for the feature.
func (f Feature) Value() (driver.Value, error) {
	return f.Geometry.Value()
}

// Transform transforms the geometry point by point.
func (f *Feature) Transform(t Transformer) {
	f.Transform(t)
}

// Visit visits each point in the geometry.
func (f Feature) Visit(v Visitor) {
	f.Visit(v)
}

func unmarshalFeature(data []byte) (*Feature, *feature, error) {
	feat := feature{}

	// Never fails because data is always valid JSON.
	_ = json.Unmarshal(data, &feat)

	// Check the type.
	if expected, got := FeatureType, feat.Type; expected != got {
		return nil, nil, fmt.Errorf("expected type %s, got %s", expected, got)
	}

	g := geometry{}

	if err := json.Unmarshal(feat.Geometry, &g); err != nil {
		return nil, nil, err
	}

	// Unmarshal the coordinates into one of our Geometry types.
	geom, err := g.unmarshalCoordinates()
	if err != nil {
		return nil, nil, err
	}
	f := &Feature{}
	f.Geometry = geom
	f.Properties = feat.Properties
	return f, &feat, nil
}

func unmarshalFeatureBBox(data []byte) (Geometry, error) {
	f, ff, err := unmarshalFeature(data)
	if err != nil {
		return nil, err
	}
	if len(ff.BBox) > 0 {
		return WithBBox(ff.BBox, f), nil
	}
	return f, nil
}
