package jobs

import (
	"net/url"

	"github.com/theopenlane/core/pkg/openlaneclient"

	"github.com/theopenlane/riverboat/internal/olclient"
)

// CustomDomainConfig contains the configuration for the custom domain workers
type CustomDomainConfig struct {
	CloudflareAPIKey string `koanf:"cloudflareApiKey" json:"cloudflareApiKey" jsonschema:"required description=the cloudflare api key"`

	OpenLaneAPIHost  string `koanf:"openLaneAPIHost" json:"openLaneAPIHost" jsonschema:"required description=the open lane api host"`
	OpenLaneAPIToken string `koanf:"openLaneAPIToken" json:"openLaneAPIToken" jsonschema:"required description=the open lane api token"`

	DatabaseHost string `koanf:"databaseHost" json:"databaseHost" jsonschema:"required description=the database host"`
}

func getOpenlaneClient(config CustomDomainConfig) (olclient.OpenlaneClient, error) {
	olconfig := openlaneclient.NewDefaultConfig()

	baseURL, err := url.Parse(config.OpenLaneAPIHost)
	if err != nil {
		return nil, err
	}

	opts := []openlaneclient.ClientOption{openlaneclient.WithBaseURL(baseURL)}
	opts = append(opts, openlaneclient.WithCredentials(openlaneclient.Authorization{
		BearerToken: config.OpenLaneAPIToken,
	}))

	return openlaneclient.New(olconfig, opts...)
}
