package river

import (
	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"

	"github.com/theopenlane/core/pkg/corejobs"

	"github.com/theopenlane/riverboat/pkg/jobs"
)

// createWorkers creates a new workers instance
func createWorkers(c Workers) (*river.Workers, error) {
	// create workers
	workers := river.NewWorkers()

	if c.EmailWorker.Config.Enabled {
		if err := river.AddWorkerSafely(workers, &jobs.EmailWorker{
			Config: c.EmailWorker.Config,
		},
		); err != nil {
			return nil, err
		}

		log.Info().Msg("email worker enabled")
	}

	if c.SlackWorker.Config.Enabled {
		if err := river.AddWorkerSafely(workers, &jobs.SlackWorker{
			Config: c.SlackWorker.Config,
		}); err != nil {
			return nil, err
		}
		log.Info().Msg("slack worker enabled")
	}

	if c.CreateCustomDomainWorker.Config.Enabled {
		if err := river.AddWorkerSafely(workers, &corejobs.CreateCustomDomainWorker{
			Config: c.CreateCustomDomainWorker.Config,
		},
		); err != nil {
			return nil, err
		}

		log.Info().Msg("create custom domain worker enabled")
	}

	if c.ValidateCustomDomainWorker.Config.Enabled {
		if err := river.AddWorkerSafely(workers, &corejobs.ValidateCustomDomainWorker{
			Config: c.ValidateCustomDomainWorker.Config,
		},
		); err != nil {
			return nil, err
		}

		log.Info().Msg("validate custom domain worker enabled")
	}

	if c.DeleteCustomDomainWorker.Config.Enabled {
		if err := river.AddWorkerSafely(workers, &corejobs.DeleteCustomDomainWorker{
			Config: c.DeleteCustomDomainWorker.Config,
		},
		); err != nil {
			return nil, err
		}

		log.Info().Msg("delete custom domain worker enabled")
	}

	if c.ExportContentWorker.Config.Enabled {
		if err := river.AddWorkerSafely(workers, &corejobs.ExportContentWorker{
			Config: c.ExportContentWorker.Config,
		},
		); err != nil {
			return nil, err
		}

		log.Info().Msg("export content worker enabled")
	}

	if c.DeleteExportContentWorker.Config.Enabled {
		if err := river.AddWorkerSafely(workers, &corejobs.DeleteExportContentWorker{
			Config: c.DeleteExportContentWorker.Config,
		},
		); err != nil {
			return nil, err
		}

		log.Info().Msg("delete export content worker enabled")
	}

	if c.WatermarkDocWorker.Config.Enabled {
		if err := river.AddWorkerSafely(workers, &corejobs.WatermarkDocWorker{
			Config: c.WatermarkDocWorker.Config,
		},
		); err != nil {
			return nil, err
		}

		log.Info().Msg("watermark doc worker enabled")
	}

	if c.CreatePirschDomainWorker.Config.Enabled {
		if err := river.AddWorkerSafely(workers, &corejobs.CreatePirschDomainWorker{
			Config: c.CreatePirschDomainWorker.Config,
		},
		); err != nil {
			return nil, err
		}

		log.Info().Msg("create pirsch domain worker enabled")
	}

	if c.DeletePirschDomainWorker.Config.Enabled {
		if err := river.AddWorkerSafely(workers, &corejobs.DeletePirschDomainWorker{
			Config: c.DeletePirschDomainWorker.Config,
		},
		); err != nil {
			return nil, err
		}

		log.Info().Msg("delete pirsch domain worker enabled")
	}

	// add more workers here

	return workers, nil
}
