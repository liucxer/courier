package pgbuilder

import (
	"context"

	"github.com/liucxer/courier/sqlx"
	"github.com/liucxer/courier/sqlx/builder"
)

type ModalAndWithQuery interface {
	builder.Model
	WithQuery(db sqlx.DBExecutor) builder.SqlExpr
}

type WithRecursive interface {
	withRecursive()
}

func (s Stmt) With(model ModalAndWithQuery, others ...ModalAndWithQuery) *Stmt {
	if _, ok := model.(WithRecursive); ok {
		s.with = &StmtWith{stmt: &s, models: []ModalAndWithQuery{model}, recursive: true}
		return &s
	}
	s.with = &StmtWith{stmt: &s, models: append([]ModalAndWithQuery{model}, others...)}
	return &s
}

type BuildSubExpr func() builder.SqlExpr

type StmtWith struct {
	stmt      *Stmt
	recursive bool
	models    []ModalAndWithQuery
}

func (s *StmtWith) IsNil() bool {
	return s == nil || len(s.models) == 0
}

func (s *StmtWith) Ex(ctx context.Context) *builder.Ex {
	e := builder.Expr("WITH ")

	if s.recursive {
		e.WriteString("RECURSIVE ")
	}

	for i := range s.models {
		if i > 0 {
			e.WriteString(", ")
		}

		model := s.models[i]

		table := s.stmt.T(model)

		e.WriteExpr(table)
		e.WriteGroup(func(e *builder.Ex) {
			e.WriteExpr(&table.Columns)
		})

		e.WriteString(" AS ")

		e.WriteGroup(func(e *builder.Ex) {
			e.WriteByte('\n')
			e.WriteExpr(model.WithQuery(s.stmt.db))
			e.WriteByte('\n')
		})
	}

	return e.Ex(ctx)
}

func AsWithRecursiveQuery(model builder.Model, withQueryFn func(db sqlx.DBExecutor) builder.SqlExpr) ModalAndWithQuery {
	return &withRecursiveQuery{
		withQuery: withQuery{
			Model:     model,
			withQuery: withQueryFn,
		},
	}
}

type withRecursiveQuery struct {
	WithRecursive
	withQuery
}

func AsWithQuery(model builder.Model, withQueryFn func(db sqlx.DBExecutor) builder.SqlExpr) ModalAndWithQuery {
	return &withQuery{
		Model:     model,
		withQuery: withQueryFn,
	}
}

type withQuery struct {
	builder.Model
	withQuery func(db sqlx.DBExecutor) builder.SqlExpr
}

func (m *withQuery) T() *builder.Table {
	if td, ok := m.Model.(builder.TableDefinition); ok {
		return td.T()
	}

	if t, ok := m.Model.(*builder.Table); ok {
		return t
	}

	return nil
}

func (m *withQuery) WithQuery(db sqlx.DBExecutor) builder.SqlExpr {
	return m.withQuery(db)
}
