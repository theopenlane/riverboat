package jobs

import (
	"context"

	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"
)

// CreateCustomDomainArgs for the worker to process the custom domain
type CreateCustomDomainArgs struct {
	// ID of the custom domain in our system
	CustomDomainID string `json:"custom_domain_id"`
}

// Kind satisfies the river.Job interface
func (CreateCustomDomainArgs) Kind() string { return "create_custom_domain" }

// CreateCustomDomainWorker creates a custom hostname in cloudflare, and
// creates and updates the records in our system
type CreateCustomDomainWorker struct {
	river.WorkerDefaults[CreateCustomDomainArgs]

	Config CreateCustomDomainConfig
}

// CreateCustomDomainConfig contains the configuration for the worker
type CreateCustomDomainConfig struct {
	CloudflareAPIKey string `koanf:"cloudflareApiKey" json:"cloudflareApiKey" jsonschema:"required description=the cloudflare api key"`
}

// Work satisfies the river.Worker interface for the create custom domain worker
// it creates a custom domain for an organization
// todo(acookin): implement this
func (w *CreateCustomDomainWorker) Work(ctx context.Context, job *river.Job[CreateCustomDomainArgs]) error {
	log.Info().Str("custom_domain_id", job.Args.CustomDomainID).Msg("creating custom domain")

	return nil
}
