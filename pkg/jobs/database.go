package jobs

import (
	"context"

	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"
	"github.com/theopenlane/core/common/jobspec"
	dbx "github.com/theopenlane/dbx/pkg/dbxclient"
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

// InsertOpts provides the default configuration when processing this job.
func (DatabaseArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{MaxAttempts: 3, Queue: jobspec.QueueDefault} //nolint:mnd
}

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
