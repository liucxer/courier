package pgbuilder_test

import (
	"testing"
	"time"

	"github.com/liucxer/courier/sqlx-pg/pgbuilder"
	"github.com/liucxer/courier/sqlx"
	"github.com/liucxer/courier/sqlx/builder"
	"github.com/liucxer/courier/sqlx/datatypes"
	"github.com/liucxer/courier/sqlx/migration"
	"github.com/liucxer/courier/sqlx/postgresqlconnector"
	"github.com/google/uuid"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var (
	postgresConnector = &postgresqlconnector.PostgreSQLConnector{
		Host:       "postgres://postgres@0.0.0.0:5432",
		Extra:      "sslmode=disable",
		Extensions: []string{"postgis"},
	}

	DB            = sqlx.NewFeatureDatabase("test_for_pg_builder")
	TableUser     = DB.Register(&User{})
	TableUserRole = DB.Register(&UserRole{})
	db            sqlx.DBExecutor
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)

	db = DB.OpenDB(postgresConnector)

	if err := migration.Migrate(db, nil); err != nil {
		panic(err)
	}
}

func TestStmt(t *testing.T) {
	t.Run("insert single", func(t *testing.T) {
		user := &User{
			Name: "test",
		}

		err := pgbuilder.Use(db).
			Insert().Into(&User{}).
			ValuesFrom(user).
			OnConflictDoUpdateSet("i_name", TableUser.F("Name")).
			Returning(builder.Expr("*")).
			Scan(user)

		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(user.ID).NotTo(Equal(0))
	})

	t.Run("update simple", func(t *testing.T) {
		err := pgbuilder.Use(db).
			Update(&User{}).
			SetWith(pgbuilder.RecordValues{uuid.New()}, TableUser.F("Name")).
			Where(TableUser.F("Name").Eq("test")).
			Do()

		NewWithT(t).Expect(err).To(BeNil())
	})

	t.Run("insert multi", func(t *testing.T) {
		err := pgbuilder.Use(db).
			Insert().Into(&User{}).
			ValuesBy(
				func(vc *pgbuilder.RecordCollection) {
					for i := 0; i < 100; i++ {
						vc.SetRecordValues(
							uuid.New(),
						)
					}
				},
				TableUser.F("Name"),
			).
			OnConflictDoNothing(pgbuilder.PrimaryKey).
			Do()

		NewWithT(t).Expect(err).To(BeNil())
	})

	t.Run("insert from select", func(t *testing.T) {
		stmt := pgbuilder.Use(db).
			Insert().Into(&User{}).
			ValuesWith(
				pgbuilder.RecordValues{
					pgbuilder.Use(db).Select(TableUser.MustFields("Name", "Age")).From(&User{}),
				},
				TableUser.F("Name"),
				TableUser.F("Age"),
			).
			OnConflictDoNothing("i_name")

		NewWithT(t).Expect(stmt.Do()).To(BeNil())
	})

	t.Run("list", func(t *testing.T) {
		dataList := &UserDataList{}

		err := dataList.DoList(db, &pgbuilder.Pager{Size: 10})

		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(len(dataList.Data) >= 1).To(BeTrue())
		NewWithT(t).Expect(dataList.Total >= 1).To(BeTrue())
	})

	t.Run("with Select", func(t *testing.T) {
		count := 0

		vUser := builder.T("v_user", builder.Col("f_name"), builder.Col("f_age"))

		err := pgbuilder.Use(db).
			With(pgbuilder.AsWithQuery(vUser, func(db sqlx.DBExecutor) builder.SqlExpr {
				return pgbuilder.Use(db).
					Select(builder.MultiMayAutoAlias(
						vUser.Col("f_name"),
						vUser.Col("f_age"),
					)).From(&User{}).
					Where(vUser.Col("f_age").Gt(1))
			})).
			Select(builder.Count()).From(vUser).
			Scan(&count)

		NewWithT(t).Expect(err).To(BeNil())
		NewWithT(t).Expect(count > 0).To(BeTrue())
	})

	t.Run("update from", func(t *testing.T) {
		err := pgbuilder.Use(db).
			Update(&User{}).
			From(&UserRole{}).
			SetWith(
				pgbuilder.RecordValues{
					(&UserRole{}).FieldUpdatedAt(),
				},
				(&User{}).FieldUpdatedAt(),
			).
			Where((&UserRole{}).FieldCreatedAt().Eq((&User{}).FieldCreatedAt())).
			Do()

		NewWithT(t).Expect(err).To(BeNil())
	})

	t.Run("delete soft", func(t *testing.T) {
		err := pgbuilder.
			Use(db).
			Delete(&User{}).
			Do()

		NewWithT(t).Expect(err).To(BeNil())
	})

	t.Run("delete ignore deletedAt", func(t *testing.T) {
		err := pgbuilder.Use(db.WithContext(pgbuilder.ContextWithIgnoreDeletedAt(db.Context()))).
			Delete(&User{}).
			Do()

		NewWithT(t).Expect(err).To(BeNil())
	})
}

