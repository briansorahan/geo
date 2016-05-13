package geo

import (
	"encoding/json"
)

const (
	featureJSONPrefix = `{"type":"Feature","geometry":`
)

// Feature is a GeoJSON feature.
type Feature struct {
	Geometry   Geometry    `json:"geometry"`
	Properties interface{} `json:"properties,omitempty"`
}

// MarshalJSON marshals the feature to GeoJSON.
func (feature Feature) MarshalJSON() ([]byte, error) {
	geom, err := feature.Geometry.MarshalJSON()
	if err != nil {
		return nil, err
	}
	props, err := json.Marshal(feature.Properties)
	if err != nil {
		return nil, err
	}
	buf := append([]byte(featureJSONPrefix), geom...)
	buf = append(buf, []byte(`,"properties":`)...)
	buf = append(buf, props...)
	return append(buf, '}'), nil
}
