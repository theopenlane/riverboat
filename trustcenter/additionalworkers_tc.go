//go:build trustcenter

package trustcenter

import (
	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"

	"github.com/theopenlane/corejobs"
	"github.com/theopenlane/riverboat/pkg/jobs"
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
		if err := setAndValidateOpenlaneConfigDefaults(&w.WatermarkDocWorker.Config.OpenlaneConfig, w.OpenlaneConfig); err != nil {
			log.Error().Err(err).Msg("failed to set and validate openlane config defaults for watermark doc worker")
			return nil, err
		}

		if err := river.AddWorkerSafely(workers, &w.WatermarkDocWorker); err != nil {
			return nil, err
		}

		log.Info().Msg("worker enabled: trust center watermark doc")
	}

	if w.ClearTrustCenterCacheWorker.Config.Enabled {
		if err := river.AddWorkerSafely(workers, &w.ClearTrustCenterCacheWorker); err != nil {
			return nil, err
		}

		log.Info().Msg("worker enabled: clear trust center cache")
	}

	if err := createPirschDomainWorkers(w, insertOnlyClient, workers); err != nil {
		return nil, err
	}

	if err := createPreviewDomainWorkers(w, insertOnlyClient, workers); err != nil {
		return nil, err
	}

	if w.AttestNDARequestWorker.Config.Enabled {
		if err := setAndValidateOpenlaneConfigDefaults(&w.AttestNDARequestWorker.Config.OpenlaneConfig, w.OpenlaneConfig); err != nil {
			log.Error().Err(err).Msg("failed to set and validate openlane config defaults for attest NDA request worker")
			return nil, err
		}

		if err := river.AddWorkerSafely(workers, &w.AttestNDARequestWorker); err != nil {
			return nil, err
		}

		log.Info().Msg("worker enabled: attest NDA request")
	}

	// add more workers here

	return workers, nil
}

// setAndValidateOpenlaneConfigDefaults sets OpenlaneConfig fields from target if they are not already set in input
func setAndValidateOpenlaneConfigDefaults(input *corejobs.OpenlaneConfig, target jobs.OpenlaneConfig) error {
	input.SetAPIHost(target.GetAPIHost())
	input.SetAPIToken(target.GetAPIToken())

	return input.Validate()
}

func createCustomDomainWorkers(w Workers, insertOnlyClient *riverqueue.Client, workers *river.Workers) error {
	if w.CreateCustomDomainWorker.Config.Enabled {
		if err := setAndValidateOpenlaneConfigDefaults(&w.CreateCustomDomainWorker.Config.OpenlaneConfig, w.OpenlaneConfig); err != nil {
			log.Error().Err(err).Msg("failed to set and validate openlane config defaults for create custom domain worker")
			return err
		}

		w.CreateCustomDomainWorker.WithRiverClient(insertOnlyClient)

		if err := river.AddWorkerSafely(workers, &w.CreateCustomDomainWorker); err != nil {
			return err
		}

		log.Info().Msg("worker enabled: create custom domain")
	}

	if w.ValidateCustomDomainWorker.Config.Enabled {
		if err := setAndValidateOpenlaneConfigDefaults(&w.ValidateCustomDomainWorker.Config.OpenlaneConfig, w.OpenlaneConfig); err != nil {
			log.Error().Err(err).Msg("failed to set and validate openlane config defaults for validate custom domain worker")
			return err
		}

		if err := river.AddWorkerSafely(workers, &w.ValidateCustomDomainWorker); err != nil {
			return err
		}

		log.Info().Msg("worker enabled: validate custom domain")
	}

	if w.DeleteCustomDomainWorker.Config.Enabled {
		if err := setAndValidateOpenlaneConfigDefaults(&w.DeleteCustomDomainWorker.Config.OpenlaneConfig, w.OpenlaneConfig); err != nil {
			log.Error().Err(err).Msg("failed to set and validate openlane config defaults for delete custom domain worker")
			return err
		}

		if err := river.AddWorkerSafely(workers, &w.DeleteCustomDomainWorker); err != nil {
			return err
		}

		log.Info().Msg("worker enabled: delete custom domain")
	}

	return nil
}

func createPirschDomainWorkers(w Workers, insertOnlyClient *riverqueue.Client, workers *river.Workers) error {
	if w.CreatePirschDomainWorker.Config.Enabled {
		if err := setAndValidateOpenlaneConfigDefaults(&w.CreatePirschDomainWorker.Config.OpenlaneConfig, w.OpenlaneConfig); err != nil {
			log.Error().Err(err).Msg("failed to set and validate openlane config defaults for create pirsch domain worker")
			return err
		}

		w.CreatePirschDomainWorker.WithRiverClient(insertOnlyClient)

		if err := river.AddWorkerSafely(workers, &w.CreatePirschDomainWorker); err != nil {
			return err
		}

		log.Info().Msg("worker enabled: create pirsch domain")
	}

	if w.DeletePirschDomainWorker.Config.Enabled {
		if err := setAndValidateOpenlaneConfigDefaults(&w.DeletePirschDomainWorker.Config.OpenlaneConfig, w.OpenlaneConfig); err != nil {
			log.Error().Err(err).Msg("failed to set and validate openlane config defaults for delete pirsch domain worker")
			return err
		}

		if err := river.AddWorkerSafely(workers, &w.DeletePirschDomainWorker); err != nil {
			return err
		}

		log.Info().Msg("delete pirsch domain worker enabled")
	}

	if w.UpdatePirschDomainWorker.Config.Enabled {
		if err := river.AddWorkerSafely(workers, &w.UpdatePirschDomainWorker); err != nil {
			return err
		}

		log.Info().Msg("worker enabled: update pirsch domain")
	}

	return nil
}

func createPreviewDomainWorkers(w Workers, insertOnlyClient *riverqueue.Client, workers *river.Workers) error {
	if w.CreatePreviewDomainWorker.Config.Enabled {
		if err := setAndValidateOpenlaneConfigDefaults(&w.CreatePreviewDomainWorker.Config.OpenlaneConfig, w.OpenlaneConfig); err != nil {
			log.Error().Err(err).Msg("failed to set and validate openlane config defaults for create preview domain worker")
			return err
		}

		w.CreatePreviewDomainWorker.WithRiverClient(insertOnlyClient)

		if err := river.AddWorkerSafely(workers, &w.CreatePreviewDomainWorker); err != nil {
			return err
		}

		log.Info().Str("worker", "create preview domain").Msg("worker enabled: create preview domain")
	}

	if w.DeletePreviewDomainWorker.Config.Enabled {
		if err := setAndValidateOpenlaneConfigDefaults(&w.DeletePreviewDomainWorker.Config.OpenlaneConfig, w.OpenlaneConfig); err != nil {
			log.Error().Err(err).Msg("failed to set and validate openlane config defaults for delete preview domain worker")
			return err
		}

		if err := river.AddWorkerSafely(workers, &w.DeletePreviewDomainWorker); err != nil {
			return err
		}

		log.Info().Msg("worker enabled: delete preview domain")
	}

	if w.ValidatePreviewDomainWorker.Config.Enabled {
		if err := setAndValidateOpenlaneConfigDefaults(&w.ValidatePreviewDomainWorker.Config.OpenlaneConfig, w.OpenlaneConfig); err != nil {
			log.Error().Err(err).Msg("failed to set and validate openlane config defaults for validate preview domain worker")
			return err
		}

		if err := river.AddWorkerSafely(workers, &w.ValidatePreviewDomainWorker); err != nil {
			return err
		}

		log.Info().Msg("worker enabled: validate preview domain")
	}

	return nil
}
