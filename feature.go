package geo

import (
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

type feature struct {
	Type     string          `json:"type"`
	Geometry json.RawMessage `json:"geometry"`
}

// UnmarshalJSON unmarshals a feature from JSON.
func (f *Feature) UnmarshalJSON(data []byte) error {
	feat := feature{}
	if err := json.Unmarshal(data, &feat); err != nil {
		return err
	}
	if feat.Type != "Feature" {
		return fmt.Errorf("could not unmarshal feature from %s", string(data))
	}
	g := geometry{}
	if err := json.Unmarshal(feat.Geometry, &g); err != nil {
		return err
	}
	geom, err := g.Geometry()
	if err != nil {
		return err
	}
	f.Geometry = geom
	return nil
}
