package riverqueue

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivermigrate"
)

// RunMigrations runs the migrations for the river server
// see https://riverqueue.com/docs/migrations for more information
func RunMigrations(ctx context.Context, dbPool *pgxpool.Pool) error {
	// run migrations here
	migrator, err := rivermigrate.New(riverpgxv5.New(dbPool), nil)
	if err != nil {
		return err
	}

	if _, err := migrator.Migrate(ctx, rivermigrate.DirectionUp, nil); err != nil {
		return err
	}

	return nil
}
