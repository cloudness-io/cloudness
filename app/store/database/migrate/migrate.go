package migrate

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"

	"github.com/jmoiron/sqlx"
	"github.com/maragudk/migrate"
	"github.com/rs/zerolog/log"
)

//go:embed postgres/*.sql
var postgres embed.FS

//go:embed sqlite/*.sql
var sqlite embed.FS

const (
	tableName = "migrations"

	postgresDriverName = "postgres"
	postgresSourceDir  = "postgres"

	sqliteDriverName = "sqlite3"
	sqliteSourceDir  = "sqlite"
)

// Migrate performs the database migration.
func Migrate(ctx context.Context, db *sqlx.DB) error {
	opts, err := getMigrator(db)
	if err != nil {
		return fmt.Errorf("failed to get migrator: %w", err)
	}
	return migrate.New(opts).MigrateUp(ctx)
}

func getMigrator(db *sqlx.DB) (migrate.Options, error) {
	before := func(ctx context.Context, _ *sql.Tx, version string) error {
		ctx = log.Ctx(ctx).With().
			Str("migrate.version", version).
			Str("migrate.phase", "before").
			Logger().WithContext(ctx)
		log := log.Ctx(ctx)

		log.Info().Msg("[START]")
		defer log.Info().Msg("[DONE]")

		return nil
	}

	after := func(ctx context.Context, dbtx *sql.Tx, version string) error {
		ctx = log.Ctx(ctx).With().
			Str("migrate.version", version).
			Str("migrate.phase", "after").
			Logger().WithContext(ctx)
		log := log.Ctx(ctx)

		log.Info().Msg("[START]")
		defer log.Info().Msg("[DONE]")

		return nil
	}

	opts := migrate.Options{
		After:  after,
		Before: before,
		DB:     db.DB,
		FS:     sqlite,
		Table:  tableName,
	}

	switch db.DriverName() {
	case sqliteDriverName:
		folder, _ := fs.Sub(sqlite, sqliteSourceDir)
		opts.FS = folder
	case postgresDriverName:
		folder, _ := fs.Sub(postgres, postgresSourceDir)
		opts.FS = folder

	default:
		return migrate.Options{}, fmt.Errorf("unsupported driver '%s'", db.DriverName())
	}

	return opts, nil
}
