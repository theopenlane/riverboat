package river

import (
	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"

	"github.com/theopenlane/corejobs"
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

		log.Info().Msg("worker enabled: email")
	}

	if w.SlackWorker.Config.Enabled {
		if err := river.AddWorkerSafely(workers, &jobs.SlackWorker{
			Config: w.SlackWorker.Config,
		}); err != nil {
			return nil, err
		}

		log.Info().Msg("worker enabled: slack")
	}

	if err := createExportWorkers(w, workers); err != nil {
		return nil, err
	}

	if err := createOrganizationDeletionWorkers(w, insertOnlyClient, workers); err != nil {
		return nil, err
	}

	// add more workers here

	return workers, nil
}

func createExportWorkers(w Workers, workers *river.Workers) error {
	if w.ExportContentWorker.Config.Enabled {
		if err := w.ExportContentWorker.Config.SetDefaultsIfUnset(w.OpenlaneConfig); err != nil {
			return err
		}

		if err := river.AddWorkerSafely(workers, &w.ExportContentWorker); err != nil {
			return err
		}

		log.Info().Msg("worker enabled: export content")
	}

	if w.DeleteExportContentWorker.Config.Enabled {
		if err := w.DeleteExportContentWorker.Config.SetDefaultsIfUnset(w.OpenlaneConfig); err != nil {
			return err
		}

		if err := river.AddWorkerSafely(workers, &w.DeleteExportContentWorker); err != nil {
			return err
		}

		log.Info().Msg("worker enabled: delete export content")
	}

	return nil
}

func setAndValidateOpenlaneConfigDefaults(input *corejobs.OpenlaneConfig, target jobs.OpenlaneConfig) error {
	input.SetAPIHost(target.GetAPIHost())
	input.SetAPIToken(target.GetAPIToken())

	return input.Validate()
}

func createOrganizationDeletionWorkers(w Workers, insertOnlyClient *riverqueue.Client, workers *river.Workers) error {
	if w.OrganizationDeletionReminderWorker.Config.Enabled {
		if err := setAndValidateOpenlaneConfigDefaults(&w.OrganizationDeletionReminderWorker.Config.OpenlaneConfig, w.OpenlaneConfig); err != nil {
			log.Error().Err(err).Msg("failed to set and validate openlane config defaults for organization deletion reminder worker")
			return err
		}

		w.EmailConfig.SetDefaultsIfUnset(&w.OrganizationDeletionReminderWorker.Config.Email.Config)

		w.OrganizationDeletionReminderWorker.WithRiverClient(insertOnlyClient)

		if err := river.AddWorkerSafely(workers, &w.OrganizationDeletionReminderWorker); err != nil {
			return err
		}

		log.Info().Msg("worker enabled: organization deletion reminder")
	}

	if w.OrganizationDeletionWorker.Config.Enabled {
		if err := setAndValidateOpenlaneConfigDefaults(&w.OrganizationDeletionWorker.Config.OpenlaneConfig, w.OpenlaneConfig); err != nil {
			log.Error().Err(err).Msg("failed to set and validate openlane config defaults for organization deletion worker")
			return err
		}

		if err := river.AddWorkerSafely(workers, &w.OrganizationDeletionWorker); err != nil {
			return err
		}

		log.Info().Msg("worker enabled: organization deletion")
	}

	return nil
}
