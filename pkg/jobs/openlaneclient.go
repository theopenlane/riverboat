package jobs

import (
	"net/url"

	"github.com/Yamashou/gqlgenc/clientv2"
	"github.com/theopenlane/core/pkg/openlaneclient"
)

// NewAPIClient creates a new Openlane API client with the provided token
func NewAPIClient(baseURL url.URL, token string) (*openlaneclient.OpenlaneClient, error) {
	config := openlaneclient.Config{
		BaseURL:         &baseURL,
		GraphQLPath:     "/query",
		Clientv2Options: clientv2.Options{ParseDataAlongWithErrors: false},
	}

	return openlaneclient.New(config, openlaneclient.WithCredentials(openlaneclient.Authorization{
		BearerToken: token,
	}))
}
