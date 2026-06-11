package jobs_test

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/samber/lo"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/theopenlane/core/common/enums"
	"github.com/theopenlane/core/common/jobspec"
	"github.com/theopenlane/core/common/models"
	"github.com/theopenlane/go-client/graphclient"
	"github.com/theopenlane/riverboat/pkg/jobs"
	olmocks "github.com/theopenlane/riverboat/pkg/jobs/openlane/mocks"
	rivermocks "github.com/theopenlane/riverboat/pkg/riverqueue/mocks"
)

func TestOrganizationDeleteWorker(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	olMock := olmocks.NewMockGraphClient(t)

	pendingDeletionAt := models.DateTime(time.Now().Add(-24 * time.Hour))
	recoveredSettings := organizationSettingsResponse(
		organizationSettingEdge("setting-paid", "org-paid", "paid-org", true, nil, &pendingDeletionAt),
	)
	deletableSettings := organizationSettingsResponse(
		organizationSettingEdge("setting-unpaid", "org-unpaid", "unpaid-org", false, nil, &pendingDeletionAt),
	)

	olMock.EXPECT().
		GetOrganizationSettings(mock.Anything, mock.Anything, (*int64)(nil), (*string)(nil), (*string)(nil), mock.MatchedBy(func(where *graphclient.OrganizationSettingWhereInput) bool {
			return where != nil &&
				where.PendingDeletionAtNotNil != nil && *where.PendingDeletionAtNotNil &&
				where.PendingDeletionAtLte == nil &&
				len(where.HasOrganizationWith) == 1 &&
				where.HasOrganizationWith[0].IDNeq != nil && *where.HasOrganizationWith[0].IDNeq == "system-admin" &&
				where.HasOrganizationWith[0].PersonalOrg != nil &&
				!*where.HasOrganizationWith[0].PersonalOrg &&
				len(where.HasOrganizationWith[0].HasOrgSubscriptionsWith) == 1 &&
				where.HasOrganizationWith[0].HasOrgSubscriptionsWith[0].Active != nil &&
				*where.HasOrganizationWith[0].HasOrgSubscriptionsWith[0].Active
		}), ([]*graphclient.OrganizationSettingOrder)(nil)).
		Return(recoveredSettings, nil).
		Once()

	olMock.EXPECT().
		GetOrganizationSettings(mock.Anything, (*int64)(nil), mock.MatchedBy(func(last *int64) bool {
			return last != nil && *last == 2
		}), (*string)(nil), (*string)(nil), mock.MatchedBy(func(where *graphclient.OrganizationSettingWhereInput) bool {
			return where != nil &&
				where.PendingDeletionAtNotNil != nil && *where.PendingDeletionAtNotNil &&
				where.PendingDeletionAtLte != nil &&
				len(where.HasOrganizationWith) == 1 &&
				where.HasOrganizationWith[0].IDNeq != nil && *where.HasOrganizationWith[0].IDNeq == "system-admin" &&
				where.HasOrganizationWith[0].Not != nil &&
				len(where.HasOrganizationWith[0].Not.HasOrgSubscriptionsWith) == 1 &&
				where.HasOrganizationWith[0].Not.HasOrgSubscriptionsWith[0].Active != nil &&
				*where.HasOrganizationWith[0].Not.HasOrgSubscriptionsWith[0].Active
		}), mock.MatchedBy(func(orderBy []*graphclient.OrganizationSettingOrder) bool {
			return len(orderBy) == 1 &&
				orderBy[0] != nil &&
				orderBy[0].Field == graphclient.OrganizationSettingOrderFieldUpdatedAt &&
				orderBy[0].Direction == graphclient.OrderDirectionAsc
		})).
		Return(deletableSettings, nil).
		Once()

	olMock.EXPECT().
		UpdateOrganizationSetting(mock.Anything, "setting-paid", mock.MatchedBy(func(input graphclient.UpdateOrganizationSettingInput) bool {
			return input.ClearPendingDeletionAt != nil && *input.ClearPendingDeletionAt
		})).
		Return(&graphclient.UpdateOrganizationSetting{}, nil).
		Once()

	olMock.EXPECT().
		DeleteOrganization(mock.Anything, "org-unpaid").
		Return(&graphclient.DeleteOrganization{}, nil).
		Once()

	worker := &jobs.OrganizationDeleteWorker{
		Config: jobs.OrganizationDeleteConfig{
			MaxDeletesPerRun: 2,
			SystemAdminOrgID: "system-admin",
		},
	}
	worker.WithOpenlaneClient(olMock)

	err := worker.Work(ctx, &river.Job[jobspec.OrganizationDeletionArgs]{
		JobRow: &rivertype.JobRow{Kind: jobspec.OrganizationDeletionArgs{}.Kind()},
	})
	require.NoError(t, err)
}

