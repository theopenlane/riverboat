package jobs

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/riverqueue/river"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/theopenlane/core/common/enums"
	"github.com/theopenlane/core/common/jobspec"
	"github.com/theopenlane/core/common/models"
	"github.com/theopenlane/emailtemplates"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/theopenlane/riverboat/pkg/jobs/openlane"
	"github.com/theopenlane/riverboat/pkg/riverqueue"
)

const reminderStaggerDifference = 30 * time.Second

var (
	errDeletionDaysTooLow            = errors.New("deletion days must be at least 1 day")
	errOrgDeletionAfterCancelDaysLow = errors.New("org deletion after cancel days must be at least 1 day")
	errRiverClientRequired           = errors.New("river client is not set on worker")
	errSlackChannelRequired          = errors.New("slack channel is required")
	errSystemAdminOrgIDRequired      = errors.New("system admin org id is required")
)

// OrganizationPaymentReminderConfig contains the configuration for the organization payment reminder worker.
type OrganizationPaymentReminderConfig struct {
	OpenlaneConfig `koanf:",squash" jsonschema:"description=the openlane API configuration for organization payment reminders"`

	// OrgDeletionAfterCancelDays is the number of days after a previously active organization's subscription is canceled before it is queued for deletion
	OrgDeletionAfterCancelDays uint8 `koanf:"orgdeletionaftercanceldays" json:"orgdeletionaftercanceldays" jsonschema:"required,default=30 description=the number of days after a previously active organization's subscription cancellation before deletion is queued"`

	// DeletionDays is the number of days an org has before the deletion actually occurs. Once an org is earmarked for
	// deletion, we do not delete immediately, instead we send them an email and update "pending_deletion_at". SO if
	// DeletionDays is set to 30, the org will be deleted at in 30 days ( pending_deletion_at set to today + 30 days)
	DeletionDays uint8 `koanf:"deletiondays" json:"deletiondays" jsonschema:"required,default=7 description=the number of days before an organization pending deletion is executed"`

	// SystemAdminOrgID is the organization ID that should never be queued for deletion
	SystemAdminOrgID string `koanf:"systemadminorgid" json:"systemadminorgid" default:"01101101011010010111010001100010" jsonschema:"description=organization ID that must never be queued for deletion,default=01101101011010010111010001100010"`

	// SlackChannel is the channel name where organization deletion reminders summaries will be sent to
	SlackChannel string `koanf:"slackchannel" json:"slackchannel" default:"customers" jsonschema:"description=slack channel for organization deletion summaries,default=customers"`

	// Enabled is used to determine if to register this worker or not
	Enabled bool `koanf:"enabled" json:"enabled" jsonschema:"required description=whether the organization payment reminder worker is enabled"`

	// DryRun prints matching organization IDs instead of mutating state or queueing emails - if enabled.
	DryRun bool `koanf:"dryrun" json:"dryrun" jsonschema:"description=if true, only print organization IDs that would be processed"`

	Email struct {
		Enabled bool                  `koanf:"enabled" json:"enabled" jsonschema:"required description=whether to send emails about the scheduled delution"`
		Config  emailtemplates.Config `json:"config" koanf:"config" jsonschema:"required description=email settings"`
	} `json:"email" koanf:"email" jsonschema:"required description=email configuration"`
}

// OrganizationPaymentReminderWorker fetches organizations for payment reminder processing.
type OrganizationPaymentReminderWorker struct {
	river.WorkerDefaults[jobspec.OrganizationDeletionReminderArgs]

	Config OrganizationPaymentReminderConfig `koanf:"config" json:"config" jsonschema:"description=the configuration for organization payment reminders"`

	olClient    openlane.GraphClient
	riverClient riverqueue.JobClient
}

// WithOpenlaneClient sets the Openlane client for the worker.
func (w *OrganizationPaymentReminderWorker) WithOpenlaneClient(cl openlane.GraphClient) *OrganizationPaymentReminderWorker {
	w.olClient = cl
	return w
}

// WithRiverClient sets the River client for the worker.
func (w *OrganizationPaymentReminderWorker) WithRiverClient(cl riverqueue.JobClient) *OrganizationPaymentReminderWorker {
	w.riverClient = cl
	return w
}

