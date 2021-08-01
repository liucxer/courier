package pgbuilder

import (
	"github.com/liucxer/courier/sqlx"
	"github.com/liucxer/courier/sqlx/builder"
)

type Pager struct {
	Size   int64 `name:"size,omitempty" in:"query" default:"10" validate:"@int64[-1,]"`
	Offset int64 `name:"offset,omitempty" in:"query" default:"0" validate:"@int64[0,]"`
}

type Counter interface {
	SetCount(i int)
}

type WithCountExpr interface {
	CountExpr(db sqlx.DBExecutor) builder.SqlExpr
}

type WithTotal struct {
	Total int `json:"total"`
}

func (c *WithTotal) SetCount(i int) {
	c.Total = i
}
