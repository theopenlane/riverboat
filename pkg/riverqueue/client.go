package riverqueue

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivertype"
	"github.com/rs/zerolog/log"
)

// JobClient is an interface for the river client to insert jobs
// this interface is only used for inserting new jobs and will not contain any other methods
type JobClient interface {
	// InsertMany inserts many jobs at once. Each job is inserted as an InsertManyParams tuple, which takes job args along with an optional set of insert options, which override insert options provided
	// by an JobArgsWithInsertOpts.InsertOpts implementation or any client-level defaults. The provided context is used for the underlying Postgres inserts and can be used to cancel the operation or apply a timeout.
	InsertMany(ctx context.Context, params []river.InsertManyParams) ([]*rivertype.JobInsertResult, error)
	// InsertManyTx inserts many jobs at once. Each job is inserted as an InsertManyParams tuple, which takes job args along with an optional set of insert options, which override insert options provided
	// by an JobArgsWithInsertOpts.InsertOpts implementation or any client-level defaults. The provided context is used for the underlying Postgres inserts and can be used to cancel the operation or apply a timeout.
	InsertManyTx(ctx context.Context, tx pgx.Tx, params []river.InsertManyParams) ([]*rivertype.JobInsertResult, error)
	// Insert inserts a new job with the provided args. Job opts can be used to override any defaults that may have been provided by an implementation of JobArgsWithInsertOpts.InsertOpts,
	// as well as any global defaults. The provided context is used for the underlying Postgres insert and can be used to cancel the operation or apply a timeout.
	Insert(ctx context.Context, args river.JobArgs, opts *river.InsertOpts) (*rivertype.JobInsertResult, error)
	// InsertTx inserts a new job with the provided args on the given transaction. Job opts can be used to override any defaults that may have been provided by an implementation of JobArgsWithInsertOpts.InsertOpts,
	// as well as any global defaults. The provided context is used for the underlying Postgres insert and can be used to cancel the operation or apply a timeout.
	InsertTx(ctx context.Context, tx pgx.Tx, args river.JobArgs, opts *river.InsertOpts) (*rivertype.JobInsertResult, error)
	// InsertManyFast inserts many jobs at once using Postgres' `COPY FROM` mechanism, making the operation quite fast and memory efficient. Each job is inserted as an InsertManyParams tuple,
	// which takes job args along with an optional set of insert options, which override insert options provided by an JobArgsWithInsertOpts.InsertOpts implementation or any client-level defaults.
	// The provided context is used for the underlying Postgres inserts and can be used to cancel the operation or apply a timeout.
	InsertManyFast(ctx context.Context, params []river.InsertManyParams) (int, error)
	// InsertManyTx inserts many jobs at once using Postgres' `COPY FROM` mechanism, making the operation quite fast and memory efficient. Each job is inserted as an InsertManyParams tuple,
	// which takes job args along with an optional set of insert options, which override insert options provided by an JobArgsWithInsertOpts.InsertOpts implementation or any client-level defaults.
	// The provided context is used for the underlying Postgres inserts and can be used to cancel the operation or apply a timeout.
	InsertManyFastTx(ctx context.Context, tx pgx.Tx, params []river.InsertManyParams) (int, error)
	// JobCancel cancels the job with the given ID. If possible, the job is cancelled immediately and will not be retried.
	// The provided context is used for the underlying Postgres update and can be used to cancel the operation or apply a timeout.
	JobCancel(ctx context.Context, jobID int64) (*rivertype.JobRow, error)
	// JobCancelTx cancels the job with the given ID within the specified transaction. This variant lets a caller cancel a job atomically alongside other database changes.
	// A cancelled job doesn't take effect until the transaction commits, and if the transaction rolls back, so too is the cancelled job.
	JobCancelTx(ctx context.Context, tx pgx.Tx, jobID int64) (*rivertype.JobRow, error)

	// GetPool returns the underlying pgx pool
	GetPool() *pgxpool.Pool
	// TruncateRiverTables truncates River tables in the target database. This is for test cleanup and should obviously only be used in tests.
	TruncateRiverTables(ctx context.Context) error
	// GetRiverClient returns the underlying river client
	// this can be used to interact directly with the river client for more advanced use cases (e.g. starting the river server)
	// which are outside the scope of the insert-only client interface
	GetRiverClient() *river.Client[pgx.Tx]
	// Close closes the underlying pgx pool
	Close() error
}

// Config settings for the river client
type Config struct {
	// ConnectionURI is the connection URI for the database
	ConnectionURI string `koanf:"connectionURI" json:"connectionURI" sensitive:"true" default:"postgres://postgres:password@0.0.0.0:5432/jobs?sslmode=disable"`
	// RunMigrations is a flag to determine if migrations should be run
	RunMigrations bool `koanf:"runMigrations" json:"runMigrations" default:"false"`
	// RiverConf is the river configuration
	RiverConf river.Config `koanf:"riverConf" json:"riverConf"`
	// Metrics configuration settings to enable and configure the metrics middleware
	Metrics MetricsConfig `koanf:"metrics" json:"metrics"`
}

// MetricsConfig allows the configuration of the metrics middleware
type MetricsConfig struct {
	// Enable toggles otel metrics middleware
	EnableMetrics bool `koanf:"enableMetrics" json:"enableMetrics" default:"false"`
	// DurationUnit sets the duration unit for metrics
	MetricsDurationUnit string `koanf:"metricsDurationUnit" json:"metricsDurationUnit" default:"ms"`
	// EnableSemanticMetrics toggles semantic metrics
	EnableSemanticMetrics bool `koanf:"enableSemanticMetrics" json:"enableSemanticMetrics" default:"true"`
}
// Client is a river Client that implements the JobClient interface
type Client struct {
	config Config

	pool *pgxpool.Pool

	// riverClient is the river client that is used to interact with the river server
	// using the pgx driver
	riverClient *river.Client[pgx.Tx]
}

