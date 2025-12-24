package river

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"

	"github.com/theopenlane/riverboat/pkg/riverqueue"
)

const (
	defaultMaxWorkers = 100
)

// Start the river server with the given configuration
func Start(ctx context.Context, c Config) error {
	logger := createLogger(c.Logger)

	// setup metrics exporting for river
	if c.Metrics.EnableMetrics {
		log.Info().Msg("setting up OpenTelemetry metrics exporter for river")

		if err := setupMetricsExporter(); err != nil {
			log.Error().Err(err).Msg("failed to setup otel metrics exporter")
		}
	}

	insertOnlyClient, err := riverqueue.New(
		ctx, riverqueue.WithConnectionURI(c.DatabaseHost),
		riverqueue.WithLogger(logger),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create insert-only river client")
	}
	defer insertOnlyClient.Close() // nolint:errcheck

	// Create workers based on the configuration
	worker, err := createWorkers(c.Workers, insertOnlyClient)
	if err != nil {
		log.Error().Err(err).Msg("failed to create workers")
	}

	worker, err = addConditionalWorkers(worker, nil, insertOnlyClient)
	if err != nil {
		log.Error().Err(err).Msg("failed to add conditional workers")
	}

	// create periodic jobs
	periodicJobs, err := createPeriodicJobs(c.Workers)
	if err != nil {
		log.Error().Err(err).Msg("failed to create periodic jobs schedules")
	}

	additionalPeriodicJobs, err := createAdditionalPeriodicJobs(nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to create additional periodic jobs schedules")
	}

	// append additional periodic jobs
	periodicJobs = append(periodicJobs, additionalPeriodicJobs...)

	log.Debug().Msg("workers created")

	// create queues
	queues := createQueueConfig(c.Queues)

	log.Debug().Interface("queues", queues).Msg("queues created")

	// create a new river client
	client, err := riverqueue.New(
		ctx,
		riverqueue.WithConnectionURI(c.DatabaseHost),
		riverqueue.WithRunMigrations(true),
		riverqueue.WithLogger(logger),
		riverqueue.WithWorkers(worker),
		riverqueue.WithQueues(queues),
		riverqueue.WithPeriodicJobs(periodicJobs),
		riverqueue.WithMaxRetries(c.DefaultMaxRetries),
		riverqueue.WithMetrics(&c.Metrics),
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to create river client")
	}

	// get the underlying river client
	rc := client.GetRiverClient()

	log.Info().Msg(startBlock)

	// run the client
	if err := rc.Start(ctx); err != nil {
		log.Error().Err(err).Msg("failed to start river client")
	}

	// start the metrics server
	go func() {
		if err := registerMetricsServer(ctx); err != nil {
			log.Error().Err(err).Msg("failed to start metrics server")
		}
	}()

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

		err := rc.Stop(softStopCtx)

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
		err = rc.StopAndCancel(hardStopCtx)
		if err != nil && errors.Is(err, context.DeadlineExceeded) {
			log.Info().Msg("hard stop timeout; ignoring stop procedure and exiting unsafely")
		} else if err != nil {
			log.Panic().Err(err).Msg("hard stop failed")
		}
	}()

	<-rc.Stopped()

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

	// if no queues are defined, just use the default queue
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
