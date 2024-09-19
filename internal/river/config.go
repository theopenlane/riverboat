package river

import "github.com/theopenlane/riverboat/pkg/jobs"

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
	// MaxWorkers alloted for the queue
	MaxWorkers int `koanf:"maxWorkers" json:"maxWorkers" default:"100"`
}

// Logger is the configuration for the logger used in the river server
type Logger struct {
	// Debug enables debug logging
	Debug bool `koanf:"-" json:"-"`
	// Pretty enables pretty logging
	Pretty bool `koanf:"-" json:"-"`
}

type Workers struct {
	EmailWorker jobs.EmailWorker `koanf:"emailWorker" json:"emailWorker"`

	// add more workers here
}
