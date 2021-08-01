package pgbuilder

import (
	"context"

	"github.com/liucxer/courier/sqlx/builder"
)

func (s *Stmt) Insert() *StmtInsert {
	return &StmtInsert{
		stmt: s,
	}
}

/**
[ WITH [ RECURSIVE ] with_query [, ...] ]
INSERT INTO table_name [ AS alias ] [ ( column_name [, ...] ) ]
    [ OVERRIDING { SYSTEM | USER} VALUE ]
    { DEFAULT VALUES | VALUES ( { expression | DEFAULT } [, ...] ) [, ...] | query }
    [ ON CONFLICT [ conflict_target ] conflict_action ]
    [ RETURNING * | output_expression [ [ AS ] output_name ] [, ...] ]

where conflict_target can be one of:

    ( { index_column_name | ( index_expression ) } [ COLLATE collation ] [ opclass ] [, ...] ) [ WHERE index_predicate ]
    ON CONSTRAINT constraint_name

and conflict_action is one of:

    DO NOTHING
    DO UPDATE SET { column_name = { expression | DEFAULT } |
                    ( column_name [, ...] ) = [ ROW ] ( { expression | DEFAULT } [, ...] ) |
                    ( column_name [, ...] ) = ( sub-SELECT )
                  } [, ...]
              [ WHERE condition ]
*/
type StmtInsert struct {
	stmt *Stmt

	model builder.Model

	additions []builder.Addition

	rc *RecordCollection
}

func (s *StmtInsert) Do() error {
	if s.IsNil() {
		return nil
	}
	_, err := s.stmt.db.ExecExpr(s)
	return err
}

func (s *StmtInsert) IsNil() bool {
	return s == nil || s.stmt == nil || s.model == nil || s.rc == nil
}

func (s StmtInsert) Into(model builder.Model, additions ...builder.Addition) *StmtInsert {
	s.model = model
	s.additions = additions
	return &s
}

func (s *StmtInsert) Returning(target builder.SqlExpr) CouldScan {
	return s.stmt.ReturningOf(s, target)
}

func (s StmtInsert) ValuesBy(collect func(vc *RecordCollection), columns ...*builder.Column) *StmtInsert {
	s.rc = RecordCollectionBy(collect, columns...)
	return &s
}

func (s StmtInsert) ValuesWith(recordValues RecordValues, columns ...*builder.Column) *StmtInsert {
	s.rc = RecordCollectionWith(recordValues, columns...)
	return &s
}

func (s StmtInsert) ValuesFrom(model builder.Model, columnsCouldBeZeroValue ...*builder.Column) *StmtInsert {
	s.rc = RecordCollectionFrom(s.stmt.db, model, columnsCouldBeZeroValue...)
	return &s
}

func (s StmtInsert) OnConflictDoNothing(indexKey string) *StmtInsert {
	conflictColumns := s.stmt.conflictColumns(s.model, indexKey)
	if conflictColumns == nil {
		return &s
	}
	onConflict := builder.OnConflict(conflictColumns).DoNothing()
	s.additions = append(append(builder.Additions{}, s.additions...), onConflict)
	return &s
}

func (s StmtInsert) OnConflictDoUpdateSet(indexKey string, excludedColumns ...*builder.Column) *StmtInsert {
	conflictColumns := s.stmt.conflictColumns(s.model, indexKey)
	if conflictColumns == nil {
		return &s
	}
	onConflict := builder.OnConflict(conflictColumns).DoUpdateSet(ExcludedFields(excludedColumns...)...)
	s.additions = append(append(builder.Additions{}, s.additions...), onConflict)
	return &s
}

func (s *StmtInsert) Ex(ctx context.Context) *builder.Ex {
	return s.stmt.ExprBy(func(ctx context.Context) *builder.Ex {
		return builder.Insert().
			Into(s.stmt.T(s.model), s.additions...).
			Values(s.rc.Columns, s.rc.Values()...).
			Ex(ctx)
	}).Ex(ctx)
}

func Excluded(f *builder.Column) builder.SqlExpr {
	return builder.Expr("EXCLUDED.?", f)
}

func ExcludedFields(fields ...*builder.Column) builder.Assignments {
	assignments := builder.Assignments{}

	for i := range fields {
		f := fields[i]
		if f != nil {
			assignments = append(assignments, f.ValueBy(Excluded(f)))
		}
	}

	return assignments
}
