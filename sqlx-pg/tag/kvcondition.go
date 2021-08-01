package tag

import (
	"github.com/liucxer/courier/kvcondition"
	"github.com/liucxer/courier/sqlx-pg/pgbuilder"
	"github.com/liucxer/courier/sqlx"
	"github.com/liucxer/courier/sqlx/builder"
)

func SelectByKVCondition(db sqlx.DBExecutor, kvc kvcondition.KVCondition, modelWithTag ModelWithTag, idCol *builder.Column) *pgbuilder.StmtSelect {
	if kvc.Node == nil {
		return nil
	}

	s := pgbuilder.Use(db)

	keyCol := modelWithTag.FieldKey()
	valueCol := modelWithTag.FieldValue()

	var visit func(node kvcondition.Node) builder.SqlCondition

	visit = func(node kvcondition.Node) builder.SqlCondition {
		switch n := node.(type) {
		case *kvcondition.Condition:
			switch n.Operator {
			case kvcondition.ConditionOperatorOR:
				return builder.Or(
					idCol.In(s.Select(idCol).From(modelWithTag).Where(visit(n.Left))),
					visit(n.Right),
				)
			case kvcondition.ConditionOperatorAND:
				return builder.And(
					idCol.In(s.Select(idCol).From(modelWithTag).Where(visit(n.Left))),
					visit(n.Right),
				)
			}
			return nil
		case *kvcondition.Rule:
			k := string(n.Key)
			v := string(n.Value)

			switch n.Operator {
			case kvcondition.OperatorExists:
				return keyCol.Eq(k)
			case kvcondition.OperatorEqual:
				return builder.And(
					keyCol.Eq(k),
					valueCol.Eq(v),
				)
			case kvcondition.OperatorNotEqual:
				return builder.And(
					keyCol.Eq(k),
					valueCol.Neq(v),
				)
			case kvcondition.OperatorContains:
				return builder.And(
					keyCol.Eq(k),
					valueCol.Like(v),
				)
			case kvcondition.OperatorStartsWith:
				return builder.And(
					keyCol.Eq(k),
					valueCol.RightLike(v),
				)
			case kvcondition.OperatorEndsWith:
				return builder.And(
					keyCol.Eq(k),
					valueCol.LeftLike(v),
				)
			}
			return nil
		}
		return nil
	}

	return s.Select(idCol).From(modelWithTag).Where(visit(kvc.Node))
}
