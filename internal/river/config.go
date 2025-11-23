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
	DatabaseHost string `koanf:"databasehost" json:"databasehost" sensitive:"true" default:"postgres://postgres:password@0.0.0.0:5432/jobs?sslmode=disable"`
	// Queues to be enabled on the server, if not provided, a default queue is created
	Queues []Queue `koanf:"queues" json:"queues" default:""`
	// Workers to be enabled on the server
	Workers Workers `koanf:"workers" json:"workers"`
	// DefaultMaxRetries is the maximum number of retries for failed jobs, this can be set differently per job
	DefaultMaxRetries int `koanf:"defaultmaxretries" json:"defaultmaxretries" default:"10"`
	// Metrics enables or disables metrics collection
	Metrics riverqueue.MetricsConfig `koanf:"metrics" json:"metrics"`
}

// Queue is the configuration for a queue
type Queue struct {
	// Name of the queue
	Name string `koanf:"name" json:"name" default:"default"`
	// MaxWorkers allotted for the queue
	MaxWorkers int `koanf:"maxworkers" json:"maxworkers" default:"100"`
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
	OpenlaneConfig corejobs.OpenlaneConfig `koanf:"openlaneconfig" json:"openlaneconfig"`

	// EmailWorker configuration for sending emails
	EmailWorker jobs.EmailWorker `koanf:"emailworker" json:"emailworker"`

	// DatabaseWorker configuration for creating databases using openlane/dbx
	DatabaseWorker jobs.DatabaseWorker `koanf:"databaseworker" json:"databaseworker"`

	// CreateCustomDomainWorker configuration for creating custom domains
	CreateCustomDomainWorker corejobs.CreateCustomDomainWorker `koanf:"createcustomdomainworker" json:"createcustomdomainworker"`

	// ValidateCustomDomainWorker configuration for validating custom domains
	ValidateCustomDomainWorker corejobs.ValidateCustomDomainWorker `koanf:"validatecustomdomainworker" json:"validatecustomdomainworker"`

	// DeleteCustomDomainWorker configuration for deleting custom domains
	DeleteCustomDomainWorker corejobs.DeleteCustomDomainWorker `koanf:"deletecustomdomainworker" json:"deletecustomdomainworker"`

	// ExportContentWorker configuration for exporting content
	ExportContentWorker corejobs.ExportContentWorker `koanf:"exportcontentworker" json:"exportcontentworker"`

	// DeleteExportContentWorker configuration for batch deleting exports and clogging object storage
	DeleteExportContentWorker corejobs.DeleteExportContentWorker `koanf:"deleteexportcontentworker" json:"deleteexportcontentworker"`

	// WatermarkDocWorker configuration for watermarking documents
	WatermarkDocWorker corejobs.WatermarkDocWorker `koanf:"watermarkdocworker" json:"watermarkdocworker"`

	// CreatePirschDomainWorker configuration for creating Pirsch domains
	CreatePirschDomainWorker corejobs.CreatePirschDomainWorker `koanf:"createpirschdomainworker" json:"createpirschdomainworker"`

	// DeletePirschDomainWorker configuration for deleting Pirsch domains
	DeletePirschDomainWorker corejobs.DeletePirschDomainWorker `koanf:"deletepirschdomainworker" json:"deletepirschdomainworker"`

	// UpdatePirschDomainWorker configuration for updating Pirsch domains
	UpdatePirschDomainWorker corejobs.UpdatePirschDomainWorker `koanf:"updatepirschdomainworker" json:"updatepirschdomainworker"`

	// SlackWorker configuration for sending Slack messages
	SlackWorker jobs.SlackWorker `koanf:"slackworker" json:"slackworker"`

	// PreviewDomainWorkers
	CreatePreviewDomainWorker   corejobs.CreatePreviewDomainWorker   `koanf:"createpreviewdomainworker" json:"createpreviewdomainworker"`
	DeletePreviewDomainWorker   corejobs.DeletePreviewDomainWorker   `koanf:"deletepreviewdomainworker" json:"deletepreviewdomainworker"`
	ValidatePreviewDomainWorker corejobs.ValidatePreviewDomainWorker `koanf:"validatepreviewdomainworker" json:"validatepreviewdomainworker"`

	// add more workers here
}
