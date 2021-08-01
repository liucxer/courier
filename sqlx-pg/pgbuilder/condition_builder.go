package pgbuilder

import (
	"github.com/liucxer/courier/sqlx"
	"github.com/liucxer/courier/sqlx/builder"
)

func ToCondition(db sqlx.DBExecutor, b ConditionBuilder) builder.SqlCondition {
	if b == nil {
		return builder.EmptyCond()
	}
	return b.ToCondition(db)
}

type ConditionBuilder interface {
	ToCondition(db sqlx.DBExecutor) builder.SqlCondition
}

func ConditionBuilderFromCondition(c builder.SqlCondition) ConditionBuilder {
	return ConditionBuilderBy(func(db sqlx.DBExecutor) builder.SqlCondition {
		return c
	})
}

func ConditionBuilderBy(build func(db sqlx.DBExecutor) builder.SqlCondition) ConditionBuilder {
	return &conditionBuilder{build: build}
}

type conditionBuilder struct {
	build func(db sqlx.DBExecutor) builder.SqlCondition
}

func (c *conditionBuilder) ToCondition(db sqlx.DBExecutor) builder.SqlCondition {
	return c.build(db)
}

func OneOf(builders ...ConditionBuilder) ConditionBuilder {
	return &conditionBuilderCompose{
		typ:      "or",
		builders: builders,
	}
}

func AllOf(builders ...ConditionBuilder) ConditionBuilder {
	return &conditionBuilderCompose{
		typ:      "all",
		builders: builders,
	}
}

type conditionBuilderCompose struct {
	typ      string
	builders []ConditionBuilder
}

func (c *conditionBuilderCompose) ToCondition(db sqlx.DBExecutor) builder.SqlCondition {
	where := builder.EmptyCond()

	for i := range c.builders {
		b := c.builders[i]
		if b == nil {
			continue
		}

		sub := b.ToCondition(db)

		if builder.IsNilExpr(sub) {
			continue
		}

		switch c.typ {
		case "or":
			where = where.Or(sub)
		case "all":
			where = where.And(sub)
		}
	}

	return where
}

type SubConditionBuilder interface {
	ConditionBuilder
	SelectFrom(db sqlx.DBExecutor) *StmtSelect
}

func SubSelect(target *builder.Column, subConditionBuilder SubConditionBuilder) ConditionBuilder {
	return ConditionBuilderBy(func(db sqlx.DBExecutor) builder.SqlCondition {
		if subConditionBuilder == nil {
			return nil
		}

		where := subConditionBuilder.ToCondition(db)
		if builder.IsNilExpr(where) {
			return nil
		}

		return target.In(subConditionBuilder.SelectFrom(db).Where(where))
	})
}
