package er_test

import (
	"encoding/json"

	"github.com/liucxer/courier/sqlx/er"
	"github.com/liucxer/courier/sqlx/generator/__examples__/database"
	"github.com/liucxer/courier/sqlx/postgresqlconnector"
)

func ExampleDatabaseERFromDB() {
	ers := er.DatabaseERFromDB(database.DBTest, &postgresqlconnector.PostgreSQLConnector{})
	_, _ = json.MarshalIndent(ers, "", "  ")
	// Output:
}
