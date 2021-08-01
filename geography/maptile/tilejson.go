package maptile

const TileJSONVersion = "2.2.0"

type SchemeType string

const (
	SchemeTypeXYZ  SchemeType = "xyz"
	SchemeTypeTMLS SchemeType = "tms"
)

func NewTileBounds(minLon, minLat, maxLon, maxLat float64) *[4]float64 {
	v := [4]float64{minLon, minLat, maxLon, maxLat}
	return &v
}

func NewTileCenter(lon, lat float64, zoom float64) *[3]float64 {
	v := [3]float64{lon, lat, zoom}
	return &v
}

// https://github.com/mapbox/tilejson-spec
type TileJSON struct {
	// REQUIRED. A semver.org style version number. Describes the version of
	// the TileJSON spec that is implemented by this JSON object.
	TileJSON string `json:"tilejson"`
	// REQUIRED. An array of tile endpoints. {z}, {x} and {y}, if present,
	// are replaced with the corresponding integers. If multiple endpoints are specified, clients
	// may use any combination of endpoints. All endpoints MUST return the same
	// content for the same URL. The array MUST contain at least one endpoint.
	Tiles []string `json:"tiles"`
	// OPTIONAL. Default: null. Contains an attribution to be displayed
	// when the map is shown to a user. Implementations MAY decide to treat this
	// as HTML or literal text. For security reasons, make absolutely sure that
	// this field can't be abused as a vector for XSS or beacon tracking.
	Attribution *string `json:"attribution,omitempty"`
	// OPTIONAL. Default: [-180, -90, 180, 90].
	// The maximum extent of available map tiles. Bounds MUST define an area
	// covered by all zoom levels. The bounds are represented in WGS:84
	// latitude and longitude values, in the order left, bottom, right, top.
	// Values may be integers or floating point numbers.
	Bounds *[4]float64 `json:"bounds,omitempty"`
	// OPTIONAL. Default: null.
	// The first value is the longitude, the second is latitude (both in
	// WGS:84 values), the third value is the zoom level as an integer.
	// Longitude and latitude MUST be within the specified bounds.
	// The zoom level MUST be between minzoom and maxzoom.
	// Implementations can use this value to set the default location. If the
	// value is null, implementations may use their own algorithm for
	// determining a default location.
	Center *[3]float64 `json:"center,omitempty"`
	// pbf - protocol buffer
	Format string `json:"format"`
	// OPTIONAL. Default: 0. >= 0, <= 22.
	// A positive integer specifying the minimum zoom level.
	MinZoom uint `json:"minzoom,omitempty"`
	// OPTIONAL. Default: 22. >= 0, <= 22.
	// An positive integer specifying the maximum zoom level. MUST be >= minzoom.
	MaxZoom uint `json:"maxzoom,omitempty"`
	// OPTIONAL. Default: null. A name describing the tileset. The name can
	// contain any legal character. Implementations SHOULD NOT interpret the
	// name as HTML.
	Name string `json:"name,omitempty"`
	// OPTIONAL. Default: null. A text description of the tileset. The
	// description can contain any legal character. Implementations SHOULD NOT
	// interpret the description as HTML.
	Description string `json:"description,omitempty"`
	// OPTIONAL. Default: "xyz". Either "xyz" or "tms". Influences the y
	// direction of the tile coordinates.
	// The global-mercator (aka Spherical Mercator) profile is assumed.
	Scheme SchemeType `json:"scheme,omitempty"`
	// OPTIONAL. Default: []. An array of interactivity endpoints. {z}, {x}
	// and {y}, if present, are replaced with the corresponding integers. If multiple
	// endpoints are specified, clients may use any combination of endpoints.
	// All endpoints MUST return the same content for the same URL.
	// If the array doesn't contain any entries, interactivity is not supported
	// for this tileset.
	// See https://github.com/mapbox/utfgrid-spec/tree/master/1.2
	// for the interactivity specification.
	Grids []string `json:"grids,omitempty"`
	// OPTIONAL. Default: []. An array of data files in GeoJSON format.
	// {z}, {x} and {y}, if present,
	// are replaced with the corresponding integers. If multiple
	// endpoints are specified, clients may use any combination of endpoints.
	// All endpoints MUST return the same content for the same URL.
	// If the array doesn't contain any entries, then no data is present in
	// the map.
	Data []string `json:"data,omitempty"`
	// OPTIONAL. Default: "1.0.0". A semver.org style version number. When
	// changes across tiles are introduced, the minor version MUST change.
	// This may lead to cut off labels. Therefore, implementors can decide to
	// clean their cache when the minor version changes. Changes to the patch
	// level MUST only have changes to tiles that are contained within one tile.
	// When tiles change significantly, the major version MUST be increased.
	// Implementations MUST NOT use tiles with different major versions.
	Version string `json:"version,omitempty"`
	// OPTIONAL. Default: null. Contains a mustache template to be used to
	// format data from grids for interaction.
	// See https://github.com/mapbox/utfgrid-spec/tree/master/1.2
	// for the interactivity specification.
	Template string `json:"template,omitempty"`
	// OPTIONAL. Default: null. Contains a legend to be displayed with the map.
	// Implementations MAY decide to treat this as HTML or literal text.
	// For security reasons, make absolutely sure that this field can't be
	// abused as a vector for XSS or beacon tracking.
	Legend string `json:"legend,omitempty"`

	// vector layer details. This is not part of the tileJSON spec
	VectorLayers []VectorLayer `json:"vector_layers"`
}

type FieldType string

const (
	FieldTypeString  FieldType = "String"
	FieldTypeNumber  FieldType = "Number"
	FieldTypeBoolean FieldType = "Boolean"
)

// follow https://github.com/mapbox/mbtiles-spec/blob/master/1.3/spec.md
type VectorLayer struct {
	ID          string               `json:"id"`
	Fields      map[string]FieldType `json:"name"`
	Description string               `json:"description"`
	MinZoom     uint                 `json:"minzoom,omitempty"`
	MaxZoom     uint                 `json:"maxzoom,omitempty"`
}