// ensure the client implements the JobClient interface
var _ JobClient = &Client{}

// New creates a new river client with the options provided
func New(ctx context.Context, opts ...Option) (c *Client, err error) {
	// Initialize the Client struct
	c = &Client{}

	// apply the options to the client
	for _, opt := range opts {
		opt(c)
	}

	if c.config.ConnectionURI == "" {
		return nil, ErrConnectionURIRequired
	}

	// create a new river client with the given connection URI
	c.pool, err = pgxpool.New(ctx, c.config.ConnectionURI)
	if err != nil {
		log.Error().Err(err).Msg("error creating job queue database connection")
		return nil, err
	}

	// run migrations if the flag is set
	if c.config.RunMigrations {
		log.Info().Msg("running migrations")

		if err := RunMigrations(ctx, c.pool); err != nil {
			log.Error().Err(err).Msg("error running migrations")
			return nil, err
		}
	}


	// Setup otel metrics if enabled
	opts = []river.ClientOption{}
	if c.config.Metrics.EnableMetrics {
		opts = append(opts, river.WithOtelMetrics(
			c.config.Metrics.EnableMetrics,
			c.config.Metrics.MetricsDurationUnit,
			c.config.Metrics.EnableSemanticMetrics,
		))
	}
	c.riverClient, err = river.NewClient(riverpgxv5.New(c.pool), &c.config.RiverConf, opts...)
	if err != nil {
		log.Error().Err(err).Msg("error creating river client")
		return nil, err
	}

	return c, nil
}

// GetPool returns the underlying pgx pool
func (c *Client) GetPool() *pgxpool.Pool {
	return c.pool
}

// Close satisfies the JobClient interface
func (c *Client) Close() error {
	c.pool.Close()

	return nil
}

// GetRiverClient returns the underlying river client
func (c *Client) GetRiverClient() *river.Client[pgx.Tx] {
	return c.riverClient
}

// Insert satisfies the JobClient interface and inserts a new job with the provided args and opts
// it will start a new transaction and insert the job
func (c *Client) Insert(ctx context.Context, args river.JobArgs, opts *river.InsertOpts) (*rivertype.JobInsertResult, error) {
	tx, err := c.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	res, err := c.riverClient.InsertTx(ctx, tx, args, opts)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return nil, err
		}

		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return res, nil
}

// InsertMany satisfies the JobClient interface and inserts many jobs at once
// it will start a new transaction and insert the jobs
func (c *Client) InsertMany(ctx context.Context, params []river.InsertManyParams) ([]*rivertype.JobInsertResult, error) {
	tx, err := c.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	res, err := c.riverClient.InsertManyTx(ctx, tx, params)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return nil, err
		}

		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return res, nil
}

// InsertManyTx satisfies the JobClient interface
func (c *Client) InsertManyTx(ctx context.Context, tx pgx.Tx, params []river.InsertManyParams) ([]*rivertype.JobInsertResult, error) {
	return c.riverClient.InsertManyTx(ctx, tx, params)
}

// InsertTx satisfies the JobClient interface
func (c *Client) InsertTx(ctx context.Context, tx pgx.Tx, args river.JobArgs, opts *river.InsertOpts) (*rivertype.JobInsertResult, error) {
	return c.riverClient.InsertTx(ctx, tx, args, opts)
}

// InsertManyFast satisfies the JobClient interface and inserts many jobs at once using Postgres' `COPY FROM` mechanism
// it will start a new transaction and insert the jobs and commit the transaction after the insert
func (c *Client) InsertManyFast(ctx context.Context, params []river.InsertManyParams) (int, error) {
	tx, err := c.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, err
	}

	res, err := c.riverClient.InsertManyFastTx(ctx, tx, params)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return 0, err
		}

		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return res, nil
}

// InsertManyFastTx satisfies the JobClient interface
func (c *Client) InsertManyFastTx(ctx context.Context, tx pgx.Tx, params []river.InsertManyParams) (int, error) {
	return c.riverClient.InsertManyFastTx(ctx, tx, params)
}

// JobCancel satisfies the JobClient interface and cancels the job with the given ID
// it will start a new transaction and cancel the job and commit the transaction
func (c *Client) JobCancel(ctx context.Context, jobID int64) (*rivertype.JobRow, error) {
	tx, err := c.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	res, err := c.riverClient.JobCancelTx(ctx, tx, jobID)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return nil, err
		}

		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return res, nil
}

// JobCancelTx satisfies the JobClient interface
func (c *Client) JobCancelTx(ctx context.Context, tx pgx.Tx, jobID int64) (*rivertype.JobRow, error) {
	return c.riverClient.JobCancelTx(ctx, tx, jobID)
}

// TruncateRiverTables truncates River tables in the target database. This is
// for test cleanup and should obviously only be used in tests.
func (c *Client) TruncateRiverTables(ctx context.Context) error {
	pool := c.GetPool()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second) // nolint:mnd
	defer cancel()

	tables := []string{"river_job", "river_leader", "river_queue"}

	for _, table := range tables {
		if _, err := pool.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s;", table)); err != nil {
			return fmt.Errorf("error truncating %q: %w", table, err)
		}
	}

	return nil
}

// Healthcheck pings the DB to check if the connection is working
func Healthcheck(client *Client) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		if err := client.GetPool().Ping(ctx); err != nil {
			return fmt.Errorf("db connection failed: %w", err)
		}

		return nil
	}
}
