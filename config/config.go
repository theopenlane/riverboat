package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/mcuadros/go-defaults"

	"github.com/theopenlane/riverboat/internal/river"
)

var (
	DefaultConfigFilePath = "./config/.config.yaml"
	envPrefix             = "RIVERBOAT_"
)

// Config contains the configuration for the server
type Config struct {
	// RefreshInterval determines how often to reload the config
	RefreshInterval time.Duration `json:"refreshInterval" koanf:"refreshInterval" default:"10m"`
	// JobQueue is the configuration for the job queue
	JobQueue river.Config `koanf:"jobQueue" json:"jobQueue"`
}

// Load is responsible for loading the configuration from a YAML file and environment variables.
// If the `cfgFile` is empty or nil, it sets the default configuration file path.
// Config settings are taken from default values, then from the config file, and finally from environment
// the later overwriting the former.
func Load(cfgFile *string) (*Config, error) {
	k := koanf.New(".")

	if cfgFile == nil || *cfgFile == "" {
		*cfgFile = DefaultConfigFilePath
	}

	// load defaults
	conf := &Config{}
	defaults.SetDefaults(conf)

	// parse yaml config
	if err := k.Load(file.Provider(*cfgFile), yaml.Parser()); err != nil {
		fmt.Println("failed to load config file", err)
	} else {
		// unmarshal the config
		if err := k.Unmarshal("", &conf); err != nil {
			panic(err)
		}
	}

	// load env vars
	if err := k.Load(env.Provider(envPrefix, ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(
			strings.TrimPrefix(s, envPrefix)), "_", ".")
	}), nil); err != nil {
		panic(err)
	}

	// unmarshal the env vars
	if err := k.Unmarshal("", &conf); err != nil {
		panic(err)
	}

	return conf, nil
}