type UserParams struct {
	Names []string `name:"name" in:"query"`
	Ages  []int    `name:"age" in:"query"`
}

func (u *UserParams) ToCondition(db sqlx.DBExecutor) builder.SqlCondition {
	where := builder.EmptyCond()

	if len(u.Names) > 0 {
		where = where.And(TableUser.F("Name").In(u.Names))
	}

	if len(u.Ages) > 0 {
		where = where.And(TableUser.F("Age").In(u.Ages))
	}

	return where
}

type UserDataList struct {
	UserParams `json:"-"`
	Data       []*User `json:"data"`
	pgbuilder.WithTotal
}

func (UserDataList) New() interface{} {
	return &User{}
}

func (u *UserDataList) Next(v interface{}) error {
	u.Data = append(u.Data, v.(*User))
	return nil
}

func (u *UserDataList) DoList(db sqlx.DBExecutor, pager *pgbuilder.Pager, additions ...builder.Addition) error {
	return pgbuilder.Use(db).Select(nil).
		From(&User{}).
		Where(u.ToCondition(db), additions...).
		List(u, pager)
}

type User struct {
	ID   uint64 `db:"f_id,autoincrement"`
	Name string `db:"f_name,size=255,default=''"`
	Age  int    `db:"f_age,default='18'"`

	OperationTimesWithDeletedAt
}

func (user *User) TableName() string {
	return "t_user"
}

func (user *User) PrimaryKey() []string {
	return []string{"ID"}
}

func (user *User) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"i_name": {"Name"},
	}
}

func (User) FieldDeletedAt() *builder.Column {
	return TableUser.F("DeletedAt")
}

func (User) FieldCreatedAt() *builder.Column {
	return TableUser.F("CreatedAt")
}

func (User) FieldUpdatedAt() *builder.Column {
	return TableUser.F("UpdatedAt")
}

type UserRole struct {
	ID     uint64 `db:"f_id,autoincrement"`
	UserID uint64 `db:"f_user_id"`
	OperationTimes
}

func (UserRole) TableName() string {
	return "t_user_role"
}

func (UserRole) PrimaryKey() []string {
	return []string{"ID"}
}

func (UserRole) Indexes() builder.Indexes {
	return builder.Indexes{
		"i_user_role": {"UserID"},
	}
}

func (UserRole) FieldUpdatedAt() *builder.Column {
	return TableUserRole.F("UpdatedAt")
}

func (UserRole) FieldCreatedAt() *builder.Column {
	return TableUserRole.F("CreatedAt")
}

type OperationTimes struct {
	CreatedAt datatypes.Timestamp `db:"f_created_at,default='0'" json:"createdAt" `
	UpdatedAt datatypes.Timestamp `db:"f_updated_at,default='0'" json:"updatedAt"`
}

func (times *OperationTimes) MarkUpdatedAt() {
	times.UpdatedAt = datatypes.Timestamp(time.Now())
}

func (times *OperationTimes) MarkCreatedAt() {
	times.MarkUpdatedAt()
	times.CreatedAt = times.UpdatedAt
}

type OperationTimesWithDeletedAt struct {
	OperationTimes
	DeletedAt datatypes.Timestamp `db:"f_deleted_at,default='0'" json:"-"`
}
