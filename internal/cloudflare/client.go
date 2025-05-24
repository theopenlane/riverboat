package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/custom_hostnames"
	"github.com/cloudflare/cloudflare-go/v4/option"
)

type CustomHostnamesService interface {
	New(context.Context, custom_hostnames.CustomHostnameNewParams, ...option.RequestOption) (*custom_hostnames.CustomHostnameNewResponse, error)
	Delete(context.Context, string, custom_hostnames.CustomHostnameDeleteParams, ...option.RequestOption) (*custom_hostnames.CustomHostnameDeleteResponse, error)
}

type Client interface {
	CustomHostnames() CustomHostnamesService
}

type cloudflareClientImpl struct {
	client *cloudflare.Client
}

func NewClient(apiKey string) Client {
	return &cloudflareClientImpl{
		client: cloudflare.NewClient(
			option.WithAPIToken(apiKey),
		),
	}
}

func (c *cloudflareClientImpl) CustomHostnames() CustomHostnamesService {
	return c.client.CustomHostnames
}