// Work satisfies the river.Worker interface for the organization payment reminder worker.
func (w *OrganizationPaymentReminderWorker) Work(ctx context.Context, job *river.Job[jobspec.OrganizationDeletionReminderArgs]) error {
	logger := log.Ctx(ctx).With().Str("job_kind", job.Kind).Logger()
	logger.Info().Msg("starting organization payment reminder job")

	cancelDays, err := w.validateConfig()
	if err != nil {
		return err
	}

	if w.olClient == nil {
		cl, err := w.Config.getOpenlaneClient()
		if err != nil {
			return err
		}

		w.olClient = cl
	}

	if w.riverClient == nil {
		logger.Error().Msg("river client is not set on worker, cannot insert organization deletion jobs")
		return errRiverClientRequired
	}

	var (
		emailQueueOffset    int
		now                 = time.Now()
		pendingOrgsToDelete []string
	)

	where := &graphclient.OrganizationSettingWhereInput{
		PendingDeletionAtIsNil: lo.ToPtr(true),
		HasOrganizationWith: []*graphclient.OrganizationWhereInput{
			{
				IDNeq:       lo.ToPtr(w.Config.SystemAdminOrgID),
				PersonalOrg: lo.ToPtr(false),
				Not: &graphclient.OrganizationWhereInput{
					HasOrgSubscriptionsWith: []*graphclient.OrgSubscriptionWhereInput{
						{
							Active: lo.ToPtr(true),
						},
					},
				},
				HasOrgSubscriptionsWith: []*graphclient.OrgSubscriptionWhereInput{
					{
						Active:       lo.ToPtr(false),
						UpdatedAtLte: lo.ToPtr(now.Add(-time.Duration(cancelDays) * 24 * time.Hour)),
					},
				},
			},
		},
	}

	var after *string

	for {
		settings, err := w.olClient.GetOrganizationSettings(ctx, &defaultPageSize, nil, after, nil, where, nil)
		if err != nil {
			logger.Error().Err(err).Msg("failed to fetch organizations")
			return err
		}

		if len(settings.OrganizationSettings.Edges) == 0 {
			break
		}

		for _, edge := range settings.OrganizationSettings.Edges {
			summary, isProcessed, err := w.processOrganization(ctx, logger, edge, &emailQueueOffset)
			if err != nil {
				return err
			}

			if isProcessed {
				pendingOrgsToDelete = append(pendingOrgsToDelete, summary)
			}
		}

		if !settings.OrganizationSettings.PageInfo.HasNextPage {
			break
		}

		after = settings.OrganizationSettings.PageInfo.EndCursor
	}

	return w.createReminderSummary(ctx, logger, pendingOrgsToDelete)
}

func (w *OrganizationPaymentReminderWorker) validateConfig() (uint8, error) {
	if w.Config.DeletionDays <= 0 {
		return 0, errDeletionDaysTooLow
	}

	w.Config.SystemAdminOrgID = strings.TrimSpace(w.Config.SystemAdminOrgID)
	if w.Config.SystemAdminOrgID == "" {
		return 0, errSystemAdminOrgIDRequired
	}

	w.Config.SlackChannel = strings.TrimSpace(w.Config.SlackChannel)
	if w.Config.SlackChannel == "" {
		return 0, errSlackChannelRequired
	}

	cancelDays := w.Config.OrgDeletionAfterCancelDays
	if cancelDays <= 0 {
		return 0, errOrgDeletionAfterCancelDaysLow
	}

	return cancelDays, nil
}

func (w *OrganizationPaymentReminderWorker) processOrganization(ctx context.Context, logger zerolog.Logger, edge *graphclient.GetOrganizationSettings_OrganizationSettings_Edges, emailQueueOffset *int) (string, bool, error) {
	if edge == nil || edge.Node == nil || edge.Node.Organization == nil {
		return "", false, nil
	}

	summary := fmt.Sprintf("%s (%s)", edge.Node.Organization.Name, edge.Node.Organization.ID)

	logger.Info().
		Str("organization_id", edge.Node.Organization.ID).
		Str("setting_id", edge.Node.ID).
		Msg("processing organization setting")

	if w.Config.DryRun {
		logger.Info().
			Str("organization_id", edge.Node.Organization.ID).
			Str("organization_name", edge.Node.Organization.Name).
			Msg("dry run: this organization would be scheduled for deletion - Dry Run")

		return summary, true, nil
	}

	pendingDeletionAt := time.Now().AddDate(0, 0, int(w.Config.DeletionDays))

	_, err := w.olClient.UpdateOrganizationSetting(ctx, edge.Node.ID, graphclient.UpdateOrganizationSettingInput{
		PendingDeletionAt: lo.ToPtr(models.DateTime(pendingDeletionAt)),
	})
	if err != nil {
		logger.Error().
			Err(err).
			Str("organization_id", edge.Node.Organization.ID).
			Str("setting_id", edge.Node.ID).
			Msg("failed to update organization settings pending deletion state")

		return "", false, err
	}

	if !w.Config.Email.Enabled {
		return summary, true, nil
	}

	if err := w.sendReminder(ctx, logger, edge, pendingDeletionAt, emailQueueOffset); err != nil {
		return "", false, err
	}

	logger.Info().
		Str("organization_id", edge.Node.Organization.ID).
		Str("setting_id", edge.Node.ID).
		Msg("sent notification to owners and admins about deletion")

	return summary, true, nil
}

