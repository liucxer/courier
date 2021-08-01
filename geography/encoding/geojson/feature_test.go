package geojson_test

import (
	"bytes"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/liucxer/courier/geography/encoding/geojson"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMarshalFeatureID(t *testing.T) {
	f := &geojson.Feature{
		ID:       "snail",
		Geometry: &geojson.Geometry{},
	}

	data, err := f.MarshalJSON()
	if err != nil {
		t.Fatalf("should marshal, %v", err)
	}

	fmt.Println(string(data))
	if !bytes.Equal(data, []byte(`{"id":"snail","type":"Feature","geometry":{"type":""},"properties":null}`)) {
		t.Errorf("data not correct")
		t.Logf("%v", string(data))
	}
	f.ID = 123
	f.Geometry = &geojson.Geometry{}
	data, err = f.MarshalJSON()
	if err != nil {
		t.Fatalf("should marshal, %v", err)

	}

	if !bytes.Equal(data, []byte(`{"id":123,"type":"Feature","geometry":{"type":""},"properties":null}`)) {
		t.Errorf("data not correct")
		t.Logf("%v", string(data))
	}
}

func TestFeature_UnmarshalText(t *testing.T) {
	var rawJSON = `{
   "type": "Feature",
   "properties": {
       "Id": 111
   },
   "geometry": {
       "type": "Polygon",
       "coordinates": [[
               [11.11, 11.12],
               [12.11, 12.12],
               [13.11, 13.12],
               [14.11, 14.12]
		]]
   },
	"crs":{
		"name":"snail"
	}
}`

	feat := &geojson.Feature{}
	err := feat.UnmarshalText([]byte(rawJSON))
	require.NoError(t, err)
	spew.Dump(feat)
}
