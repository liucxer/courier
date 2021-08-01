package mvt

import (
	"strings"

	"github.com/liucxer/courier/geography/encoding/mvt/vector_tile"
)

type FeatureMarshaller interface {
	Type() string
	Geometry() []uint32
}

type Coord interface {
	X() float64
	Y() float64
}

type Feature struct {
	ID         uint64
	Type       string
	Geometry   []uint32
	Properties map[string]interface{}
}

func (f Feature) GeomType() vector_tile.Tile_GeomType {
	switch strings.ToUpper(f.Type) {
	case "POINT", "MULTIPOINT":
		return vector_tile.Tile_POINT
	case "POLYGON", "MULTIPOLYGON":
		return vector_tile.Tile_POLYGON
	case "LINESTRING", "MULTILINESTRING":
		return vector_tile.Tile_LINESTRING
	}
	return vector_tile.Tile_UNKNOWN
}

const (
	moveTo    = 1
	lineTo    = 2
	closePath = 7
)

func NewFeatureWriter(cap int) *FeatureWriter {
	return &FeatureWriter{
		data: make([]uint32, 0, cap),
	}
}

type FeatureWriter struct {
	data         []uint32
	prevX, prevY int32
}

func (w *FeatureWriter) Data() []uint32 {
	return w.data
}

func (w *FeatureWriter) MoveTo(l int, getCoord func(i int) Coord) {
	w.data = append(w.data, (uint32(l)<<3)|moveTo)
	for i := 0; i < l; i++ {
		w.addCoord(getCoord(i))
	}
}

func (w *FeatureWriter) LineTo(l int, getCoord func(i int) Coord) {
	w.data = append(w.data, (uint32(l)<<3)|lineTo)
	for i := 0; i < l; i++ {
		w.addCoord(getCoord(i))
	}
}

func (w *FeatureWriter) ClosePath() {
	w.data = append(w.data, (1<<3)|closePath)
}

func (w *FeatureWriter) addCoord(coord Coord) {
	x0 := int32(coord.X())
	y0 := int32(coord.Y())

	dx := x0 - w.prevX
	dy := y0 - w.prevY

	w.prevX = x0
	w.prevY = y0

	w.data = append(w.data,
		uint32((dx<<1)^(dx>>31)),
		uint32((dy<<1)^(dy>>31)),
	)
}
