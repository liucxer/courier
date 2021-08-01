package pgbuilder

import (
	"github.com/liucxer/courier/sqlx"
	"github.com/liucxer/courier/sqlx/builder"
)

type Receiver func(v interface{})

type Collection struct {
	new       func() interface{}
	id        func(v interface{}) interface{}
	condition func(db sqlx.DBExecutor, ids []interface{}) builder.SqlCondition
	ids       []interface{}
	receivers map[interface{}][]Receiver
}

func (c *Collection) New() interface{} {
	return c.new()
}

func (c *Collection) Init(new func() interface{}, id func(v interface{}) interface{}, condition func(db sqlx.DBExecutor, ids []interface{}) builder.SqlCondition) {
	if c.receivers == nil {
		c.receivers = map[interface{}][]Receiver{}
		c.id = id
		c.condition = condition
		c.new = new
	}
}

func (c *Collection) OnNext(id interface{}, receivers ...Receiver) {
	if len(receivers) > 0 {
		c.ids = append(c.ids, id)
		c.receivers[id] = append(c.receivers[id], receivers...)
	}
}

func (c *Collection) ToCondition(db sqlx.DBExecutor) builder.SqlCondition {
	if c.receivers == nil || len(c.ids) == 0 {
		return builder.EmptyCond()
	}

	return c.condition(db, c.ids)
}

func (c *Collection) Next(v interface{}) error {
	id := c.id(v)

	if receivers, ok := c.receivers[id]; ok {
		for i := range receivers {
			receivers[i](v)
		}
	}

	return nil
}

func (c *Collection) DoList(db sqlx.DBExecutor, pager *Pager, additions ...builder.Addition) error {
	if c.receivers == nil {
		return nil
	}

	return Use(db).
		Select(nil).
		From(c.new().(builder.Model)).
		Where(c.ToCondition(db), additions...).
		List(c, pager)
}
