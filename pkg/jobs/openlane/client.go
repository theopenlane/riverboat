package openlane

import (
	openlane "github.com/theopenlane/go-client"
	"github.com/theopenlane/go-client/graphclient"
)

// GraphClient is a type alias for the graphclient.GraphClient interface
type GraphClient = graphclient.GraphClient

// Client is a type alias for openlane.Client
// which provides access to the Openlane API
type Client = *openlane.Client

// New creates a new Openlane client with the provided configuration and options.
// It returns an implementation of the OpenlaneClient interface and any error encountered.
func New(opts ...openlane.ClientOption) (Client, error) {
	return openlane.New(opts...)
}
