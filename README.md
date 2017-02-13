# geo

[![Go Report Card](https://goreportcard.com/badge/github.com/briansorahan/geo)](https://goreportcard.com/report/github.com/briansorahan/geo)
[![wercker status](https://app.wercker.com/status/deafc383e082c1a3fd05f5550383592e/s/master "wercker status")](https://app.wercker.com/project/byKey/deafc383e082c1a3fd05f5550383592e)

Simple package for converting geometrical primitives to/from [GeoJSON](http://geojson.org) and [Well Known Text](https://en.wikipedia.org/wiki/Well-known_text).

This package aims to be simple and high quality.
If test coverage is not 100% feel free to open an issue.

Note this package is not [RFC 7946](https://tools.ietf.org/html/rfc7946) compliant and includes the non-standard "Circle" type. The circle implementation seeks to adhere to this https://github.com/geojson/geojson-spec/wiki/Proposal---Circles-and-Ellipses-Geoms
