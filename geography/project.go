package geography

import (
	"math"
)

const (
	// https://epsg.io/3857
	SRS3857 = 3857
	// https://epsg.io/4326
	SRS4326 = 4326
)

const (
	EPSLN     = 1.0e-10
	D2R       = math.Pi / 180.0
	R2D       = 180.0 / math.Pi
	A         = 6378137.0
	MAXEXTENT = 20037508.342789244
)

func PseudoMercatorToWGS84(x, y float64) (lon, lat float64) {
	return x * R2D / A, ((math.Pi * 0.5) - 2.0*math.Atan(math.Exp(-y/A))) * R2D
}

func WGS84ToPseudoMercator(lon, lat float64) (x, y float64) {
	x = A * lon * D2R
	y = A * math.Log(math.Tan((math.Pi*0.25)+(0.5*lat*D2R)))

	// if xy value is beyond maxextent (e.g. poles), return maxextent.
	if x > MAXEXTENT {
		x = MAXEXTENT
	}
	if x < -MAXEXTENT {
		x = -MAXEXTENT
	}
	if y > MAXEXTENT {
		y = MAXEXTENT
	}
	if y < -MAXEXTENT {
		y = -MAXEXTENT
	}
	return
}

func TileXYToLonLat(x, y float64, zoom uint32) (lon, lat float64) {
	maxTiles := float64(uint32(1 << zoom))
	return 360.0 * (x/maxTiles - 0.5), 2.0*math.Atan(math.Exp(math.Pi-(2*math.Pi)*(y/maxTiles)))*R2D - 90.0
}
