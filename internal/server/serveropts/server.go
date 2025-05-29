package serveropts

import (
	"github.com/theopenlane/riverboat/config"
	serverconfig "github.com/theopenlane/riverboat/internal/server/config"
)

// ServerOptions holds the configuration and provider for the server
type ServerOptions struct {
	// ConfigProvider is the provider for the server configuration
	ConfigProvider serverconfig.Provider
	// Config holds the server configuration settings
	Config serverconfig.Config
}

// NewServerOptions creates a new ServerOptions instance with the provided options and configuration location
func NewServerOptions(opts []ServerOption, cfgLoc string) *ServerOptions {
	// load koanf config
	c, err := config.Load(&cfgLoc)
	if err != nil {
		panic(err)
	}

	so := &ServerOptions{
		Config: serverconfig.Config{
			Settings: *c,
		},
	}

	for _, opt := range opts {
		opt.apply(so)
	}

	return so
}

// AddServerOptions applies a server option after the initial setup
// this should be used when information is not available on NewServerOptions
func (so *ServerOptions) AddServerOptions(opt ServerOption) {
	opt.apply(so)
}
