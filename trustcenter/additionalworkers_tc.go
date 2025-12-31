//go:build trustcenter

package trustcenter

import (
	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"

	"github.com/theopenlane/corejobs"

	"github.com/theopenlane/riverboat/pkg/riverqueue"
)

// AddConditionalWorkers adds trust center specific workers when the trustcenter build tag is present
func AddConditionalWorkers(workers *river.Workers, w Workers, insertOnlyClient *riverqueue.Client) (*river.Workers, error) {
	log.Info().Msg("adding additional trust center workers")

	// create workers
	if err := createCustomDomainWorkers(w, insertOnlyClient, workers); err != nil {
		return nil, err
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

		log.Info().Msg("worker enabled: trust center watermark doc")
	}

	if w.ClearTrustCenterCacheWorker.Config.Enabled {
		if err := river.AddWorkerSafely(workers, &corejobs.ClearTrustCenterCacheWorker{
			Config: w.ClearTrustCenterCacheWorker.Config,
		}); err != nil {
			return nil, err
		}

		log.Info().Msg("ClearTrustCenterCacheWorker worker enabled")
	}

	if err := createPirschDomainWorkers(w, insertOnlyClient, workers); err != nil {
		return nil, err
	}

	if err := createPreviewDomainWorkers(w, insertOnlyClient, workers); err != nil {
		return nil, err
	}

	// add more workers here

	return workers, nil
}

func createCustomDomainWorkers(w Workers, insertOnlyClient *riverqueue.Client, workers *river.Workers) error {
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
			return err
		}

		log.Info().Msg("worker enabled: create custom domain")
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
			return err
		}

		log.Info().Msg("worker enabled: validate custom domain")
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
			return err
		}

		log.Info().Msg("worker enabled: delete custom domain")
	}

	return nil
}

func createPirschDomainWorkers(w Workers, insertOnlyClient *riverqueue.Client, workers *river.Workers) error {
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
			return err
		}

		log.Info().Msg("worker enabled: create pirsch domain")
	}

	if w.DeletePirschDomainWorker.Config.Enabled {
		if err := river.AddWorkerSafely(workers, &corejobs.DeletePirschDomainWorker{
			Config: w.DeletePirschDomainWorker.Config,
		},
		); err != nil {
			return err
		}

		log.Info().Msg("delete pirsch domain worker enabled")
	}

	if w.UpdatePirschDomainWorker.Config.Enabled {
		if err := river.AddWorkerSafely(workers, &corejobs.UpdatePirschDomainWorker{
			Config: w.UpdatePirschDomainWorker.Config,
		},
		); err != nil {
			return err
		}

		log.Info().Msg("worker enabled: update pirsch domain")
	}

	return nil
}

func createPreviewDomainWorkers(w Workers, insertOnlyClient *riverqueue.Client, workers *river.Workers) error {
	if w.CreatePreviewDomainWorker.Config.Enabled {
		previewDomainConfig := &corejobs.CreatePreviewDomainWorker{
			Config: w.CreatePreviewDomainWorker.Config,
		}

		// set Openlane config defaults if not set
		if previewDomainConfig.Config.OpenlaneAPIHost == "" {
			previewDomainConfig.Config.OpenlaneAPIHost = w.OpenlaneConfig.OpenlaneAPIHost
		}

		if previewDomainConfig.Config.OpenlaneAPIToken == "" {
			previewDomainConfig.Config.OpenlaneAPIToken = w.OpenlaneConfig.OpenlaneAPIToken
		}

		previewDomainConfig.WithRiverClient(insertOnlyClient)

		if err := river.AddWorkerSafely(workers, previewDomainConfig); err != nil {
			return err
		}

		log.Info().Msg("worker enabled: create preview domain")
	}

	if w.DeletePreviewDomainWorker.Config.Enabled {
		if err := river.AddWorkerSafely(workers, &corejobs.DeletePreviewDomainWorker{
			Config: w.DeletePreviewDomainWorker.Config,
		},
		); err != nil {
			return err
		}

		log.Info().Msg("worker enabled: delete preview domain")
	}

	if w.ValidatePreviewDomainWorker.Config.Enabled {
		if err := river.AddWorkerSafely(workers, &corejobs.ValidatePreviewDomainWorker{
			Config: w.ValidatePreviewDomainWorker.Config,
		},
		); err != nil {
			return err
		}

		log.Info().Msg("worker enabled: validate preview domain")
	}

	return nil
}
