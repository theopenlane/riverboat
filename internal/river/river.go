package river

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/rs/zerolog/log"

	"github.com/theopenlane/riverboat/pkg/jobs"
)

const (
	defaultMaxWorkers = 100
)

// Start the river server
func Start(ctx context.Context, c Config) error {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Create a new database connection pool
	dbPool, err := pgxpool.New(ctx, c.DatabaseHost)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")

		return err
	}

	// Create workers
	worker, err := createWorkers(c.Workers)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create workers")

		return err
	}

	log.Debug().Msg("workers created")

	// create queues
	queues := createQueueConfig(c.Queues)

	log.Debug().Interface("queues", queues).Msg("queues created")

	client, err := river.NewClient(
		riverpgxv5.New(dbPool),
		&river.Config{
			Workers: worker,
			Queues:  queues,
			Logger:  createLogger(c.Logger),
		},
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create river client")

		return err
	}

	go func() {
		// Run the client inline. All executed jobs will inherit from ctx:
		if err := client.Start(ctx); err != nil {
			log.Panic().Err(err).Msg("Failed to start river client")
		}
	}()

	<-ch

	return nil
}

// createLogger creates a new logger based on the configuration
func createLogger(c Logger) *slog.Logger {
	level := slog.LevelInfo

	if c.Debug {
		level = slog.LevelDebug
	}

	// create a new pretty logger
	opts := slog.HandlerOptions{
		Level: level,
	}

	if c.Pretty {
		return slog.New(slog.NewTextHandler(os.Stderr, &opts))
	}

	// create a new logger
	return slog.New(slog.NewJSONHandler(os.Stderr, &opts))
}

// createQueueConfig creates a map of queue configurations
func createQueueConfig(queues []Queue) map[string]river.QueueConfig {
	qc := map[string]river.QueueConfig{
		river.QueueDefault: {MaxWorkers: defaultMaxWorkers},
	}

	// if no queues are defined, create a default queue
	if len(queues) == 0 {
		log.Debug().Msg("using default queues")

		return qc
	}

	for _, q := range queues {
		qc[q.Name] = river.QueueConfig{
			MaxWorkers: q.MaxWorkers,
		}
	}

	return qc
}

// createWorkers creates a new workers instance
func createWorkers(c Workers) (*river.Workers, error) {
	// create workers
	workers := river.NewWorkers()

	if err := river.AddWorkerSafely(workers, &jobs.EmailWorker{
		EmailConfig: c.EmailWorker.EmailConfig,
	},
	); err != nil {
		return nil, err
	}

	// add more workers here

	return workers, nil
}
