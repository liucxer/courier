package geography

import (
	"database/sql/driver"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"github.com/liucxer/courier/sqlx"
	"github.com/liucxer/courier/sqlx/builder"
	"github.com/liucxer/courier/sqlx/migration"
	"github.com/liucxer/courier/sqlx/mysqlconnector"
	"github.com/liucxer/courier/sqlx/postgresqlconnector"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

var (
	mysqlConnector = &mysqlconnector.MysqlConnector{
		Host:  "root@tcp(0.0.0.0:3306)",
		Extra: "charset=utf8mb4&parseTime=true&interpolateParams=true&autocommit=true&loc=Local",
	}

	postgresConnector = &postgresqlconnector.PostgreSQLConnector{
		Host:       "postgres://postgres@0.0.0.0:5432",
		Extra:      "sslmode=disable",
		Extensions: []string{"postgis", "hstore"},
	}
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

type Geometries struct {
	ID              string          `db:"F_id"`
	Point           Point           `db:"F_point"`
	MultiPoint      MultiPoint      `db:"F_multi_point,null"`
	LineString      LineString      `db:"F_line_string,null"`
	MultiLineString MultiLineString `db:"F_multi_line_string,null"`
	Polygon         Polygon         `db:"f_polygon,null"`
	MultiPolygon    MultiPolygon    `db:"f_multi_polygon,null"`
	Geometry        Geometry        `db:"f_geometry"`
}

func (Geometries) PrimaryKey() []string {
	return []string{"ID"}
}

func (Geometries) TableName() string {
	return "t_geom"
}

func TestGeomRW(t *testing.T) {
	tt := require.New(t)

	dbGeom := sqlx.NewDatabase("geo")

	for _, connector := range []driver.Connector{
		postgresConnector,
		mysqlConnector,
	} {
		db := dbGeom.OpenDB(connector)

		userTable := dbGeom.Register(&Geometries{})
		err := migration.Migrate(db, nil)
		tt.NoError(err)

		g := Geometries{
			ID:              uuid.New().String(),
			Point:           Point{-1, 2},
			MultiPoint:      MultiPoint{{-1, 2}, {2, 1}},
			LineString:      LineString{{-1, 2}, {2, 1}},
			MultiLineString: MultiLineString{{{-1, 2}, {2, 1}}, {{-1, 2}, {2, 1}}},
			Polygon: Polygon{
				{{-1, 2}, {2, 1}, {2, 2}, {2, 3}, {-1, 2}},
				{{-1, 2}, {2, 1}, {2, 2}, {-1, 2}},
				{{-1, 2}, {2, 3}, {2, 4}, {-1, 2}},
			},
			MultiPolygon: MultiPolygon{{{{-1, 2}, {2, 1}, {-1, 2}, {-1, 2}}}, {{{-1, 2}, {2, 1}, {-1, 2}, {-1, 2}}}},
			Geometry:     ToGeometry(Point{-1, 2}),
		}

		_, errInsert := db.ExecExpr(sqlx.InsertToDB(db, g, nil))
		tt.NoError(errInsert)

		{
			gForSelect := Geometries{}
			err := db.QueryExprAndScan(
				builder.Select(nil).From(
					userTable,
					builder.Where(userTable.F("ID").Eq(g.ID)),
				),
				&gForSelect,
			)
			tt.NoError(err)
			spew.Dump(gForSelect)
		}

		dbGeom.Tables.Range(func(t *builder.Table, idx int) {
			_, err := db.ExecExpr(db.Dialect().DropTable(t))
			tt.NoError(err)
		})
	}

}
