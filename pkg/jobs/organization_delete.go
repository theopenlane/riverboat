package jobs

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/riverqueue/river"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"

	"github.com/theopenlane/core/common/jobspec"
	"github.com/theopenlane/core/common/models"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/theopenlane/riverboat/pkg/jobs/openlane"
	"github.com/theopenlane/riverboat/pkg/riverqueue"
)

// OrganizationDeleteConfig contains the configuration for the organization deletion worker.
type OrganizationDeleteConfig struct {
	OpenlaneConfig `koanf:",squash" jsonschema:"description=the openlane API configuration for organization deletion"`

	RunInterval time.Duration `koanf:"runinterval" json:"runinterval" jsonschema:"required,default=24h description=how often to run the organization deletion worker"`

	MaxDeletesPerRun int64 `koanf:"maxdeletesperrun" json:"maxdeletesperrun" jsonschema:"required,default=25 description=maximum number of overdue organizations to delete in a single run"`

	// SystemAdminOrgID is the organization ID that should never be deleted
	SystemAdminOrgID string `koanf:"systemadminorgid" json:"systemadminorgid" default:"01101101011010010111010001100010" jsonschema:"description=organization ID that must never be deleted,default=01101101011010010111010001100010"`

	// SlackChannel is the channel name or ID where organization deletion summaries are sent
	SlackChannel string `koanf:"slackchannel" json:"slackchannel" default:"customers" jsonschema:"description=slack channel for organization deletion summaries,default=customers"`

	Enabled bool `koanf:"enabled" json:"enabled" jsonschema:"required description=whether the organization deletion worker is enabled"`
}

// OrganizationDeleteWorker deletes organizations in Openlane.
type OrganizationDeleteWorker struct {
	river.WorkerDefaults[jobspec.OrganizationDeletionArgs]

	Config OrganizationDeleteConfig `koanf:"config" json:"config" jsonschema:"description=the configuration for organization deletion"`

	olClient    openlane.GraphClient
	riverClient riverqueue.JobClient
}

// WithOpenlaneClient sets the Openlane client for the worker.
func (w *OrganizationDeleteWorker) WithOpenlaneClient(cl openlane.GraphClient) *OrganizationDeleteWorker {
	w.olClient = cl
	return w
}

// WithRiverClient sets the River client for the worker.
func (w *OrganizationDeleteWorker) WithRiverClient(cl riverqueue.JobClient) *OrganizationDeleteWorker {
	w.riverClient = cl
	return w
}

