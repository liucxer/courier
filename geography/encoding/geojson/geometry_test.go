package geojson_test

import (
	"testing"

	"github.com/liucxer/courier/geography/encoding/geojson"
	"github.com/stretchr/testify/require"
)

func TestGeometry_UnmarshalText(t *testing.T) {
	rawJSON := `{"type": "GeometryCollection", "geometries": [
		{"type": "Point", "coordinates": [102.0, 0.5]},
		{"type": "MultiPoint", "coordinates": [[102.0, 0.5],[111,222]]},
		{"type": "LineString", "coordinates": [[102.0, 0.5],[111,222]]},
		{"type": "MultiLineString", "coordinates": [[[11,22],[33,44]],[[55,66],[77,888]]]},
		{"type": "Polygon", "coordinates": [[[11,22],[33,44]],[[55,66],[77,888]]]},
		{"type": "MultiPolygon", "coordinates": [[[[11,22],[33,44]],[[55,66],[77,888]]]]}

	]}`

	geo := &geojson.Geometry{}
	err := geo.UnmarshalText([]byte(rawJSON))
	require.NoError(t, err)
	//spew.Dump(geo)
}
