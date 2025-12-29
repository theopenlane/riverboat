package jobs

import (
	"context"

	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"

	"github.com/theopenlane/core/common/jobspec"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/theopenlane/riverboat/pkg/jobs/openlane"
)

// InsertOpts configures job insertion options including retry behavior and scheduling
func InsertOpts(a jobspec.CreateTaskArgs) river.InsertOpts {
	opts := river.InsertOpts{
		MaxAttempts: 3, // nolint:mnd
	}

	// If scheduled time is specified, set it in the options
	if a.ScheduledAt != nil {
		opts.ScheduledAt = *a.ScheduledAt
	}

	return opts
}

// TaskWorkerConfig contains configuration for the create task worker
type TaskWorkerConfig struct {
	// embed OpenlaneConfig to reuse validation and client creation logic
	OpenlaneConfig `koanf:",squash" jsonschema:"description=the openlane API configuration for task creation"`

	Enabled bool `koanf:"enabled" json:"enabled" jsonschema:"required,description=whether the task worker is enabled"`
}

// CreateTaskWorker processes create task jobs
type CreateTaskWorker struct {
	river.WorkerDefaults[jobspec.CreateTaskArgs]

	Config TaskWorkerConfig `koanf:"config" json:"config"`

	olClient openlane.GraphClient
}

// WithOpenlaneClient sets the Openlane client (used for testing)
func (w *CreateTaskWorker) WithOpenlaneClient(cl openlane.GraphClient) *CreateTaskWorker {
	w.olClient = cl
	return w
}

// Work satisfies the river.Worker interface
func (w *CreateTaskWorker) Work(ctx context.Context, job *river.Job[jobspec.CreateTaskArgs]) error {
	logger := log.With().
		Str("organization_id", job.Args.OrganizationID).
		Str("task_title", job.Args.Title).
		Logger()

	logger.Info().Msg("starting create task job")

	// Validate required fields
	if err := w.validateArgs(job.Args); err != nil {
		logger.Error().Err(err).Msg("invalid job arguments")
		return err
	}

	// Initialize Openlane client if not already set
	if w.olClient == nil {
		cl, err := w.Config.getOpenlaneClient()
		if err != nil {
			logger.Error().Err(err).Msg("failed to create openlane client")
			return err
		}

		w.olClient = cl
	}

	// Build task input from provided arguments
	taskInput := graphclient.CreateTaskInput{
		Title:   job.Args.Title,
		Details: &job.Args.Description,
		OwnerID: &job.Args.OrganizationID,
	}

	// Add optional fields if provided
	if job.Args.Category != nil {
		taskInput.Category = job.Args.Category
	}

	if job.Args.AssigneeID != nil {
		taskInput.AssigneeID = job.Args.AssigneeID
	}

	if job.Args.AssignerID != nil {
		taskInput.AssignerID = job.Args.AssignerID
	}

	if job.Args.DueDate != nil {
		taskInput.Due = job.Args.DueDate
	}

	if len(job.Args.InternalPolicyIDs) > 0 {
		taskInput.InternalPolicyIDs = job.Args.InternalPolicyIDs
	}

	if len(job.Args.Tags) > 0 {
		taskInput.Tags = job.Args.Tags
	}

	// Create the task using Openlane client
	createdTask, err := w.olClient.CreateTask(ctx, taskInput)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create task")
		return err
	}

	logger.Info().
		Str("task_id", createdTask.CreateTask.Task.ID).
		Str("task_title", createdTask.CreateTask.Task.Title).
		Msg("task created successfully")

	return nil
}

// validateArgs validates the job arguments
func (w *CreateTaskWorker) validateArgs(args jobspec.CreateTaskArgs) error {
	// Validate required fields
	if args.OrganizationID == "" {
		return newMissingRequiredArg("organization_id", jobspec.CreateTaskArgs{}.Kind())
	}

	if args.Title == "" {
		return newMissingRequiredArg("title", jobspec.CreateTaskArgs{}.Kind())
	}

	if args.Description == "" {
		return newMissingRequiredArg("description", jobspec.CreateTaskArgs{}.Kind())
	}

	return nil
}
