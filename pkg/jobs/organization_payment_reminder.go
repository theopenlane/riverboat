package jobs

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/riverqueue/river"
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
	errDeletionDaysTooLow          = errors.New("deletion days must be at least 1 day")
	errPaymentMethodIntervalTooLow = errors.New("payment method interval must be atleast 1 day")
)

// OrganizationPaymentReminderConfig contains the configuration for the organization payment reminder worker.
type OrganizationPaymentReminderConfig struct {
	OpenlaneConfig `koanf:",squash" jsonschema:"description=the openlane API configuration for organization payment reminders"`

	// PaymentMethodInterval is the amount of days an org must have a payment method attached or else it will be earmarked for deletion
	// This is after org creation. So if an org is created 7 days ago and this is set to 6 days, the org will be marked
	// as pending deletion. But if set to say 8 days, nothing happens
	PaymentMethodInterval uint8 `koanf:"paymentmethodinterval" json:"paymentmethodinterval" jsonschema:"required,default=30 description=the number of days after organization creation before deletion is queued"`

	// DeletionDays is the number of days an org has before the deletion actually occurs. Once an org is earmarked for
	// deletion, we do not delete immediately, instead we send them an email and update "pending_deletion_at". SO if
	// DeletionDays is set to 30, the org will be deleted at in 30 days ( pending_deletion_at set to today + 30 days)
	DeletionDays uint8 `koanf:"deletiondays" json:"deletiondays" jsonschema:"required,default=7 description=the number of days before an organization pending deletion is executed"`

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

	if w.Config.DeletionDays <= 0 {
		return errDeletionDaysTooLow
	}

	if w.Config.PaymentMethodInterval <= 0 {
		return errPaymentMethodIntervalTooLow
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
		return errors.New("river client is not set on worker") //nolint:err113
	}

	var (
		after            *string
		emailQueueOffset int
	)

	for {
		settings, err := w.olClient.GetOrganizationSettings(ctx, &defaultPageSize, nil, after, nil,
			&graphclient.OrganizationSettingWhereInput{
				PendingDeletionAtIsNil: lo.ToPtr(true),
				PaymentMethodAdded:     lo.ToPtr(false),
				HasOrganizationWith: []*graphclient.OrganizationWhereInput{
					{
						PersonalOrg: lo.ToPtr(false),
						Not: &graphclient.OrganizationWhereInput{
							HasOrgSubscriptionsWith: []*graphclient.OrgSubscriptionWhereInput{
								{
									Active: lo.ToPtr(true),
								},
							},
						},
					},
				},
			}, nil)
		if err != nil {
			logger.Error().Err(err).Msg("failed to fetch organizations")
			return err
		}

		if len(settings.OrganizationSettings.Edges) == 0 {
			break
		}

		for _, edge := range settings.OrganizationSettings.Edges {
			if edge == nil || edge.Node == nil {
				continue
			}

			logger.Info().
				Str("organization_id", edge.Node.Organization.ID).
				Str("setting_id", edge.Node.ID).
				Msg("processing organization setting")

			if !isPastPaymentIntervalTimeline(edge.Node.Organization.CreatedAt, w.Config.PaymentMethodInterval) {
				logger.Info().
					Str("organization_id", edge.Node.Organization.ID).
					Str("setting_id", edge.Node.ID).
					Uint8("payment_method_interval_days", w.Config.PaymentMethodInterval).
					Msg("skipping organization before payment reminder interval")

				continue
			}

			if w.Config.DryRun {
				logger.Info().
					Str("organization_id", edge.Node.Organization.ID).
					Str("organization_name", edge.Node.Organization.Name).
					Msg("dry run: this organization would be scheduled for deletion")

				continue
			}

			pendingDeletionAt := time.Now().AddDate(0, 0, int(w.Config.DeletionDays))

			_, err = w.olClient.UpdateOrganizationSetting(ctx, edge.Node.ID, graphclient.UpdateOrganizationSettingInput{
				PendingDeletionAt: lo.ToPtr(models.DateTime(pendingDeletionAt)),
			})
			if err != nil {
				logger.Error().
					Err(err).
					Str("organization_id", edge.Node.Organization.ID).
					Str("setting_id", edge.Node.ID).
					Msg("failed to update organization settings pending deletion state")

				return err
			}

			if !w.Config.Email.Enabled {
				continue
			}

			members, err := w.olClient.GetOrgMembersByOrgID(ctx, &graphclient.OrgMembershipWhereInput{
				OrganizationID: lo.ToPtr(edge.Node.Organization.ID),
				RoleIn:         []enums.Role{enums.RoleAdmin, enums.RoleOwner},
			})
			if err != nil {
				logger.Error().
					Err(err).
					Str("organization_id", edge.Node.Organization.ID).
					Str("setting_id", edge.Node.ID).
					Msg("failed to fetch organization admin members")

				return err
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

			recipients = lo.UniqBy(recipients, func(r emailtemplates.Recipient) string {
				return strings.ToLower(strings.TrimSpace(r.Email))
			})

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
					ScheduledAt: time.Now().Add(time.Duration(emailQueueOffset) * reminderStaggerDifference),
				})
				if err != nil {
					logger.Error().
						Err(err).
						Str("organization_id", edge.Node.Organization.ID).
						Str("setting_id", edge.Node.ID).
						Msg("failed to insert email job for organization deletion notice")

					return err
				}

				emailQueueOffset++
			}

			logger.Info().
				Str("organization_id", edge.Node.Organization.ID).
				Str("setting_id", edge.Node.ID).
				Msg("sent notification to owners and admins about deletion")
		}

		if !settings.OrganizationSettings.PageInfo.HasNextPage {
			break
		}

		after = settings.OrganizationSettings.PageInfo.EndCursor
	}

	return nil
}

func isPastPaymentIntervalTimeline(createdAt *time.Time, intervalDays uint8) bool {
	if createdAt == nil {
		return false
	}

	if intervalDays == 0 {
		return true
	}

	return time.Since(*createdAt) >= time.Duration(intervalDays)*24*time.Hour
}
