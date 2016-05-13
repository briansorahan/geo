package geo

const featureCollectionJSONPrefix = `{"type":"FeatureCollection","features":[`

// FeatureCollection represents a GeoJSON feature collection.
type FeatureCollection []Feature

// MarshalJSON marshals the feature collection to JSON.
func (coll FeatureCollection) MarshalJSON() ([]byte, error) {
	buf := []byte(featureCollectionJSONPrefix)
	for i, feature := range coll {
		feat, err := feature.MarshalJSON()
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
