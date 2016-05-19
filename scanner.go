package geo

import "fmt"

// scanner is an interface for scanning Well Known Text.
type scanner interface {
	scan(s string) error
}

// scan scans an interface{} with a scanner
func scan(s scanner, data interface{}) error {
	switch v := data.(type) {
	case []byte:
		return s.scan(string(v))
	case string:
		return s.scan(v)
	default:
		return fmt.Errorf("could not scan from %T", data)
	}
}
