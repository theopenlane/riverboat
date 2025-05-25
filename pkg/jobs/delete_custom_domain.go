package jobs

import (
	"context"

	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"
)

// DeleteCustomDomainArgs for the worker to process the custom domain
type DeleteCustomDomainArgs struct {
	// ID of the custom domain in our system
	CustomDomainID string `json:"custom_domain_id"`
}

// Kind satisfies the river.Job interface
func (DeleteCustomDomainArgs) Kind() string { return "delete_custom_domain" }

// DeleteCustomDomainWorker delete the custom hostname from cloudflare and
// updates the records in our system
type DeleteCustomDomainWorker struct {
	river.WorkerDefaults[DeleteCustomDomainArgs]

	Config DeleteCustomDomainConfig
}

// DeleteCustomDomainConfig contains the configuration for the example worker
type DeleteCustomDomainConfig struct {
	CloudflareAPIKey string `koanf:"cloudflareApiKey" json:"cloudflareApiKey" jsonschema:"required description=the cloudflare api key"`
}

// Work satisfies the river.Worker interface for the delete custom domain worker
// it deletes a custom domain for an organization
// todo(acookin): implement this
func (w *DeleteCustomDomainWorker) Work(_ context.Context, job *river.Job[DeleteCustomDomainArgs]) error {
	log.Info().Str("custom_domain_id", job.Args.CustomDomainID).Msg("creating custom domain")

	return nil
}
