package jobs

import (
	"context"
	"net/url"
	"time"

	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"
	"github.com/theopenlane/core/pkg/openlaneclient"
)

// OnboardingArgs for the worker to process the job
type OnboardingArgs struct {
	// OrganizationID is the id of the organization the tasks should be created for
	OrganizationID string `json:"organizationID"`
	// OwnerID is the user ID that owns the new organization
	OwnerID string `json:"ownerID"`
	// AdminGroupID is the group ID of the managed Admin group
	AdminGroupID string `json:"adminGroupID"`
	// Token is a token to use for the API that has access to the organization
	Token string `json:"token"`
}

// Kind satisfies the river.Job interface
func (OnboardingArgs) Kind() string { return "onboarding" }

type OnboardingWorker struct {
	river.WorkerDefaults[OnboardingArgs]

	Config OnboardingConfig `koanf:"config" json:"config" jsonschema:"description=the email configuration"`
}

// OnboardingConfig contains the configuration for the onboarding worker
type OnboardingConfig struct {
	// StarterTasks are the tasks to create for the organization after signup
	StarterTasks []Task `json:"starterTasks"`
	// APIBaseURL is the base URL for the Openlane API
	APIBaseURL url.URL `json:"apiBaseURL"`
}

type Task struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Details     map[string]any `json:"details"`
}

// validateOnboardingConfig validates the subscription configuration settings
func (w *OnboardingWorker) validateOnboardingConfig() error {
	// default to the dev URL if none is provided in the config
	if w.Config.APIBaseURL.String() == "" {
		w.Config.APIBaseURL = url.URL{
			Scheme: "http",
			Host:   "localhost:17608",
		}
	}

	return nil
}

// Work satisfies the river.Worker interface for the email worker
func (w *OnboardingWorker) Work(ctx context.Context, job *river.Job[OnboardingArgs]) error {
	// validate the email configuration
	if err := w.validateOnboardingConfig(); err != nil {
		return err
	}

	if job.Args.Token == "" {
		return newMissingRequiredArg("token", job.Args.Kind())
	}

	client, err := NewAPIClient(w.Config.APIBaseURL, job.Args.Token)
	if err != nil {
		return err
	}

	taskInput := []*openlaneclient.CreateTaskInput{}
	// create starter tasks due 30 days from now
	dueDate := time.Now().AddDate(0, 0, 30) //nolint:mnd

	for _, t := range w.Config.StarterTasks {
		t := openlaneclient.CreateTaskInput{
			Title:       t.Title,
			Description: &t.Description,
			Details:     t.Details,
			Tags:        []string{"onboarding", "setup"},
			AssigneeID:  &job.Args.OwnerID,
			GroupIDs:    []string{job.Args.AdminGroupID}, // allow all admins to see starter tasks
			OwnerID:     job.Args.OrganizationID,
			Due:         &dueDate,
		}

		taskInput = append(taskInput, &t)
	}

	_, err = client.CreateBulkTask(ctx, taskInput)
	if err != nil {
		log.Error().Err(err).Msg("error creating onboarding tasks")

		return err
	}

	return nil
}
