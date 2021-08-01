package tag

import (
	"time"

	"github.com/liucxer/courier/kvcondition"
	"github.com/liucxer/courier/sqlx-pg/pgbuilder"
	"github.com/liucxer/courier/sqlx"
	"github.com/liucxer/courier/sqlx/builder"
	"github.com/liucxer/courier/sqlx/datatypes"
)

func TaggerFor(modelWithTag ModelWithTag) *Tagger {
	return &Tagger{modelWithTag: modelWithTag, targetFieldName: modelWithTag.UniqueIndexes()[modelWithTag.UniqueIndexITag()][0]}
}

type Tagger struct {
	modelWithTag    ModelWithTag
	targetFieldName string
}

func (t *Tagger) Tag(db sqlx.DBExecutor, targetID uint64, tags Tags) error {
	modelWithTag := t.modelWithTag

	table := db.T(modelWithTag)

	now := datatypes.Timestamp(time.Now())

	expr := pgbuilder.Use(db).
		Insert().Into(modelWithTag).
		OnConflictDoNothing(modelWithTag.UniqueIndexITag()).
		ValuesBy(
			func(vc *pgbuilder.RecordCollection) {
				for key, values := range tags {
					for i := range values {
						vc.SetRecordValues(
							targetID,
							key,
							values[i],
							now,
							now,
						)
					}
				}
			},
			table.F(t.targetFieldName),
			modelWithTag.FieldKey(),
			modelWithTag.FieldValue(),
			modelWithTag.FieldCreatedAt(),
			modelWithTag.FieldUpdatedAt(),
		)

	return expr.Do()
}

func (t *Tagger) UnTag(db sqlx.DBExecutor, targetID uint64, keyAndValues ...string) error {
	modelWithTag := t.modelWithTag

	table := db.T(modelWithTag)

	where := table.F(t.targetFieldName).Eq(targetID)

	if len(keyAndValues) > 0 {
		where = builder.And(where, modelWithTag.FieldKey().Eq(keyAndValues[0]))

		if len(keyAndValues) > 1 {
			where = builder.And(where, modelWithTag.FieldValue().In(keyAndValues[1:]))
		}
	}

	return pgbuilder.Use(db).Delete(modelWithTag).Where(where).Do()
}

func (t *Tagger) GetTags(db sqlx.DBExecutor, targetID uint64) (Tags, error) {
	tags := Tags{}

	where := db.T(t.modelWithTag).F(t.targetFieldName).Eq(targetID)

	err := pgbuilder.Use(db).Select(nil).From(t.modelWithTag).Where(where).List(tags, nil)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (t *Tagger) SelectFor(db sqlx.DBExecutor, tagql kvcondition.KVCondition) *pgbuilder.StmtSelect {
	modelWithTag := t.modelWithTag

	stmtSelect := SelectByKVCondition(db, tagql, modelWithTag, db.T(modelWithTag).F(t.targetFieldName))
	if stmtSelect == nil {
		return nil
	}

	return stmtSelect
}

type ModelWithTag interface {
	builder.Model

	FieldKey() *builder.Column
	FieldValue() *builder.Column

	FieldCreatedAt() *builder.Column
	FieldUpdatedAt() *builder.Column

	UniqueIndexes() builder.Indexes
	UniqueIndexITag() string
}
