package maptile

import (
	"math"
	"sync"

	"github.com/liucxer/courier/geography/encoding/mvt"

	"github.com/liucxer/courier/geography"
)

func NewMapTile(z, x, y uint32) *MapTile {
	return &MapTile{
		Z: z,
		X: x,
		Y: y,
	}
}

type MapTile struct {
	coordsTransform CoordsTransform
	Z               uint32
	X               uint32
	Y               uint32
	Layers          []*Layer
}

type CoordsTransform interface {
	ToEarth(point geography.Point) geography.Point
	ToMars(point geography.Point) geography.Point
}

func (t *MapTile) SetCoordsTransform(coordsTransform CoordsTransform) {
	t.coordsTransform = coordsTransform
}

func (t *MapTile) MarshalMVT(w *mvt.MVTWriter) error {
	for i := range t.Layers {
		layer := t.Layers[i]
		if layer == nil || len(layer.Features) == 0 {
			continue
		}

		features := make([]*mvt.Feature, 0)

		for i := range layer.Features {
			feat := layer.Features[i]
			if feat == nil {
				continue
			}

			geo := feat.ToGeom()
			if geo == nil {
				continue
			}
			g := geo.Project(t.NewTransform(layer.Extent))

			f := &mvt.Feature{
				Type:       g.Type(),
				Geometry:   g.Geometry(),
				Properties: feat.Properties(),
			}
			if f == nil || len(f.Geometry) == 0 {
				continue
			}

			if fid, ok := feat.(FeatureID); ok {
				f.ID = fid.ID()
			}

			features = append(features, f)
		}

		w.WriteLayer(layer.Name, layer.Extent, features...)
	}
	return nil
}

func (t *MapTile) NewTransform(extent uint32) geography.Transform {
	n := uint32(TrailingZeros32(extent))
	z := uint32(t.Z) + n

	minx := float64(t.X << n)
	miny := float64(t.Y << n)

	return func(p geography.Point) geography.Point {
		if t.coordsTransform != nil {
			p = t.coordsTransform.ToMars(p)
		}
		x, y := lonLatToPixelXY(p[0], p[1], z)
		return geography.Point{
			math.Floor(x - minx),
			math.Floor(y - miny),
		}
	}
}

func (t *MapTile) BBox() geography.Bound {
	buffer := 0.0
	x := float64(t.X)
	y := float64(t.Y)

	minx := x - buffer

	miny := y - buffer
	if miny < 0 {
		miny = 0
	}

	lon1, lat1 := geography.TileXYToLonLat(minx, miny, uint32(t.Z))

	maxX := x + 1 + buffer

	maxTiles := float64(uint32(1 << t.Z))
	maxY := y + 1 + buffer
	if maxY > maxTiles {
		maxY = maxTiles
	}

	lon2, lat2 := geography.TileXYToLonLat(maxX, maxY, uint32(t.Z))

	if t.coordsTransform != nil {
		return geography.Bound{
			Min: t.coordsTransform.ToEarth(geography.Point{lon1, lat2}),
			Max: t.coordsTransform.ToEarth(geography.Point{lon2, lat1}),
		}
	}

	return geography.Bound{
		Min: geography.Point{lon1, lat2},
		Max: geography.Point{lon2, lat1},
	}
}

func (t *MapTile) AddLayers(layers ...*Layer) {
	t.Layers = append(t.Layers, layers...)
}

func (t *MapTile) AddTileLayers(tileLayers ...TileLayer) (e error) {
	wg := sync.WaitGroup{}

	result := make(chan interface{})

	for i := range tileLayers {
		wg.Add(1)
		go func(tileLayer TileLayer) {
			defer wg.Done()
			features, err := tileLayer.Features(t)
			if err != nil {
				result <- err
				return
			}
			extend := uint32(0)

			if tileLayerExtentConf, ok := tileLayer.(TileLayerExtentConf); ok {
				extend = tileLayerExtentConf.Extent()
			}

			result <- NewLayer(tileLayer.Name(), extend, features...)

		}(tileLayers[i])
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	for r := range result {
		switch v := r.(type) {
		case error:
			e = v
		case *Layer:
			t.AddLayers(v)
		}
	}
	return
}
