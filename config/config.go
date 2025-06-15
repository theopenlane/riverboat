package config

import (
	"strings"
	"time"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
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
	RefreshInterval time.Duration `json:"refreshInterval" koanf:"refreshInterval" default:"10m"`
	// River is the configuration for the job queue
	River river.Config `koanf:"river" json:"river"`
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

	// load defaults
	conf := &Config{}
	defaults.SetDefaults(conf)

	// parse yaml config
	if err := k.Load(file.Provider(*cfgFile), yaml.Parser()); err != nil {
		log.Warn().Err(err).Msg("failed to load config file - ensure the .config.yaml is present and valid or use environment variables to set the configuration")
	}

	// unmarshal the config
	if err := k.Unmarshal("", &conf); err != nil {
		log.Fatal().Err(err).Msg("failed to unmarshal config file")
	}

	// load env vars
	if err := k.Load(env.ProviderWithValue(envPrefix, ".", func(s, v string) (string, interface{}) {
		key := strings.ReplaceAll(strings.ToLower(strings.TrimPrefix(s, envPrefix)), "_", ".")

		if strings.Contains(v, ",") {
			return key, strings.Split(v, ",")
		}

		return key, v
	}), nil); err != nil {
		log.Warn().Err(err).Msg("failed to load env vars, some settings may not be applied")
	}

	// unmarshal the env vars
	if err := k.Unmarshal("", &conf); err != nil {
		log.Fatal().Err(err).Msg("failed to unmarshal env vars")
	}

	return conf, nil
}