func (w *OrganizationPaymentReminderWorker) sendReminder(ctx context.Context, logger zerolog.Logger, edge *graphclient.GetOrganizationSettings_OrganizationSettings_Edges, pendingDeletionAt time.Time, emailQueueOffset *int) error {
	recipients, err := w.getReminderRecipients(ctx, edge)
	if err != nil {
		logger.Error().
			Err(err).
			Str("organization_id", edge.Node.Organization.ID).
			Str("setting_id", edge.Node.ID).
			Msg("failed to fetch organization admin members")

		return err
	}

	for _, recipient := range recipients {
		email, err := w.Config.Email.Config.NewOrgDeletionNoticeEmail(recipient, emailtemplates.OrgDeletionNoticeTemplateData{
			OrganizationName: edge.Node.Organization.Name,
			Date:             pendingDeletionAt,
		})
		if err != nil {
			logger.Error().
				Err(err).
				Str("organization_id", edge.Node.Organization.ID).
				Str("setting_id", edge.Node.ID).
				Msg("failed to create organization deletion notice email")

			return err
		}

		_, err = w.riverClient.Insert(ctx, EmailArgs{
			Message: *email,
		}, &river.InsertOpts{
			ScheduledAt: time.Now().Add(time.Duration(*emailQueueOffset) * reminderStaggerDifference),
		})
		if err != nil {
			logger.Error().
				Err(err).
				Str("organization_id", edge.Node.Organization.ID).
				Str("setting_id", edge.Node.ID).
				Msg("failed to insert email job for organization deletion notice")

			return err
		}

		*emailQueueOffset++
	}

	return nil
}

func (w *OrganizationPaymentReminderWorker) getReminderRecipients(ctx context.Context, edge *graphclient.GetOrganizationSettings_OrganizationSettings_Edges) ([]emailtemplates.Recipient, error) {
	members, err := w.olClient.GetOrgMembersByOrgID(ctx, &graphclient.OrgMembershipWhereInput{
		OrganizationID: lo.ToPtr(edge.Node.Organization.ID),
		RoleIn:         []enums.Role{enums.RoleAdmin, enums.RoleOwner},
	})
	if err != nil {
		return nil, err
	}

	recipients := make([]emailtemplates.Recipient, 0, len(members.OrgMemberships.Edges)+1)
	for _, member := range members.OrgMemberships.Edges {
		if member == nil || member.Node == nil || member.Node.User.Email == "" {
			continue
		}

		recipients = append(recipients, emailtemplates.Recipient{
			Email:     member.Node.User.Email,
			FirstName: lo.FromPtr(member.Node.User.FirstName),
			LastName:  lo.FromPtr(member.Node.User.LastName),
		})
	}

	orgBillingEmail := strings.TrimSpace(lo.FromPtr(edge.Node.BillingEmail))
	if orgBillingEmail != "" {
		recipients = append(recipients, emailtemplates.Recipient{
			Email: orgBillingEmail,
			// add default names
			FirstName: "Billing",
			LastName:  "Admin",
		})
	}

	return lo.UniqBy(recipients, func(r emailtemplates.Recipient) string {
		return strings.ToLower(strings.TrimSpace(r.Email))
	}), nil
}

func (w *OrganizationPaymentReminderWorker) createReminderSummary(ctx context.Context, logger zerolog.Logger, pendingOrgsToDelete []string) error {
	action := "set to be deleted"
	if w.Config.DryRun {
		action = "would be deleted"
	}

	message := fmt.Sprintf("Organization deletion reminder summary: %d orgs %s", len(pendingOrgsToDelete), action)
	if len(pendingOrgsToDelete) > 0 {
		message = fmt.Sprintf("%s:\n- %s", message, strings.Join(pendingOrgsToDelete, "\n- "))
	}

	if w.Config.DryRun {
		message += " - Dry Run"
	}

	if _, err := w.riverClient.Insert(ctx, SlackArgs{
		Channel: w.Config.SlackChannel,
		Message: message,
	}, nil); err != nil {
		logger.Error().Err(err).Msg("failed to insert slack job for organization deletion reminder summary")
		return err
	}

	return nil
}
