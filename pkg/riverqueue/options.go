package riverqueue

import (
	"log/slog"
	"time"

	"github.com/riverqueue/river"
)

// Option is a function that configures a client
type Option func(*Client)

// WithConnectionURI sets the connection URI for the client
func WithConnectionURI(uri string) Option {
	return func(c *Client) {
		c.config.ConnectionURI = uri
	}
}

// WithLogger sets the logger for the client
func WithLogger(l *slog.Logger) Option {
	return func(c *Client) {
		c.config.RiverConf.Logger = l
	}
}

// WithMaxRetries sets the maximum number of retries for the client
func WithMaxRetries(maxRetries int) Option {
	return func(c *Client) {
		c.config.RiverConf.MaxAttempts = maxRetries
	}
}

// WithJobTimeout sets the job timeout for the client
func WithJobTimeout(jobTimeout time.Duration) Option {
	return func(c *Client) {
		c.config.RiverConf.JobTimeout = jobTimeout
	}
}

// WithRunMigrations sets the run migrations flag for the client
func WithRunMigrations(runMigrations bool) Option {
	return func(c *Client) {
		c.config.RunMigrations = runMigrations
	}
}

// WithWorkers sets the workers for the client
// this should be omitted when creating an insert only client
func WithWorkers(workers *river.Workers) Option {
	return func(c *Client) {
		c.config.RiverConf.Workers = workers
	}
}

// WithQueues sets the queues for the client
// this should be omitted when creating an insert only client
func WithQueues(q map[string]river.QueueConfig) Option {
	return func(c *Client) {
		c.config.RiverConf.Queues = q
	}
}

// WithRiverConfig sets the entire river configuration for the client
// prefer using the other options when possible
func WithRiverConfig(conf river.Config) Option {
	return func(c *Client) {
		c.config.RiverConf = conf
	}
}
