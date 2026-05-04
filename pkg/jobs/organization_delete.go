package jobs

import (
	"context"
	"time"

	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"

	"github.com/theopenlane/core/common/jobspec"
	"github.com/theopenlane/core/common/models"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/theopenlane/riverboat/pkg/jobs/openlane"
)

// OrganizationDeleteConfig contains the configuration for the organization deletion worker.
type OrganizationDeleteConfig struct {
	OpenlaneConfig `koanf:",squash" jsonschema:"description=the openlane API configuration for organization deletion"`

	RunInterval time.Duration `koanf:"runinterval" json:"runinterval" jsonschema:"required,default=24h description=how often to run the organization deletion worker"`

	MaxDeletesPerRun int64 `koanf:"maxdeletesperrun" json:"maxdeletesperrun" jsonschema:"required,default=25 description=maximum number of overdue organizations to delete in a single run"`

	Enabled bool `koanf:"enabled" json:"enabled" jsonschema:"required description=whether the organization deletion worker is enabled"`
}

// OrganizationDeleteWorker deletes organizations in Openlane.
type OrganizationDeleteWorker struct {
	river.WorkerDefaults[jobspec.OrganizationDeletionArgs]

	Config OrganizationDeleteConfig `koanf:"config" json:"config" jsonschema:"description=the configuration for organization deletion"`

	olClient openlane.GraphClient
}

// WithOpenlaneClient sets the Openlane client for the worker.
func (w *OrganizationDeleteWorker) WithOpenlaneClient(cl openlane.GraphClient) *OrganizationDeleteWorker {
	w.olClient = cl
	return w
}

// Work satisfies the river.Worker interface for the organization deletion worker.
func (w *OrganizationDeleteWorker) Work(ctx context.Context, job *river.Job[jobspec.OrganizationDeletionArgs]) error {
	logger := log.Ctx(ctx).With().Str("job_kind", job.Kind).Logger()
	logger.Info().Msg("starting organization deletion")

	if w.olClient == nil {
		cl, err := w.Config.getOpenlaneClient()
		if err != nil {
			return err
		}

		w.olClient = cl
	}

	now, err := models.ToDateTime(time.Now().Format("2006-01-02"))
	if err != nil {
		logger.Error().Err(err).Msg("failed to format current date for pending deletion filter")
		return err
	}

	settings, err := w.olClient.GetOrganizationSettings(ctx, nil, &w.Config.MaxDeletesPerRun, nil, nil,
		&graphclient.OrganizationSettingWhereInput{
			PendingDeletionAtNotNil: lo.ToPtr(true),
			PendingDeletionAtLte:    now,
		}, []*graphclient.OrganizationSettingOrder{
			{
				Field:     graphclient.OrganizationSettingOrderFieldUpdatedAt,
				Direction: graphclient.OrderDirectionAsc,
			},
		})
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch organizations pending deletion")
		return err
	}

	for _, edge := range settings.OrganizationSettings.Edges {
		if edge == nil || edge.Node == nil || edge.Node.Organization == nil {
			continue
		}

		orgLogger := logger.With().
			Str("organization_id", edge.Node.Organization.ID).
			Str("setting_id", edge.Node.ID).
			Logger()

		if edge.Node.PaymentMethodAdded {
			if _, err := w.olClient.UpdateOrganizationSetting(ctx, edge.Node.ID, graphclient.UpdateOrganizationSettingInput{
				ClearPendingDeletionAt: lo.ToPtr(true),
			}); err != nil {
				orgLogger.Error().Err(err).Msg("failed to clear organization pending deletion state")
				return err
			}

			orgLogger.Info().Msg("cleared organization pending deletion state because payment method was added")

			continue
		}

		if _, err := w.olClient.DeleteOrganization(ctx, edge.Node.Organization.ID); err != nil {
			orgLogger.Error().Err(err).Msg("failed to delete organization")
			return err
		}

		orgLogger.Info().Msg("successfully deleted organization")
	}

	logger.Info().Msg("finished organization deletion")

	return nil
}
