package river

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/rs/zerolog/log"
)

const (
	defaultMaxWorkers = 100
)

// Start the river server with the given configuration
func Start(ctx context.Context, c Config) error {
	// Create a new database connection pool
	dbPool, err := pgxpool.New(ctx, c.DatabaseHost)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	// Run migrations on startup
	if err := runMigrations(ctx, dbPool); err != nil {
		log.Fatal().Err(err).Msg("failed to run migrations")
	}

	// Create workers based on the configuration
	worker, err := createWorkers(c.Workers)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create workers")
	}

	log.Debug().Msg("workers created")

	// create queues
	queues := createQueueConfig(c.Queues)

	log.Debug().Interface("queues", queues).Msg("queues created")

	// create a new river client
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
	}

	log.Info().Msg(startBlock)

	// run the client
	if err := client.Start(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to start river client")
	}

	sigintOrTerm := make(chan os.Signal, 1)
	signal.Notify(sigintOrTerm, syscall.SIGINT, syscall.SIGTERM)

	// this waits for SIGINT/SIGTERM and when received, tries to stop
	// gracefully by allowing a chance for jobs to finish. But if that isn't
	// working, a second SIGINT/SIGTERM will tell it to terminate with prejudice and
	// it'll issue a hard stop that cancels the context of all active jobs. In
	// case that doesn't work, a third SIGINT/SIGTERM ignores River's stop procedure
	// completely and exits uncleanly.
	go func() {
		<-sigintOrTerm
		log.Info().Msg("Received SIGINT/SIGTERM; initiating soft stop (try to wait for jobs to finish)")

		softStopCtx, softStopCtxCancel := context.WithTimeout(ctx, 10*time.Second) // nolint:mnd
		defer softStopCtxCancel()

		go func() {
			select {
			case <-sigintOrTerm:
				log.Info().Msg("Received SIGINT/SIGTERM again; initiating hard stop (cancel everything)")
				softStopCtxCancel()
			case <-softStopCtx.Done():
				log.Info().Msg("soft stop timeout; initiating hard stop (cancel everything)")
			}
		}()

		err := client.Stop(softStopCtx)

		if err != nil && !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled) {
			panic(err)
		}

		if err == nil {
			log.Info().Msg("soft stop succeeded")

			return
		}

		hardStopCtx, hardStopCtxCancel := context.WithTimeout(ctx, 10*time.Second) // nolint:mnd
		defer hardStopCtxCancel()

		// As long as all jobs respect context cancellation, StopAndCancel will
		// always work. However, in the case of a bug where a job blocks despite
		// being cancelled, it may be necessary to either ignore River's stop
		// result (what's shown here) or have a supervisor kill the process.
		err = client.StopAndCancel(hardStopCtx)
		if err != nil && errors.Is(err, context.DeadlineExceeded) {
			log.Info().Msg("hard stop timeout; ignoring stop procedure and exiting unsafely")
		} else if err != nil {
			log.Panic().Err(err).Msg("hard stop failed")
		}
	}()
	<-client.Stopped()

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

var startBlock = `
          $$\                                $$\                            $$\
          \__|                               $$ |                           $$ |
 $$$$$$\  $$\ $$\    $$\  $$$$$$\   $$$$$$\  $$$$$$$\   $$$$$$\   $$$$$$\ $$$$$$\
$$  __$$\ $$ |\$$\  $$  |$$  __$$\ $$  __$$\ $$  __$$\ $$  __$$\  \____$$\\_$$  _|
$$ |  \__|$$ | \$$\$$  / $$$$$$$$ |$$ |  \__|$$ |  $$ |$$ /  $$ | $$$$$$$ | $$ |
$$ |      $$ |  \$$$  /  $$   ____|$$ |      $$ |  $$ |$$ |  $$ |$$  __$$ | $$ |$$\
$$ |      $$ |   \$  /   \$$$$$$$\ $$ |      $$$$$$$  |\$$$$$$  |\$$$$$$$ | \$$$$  |
\__|      \__|    \_/     \_______|\__|      \_______/  \______/  \_______|  \____/

`
