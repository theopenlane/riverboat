package river

import (
	"github.com/riverqueue/river"

	"github.com/theopenlane/core/pkg/corejobs"

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

	if err := river.AddWorkerSafely(workers, &corejobs.CreateCustomDomainWorker{
		Config: c.CreateCustomDomainWorker.Config,
	},
	); err != nil {
		return nil, err
	}

	if err := river.AddWorkerSafely(workers, &corejobs.ValidateCustomDomainWorker{
		Config: c.ValidateCustomDomainWorker.Config,
	},
	); err != nil {
		return nil, err
	}

	if err := river.AddWorkerSafely(workers, &corejobs.DeleteCustomDomainWorker{
		Config: c.DeleteCustomDomainWorker.Config,
	},
	); err != nil {
		return nil, err
	}

	if err := river.AddWorkerSafely(workers, &corejobs.ExportContentWorker{
		Config: c.ExportContentWorker.Config,
	},
	); err != nil {
		return nil, err
	}

	// add more workers here

	return workers, nil
}
