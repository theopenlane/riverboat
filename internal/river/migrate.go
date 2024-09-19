package river

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivermigrate"
)

// runMigrations runs the migrations for the river server
// see https://riverqueue.com/docs/migrations for more information
func runMigrations(ctx context.Context, dbPool *pgxpool.Pool) error {
	// run migrations here
	migrator := rivermigrate.New(riverpgxv5.New(dbPool), nil)

	if _, err := migrator.Migrate(ctx, rivermigrate.DirectionUp, nil); err != nil {
		return err
	}

	return nil
}
