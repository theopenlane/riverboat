package river

import (
	"context"

	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"
)

type WaitsForCancelOnlyArgs struct{}

func (WaitsForCancelOnlyArgs) Kind() string { return "waits_for_cancel_only" }

// WaitsForCancelOnlyWorker is a worker that will never finish jobs until its
// context is cancelled.
type WaitsForCancelOnlyWorker struct {
	river.WorkerDefaults[WaitsForCancelOnlyArgs]

	jobStarted chan struct{}
}

func (w *WaitsForCancelOnlyWorker) Work(ctx context.Context, job *river.Job[WaitsForCancelOnlyArgs]) error {
	log.Info().Msg("creating working job that doesn't finish until cancelled")
	close(w.jobStarted)

	<-ctx.Done()
	log.Info().Msg("job cancelled")

	// In the event of cancellation, an error should be returned so that the job
	// goes back in the retry queue.
	return ctx.Err()
}
