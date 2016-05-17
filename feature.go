package geo

import (
	"bytes"
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
func (feature Feature) MarshalJSON() ([]byte, error) {
	geom, err := feature.Geometry.MarshalJSON()
	if err != nil {
		return nil, err
	}
	props, err := json.Marshal(feature.Properties)
	if err != nil {
		return nil, err
	}
	buf := append(featureJSONPrefix, geom...)
	buf = append(buf, propertiesJSONKey...)
	buf = append(buf, props...)
	return append(buf, '}'), nil
}

// UnmarshalJSON unmarshals a feature from JSON.
func (feature *Feature) UnmarshalJSON(data []byte) error {
	idx := bytes.Index(data, featureJSONPrefix)
	if idx == -1 {
		return fmt.Errorf("could not unmarshal feature from %s", string(data))
	}
	// Get the index of the properties key/value pair.
	// If there isn't one, set the index to the last position in the data,
	// so we can still use it to slice the data.
	pidx := bytes.Index(data, propertiesJSONKey)
	if pidx == -1 {
		pidx = len(data)
	}
	if err := feature.Geometry.UnmarshalJSON(data[idx+len(featureJSONPrefix) : pidx]); err != nil {
		return fmt.Errorf("could not unmarshal feature from %s", string(data))
	}
	// pidx is equal to the end of the data if there are no properties
	if pidx == len(data) {
		return nil
	}
	props := map[string]interface{}{}
	// Slice between the ':' after the properties key and the last '}'
	if err := json.Unmarshal(data[pidx+len(propertiesJSONKey):len(data)-1], &props); err != nil {
		return err
	}
	feature.Properties = props
	return nil
}
