package coordstransform

import (
	"fmt"
	"github.com/liucxer/courier/geography"
)

func ExampleWGS84toGCJ02() {
	lngLat := (CoordsTransform{}).ToMars(geography.Point{116.404, 39.915})
	fmt.Println(lngLat[0], lngLat[1])
	// Output:
	// 116.41024449916938 39.91640428150164
}

func ExampleGCJ02toWGS84() {
	lngLat := (CoordsTransform{}).ToEarth(geography.Point{116.404, 39.915})
	fmt.Println(lngLat[0], lngLat[1])
	// Output:
	// 116.39775550083061 39.91359571849836
}
