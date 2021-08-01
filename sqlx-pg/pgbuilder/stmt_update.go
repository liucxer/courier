package pgbuilder

import (
	"context"
	"time"

	"github.com/liucxer/courier/sqlx/builder"
	"github.com/liucxer/courier/sqlx/datatypes"
)

func (s *Stmt) Update(model builder.Model, modifiers ...string) *StmtUpdate {
	return &StmtUpdate{
		stmt:      s,
		model:     model,
		modifiers: modifiers,
		rc:        &RecordCollection{},
	}
}

/**
[ WITH [ RECURSIVE ] with_query [, ...] ]
UPDATE [ ONLY ] table_name [ * ] [ [ AS ] alias ]
    SET { column_name = { expression | DEFAULT } |
          ( column_name [, ...] ) = [ ROW ] ( { expression | DEFAULT } [, ...] ) |
          ( column_name [, ...] ) = ( sub-SELECT )
        } [, ...]
    [ FROM from_list ]
    [ WHERE condition | WHERE CURRENT OF cursor_name ]
    [ RETURNING * | output_expression [ [ AS ] output_name ] [, ...] ]
*/
type StmtUpdate struct {
	stmt      *Stmt
	modifiers []string
	model     builder.Model
	from      builder.Model
	where     builder.SqlCondition
	rc        *RecordCollection
}

func (s *StmtUpdate) Do() error {
	if s.IsNil() {
		return nil
	}
	_, err := s.stmt.db.ExecExpr(s)
	return err
}

func (s *StmtUpdate) IsNil() bool {
	return s.stmt == nil || s.model == nil || s.rc == nil
}

func (s StmtUpdate) From(from builder.Model) *StmtUpdate {
	s.from = from
	return &s
}

func (s StmtUpdate) Where(where builder.SqlCondition) *StmtUpdate {
	s.where = where
	return &s
}

func (s StmtUpdate) SetBy(collect func(vc *RecordCollection), columns ...*builder.Column) *StmtUpdate {
	s.rc = RecordCollectionBy(collect, columns...)
	return &s
}

func (s StmtUpdate) SetWith(recordValues RecordValues, columns ...*builder.Column) *StmtUpdate {
	s.rc = RecordCollectionWith(recordValues, columns...)
	return &s
}

func (s StmtUpdate) SetFrom(model builder.Model, columnsCouldBeZeroValue ...*builder.Column) *StmtUpdate {
	s.rc = RecordCollectionFrom(s.stmt.db, model, columnsCouldBeZeroValue...)
	return &s
}

func (s *StmtUpdate) Returning(target builder.SqlExpr) CouldScan {
	return s.stmt.ReturningOf(s, target)
}

func (s *StmtUpdate) Ex(ctx context.Context) *builder.Ex {
	where := s.where

	// 已经删除不应该再次处理
	if modelWithDeleted, ok := s.model.(ModelWithDeletedAt); ok {
		where = builder.And(
			where,
			modelWithDeleted.FieldDeletedAt().Eq(0),
		)
	}

	rc := s.rc

	if modelWithUpdatedAt, ok := s.model.(ModelWithUpdatedAt); ok {
		// 补全更新时间
		if rc.Columns.F(modelWithUpdatedAt.FieldUpdatedAt().FieldName) == nil {
			rc = s.rc.WithExtendCol(modelWithUpdatedAt.FieldUpdatedAt(), datatypes.Timestamp(time.Now()))
		}
	}

	return s.stmt.ExprBy(func(ctx context.Context) *builder.Ex {
		e := builder.Expr("UPDATE")

		if len(s.modifiers) > 0 {
			for i := range s.modifiers {
				e.WriteByte(' ')
				e.WriteString(s.modifiers[i])
			}
		}

		e.WriteByte(' ')
		e.WriteExpr(s.stmt.db.T(s.model))

		e.WriteString(" SET ")

		builder.WriteAssignments(e, s.rc.AsAssignments()...)

		if s.from != nil {
			ctx = builder.ContextWithToggles(ctx, builder.Toggles{
				builder.ToggleMultiTable: true,
			})

			e.WriteString(" FROM ")
			e.WriteExpr(s.stmt.db.T(s.from))
		}

		builder.WriteAdditions(e, builder.Where(where))

		return e.Ex(ctx)
	}).Ex(ctx)
}
