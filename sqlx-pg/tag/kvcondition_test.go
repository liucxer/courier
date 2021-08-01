package tag

import (
	"context"
	"fmt"
	"testing"

	"github.com/liucxer/courier/kvcondition"
	"github.com/liucxer/courier/sqlx-pg/pgutils"
	"github.com/liucxer/courier/sqlx"
	"github.com/liucxer/courier/sqlx/builder"
	"github.com/liucxer/courier/sqlx/datatypes"
	"github.com/liucxer/courier/sqlx/postgresqlconnector"
	. "github.com/onsi/gomega"
)

var (
	postgresConnector = &postgresqlconnector.PostgreSQLConnector{
		Host:       "postgres://postgres@0.0.0.0:5432",
		Extra:      "sslmode=disable",
		Extensions: []string{"postgis"},
	}
	DBTest = sqlx.NewDatabase("test")
)

func TestFromKVCondition(t *testing.T) {
	db := DBTest.OpenDB(postgresConnector)

	ql, err := kvcondition.ParseKVCondition([]byte(`areaLevel = CITY & areaCode = 140100000 | 12312 = 123`))
	NewWithT(t).Expect(err).To(BeNil())

	p := &ProjectTag{}

	s := SelectByKVCondition(db, *ql, p, p.FieldProjectID())

	e := s.Ex(context.Background())

	sql, err := pgutils.InterpolateParams(e)
	NewWithT(t).Expect(err).To(BeNil())
	fmt.Println(sql)
}

var ProjectTagTable = DBTest.Register(&ProjectTag{})

type ProjectTag struct {
	ProjectID uint64 `db:"f_project_id"`
	Tag

	UpdatedAt datatypes.Timestamp `db:"f_updated_at"`
	CreatedAt datatypes.Timestamp `db:"f_created_at"`
}

func (ProjectTag) UniqueIndexITag() string {
	return "i_tag"
}

func (ProjectTag) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"i_tag": []string{
			"ProjectID",
			"Key",
			"Value",
		},
	}
}

func (ProjectTag) TableName() string {
	return "t_project_tag"
}

func (ProjectTag) FieldProjectID() *builder.Column {
	return ProjectTagTable.F("ProjectID")
}

func (ProjectTag) FieldKey() *builder.Column {
	return ProjectTagTable.F("Key")
}

func (ProjectTag) FieldValue() *builder.Column {
	return ProjectTagTable.F("Value")
}

func (ProjectTag) FieldCreatedAt() *builder.Column {
	return ProjectTagTable.F("CreatedAt")
}

func (p ProjectTag) FieldUpdatedAt() *builder.Column {
	return ProjectTagTable.F("UpdatedAt")
}
