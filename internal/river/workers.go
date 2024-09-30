package river

import (
	"github.com/riverqueue/river"

	"github.com/theopenlane/riverboat/pkg/jobs"
)

// createWorkers creates a new workers instance
func createWorkers(c Workers) (*river.Workers, error) {
	// create workers
	workers := river.NewWorkers()

	if err := river.AddWorkerSafely(workers, &jobs.EmailWorker{
		Config: c.EmailWorker.Config,
	},
	); err != nil {
		return nil, err
	}

	if err := river.AddWorkerSafely(workers, &jobs.DatabaseWorker{
		Config: c.DatabaseWorker.Config,
	},
	); err != nil {
		return nil, err
	}

	// add more workers here

	return workers, nil
}
