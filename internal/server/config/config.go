package config

import (
	"github.com/theopenlane/riverboat/config"
)

// Config is the configuration for the http server
type Config struct {
	// add all the configuration settings for the server
	Settings config.Config
}

// Ensure that *Config implements ConfigProvider interface.
var _ ConfigProvider = &Config{}

// GetConfig implements ConfigProvider.
func (c *Config) GetConfig() (*Config, error) {
	return c, nil
}
