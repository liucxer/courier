package pgbuilder

import (
	"context"
	"strings"

	"github.com/liucxer/courier/sqlx"
	"github.com/liucxer/courier/sqlx/builder"
)

func Use(db sqlx.DBExecutor) *Stmt {
	return &Stmt{db: db}
}

type CouldScan interface {
	Scan(v interface{}) error
}

type Stmt struct {
	db   sqlx.DBExecutor
	with *StmtWith
}

func (s *Stmt) ExprBy(do func(ctx context.Context) *builder.Ex) builder.SqlExpr {
	return builder.ExprBy(func(ctx context.Context) *builder.Ex {
		e := builder.Expr("")
		if !s.with.IsNil() {
			e.WriteExpr(s.with)
			e.WriteByte('\n')
		}

		e.WriteExpr(do(ctx))

		return e.Ex(ctx)
	})
}

const PrimaryKey = "pk"

const (
	toggleKeyIgnoreDeletedAt = "$$IgnoreDeletedAt$$"
)

func ContextWithIgnoreDeletedAt(ctx context.Context) context.Context {
	return builder.ContextWithToggles(ctx, builder.Toggles{
		toggleKeyIgnoreDeletedAt: true,
	})
}

func (s *Stmt) conflictColumns(model builder.Model, indexKey string) *builder.Columns {
	if strings.ToLower(indexKey) == PrimaryKey {
		if modelWithPrimaryKey, ok := model.(ModelWithPrimaryKey); ok {
			return s.T(model).MustFields(modelWithPrimaryKey.PrimaryKey()...)
		}
		return nil
	}

	if modelWithUniqueIndexes, ok := model.(ModelWithUniqueIndexes); ok {
		if fieldNames, ok := modelWithUniqueIndexes.UniqueIndexes()[indexKey]; ok {
			return s.T(model).MustFields(fieldNames...)
		}
	}
	return nil
}

func (s *Stmt) T(model builder.Model) *builder.Table {
	return s.db.T(model)
}

type ModelWithDeletedAt interface {
	builder.Model
	FieldDeletedAt() *builder.Column
}

type ModelWithCreatedAt interface {
	builder.Model
	FieldCreatedAt() *builder.Column
}

type ModelWithUpdatedAt interface {
	builder.Model
	FieldUpdatedAt() *builder.Column
}

type ModelWithUniqueIndexes interface {
	builder.Model
	UniqueIndexes() builder.Indexes
}

type ModelWithPrimaryKey interface {
	builder.Model
	PrimaryKey() []string
}

func (s *Stmt) ReturningOf(expr builder.SqlExpr, target builder.SqlExpr) CouldScan {
	return &returning{
		stmt:   s,
		expr:   expr,
		target: target,
	}
}

type returning struct {
	stmt   *Stmt
	expr   builder.SqlExpr
	target builder.SqlExpr
}

func (r *returning) Scan(v interface{}) error {
	return r.stmt.db.QueryExprAndScan(r, v)
}

func (r *returning) IsNil() bool {
	return r == nil || builder.IsNilExpr(r.expr)
}

func (r *returning) Ex(ctx context.Context) *builder.Ex {
	e := builder.Expr("")
	e.WriteExpr(r.expr)

	e.WriteString(" RETURNING ")

	if builder.IsNilExpr(r.target) {
		e.WriteString("*")
	} else {
		e.WriteExpr(r.target)
	}

	return e.Ex(ctx)
}
