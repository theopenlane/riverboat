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

		// set Openlane config defaults if not set
		if customDomainConfig.Config.OpenlaneAPIHost == "" {
			customDomainConfig.Config.OpenlaneAPIHost = w.OpenlaneConfig.OpenlaneAPIHost
		}

		if customDomainConfig.Config.OpenlaneAPIToken == "" {
			customDomainConfig.Config.OpenlaneAPIToken = w.OpenlaneConfig.OpenlaneAPIToken
		}

		customDomainConfig.WithRiverClient(insertOnlyClient)

		if err := river.AddWorkerSafely(workers, customDomainConfig); err != nil {
			return nil, err
		}

		log.Info().Msg("create custom domain worker enabled")
	}

	if w.ValidateCustomDomainWorker.Config.Enabled {
		validateCustomDomainConfig := &corejobs.ValidateCustomDomainWorker{
			Config: w.ValidateCustomDomainWorker.Config,
		}

		// set Openlane config defaults if not set
		if validateCustomDomainConfig.Config.OpenlaneAPIHost == "" {
			validateCustomDomainConfig.Config.OpenlaneAPIHost = w.OpenlaneConfig.OpenlaneAPIHost
		}

		if validateCustomDomainConfig.Config.OpenlaneAPIToken == "" {
			validateCustomDomainConfig.Config.OpenlaneAPIToken = w.OpenlaneConfig.OpenlaneAPIToken
		}

		if err := river.AddWorkerSafely(workers, validateCustomDomainConfig); err != nil {
			return nil, err
		}

		log.Info().Msg("validate custom domain worker enabled")
	}

	if w.DeleteCustomDomainWorker.Config.Enabled {
		deleteCustomDomainConfig := &corejobs.DeleteCustomDomainWorker{
			Config: w.DeleteCustomDomainWorker.Config,
		}

		// set Openlane config defaults if not set
		if deleteCustomDomainConfig.Config.OpenlaneAPIHost == "" {
			deleteCustomDomainConfig.Config.OpenlaneAPIHost = w.OpenlaneConfig.OpenlaneAPIHost
		}

		if deleteCustomDomainConfig.Config.OpenlaneAPIToken == "" {
			deleteCustomDomainConfig.Config.OpenlaneAPIToken = w.OpenlaneConfig.OpenlaneAPIToken
		}

		if err := river.AddWorkerSafely(workers, deleteCustomDomainConfig); err != nil {
			return nil, err
		}

		log.Info().Msg("delete custom domain worker enabled")
	}

	if w.ExportContentWorker.Config.Enabled {
		exportContentConfig := &corejobs.ExportContentWorker{
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
			return nil, err
		}

		log.Info().Msg("export content worker enabled")
	}

	if w.DeleteExportContentWorker.Config.Enabled {
		deleteExportContentConfig := &corejobs.DeleteExportContentWorker{
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
			return nil, err
		}

		log.Info().Msg("delete export content worker enabled")
	}

	if w.WatermarkDocWorker.Config.Enabled {
		watermarkDocConfig := &corejobs.WatermarkDocWorker{
			Config: w.WatermarkDocWorker.Config,
		}

		// set Openlane config defaults if not set
		if watermarkDocConfig.Config.OpenlaneAPIHost == "" {
			watermarkDocConfig.Config.OpenlaneAPIHost = w.OpenlaneConfig.OpenlaneAPIHost
		}

		if watermarkDocConfig.Config.OpenlaneAPIToken == "" {
			watermarkDocConfig.Config.OpenlaneAPIToken = w.OpenlaneConfig.OpenlaneAPIToken
		}

		if err := river.AddWorkerSafely(workers, watermarkDocConfig); err != nil {
			return nil, err
		}

		log.Info().Msg("watermark doc worker enabled")
	}

	if w.CreatePirschDomainWorker.Config.Enabled {
		pirschDomainConfig := &corejobs.CreatePirschDomainWorker{
			Config: w.CreatePirschDomainWorker.Config,
		}

		// set Openlane config defaults if not set
		if pirschDomainConfig.Config.OpenlaneAPIHost == "" {
			pirschDomainConfig.Config.OpenlaneAPIHost = w.OpenlaneConfig.OpenlaneAPIHost
		}

		if pirschDomainConfig.Config.OpenlaneAPIToken == "" {
			pirschDomainConfig.Config.OpenlaneAPIToken = w.OpenlaneConfig.OpenlaneAPIToken
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
