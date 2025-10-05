package jobs

import (
	"context"

	dbx "github.com/theopenlane/dbx/pkg/dbxclient"

	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"
)

// DatabaseArgs are the arguments for the database worker
type DatabaseArgs struct {
	// OrganizationID is the organization id to create the database for
	OrganizationID string `json:"organization_id"`
	// Location is the location to create the database in, e.g. AMER
	Location string `json:"location"`
}

// Kind satisfies the river.Args interface for the database worker
func (DatabaseArgs) Kind() string { return "dedicated_database" }

// DatabaseWorker is a worker to create a dedicated database for an organization
type DatabaseWorker struct {
	river.WorkerDefaults[DatabaseArgs]

	Config dbx.Config `koanf:"config" json:"config" jsonschema:"description=the database configuration"`
}

// validateDatabaseInput validates the input for the database worker
func validateDatabaseInput(job *river.Job[DatabaseArgs]) error {
	if job.Args.OrganizationID == "" {
		return newMissingRequiredArg("organization_id", DatabaseArgs{}.Kind())
	}

	if job.Args.Location == "" {
		return newMissingRequiredArg("location", DatabaseArgs{}.Kind())
	}

	return nil
}

// Work satisfies the river.Worker interface for the database worker
// it creates a dedicated database using the dbx client for the organization
func (w *DatabaseWorker) Work(ctx context.Context, job *river.Job[DatabaseArgs]) error {
	// if its not enabled, return early
	if !w.Config.Enabled {
		return nil
	}

	if err := validateDatabaseInput(job); err != nil {
		return err
	}

	input := dbx.CreateDatabaseInput{
		OrganizationID: job.Args.OrganizationID,
		Geo:            &job.Args.Location,
	}

	log.Debug().
		Str("org", input.OrganizationID).
		Str("geo", *input.Geo).
		Msg("creating database")

	client := w.Config.NewDefaultClient()
	if _, err := client.CreateDatabase(ctx, input); err != nil {
		log.Error().Err(err).Msg("failed to create database")
		return err
	}

	return nil
}

// RegisterSlackJob registers the Slack job.
func RegisterSlackJob() {
	RegisterJob("slack", func(ctx context.Context, params map[string]interface{}) error {
		channel, _ := params["channel"].(string)
		message, _ := params["message"].(string)
		devMode, _ := params["dev_mode"].(bool)
		return SendSlackMessage(ctx, SlackJobArgs{
			Channel: channel,
			Message: message,
			DevMode: devMode,
		})
	})
}
