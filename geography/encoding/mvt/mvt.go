package mvt

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"fmt"

	"github.com/liucxer/courier/geography/encoding/mvt/vector_tile"
	"github.com/liucxer/courier/ptr"
	"github.com/golang/protobuf/proto"
)

func ToMVT(v MVTMarshaller) (*MVT, error) {
	w := &MVTWriter{}
	if err := v.MarshalMVT(w); err != nil {
		return nil, err
	}
	mvtData, err := proto.Marshal(&w.Tile)
	if err != nil {
		return nil, err
	}

	mvt := &MVT{}
	mvt.Write(mvtData)
	return mvt, nil
}

type MVTMarshaller interface {
	MarshalMVT(w *MVTWriter) error
}

type MVTWriter struct {
	vector_tile.Tile
}

func (w *MVTWriter) WriteLayer(name string, extent uint32, features ...*Feature) {
	layer := &vector_tile.Tile_Layer{
		Version: ptr.Uint32(2),
		Name:    &name,
		Extent:  &extent,
	}

	keyIdxSet := make(map[string]uint32, 0)
	valueIdxSet := make(map[string]uint32, 0)

	addKeyValue := func(f *vector_tile.Tile_Feature, key string, value interface{}) {
		tv := vectorTileValue(value)

		keyIdx, ok := keyIdxSet[key]
		if !ok {
			layer.Keys = append(layer.Keys, key)
			keyIdxSet[key] = uint32(len(layer.Keys) - 1)
			keyIdx = keyIdxSet[key]
		}

		valueKey := tv.String()
		valueIdx, ok := valueIdxSet[valueKey]
		if !ok {
			layer.Values = append(layer.Values, tv)
			valueIdxSet[valueKey] = uint32(len(layer.Values) - 1)
			valueIdx = valueIdxSet[valueKey]
		}

		f.Tags = append(f.Tags, keyIdx, valueIdx)
	}

	layer.Features = func() []*vector_tile.Tile_Feature {
		vtFeatures := make([]*vector_tile.Tile_Feature, 0)
		for i := range features {
			f := features[i]
			if f == nil {
				continue
			}
			id := f.ID
			geomType := f.GeomType()

			feat := &vector_tile.Tile_Feature{
				Id:       &id,
				Type:     &geomType,
				Geometry: f.Geometry,
			}

			for k, v := range f.Properties {
				addKeyValue(feat, k, v)
			}
			vtFeatures = append(vtFeatures, feat)
		}
		return vtFeatures
	}()

	w.Tile.Layers = append(w.Tile.Layers, layer)
}

type MVT struct {
	bytes.Buffer
}

func (MVT) ContextType() string {
	return "application/vnd.mapbox-vector-tile"
}

func vectorTileValue(i interface{}) *vector_tile.Tile_Value {
	tv := new(vector_tile.Tile_Value)
	switch t := i.(type) {
	default:
		buff := new(bytes.Buffer)
		err := binary.Write(buff, binary.BigEndian, t)
		if err == nil {
			tv.XXX_unrecognized = buff.Bytes()
		}
	case encoding.TextMarshaler:
		data, err := t.MarshalText()
		if err == nil {
			str := string(data)
			tv.StringValue = &str
		}
	case string:
		tv.StringValue = &t
	case fmt.Stringer:
		str := t.String()
		tv.StringValue = &str
	case bool:
		tv.BoolValue = &t
	case int8:
		intv := int64(t)
		tv.SintValue = &intv
	case int16:
		intv := int64(t)
		tv.SintValue = &intv
	case int32:
		intv := int64(t)
		tv.SintValue = &intv
	case int64:
		tv.IntValue = &t
	case uint8:
		intv := int64(t)
		tv.SintValue = &intv
	case uint16:
		intv := int64(t)
		tv.SintValue = &intv
	case uint32:
		intv := int64(t)
		tv.SintValue = &intv
	case uint64:
		tv.UintValue = &t
	case float32:
		tv.FloatValue = &t
	case float64:
		tv.DoubleValue = &t
	}
	return tv
}
