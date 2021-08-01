package coordstransform

import (
	"math"

	"github.com/liucxer/courier/geography"
)

type CoordsTransform struct {
}

const (
	OFFSET = 0.00669342162296594323
	AXIS   = 6378245.0
)

func (CoordsTransform) ToEarth(point geography.Point) geography.Point {
	if point.IsZero() {
		return point
	}
	mgLon, mgLat := delta(point[0], point[1])
	return geography.Point{point[0]*2 - mgLon, point[1]*2 - mgLat}
}

func (CoordsTransform) ToMars(point geography.Point) geography.Point {
	if point.IsZero() {
		return point
	}

	mgLon, mgLat := delta(point[0], point[1])
	return geography.Point{mgLon, mgLat}
}

func delta(lon, lat float64) (float64, float64) {
	dlat := transformlat(lon-105.0, lat-35.0)
	dlon := transformlng(lon-105.0, lat-35.0)

	radlat := lat / 180.0 * math.Pi
	magic := math.Sin(radlat)
	magic = 1 - OFFSET*magic*magic
	sqrtmagic := math.Sqrt(magic)

	dlat = (dlat * 180.0) / ((AXIS * (1 - OFFSET)) / (magic * sqrtmagic) * math.Pi)
	dlon = (dlon * 180.0) / (AXIS / sqrtmagic * math.Cos(radlat) * math.Pi)

	mgLat := lat + dlat
	mgLon := lon + dlon

	return mgLon, mgLat
}

func transformlat(lon, lat float64) float64 {
	var ret = -100.0 + 2.0*lon + 3.0*lat + 0.2*lat*lat + 0.1*lon*lat + 0.2*math.Sqrt(math.Abs(lon))
	ret += (20.0*math.Sin(6.0*lon*math.Pi) + 20.0*math.Sin(2.0*lon*math.Pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(lat*math.Pi) + 40.0*math.Sin(lat/3.0*math.Pi)) * 2.0 / 3.0
	ret += (160.0*math.Sin(lat/12.0*math.Pi) + 320*math.Sin(lat*math.Pi/30.0)) * 2.0 / 3.0
	return ret
}

func transformlng(lon, lat float64) float64 {
	var ret = 300.0 + lon + 2.0*lat + 0.1*lon*lon + 0.1*lon*lat + 0.1*math.Sqrt(math.Abs(lon))
	ret += (20.0*math.Sin(6.0*lon*math.Pi) + 20.0*math.Sin(2.0*lon*math.Pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(lon*math.Pi) + 40.0*math.Sin(lon/3.0*math.Pi)) * 2.0 / 3.0
	ret += (150.0*math.Sin(lon/12.0*math.Pi) + 300.0*math.Sin(lon/30.0*math.Pi)) * 2.0 / 3.0
	return ret
}
