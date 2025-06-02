// Package main generates the JSON schema and YAML schema for the riverboat configuration
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/invopop/jsonschema"
	"github.com/invopop/yaml"
	"github.com/mcuadros/go-defaults"
	"github.com/rs/zerolog/log"

	"github.com/theopenlane/utils/envparse"

	"github.com/theopenlane/riverboat/config"
)

// appName is the name of the application
const appName = "riverboat"

// const values used for the schema generator
const (
	koanfTagName   = "koanf"
	skipper        = "-"
	defaultTag     = "default"
	jsonSchemaPath = "./configgen/%s.config.json"
	yamlConfigPath = "./config/config.example.yaml"
	envConfigPath  = "./config/.env.example"
	configMapPath  = "./config/configmap.yaml"
	ownerReadWrite = 0o600
	repoRoot       = "github.com/theopenlane/%s/"
)

// includedPackages is a list of packages to include in the schema generation
// that contain Go comments to be added to the schema
// any external packages must use the jsonschema description tags to add comments
var includedPackages = []string{
	"internal/river",
	"pkg/jobs",
	"pkg/riverqueue",
}

// schemaConfig represents the configuration for the schema generator
type schemaConfig struct {
	// jsonSchemaPath represents the file path of the JSON schema to be generated
	jsonSchemaPath string
	// yamlConfigPath is the file path to the YAML configuration to be generated
	yamlConfigPath string
	// envConfigPath is the file path to the environment variable configuration to be generated
	envConfigPath string
	// configMapPath is the file path to the kubernetes config map configuration to be generated
	configMapPath string
}

func main() {
	c := schemaConfig{
		jsonSchemaPath: fmt.Sprintf(jsonSchemaPath, appName),
		yamlConfigPath: yamlConfigPath,
		envConfigPath:  envConfigPath,
		configMapPath:  configMapPath,
	}

	generateSchema(appName, c, &config.Config{})
}

// generateSchema generates a JSON schema and a YAML schema based on the provided schemaConfig and structure
func generateSchema(appName string, c schemaConfig, structure interface{}) {
	// override the default name to using the prefixed pkg name
	r := jsonschema.Reflector{
		Namer:                      namePkg,
		ExpandedStruct:             true,
		RequiredFromJSONSchemaTags: true,
		FieldNameTag:               koanfTagName,
	}

	// add go comments to the schema
	for _, pkg := range includedPackages {
		if err := r.AddGoComments(fmt.Sprintf(repoRoot, appName), pkg); err != nil {
			log.Panic().Err(err).Msg("error adding go comments to schema")
		}
	}

	s := r.Reflect(structure)

	// generate the json schema
	genJSONSchema(s, c)

	// generate yaml schema with default
	genYAMLSchema(c)

	// generate environment variables
	configMapSchema := genEnvVarSchema(c)

	// Get the configmap header
	genConfigMapSchema(configMapSchema, c)
}

func namePkg(r reflect.Type) string {
	return r.String()
}

func genJSONSchema(s interface{}, c schemaConfig) {
	// generate the json schema
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Panic().Err(err).Msg("error marshalling json schema")
	}

	if err := os.WriteFile(c.jsonSchemaPath, data, ownerReadWrite); err != nil {
		log.Panic().Err(err).Msg("error writing json schema")
	}
}

func genYAMLSchema(c schemaConfig) {
	yamlConfig := &config.Config{}
	defaults.SetDefaults(yamlConfig)

	// this uses the `json` tag to generate the yaml schema
	yamlSchema, err := yaml.Marshal(yamlConfig)
	if err != nil {
		log.Panic().Err(err).Msg("error marshalling yaml schema")
	}

	if err = os.WriteFile(c.yamlConfigPath, yamlSchema, ownerReadWrite); err != nil {
		log.Panic().Err(err).Msg("error writing yaml schema")
	}
}

func genEnvVarSchema(c schemaConfig) string {
	cp := envparse.Config{
		FieldTagName: koanfTagName,
		Skipper:      skipper,
	}

	out, err := cp.GatherEnvInfo(strings.ToUpper(appName), &config.Config{})
	if err != nil {
		log.Panic().Err(err).Msg("error gathering environment variables")
	}

	// generate the environment variables from the config
	envSchema := ""
	configMapSchema := "\n"

	for _, k := range out {
		defaultVal := k.Tags.Get(defaultTag)

		envSchema += fmt.Sprintf("%s=\"%s\"\n", k.Key, defaultVal)

		// if the default value is empty, use the value from the values.yaml
		if defaultVal == "" {
			configMapSchema += fmt.Sprintf("  %s: {{ .Values.%s }}\n", k.Key, k.FullPath)
		} else {
			switch k.Type.Kind() {
			case reflect.String, reflect.Int64:
				defaultVal = "\"" + defaultVal + "\"" // add quotes to the string
			case reflect.Slice, reflect.Array:
				defaultVal = strings.Replace(defaultVal, "[", "", 1)
				defaultVal = strings.Replace(defaultVal, "]", "", 1)
				defaultVal = "\"" + defaultVal + "\"" // add quotes to the string
			}

			configMapSchema += fmt.Sprintf("  %s: {{ .Values.%s | default %s }}\n", k.Key, k.FullPath, defaultVal)
		}
	}

	// write the environment variables to a file
	if err = os.WriteFile(c.envConfigPath, []byte(envSchema), ownerReadWrite); err != nil {
		log.Panic().Err(err).Msg("error writing environment variables to file")
	}

	return configMapSchema
}

func genConfigMapSchema(configMapSchema string, c schemaConfig) {
	cm, err := os.ReadFile("./configgen/templates/configmap.tmpl")
	if err != nil {
		log.Panic().Err(err).Msg("error reading configmap template")
	}

	// append the configmap schema to the header
	cm = append(cm, []byte(configMapSchema)...)

	// write the configmap to a file
	if err = os.WriteFile(c.configMapPath, cm, ownerReadWrite); err != nil {
		log.Panic().Err(err).Msg("error writing configmap")
	}
}
