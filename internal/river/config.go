package river

import (
	"github.com/theopenlane/core/pkg/corejobs"

	"github.com/theopenlane/riverboat/pkg/jobs"
)

// Config is the configuration for the river server
type Config struct {
	// Logger configuration, which is inherited from the core logger
	Logger Logger `koanf:"-" json:"-"`

	// DatabaseHost for connecting to the postgres database
	DatabaseHost string `koanf:"databaseHost" json:"databaseHost" default:"postgres://postgres:password@0.0.0.0:5432/jobs?sslmode=disable"`
	// Queues to be enabled on the server, if not provided, a default queue is created
	Queues []Queue `koanf:"queues" json:"queues" default:""`
	// Workers to be enabled on the server
	Workers Workers `koanf:"workers" json:"workers"`
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

	// add more workers here
}
