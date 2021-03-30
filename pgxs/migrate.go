package pgxs

import (
	"context"
	"fmt"
	"github.com/jackc/tern/migrate"
	"os"
	"strings"
)

// MigrationsTable is for saving actual schema version
var MigrationsTable = "public.schema_version"

func (db *Repo) Migrate(ctx context.Context) error {

	conn, err := newConn(ctx, db.Logger, db.Config)
	if err != nil {
		return err
	}

	migrator, err := migrate.NewMigrator(ctx, conn, MigrationsTable)
	if err != nil {
		db.Logger.Errorf("Unable to create migrator: %s", err)
		return err
	}

	migrator.OnStart = func(sequence int32, name, direction, sql string) {
		db.Logger.Infof("executing %s %s\n%s\n\n", name, direction, sql)
	}

	currentVersion, err := migrator.GetCurrentVersion(ctx)
	if err != nil {
		db.Logger.Errorf("Unable to get current version: %s", err)
		return err
	} else {
		db.Logger.Infof("Current migration version: %d", currentVersion)
	}

	if err := migrator.LoadMigrations(db.Config.MigrationsPath); err != nil {
		db.Logger.Errorf("Unable to load migrations: %s", err)
		return err
	} else {
		db.Logger.Infof("Successfully loaded migrations: %d", len(migrator.Migrations))
	}

	if err := migrator.Migrate(ctx); err != nil {
		if err, ok := err.(migrate.MigrationPgError); ok {
			if err.Detail != "" {
				fmt.Fprintln(os.Stderr, "DETAIL:", err.Detail)
			}

			if err.Position != 0 {
				ele, err := ExtractErrorLine(err.Sql, int(err.Position))
				if err != nil {
					db.Logger.Errorf("Unable to extract error line: %s", err)
				}

				prefix := fmt.Sprintf("LINE %d: ", ele.LineNum)
				db.Logger.Warnf("%s%s\n", prefix, ele.Text)

				padding := strings.Repeat(" ", len(prefix)+ele.ColumnNum-1)
				db.Logger.Warnf("%s^\n", padding)
			}
		}
		db.Logger.Errorf("Unable to migrate: %s", err)
		return err
	} else {
		actual := migrator.Migrations[len(migrator.Migrations)-1]
		db.Logger.Infow("Successfully finished migration", "name", actual.Name, "seq", actual.Sequence)
	}

	return nil
}
