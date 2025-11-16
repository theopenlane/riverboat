package river

import (
	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"

	"github.com/theopenlane/core/pkg/corejobs"

	"github.com/theopenlane/riverboat/pkg/jobs"
	"github.com/theopenlane/riverboat/pkg/riverqueue"
)

// createWorkers creates a new workers instance
func createWorkers(w Workers, insertOnlyClient *riverqueue.Client) (*river.Workers, error) {
	// create workers
	workers := river.NewWorkers()

	if w.EmailWorker.Config.Enabled {
		if err := river.AddWorkerSafely(workers, &jobs.EmailWorker{
			Config: w.EmailWorker.Config,
		},
		); err != nil {
			return nil, err
		}

		log.Info().Msg("email worker enabled")
	}

	if w.SlackWorker.Config.Enabled {
		if err := river.AddWorkerSafely(workers, &jobs.SlackWorker{
			Config: w.SlackWorker.Config,
		}); err != nil {
			return nil, err
		}
		log.Info().Msg("slack worker enabled")
	}

	if w.CreateCustomDomainWorker.Config.Enabled {
		customDomainConfig := &corejobs.CreateCustomDomainWorker{
			Config: w.CreateCustomDomainWorker.Config,
		}

		customDomainConfig.WithRiverClient(insertOnlyClient)

		if err := river.AddWorkerSafely(workers, customDomainConfig); err != nil {
			return nil, err
		}

		log.Info().Msg("create custom domain worker enabled")
	}

	if w.ValidateCustomDomainWorker.Config.Enabled {
		if err := river.AddWorkerSafely(workers, &corejobs.ValidateCustomDomainWorker{
			Config: w.ValidateCustomDomainWorker.Config,
		},
		); err != nil {
			return nil, err
		}

		log.Info().Msg("validate custom domain worker enabled")
	}

	if w.DeleteCustomDomainWorker.Config.Enabled {
		if err := river.AddWorkerSafely(workers, &corejobs.DeleteCustomDomainWorker{
			Config: w.DeleteCustomDomainWorker.Config,
		},
		); err != nil {
			return nil, err
		}

		log.Info().Msg("delete custom domain worker enabled")
	}

	if w.ExportContentWorker.Config.Enabled {
		if err := river.AddWorkerSafely(workers, &corejobs.ExportContentWorker{
			Config: w.ExportContentWorker.Config,
		},
		); err != nil {
			return nil, err
		}

		log.Info().Msg("export content worker enabled")
	}

	if w.DeleteExportContentWorker.Config.Enabled {
		if err := river.AddWorkerSafely(workers, &corejobs.DeleteExportContentWorker{
			Config: w.DeleteExportContentWorker.Config,
		},
		); err != nil {
			return nil, err
		}

		log.Info().Msg("delete export content worker enabled")
	}

	if w.WatermarkDocWorker.Config.Enabled {
		if err := river.AddWorkerSafely(workers, &corejobs.WatermarkDocWorker{
			Config: w.WatermarkDocWorker.Config,
		},
		); err != nil {
			return nil, err
		}

		log.Info().Msg("watermark doc worker enabled")
	}

	if w.CreatePirschDomainWorker.Config.Enabled {
		pirschDomainConfig := &corejobs.CreatePirschDomainWorker{
			Config: w.CreatePirschDomainWorker.Config,
		}

		pirschDomainConfig.WithRiverClient(insertOnlyClient)

		if err := river.AddWorkerSafely(workers, pirschDomainConfig); err != nil {
			return nil, err
		}

		log.Info().Msg("create pirsch domain worker enabled")
	}

	if w.DeletePirschDomainWorker.Config.Enabled {
		if err := river.AddWorkerSafely(workers, &corejobs.DeletePirschDomainWorker{
			Config: w.DeletePirschDomainWorker.Config,
		},
		); err != nil {
			return nil, err
		}

		log.Info().Msg("delete pirsch domain worker enabled")
	}

	// add more workers here

	return workers, nil
}
