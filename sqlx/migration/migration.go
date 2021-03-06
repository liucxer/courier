package migration

import (
	"context"
	"io"

	"github.com/liucxer/courier/sqlx"
	"github.com/liucxer/courier/sqlx/enummeta"
)

type contextKeyMigrationOutput int

func MigrationOutputFromContext(ctx context.Context) io.Writer {
	if opts, ok := ctx.Value(contextKeyMigrationOutput(1)).(io.Writer); ok {
		if opts != nil {
			return opts
		}
	}
	return nil
}

func MustMigrate(db sqlx.DBExecutor, w io.Writer) {
	if err := Migrate(db, w); err != nil {
		panic(err)
	}
}

func Migrate(db sqlx.DBExecutor, output io.Writer) error {
	ctx := context.WithValue(db.Context(), contextKeyMigrationOutput(1), output)

	if err := db.(sqlx.Migrator).Migrate(ctx, db); err != nil {
		return err
	}
	if output == nil {
		if err := enummeta.SyncEnum(db); err != nil {
			return err
		}
	}
	return nil
}
