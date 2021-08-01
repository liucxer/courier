## geom

[![GoDoc Widget](https://godoc.org/github.com/liucxer/courier/geography?status.svg)](https://godoc.org/github.com/liucxer/courier/geography)
[![Build Status](https://travis-ci.org/go-courier/geography.svg?branch=master)](https://travis-ci.org/go-courier/geography)
[![codecov](https://codecov.io/gh/go-courier/geography/branch/master/graph/badge.svg)](https://codecov.io/gh/go-courier/geography)
[![Go Report Card](https://goreportcard.com/badge/github.com/liucxer/courier/geography)](https://goreportcard.com/report/github.com/liucxer/courier/geography)

WIP

Geometry helpers base on .

```go
package database

import (
	"github.com/liucxer/courier/geography"
)

type Geom struct {
		ID              string                    `db:"F_id"`
    	Point           geography.Point           `db:"F_point,null"`
    	MultiPoint      geography.MultiPoint      `db:"F_multi_point,null"`
    	LineString      geography.LineString      `db:"F_line_string,null"`
    	MultiLineString geography.MultiLineString `db:"F_multi_line_string,null"`
    	Polygon         geography.Polygon         `db:"f_polygon,null"`
    	MultiPolygon    geography.MultiPolygon    `db:"f_multi_polygon,null"`
    	Geometry        geography.Geometry        `db:"f_geometry,null"`
}
```


* 3857 in db
* 4326 in go