func TestOrganizationPaymentReminderWorker(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	olMock := olmocks.NewMockGraphClient(t)
	riverMock := rivermocks.NewMockJobClient(t)
	insertScheduledAt := make([]time.Time, 0, 2)

	settings := organizationSettingsResponse(
		organizationSettingEdge("setting-1", "org-1", "acme", true, nil, nil),
	)

	olMock.EXPECT().
		GetOrganizationSettings(mock.Anything, mock.Anything, (*int64)(nil), (*string)(nil), (*string)(nil), mock.MatchedBy(func(where *graphclient.OrganizationSettingWhereInput) bool {
			return where != nil &&
				where.PendingDeletionAtIsNil != nil && *where.PendingDeletionAtIsNil &&
				where.PaymentMethodAdded == nil &&
				len(where.HasOrganizationWith) == 1 &&
				where.HasOrganizationWith[0].PersonalOrg != nil &&
				!*where.HasOrganizationWith[0].PersonalOrg &&
				where.HasOrganizationWith[0].IDNeq != nil && *where.HasOrganizationWith[0].IDNeq == "system-admin" &&
				where.HasOrganizationWith[0].Not != nil &&
				len(where.HasOrganizationWith[0].Not.HasOrgSubscriptionsWith) == 1 &&
				where.HasOrganizationWith[0].Not.HasOrgSubscriptionsWith[0].Active != nil &&
				*where.HasOrganizationWith[0].Not.HasOrgSubscriptionsWith[0].Active &&
				len(where.HasOrganizationWith[0].HasOrgSubscriptionsWith) == 1 &&
				where.HasOrganizationWith[0].HasOrgSubscriptionsWith[0].Active != nil &&
				!*where.HasOrganizationWith[0].HasOrgSubscriptionsWith[0].Active &&
				where.HasOrganizationWith[0].HasOrgSubscriptionsWith[0].UpdatedAtLte != nil &&
				where.HasOrganizationWith[0].HasOrgSubscriptionsWith[0].TrialExpiresAtNotNil == nil
		}), ([]*graphclient.OrganizationSettingOrder)(nil)).
		Return(settings, nil).
		Once()

	olMock.EXPECT().
		UpdateOrganizationSetting(mock.Anything, "setting-1", mock.MatchedBy(func(input graphclient.UpdateOrganizationSettingInput) bool {
			if input.PendingDeletionAt == nil {
				return false
			}

			pendingDeletionAt := time.Time(*input.PendingDeletionAt)
			return pendingDeletionAt.After(time.Now().Add(6*24*time.Hour)) &&
				pendingDeletionAt.Before(time.Now().Add(8*24*time.Hour))
		})).
		Return(&graphclient.UpdateOrganizationSetting{}, nil).
		Once()

	olMock.EXPECT().
		GetOrgMembersByOrgID(mock.Anything, mock.MatchedBy(func(where *graphclient.OrgMembershipWhereInput) bool {
			return where != nil &&
				where.OrganizationID != nil && *where.OrganizationID == "org-1" &&
				len(where.RoleIn) == 2 &&
				where.RoleIn[0] == enums.RoleAdmin &&
				where.RoleIn[1] == enums.RoleOwner
		})).
		Return(&graphclient.GetOrgMembersByOrgID{
			OrgMemberships: graphclient.GetOrgMembersByOrgID_OrgMemberships{
				Edges: []*graphclient.GetOrgMembersByOrgID_OrgMemberships_Edges{
					{
						Node: &graphclient.GetOrgMembersByOrgID_OrgMemberships_Edges_Node{
							ID:             "membership-1",
							OrganizationID: "org-1",
							Role:           enums.RoleAdmin,
							UserID:         "user-1",
							User: graphclient.GetOrgMembersByOrgID_OrgMemberships_Edges_Node_User{
								Email:     "admin@example.com",
								FirstName: lo.ToPtr("Ada"),
								ID:        "user-1",
								LastName:  lo.ToPtr("Lovelace"),
							},
						},
					},
					{
						Node: &graphclient.GetOrgMembersByOrgID_OrgMemberships_Edges_Node{
							ID:             "membership-2",
							OrganizationID: "org-1",
							Role:           enums.RoleOwner,
							UserID:         "user-2",
							User: graphclient.GetOrgMembersByOrgID_OrgMemberships_Edges_Node_User{
								Email: "owner@example.com",
								ID:    "user-2",
							},
						},
					},
				},
			},
		}, nil).
		Once()

	riverMock.EXPECT().
		Insert(mock.Anything, mock.MatchedBy(func(args jobs.EmailArgs) bool {
			return len(args.Message.To) == 1 &&
				args.Message.Subject == "Organization Deletion Notice for acme"
		}), mock.MatchedBy(func(opts *river.InsertOpts) bool {
			return opts != nil && !opts.ScheduledAt.IsZero()
		})).
		Run(func(_ context.Context, _ river.JobArgs, opts *river.InsertOpts) {
			insertScheduledAt = append(insertScheduledAt, opts.ScheduledAt)
		}).
		Return(&rivertype.JobInsertResult{}, nil).
		Twice()

	worker := &jobs.OrganizationPaymentReminderWorker{
		Config: jobs.OrganizationPaymentReminderConfig{
			OrgDeletionAfterCancelDays: 1,
			DeletionDays:               7,
			SystemAdminOrgID:           "system-admin",
		},
	}
	worker.Config.Email.Enabled = true
	worker.WithOpenlaneClient(olMock)
	worker.WithRiverClient(riverMock)

	err := worker.Work(ctx, &river.Job[jobspec.OrganizationDeletionReminderArgs]{
		JobRow: &rivertype.JobRow{Kind: jobspec.OrganizationDeletionReminderArgs{}.Kind()},
	})
	require.NoError(t, err)
	require.Len(t, insertScheduledAt, 2)
	require.True(t, slices.IsSortedFunc(insertScheduledAt, func(a, b time.Time) int {
		switch {
		case a.Before(b):
			return -1
		case a.After(b):
			return 1
		default:
			return 0
		}
	}), "email insert scheduled times should be monotonic")
	require.True(t, insertScheduledAt[0].After(time.Now().Add(-5*time.Second)))
	require.True(t, insertScheduledAt[1].After(insertScheduledAt[0]))
}

