package pgbuilder

import (
	"context"

	"github.com/liucxer/courier/sqlx"
	"github.com/liucxer/courier/sqlx/builder"
)

func (s *Stmt) Select(target builder.SqlExpr) *StmtSelect {
	return &StmtSelect{
		stmt:   s,
		target: target,
	}
}

/**
[ WITH [ RECURSIVE ] with_query [, ...] ]
SELECT [ ALL | DISTINCT [ ON ( expression [, ...] ) ] ]
    [ * | expression [ [ AS ] output_name ] [, ...] ]
    [ FROM from_item [, ...] ]
    [ WHERE condition ]
    [ GROUP BY grouping_element [, ...] ]
    [ HAVING condition [, ...] ]
    [ WINDOW window_name AS ( window_definition ) [, ...] ]
    [ { UNION | INTERSECT | EXCEPT } [ ALL | DISTINCT ] select ]
    [ ORDER BY expression [ ASC | DESC | USING operator ] [ NULLS { FIRST | LAST } ] [, ...] ]
    [ LIMIT { count | ALL } ]
    [ OFFSET start [ ROW | ROWS ] ]
    [ FETCH { FIRST | NEXT } [ count ] { ROW | ROWS } ONLY ]
    [ FOR { UPDATE | NO KEY UPDATE | SHARE | KEY SHARE } [ OF table_name [, ...] ] [ NOWAIT | SKIP LOCKED ] [...] ]

where from_item can be one of:

    [ ONLY ] table_name [ * ] [ [ AS ] alias [ ( column_alias [, ...] ) ] ]
                [ TABLESAMPLE sampling_method ( argument [, ...] ) [ REPEATABLE ( seed ) ] ]
    [ LATERAL ] ( select ) [ AS ] alias [ ( column_alias [, ...] ) ]
    with_query_name [ [ AS ] alias [ ( column_alias [, ...] ) ] ]
    [ LATERAL ] function_name ( [ argument [, ...] ] )
                [ WITH ORDINALITY ] [ [ AS ] alias [ ( column_alias [, ...] ) ] ]
    [ LATERAL ] function_name ( [ argument [, ...] ] ) [ AS ] alias ( column_definition [, ...] )
    [ LATERAL ] function_name ( [ argument [, ...] ] ) AS ( column_definition [, ...] )
    [ LATERAL ] ROWS FROM( function_name ( [ argument [, ...] ] ) [ AS ( column_definition [, ...] ) ] [, ...] )
                [ WITH ORDINALITY ] [ [ AS ] alias [ ( column_alias [, ...] ) ] ]
    from_item [ NATURAL ] join_type from_item [ ON join_condition | USING ( join_column [, ...] ) ]

and grouping_element can be one of:

    ( )
    expression
    ( expression [, ...] )
    ROLLUP ( { expression | ( expression [, ...] ) } [, ...] )
    CUBE ( { expression | ( expression [, ...] ) } [, ...] )
    GROUPING SETS ( grouping_element [, ...] )

and with_query is:

    with_query_name [ ( column_name [, ...] ) ] AS ( select | values | insert | update | delete )

TABLE [ ONLY ] table_name [ * ]
*/
type StmtSelect struct {
	builder.SelectStatement

	stmt *Stmt

	target builder.SqlExpr
	from   builder.Model
	where  builder.SqlCondition

	additions []builder.Addition
}

func (s *StmtSelect) IsNil() bool {
	return s == nil || s.stmt == nil || s.from == nil
}

func (s StmtSelect) From(model builder.Model) *StmtSelect {
	s.from = model
	return &s
}

func (s StmtSelect) Join(target builder.Model, joinCondition builder.SqlCondition) *StmtSelect {
	s.additions = append(append(builder.Additions{}, s.additions...), builder.Join(s.stmt.T(target)).On(joinCondition))
	return &s
}

func (s StmtSelect) CrossJoin(target builder.Model) *StmtSelect {
	s.additions = append(append(builder.Additions{}, s.additions...), builder.CrossJoin(s.stmt.T(target)))
	return &s
}

func (s StmtSelect) LeftJoin(target builder.Model, joinCondition builder.SqlCondition) *StmtSelect {
	s.additions = append(append(builder.Additions{}, s.additions...), builder.LeftJoin(s.stmt.T(target)).On(joinCondition))
	return &s
}

func (s StmtSelect) RightJoin(target builder.Model, joinCondition builder.SqlCondition) *StmtSelect {
	s.additions = append(append(builder.Additions{}, s.additions...), builder.RightJoin(s.stmt.T(target)).On(joinCondition))
	return &s
}

func (s StmtSelect) FullJoin(target builder.Model, joinCondition builder.SqlCondition) *StmtSelect {
	s.additions = append(append(builder.Additions{}, s.additions...), builder.FullJoin(s.stmt.T(target)).On(joinCondition))
	return &s
}

func (s StmtSelect) Where(where builder.SqlCondition, additions ...builder.Addition) *StmtSelect {
	s.where = where
	s.additions = append(append(builder.Additions{}, s.additions...), additions...)
	return &s
}

func (s *StmtSelect) Ex(ctx context.Context) *builder.Ex {
	where := s.where

	if !builder.TogglesFromContext(ctx).Is(toggleKeyIgnoreDeletedAt) {
		if modelWithDeleted, ok := s.from.(ModelWithDeletedAt); ok {
			where = builder.And(where, modelWithDeleted.FieldDeletedAt().Eq(0))
		}
	}

	finalAdditions := builder.Additions{
		builder.Where(where),
	}

	return s.stmt.ExprBy(func(ctx context.Context) *builder.Ex {
		return builder.
			Select(s.target).
			From(
				s.stmt.T(s.from),
				append(finalAdditions, s.additions...)...,
			).
			Ex(ctx)
	}).Ex(ctx)
}

func (s *StmtSelect) Scan(v interface{}) error {
	return s.stmt.db.QueryExprAndScan(s, v)
}

func (s *StmtSelect) List(list sqlx.ScanIterator, pager *Pager) error {
	if pager == nil {
		pager = &Pager{
			Size:   -1,
			Offset: 0,
		}
	}

	total := -1

	otherAdditions := builder.Additions{}

	if pager.Size != -1 {
		otherAdditions = append(otherAdditions, builder.Limit(pager.Size).Offset(pager.Offset))
	}

	if err := s.Where(s.where, otherAdditions...).Scan(list); err != nil {
		return err
	}

	if pager.Size != -1 {
		targetForCount := builder.SqlExpr(builder.Count())

		if withCountExpr, ok := list.(WithCountExpr); ok {
			targetForCount = withCountExpr.CountExpr(s.stmt.db)
		}

		if err := s.stmt.
			Select(targetForCount).
			From(s.from).
			Where(s.where, FilterAdditions(s.additions, func(a builder.Addition) bool {
				return !builder.IsNilExpr(a) && a.AdditionType() != builder.AdditionOrderBy
			})...).Scan(&total); err != nil {
			return err
		}
	}

	if counter, ok := list.(Counter); ok {
		counter.SetCount(total)
	}

	return nil
}

func FilterAdditions(additions builder.Additions, filter func(a builder.Addition) bool) builder.Additions {
	finalAdditions := builder.Additions{}
	for i := range additions {
		if filter(additions[i]) {
			finalAdditions = append(finalAdditions, additions[i])
		}
	}
	return finalAdditions
}
