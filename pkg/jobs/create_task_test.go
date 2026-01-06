package jobs_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/riverqueue/river"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/theopenlane/core/common/jobspec"
	"github.com/theopenlane/core/common/models"
	"github.com/theopenlane/go-client/graphclient"
	"github.com/theopenlane/riverboat/pkg/jobs"

	olmocks "github.com/theopenlane/riverboat/pkg/jobs/openlane/mocks"
)

func TestCreateTaskWorker(t *testing.T) {
	t.Parallel()

	orgID := "01JA0000000000000000000001"
	taskID := "01JB0000000000000000000001"
	assigneeID := "01JC0000000000000000000001"
	assignerID := "01JC0000000000000000000002"
	policyID := "01JD0000000000000000000001"
	title := "Information Security Policy Review"
	description := "Conduct the annual review of this internal policy"
	tag1 := "security"
	tag2 := "compliance"

	testCases := []struct {
		name               string
		input              jobspec.CreateTaskArgs
		createTaskResponse *graphclient.CreateTask
		createTaskError    error
		expectedError      string
		expectCreateTask   bool
	}{
		{
			name: "happy path - all fields provided",
			input: jobspec.CreateTaskArgs{
				OrganizationID:    orgID,
				Title:             title,
				Description:       description,
				AssigneeID:        &assigneeID,
				AssignerID:        &assignerID,
				InternalPolicyIDs: []string{policyID},
				Tags:              []string{tag1, tag2},
			},
			createTaskResponse: &graphclient.CreateTask{
				CreateTask: graphclient.CreateTask_CreateTask{
					Task: graphclient.CreateTask_CreateTask_Task{
						ID:    taskID,
						Title: title,
					},
				},
			},
			expectCreateTask: true,
		},
		{
			name: "happy path - minimal fields (only required)",
			input: jobspec.CreateTaskArgs{
				OrganizationID: orgID,
				Title:          "Simple Task",
				Description:    "A simple task with only required fields",
			},
			createTaskResponse: &graphclient.CreateTask{
				CreateTask: graphclient.CreateTask_CreateTask{
					Task: graphclient.CreateTask_CreateTask_Task{
						ID:    taskID,
						Title: "Simple Task",
					},
				},
			},
			expectCreateTask: true,
		},
		{
			name: "happy path - with internal policy link",
			input: jobspec.CreateTaskArgs{
				OrganizationID:    orgID,
				Title:             title,
				Description:       description,
				InternalPolicyIDs: []string{policyID},
			},
			createTaskResponse: &graphclient.CreateTask{
				CreateTask: graphclient.CreateTask_CreateTask{
					Task: graphclient.CreateTask_CreateTask_Task{
						ID:    taskID,
						Title: title,
					},
				},
			},
			expectCreateTask: true,
		},
		{
			name: "happy path - with assignee and assigner",
			input: jobspec.CreateTaskArgs{
				OrganizationID: orgID,
				Title:          "Assigned Task",
				Description:    "Task with assignee and assigner",
				AssigneeID:     &assigneeID,
				AssignerID:     &assignerID,
			},
			createTaskResponse: &graphclient.CreateTask{
				CreateTask: graphclient.CreateTask_CreateTask{
					Task: graphclient.CreateTask_CreateTask_Task{
						ID:    taskID,
						Title: "Assigned Task",
					},
				},
			},
			expectCreateTask: true,
		},
		{
			name: "happy path - with tags",
			input: jobspec.CreateTaskArgs{
				OrganizationID: orgID,
				Title:          "Tagged Task",
				Description:    "Task with tags",
				Tags:           []string{tag1, tag2},
			},
			createTaskResponse: &graphclient.CreateTask{
				CreateTask: graphclient.CreateTask_CreateTask{
					Task: graphclient.CreateTask_CreateTask_Task{
						ID:    taskID,
						Title: "Tagged Task",
					},
				},
			},
			expectCreateTask: true,
		},
		{
			name: "validation error - missing organization ID",
			input: jobspec.CreateTaskArgs{
				Title:       title,
				Description: description,
			},
			expectedError: "organization_id is required",
		},
		{
			name: "validation error - missing title",
			input: jobspec.CreateTaskArgs{
				OrganizationID: orgID,
				Description:    description,
			},
			expectedError: "title is required",
		},
		{
			name: "validation error - missing description",
			input: jobspec.CreateTaskArgs{
				OrganizationID: orgID,
				Title:          title,
			},
			expectedError: "description is required",
		},
		{
			name: "API error - task creation fails",
			input: jobspec.CreateTaskArgs{
				OrganizationID: orgID,
				Title:          title,
				Description:    description,
			},
			createTaskError:  errors.New("API error: database connection failed"),
			expectedError:    "API error",
			expectCreateTask: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock client
			olMock := olmocks.NewMockGraphClient(t)

			// Set up expectations
			if tc.expectCreateTask {
				olMock.EXPECT().CreateTask(
					mock.MatchedBy(func(ctx context.Context) bool {
						return ctx != nil
					}),
					mock.MatchedBy(func(input graphclient.CreateTaskInput) bool {
						// Validate basic required fields
						if input.Title != tc.input.Title {
							return false
						}
						if *input.Details != tc.input.Description {
							return false
						}
						if *input.OwnerID != tc.input.OrganizationID {
							return false
						}

						// Validate optional fields
						if tc.input.AssigneeID != nil {
							if input.AssigneeID == nil || *input.AssigneeID != *tc.input.AssigneeID {
								return false
							}
						}

						if tc.input.AssignerID != nil {
							if input.AssignerID == nil || *input.AssignerID != *tc.input.AssignerID {
								return false
							}
						}

						if len(tc.input.InternalPolicyIDs) > 0 {
							if len(input.InternalPolicyIDs) != len(tc.input.InternalPolicyIDs) {
								return false
							}
							for i, id := range tc.input.InternalPolicyIDs {
								if input.InternalPolicyIDs[i] != id {
									return false
								}
							}
						}

						if len(tc.input.Tags) > 0 {
							if len(input.Tags) != len(tc.input.Tags) {
								return false
							}
							for i, tag := range tc.input.Tags {
								if input.Tags[i] != tag {
									return false
								}
							}
						}

						return true
					}),
				).Return(tc.createTaskResponse, tc.createTaskError)
			}

			// Create worker with config
			worker := &jobs.CreateTaskWorker{
				Config: jobs.TaskWorkerConfig{
					OpenlaneConfig: jobs.OpenlaneConfig{
						OpenlaneAPIHost:  "https://api.example.com",
						OpenlaneAPIToken: "tola_test-token",
					},
					Enabled: true,
				},
			}

			// Inject mock client
			worker.WithOpenlaneClient(olMock)

			// Execute
			ctx := context.Background()
			err := worker.Work(ctx, &river.Job[jobspec.CreateTaskArgs]{
				Args: tc.input,
			})

			// Assert
			if tc.expectedError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCreateTaskWorker_WithDueDate(t *testing.T) {
	t.Parallel()

	orgID := "01JA0000000000000000000001"
	taskID := "01JB0000000000000000000001"
	title := "Task with due date"
	description := "This task has a due date"

	dueDate := models.DateTime(time.Now().Add(7 * 24 * time.Hour))

	// Create mock client
	olMock := olmocks.NewMockGraphClient(t)

	olMock.EXPECT().CreateTask(
		mock.MatchedBy(func(ctx context.Context) bool {
			return ctx != nil
		}),
		mock.MatchedBy(func(input graphclient.CreateTaskInput) bool {
			return input.Title == title &&
				*input.Details == description &&
				*input.OwnerID == orgID &&
				input.Due != nil &&
				*input.Due == dueDate
		}),
	).Return(&graphclient.CreateTask{
		CreateTask: graphclient.CreateTask_CreateTask{
			Task: graphclient.CreateTask_CreateTask_Task{
				ID:    taskID,
				Title: title,
			},
		},
	}, nil)

	// Create worker with config
	worker := &jobs.CreateTaskWorker{
		Config: jobs.TaskWorkerConfig{
			OpenlaneConfig: jobs.OpenlaneConfig{
				OpenlaneAPIHost:  "https://api.example.com",
				OpenlaneAPIToken: "tola_test-token",
			},
			Enabled: true,
		},
	}

	// Inject mock client
	worker.WithOpenlaneClient(olMock)

	// Execute
	ctx := context.Background()
	err := worker.Work(ctx, &river.Job[jobspec.CreateTaskArgs]{
		Args: jobspec.CreateTaskArgs{
			OrganizationID: orgID,
			Title:          title,
			Description:    description,
			DueDate:        &dueDate,
		},
	})

	// Assert
	require.NoError(t, err)
}

func TestCreateTaskArgs_Kind(t *testing.T) {
	t.Parallel()

	args := jobspec.CreateTaskArgs{}
	require.Equal(t, "create_task", args.Kind())
}

func TestCreateTaskArgs_InsertOpts(t *testing.T) {
	t.Parallel()

	t.Run("without scheduled time", func(t *testing.T) {
		args := jobspec.CreateTaskArgs{}
		opts := jobs.InsertOpts(args)

		require.Equal(t, 3, opts.MaxAttempts)
		require.True(t, opts.ScheduledAt.IsZero())
	})

	t.Run("with scheduled time", func(t *testing.T) {
		scheduledTime := time.Now().Add(1 * time.Hour)
		args := jobspec.CreateTaskArgs{
			ScheduledAt: &scheduledTime,
		}
		opts := jobs.InsertOpts(args)

		require.Equal(t, 3, opts.MaxAttempts)
		require.Equal(t, scheduledTime, opts.ScheduledAt)
	})
}