func organizationSettingsResponse(edges ...*graphclient.GetOrganizationSettings_OrganizationSettings_Edges) *graphclient.GetOrganizationSettings {
	return &graphclient.GetOrganizationSettings{
		OrganizationSettings: graphclient.GetOrganizationSettings_OrganizationSettings{
			Edges: edges,
			PageInfo: graphclient.GetOrganizationSettings_OrganizationSettings_PageInfo{
				HasNextPage: false,
			},
			TotalCount: int64(len(edges)),
		},
	}
}

func organizationSettingEdge(settingID, orgID, orgName string, paymentMethodAdded bool, createdAt *time.Time, pendingDeletionAt *models.DateTime) *graphclient.GetOrganizationSettings_OrganizationSettings_Edges {
	return &graphclient.GetOrganizationSettings_OrganizationSettings_Edges{
		Node: &graphclient.GetOrganizationSettings_OrganizationSettings_Edges_Node{
			ID:                 settingID,
			PaymentMethodAdded: paymentMethodAdded,
			PendingDeletionAt:  pendingDeletionAt,
			Organization: &graphclient.GetOrganizationSettings_OrganizationSettings_Edges_Node_Organization{
				CreatedAt: createdAt,
				ID:        orgID,
				Name:      orgName,
			},
		},
	}
}
