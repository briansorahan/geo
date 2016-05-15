// Package geo provides a small set of geometrical types
// and operations on them.
package geo

import "encoding/json"

// Margin is a value used to fudge equality when two floats are very close to each other.
var Margin = 1e-100

// Geometry defines the interface of every geometry type.
type Geometry interface {
	json.Marshaler
	json.Unmarshaler
	Compare(other Geometry) bool
	String() string
}
