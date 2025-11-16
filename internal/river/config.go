package river

import (
	"github.com/theopenlane/core/pkg/corejobs"

	"github.com/theopenlane/riverboat/pkg/jobs"
	"github.com/theopenlane/riverboat/pkg/riverqueue"
)

// Config is the configuration for the river server
type Config struct {
	// Logger configuration, which is inherited from the core logger
	Logger Logger `koanf:"-" json:"-"`

	// DatabaseHost for connecting to the postgres database
	DatabaseHost string `koanf:"databaseHost" json:"databaseHost" sensitive:"true" default:"postgres://postgres:password@0.0.0.0:5432/jobs?sslmode=disable"`
	// Queues to be enabled on the server, if not provided, a default queue is created
	Queues []Queue `koanf:"queues" json:"queues" default:""`
	// Workers to be enabled on the server
	Workers Workers `koanf:"workers" json:"workers"`
	// DefaultMaxRetries is the maximum number of retries for failed jobs, this can be set differently per job
	DefaultMaxRetries int `koanf:"defaultMaxRetries" json:"defaultMaxRetries" default:"10"`
	// Metrics enables or disables metrics collection
	Metrics riverqueue.MetricsConfig `koanf:"metrics" json:"metrics"`
}

// Queue is the configuration for a queue
type Queue struct {
	// Name of the queue
	Name string `koanf:"name" json:"name" default:"default"`
	// MaxWorkers allotted for the queue
	MaxWorkers int `koanf:"maxWorkers" json:"maxWorkers" default:"100"`
}

// Logger is the configuration for the logger used in the river server
type Logger struct {
	// Debug enables debug logging
	Debug bool `koanf:"-" json:"-"`
	// Pretty enables pretty logging
	Pretty bool `koanf:"-" json:"-"`
}

// Workers that will be enabled on the server
type Workers struct {
	// OpenlaneConfig configuration for openlane jobs, this is shared across multiple workers
	// if a worker needs specific configuration, it can be set in the worker's config
	OpenlaneConfig corejobs.OpenlaneConfig `koanf:"openlaneConfig" json:"openlaneConfig"`

	// EmailWorker configuration for sending emails
	EmailWorker jobs.EmailWorker `koanf:"emailWorker" json:"emailWorker"`

	// DatabaseWorker configuration for creating databases using openlane/dbx
	DatabaseWorker jobs.DatabaseWorker `koanf:"databaseWorker" json:"databaseWorker"`

	// CreateCustomDomainWorker configuration for creating custom domains
	CreateCustomDomainWorker corejobs.CreateCustomDomainWorker `koanf:"createCustomDomainWorker" json:"createCustomDomainWorker"`

	// ValidateCustomDomainWorker configuration for validating custom domains
	ValidateCustomDomainWorker corejobs.ValidateCustomDomainWorker `koanf:"validateCustomDomainWorker" json:"validateCustomDomainWorker"`

	// DeleteCustomDomainWorker configuration for deleting custom domains
	DeleteCustomDomainWorker corejobs.DeleteCustomDomainWorker `koanf:"deleteCustomDomainWorker" json:"deleteCustomDomainWorker"`

	// ExportContentWorker configuration for exporting content
	ExportContentWorker corejobs.ExportContentWorker `koanf:"exportContentWorker" json:"exportContentWorker"`

	// DeleteExportContentWorker configuration for batch deleting exports and clogging object storage
	DeleteExportContentWorker corejobs.DeleteExportContentWorker `koanf:"deleteExportContentWorker" json:"deleteExportContentWorker"`

	// WatermarkDocWorker configuration for watermarking documents
	WatermarkDocWorker corejobs.WatermarkDocWorker `koanf:"watermarkDocWorker" json:"watermarkDocWorker"`

	// CreatePirschDomainWorker configuration for creating Pirsch domains
	CreatePirschDomainWorker corejobs.CreatePirschDomainWorker `koanf:"createPirschDomainWorker" json:"createPirschDomainWorker"`

	// DeletePirschDomainWorker configuration for deleting Pirsch domains
	DeletePirschDomainWorker corejobs.DeletePirschDomainWorker `koanf:"deletePirschDomainWorker" json:"deletePirschDomainWorker"`

	// SlackWorker configuration for sending Slack messages
	SlackWorker jobs.SlackWorker `koanf:"slackWorker" json:"slackWorker"`

	// add more workers here
}