// Work satisfies the river.Worker interface for the organization deletion worker.
func (w *OrganizationDeleteWorker) Work(ctx context.Context, job *river.Job[jobspec.OrganizationDeletionArgs]) error {
	logger := log.Ctx(ctx).With().Str("job_kind", job.Kind).Logger()
	logger.Info().Msg("starting organization deletion")

	w.Config.SystemAdminOrgID = strings.TrimSpace(w.Config.SystemAdminOrgID)

	if w.Config.SystemAdminOrgID == "" {
		return errSystemAdminOrgIDRequired
	}

	w.Config.SlackChannel = strings.TrimSpace(w.Config.SlackChannel)
	if w.Config.SlackChannel == "" {
		return errSlackChannelRequired
	}

	if w.olClient == nil {
		cl, err := w.Config.getOpenlaneClient()
		if err != nil {
			return err
		}

		w.olClient = cl
	}

	if w.riverClient == nil {
		logger.Error().Msg("river client is not set on worker, cannot insert organization deletion slack summary")
		return errRiverClientRequired
	}

	if err := w.checkReactivatedSubs(ctx, logger); err != nil {
		return err
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
			HasOrganizationWith: []*graphclient.OrganizationWhereInput{
				{
					IDNeq:       lo.ToPtr(w.Config.SystemAdminOrgID),
					PersonalOrg: lo.ToPtr(false),
					Not: &graphclient.OrganizationWhereInput{
						HasOrgSubscriptionsWith: []*graphclient.OrgSubscriptionWhereInput{
							{
								Or: []*graphclient.OrgSubscriptionWhereInput{
									{
										Active: lo.ToPtr(true),
									},
									{
										StripeSubscriptionStatus: lo.ToPtr(stripeSubscriptionStatusTrialing),
									},
								},
							},
						},
					},
				},
			},
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

	deletedOrgs := []string{}

	for _, edge := range settings.OrganizationSettings.Edges {
		if edge == nil || edge.Node == nil || edge.Node.Organization == nil {
			continue
		}

		orgLogger := logger.With().
			Str("organization_id", edge.Node.Organization.ID).
			Str("setting_id", edge.Node.ID).
			Logger()

		if _, err := w.olClient.DeleteOrganization(ctx, edge.Node.Organization.ID); err != nil {
			orgLogger.Error().Err(err).Msg("failed to delete organization")
			return err
		}

		deletedOrgs = append(deletedOrgs, fmt.Sprintf("%s (%s)", edge.Node.Organization.Name, edge.Node.Organization.ID))

		orgLogger.Info().Msg("successfully deleted organization")
	}

	message := fmt.Sprintf("Organization deletion summary: %d orgs deleted", len(deletedOrgs))
	if len(deletedOrgs) > 0 {
		message = fmt.Sprintf("%s:\n- %s", message, strings.Join(deletedOrgs, "\n- "))
	}

	if _, err := w.riverClient.Insert(ctx, SlackArgs{
		Channel: w.Config.SlackChannel,
		Message: message,
	}, nil); err != nil {
		logger.Error().Err(err).Msg("failed to insert slack job for organization deletion summary")
		return err
	}

	logger.Info().Msg("finished organization deletion")

	return nil
}

// checkReactivatedSubs checks to make sure orgs that were previously earmarked for deletions and now
// have an active or trialing sub will no longer be deleted
func (w *OrganizationDeleteWorker) checkReactivatedSubs(ctx context.Context, logger zerolog.Logger) error {
	var after *string

	for {
		settings, err := w.olClient.GetOrganizationSettings(ctx, &defaultPageSize, nil, after, nil,
			&graphclient.OrganizationSettingWhereInput{
				PendingDeletionAtNotNil: lo.ToPtr(true),
				HasOrganizationWith: []*graphclient.OrganizationWhereInput{
					{
						IDNeq:       lo.ToPtr(w.Config.SystemAdminOrgID),
						PersonalOrg: lo.ToPtr(false),
						HasOrgSubscriptionsWith: []*graphclient.OrgSubscriptionWhereInput{
							{
								Or: []*graphclient.OrgSubscriptionWhereInput{
									{
										Active: lo.ToPtr(true),
									},
									{
										StripeSubscriptionStatus: lo.ToPtr(stripeSubscriptionStatusTrialing),
									},
								},
							},
						},
					},
				},
			}, nil)
		if err != nil {
			logger.Error().Err(err).Msg("failed to fetch recovered organizations pending deletion")
			return err
		}

		if len(settings.OrganizationSettings.Edges) == 0 {
			break
		}

		for _, edge := range settings.OrganizationSettings.Edges {
			if edge == nil || edge.Node == nil || edge.Node.Organization == nil {
				continue
			}

			orgLogger := logger.With().
				Str("organization_id", edge.Node.Organization.ID).
				Str("setting_id", edge.Node.ID).
				Logger()

			if _, err := w.olClient.UpdateOrganizationSetting(ctx, edge.Node.ID, graphclient.UpdateOrganizationSettingInput{
				ClearPendingDeletionAt: lo.ToPtr(true),
			}); err != nil {
				orgLogger.Error().Err(err).Msg("failed to clear organization pending deletion state")
				return err
			}

			orgLogger.Info().Msg("cleared organization pending deletion state because billing status recovered")
		}

		if !settings.OrganizationSettings.PageInfo.HasNextPage {
			break
		}

		after = settings.OrganizationSettings.PageInfo.EndCursor
	}

	return nil
}
