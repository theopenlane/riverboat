package olclient

import (
	"github.com/theopenlane/core/pkg/openlaneclient"
)

type OpenlaneClient = openlaneclient.OpenlaneGraphClient

func New(config openlaneclient.Config, opts ...openlaneclient.ClientOption) (OpenlaneClient, error) {
	return openlaneclient.New(config, opts...)
}
