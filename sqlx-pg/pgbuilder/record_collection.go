package pgbuilder

import (
	"fmt"

	"github.com/liucxer/courier/sqlx"
	"github.com/liucxer/courier/sqlx/builder"
)

func RecordCollectionFrom(db sqlx.DBExecutor, model builder.Model, columnsCouldBeZeroValue ...*builder.Column) *RecordCollection {
	t := db.T(model)

	zeroFields := make([]string, 0)
	for i := range columnsCouldBeZeroValue {
		zeroFields = append(zeroFields, columnsCouldBeZeroValue[i].FieldName)
	}

	fieldValues := sqlx.FieldValuesFromModel(t, model, zeroFields...)

	columnList := make([]*builder.Column, 0)
	recordValueList := make([]interface{}, 0)

	t.Columns.Range(func(col *builder.Column, idx int) {
		if v, ok := fieldValues[col.FieldName]; ok {
			columnList = append(columnList, col)
			recordValueList = append(recordValueList, v)
		}
	})

	return RecordCollectionBy(func(rc *RecordCollection) {
		rc.SetRecordValues(recordValueList...)
	}, columnList...)
}

func RecordCollectionWith(recordValues RecordValues, columns ...*builder.Column) *RecordCollection {
	return RecordCollectionBy(func(rc *RecordCollection) {
		rc.SetRecordValues(recordValues...)
	}, columns...)
}

func RecordCollectionBy(collect func(rc *RecordCollection), columns ...*builder.Column) *RecordCollection {
	cols := &builder.Columns{}
	for i := range columns {
		col := columns[i]
		if col == nil {
			panic(fmt.Errorf("invalid %d of columns", i))
		}
		cols.Add(col)
	}

	rc := &RecordCollection{
		Columns: cols,
		records: []RecordValues{},
	}

	collect(rc)

	return rc
}

type RecordValues []interface{}

type RecordCollection struct {
	records []RecordValues
	Columns *builder.Columns
}

func (vc *RecordCollection) IsNil() bool {
	return vc == nil || len(vc.records) == 0 || builder.IsNilExpr(vc.Columns)
}

func (vc *RecordCollection) SetRecordValues(values ...interface{}) {
	if len(values) == 1 {
		if _, ok := values[0].(builder.SelectStatement); !ok {
			if len(values) != vc.Columns.Len() {
				panic(fmt.Errorf("len of records is not matched, need %d, got %d", vc.Columns.Len(), len(values)))
			}
		}
	}

	vc.records = append(vc.records, values)
}

func (vc RecordCollection) WithExtendCol(col *builder.Column, val interface{}) *RecordCollection {
	columns := vc.Columns.Clone()
	columns.Add(col)

	records := make([]RecordValues, len(vc.records))

	for i := range records {
		records[i] = append(vc.records[i], val)
	}

	return &RecordCollection{
		Columns: columns,
		records: records,
	}
}

func (vc *RecordCollection) AsAssignments() builder.Assignments {
	if len(vc.records) == 0 {
		return nil
	}

	assignments := builder.Assignments{}

	for j := range vc.records {
		record := vc.records[j]

		vc.Columns.Range(func(col *builder.Column, idx int) {
			assignments = append(assignments, col.ValueBy(record[idx]))
		})
	}

	return assignments
}

func (vc *RecordCollection) Values() []interface{} {
	if len(vc.records) == 0 {
		return nil
	}

	values := make([]interface{}, 0)

	for j := range vc.records {
		recordValues := vc.records[j]

		values = append(values, recordValues...)
	}

	return values
}
