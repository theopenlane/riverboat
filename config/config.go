package config

import (
	"os"
	"strings"
	"time"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/mcuadros/go-defaults"
	"github.com/rs/zerolog/log"

	"github.com/theopenlane/riverboat/internal/river"
)

var (
	defaultConfigFilePath = "./config/.config.yaml"
	envPrefix             = "RIVERBOAT_"
)

// Config contains the configuration for the server
type Config struct {
	// RefreshInterval determines how often to reload the config
	RefreshInterval time.Duration `json:"refreshinterval" koanf:"refreshinterval" default:"10m"`
	// River is the configuration for the job queue
	River river.Config `koanf:"river" json:"river"`
}

// Option configures the Config
type Option func(*Config)

// New creates a Config with the supplied options applied
func New(opts ...Option) *Config {
	cfg := &Config{}
	defaults.SetDefaults(cfg)

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

// Load is responsible for loading the configuration from a YAML file and environment variables.
// If the `cfgFile` is empty or nil, it sets the default configuration file path.
// Config settings are taken from default values, then from the config file, and finally from environment
// the later overwriting the former.
func Load(cfgFile *string) (*Config, error) {
	k := koanf.New(".")

	if cfgFile == nil || *cfgFile == "" {
		*cfgFile = defaultConfigFilePath
	}

	if _, err := os.Stat(*cfgFile); err != nil {
		if os.IsNotExist(err) {
			log.Warn().Err(err).Msg("config file not found, proceeding with default configuration")
		}
	}

	// parse yaml config
	if err := k.Load(file.Provider(*cfgFile), yaml.Parser()); err != nil {
		// if it's an  unmarshal errors, panic now instead of continuing
		if strings.Contains(err.Error(), "yaml: unmarshal errors") {
			log.Fatal().Err(err).Msg("failed to unmarshal config file - ensure the .config.yaml is valid")
		} else {
			log.Warn().Err(err).Msg("failed to load config file - ensure the .config.yaml is present and valid or use environment variables to set the configuration")
		}
	}

	// load env vars
	if err := k.Load(env.Provider(".", env.Opt{
		Prefix: envPrefix,
		TransformFunc: func(key, v string) (string, interface{}) {
			key = strings.ToLower(strings.TrimPrefix(key, envPrefix))
			key = strings.ReplaceAll(key, "_", ".")

			if strings.Contains(v, ",") {
				return key, strings.Split(v, ",")
			}

			return key, v
		},
	}), nil); err != nil {
		log.Warn().Err(err).Msg("failed to load env vars, some settings may not be applied")
	}

	// create the config with defaults
	conf := New()
	if err := k.Unmarshal("", &conf); err != nil {
		log.Fatal().Err(err).Msg("failed to unmarshal config")
	}

	return conf, nil
}
