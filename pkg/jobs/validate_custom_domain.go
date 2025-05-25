package jobs

import (
	"context"

	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"
)

// ValidateCustomDomainArgs for the worker to process the custom domain
type ValidateCustomDomainArgs struct {
	// ID of the custom domain in our system this is not required. When not
	// passed, will validate all custom domains for the organization
	CustomDomainID string `json:"custom_domain_id"`
}

// Kind satisfies the river.Job interface
func (ValidateCustomDomainArgs) Kind() string { return "validate_custom_domain" }

// ValidateCustomDomainWorker checks cloudflare custom domain(s), and updates
// the status in our system
type ValidateCustomDomainWorker struct {
	river.WorkerDefaults[ValidateCustomDomainArgs]

	Config ValidateCustomDomainConfig
}

// ValidateCustomDomainConfig contains the configuration for the worker
type ValidateCustomDomainConfig struct {
	CloudflareAPIKey string `koanf:"cloudflareApiKey" json:"cloudflareApiKey" jsonschema:"required description=the cloudflare api key"`
}

// Work satisfies the river.Worker interface for the validate custom domain worker
// it validates a custom domain for an organization
// todo(acookin): implement this
func (w *ValidateCustomDomainWorker) Work(_ context.Context, job *river.Job[ValidateCustomDomainArgs]) error {
	log.Info().Str("custom_domain_id", job.Args.CustomDomainID).Msg("creating custom domain")

	return nil
}
