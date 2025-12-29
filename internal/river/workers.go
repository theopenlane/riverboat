package river

import (
	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"

	"github.com/theopenlane/corejobs"

	"github.com/theopenlane/riverboat/pkg/jobs"
	"github.com/theopenlane/riverboat/pkg/riverqueue"
)

// createWorkers creates a new workers instance
func createWorkers(w Workers, _ *riverqueue.Client) (*river.Workers, error) {
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

	if err := createExportWorkers(w, workers); err != nil {
		return nil, err
	}

	// add more workers here

	return workers, nil
}

func createExportWorkers(w Workers, workers *river.Workers) error {
	if w.ExportContentWorker.Config.Enabled {
		exportContentConfig := &jobs.ExportContentWorker{
			Config: w.ExportContentWorker.Config,
		}

		// set Openlane config defaults if not set
		if exportContentConfig.Config.OpenlaneAPIHost == "" {
			exportContentConfig.Config.OpenlaneAPIHost = w.OpenlaneConfig.OpenlaneAPIHost
		}

		if exportContentConfig.Config.OpenlaneAPIToken == "" {
			exportContentConfig.Config.OpenlaneAPIToken = w.OpenlaneConfig.OpenlaneAPIToken
		}

		if err := river.AddWorkerSafely(workers, exportContentConfig); err != nil {
			return err
		}

		log.Info().Msg("export content worker enabled")
	}

	if w.ClearTrustCenterCacheWorker.Config.Enabled {
		if err := river.AddWorkerSafely(workers, &corejobs.ClearTrustCenterCacheWorker{
			Config: w.ClearTrustCenterCacheWorker.Config,
		}); err != nil {
			return err
		}

		log.Info().Msg("ClearTrustCenterCacheWorker worker enabled")
	}

	if w.DeleteExportContentWorker.Config.Enabled {
		deleteExportContentConfig := &jobs.DeleteExportContentWorker{
			Config: w.DeleteExportContentWorker.Config,
		}

		// set Openlane config defaults if not set
		if deleteExportContentConfig.Config.OpenlaneAPIHost == "" {
			deleteExportContentConfig.Config.OpenlaneAPIHost = w.OpenlaneConfig.OpenlaneAPIHost
		}

		if deleteExportContentConfig.Config.OpenlaneAPIToken == "" {
			deleteExportContentConfig.Config.OpenlaneAPIToken = w.OpenlaneConfig.OpenlaneAPIToken
		}

		if err := river.AddWorkerSafely(workers, deleteExportContentConfig); err != nil {
			return err
		}

		log.Info().Msg("delete export content worker enabled")
	}

	return nil
}
