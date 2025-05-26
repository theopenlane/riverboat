package jobs

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/theopenlane/utils/ulids"
	"golang.org/x/sync/errgroup"
)

type ScheduledJobArgs struct {
}

func (ScheduledJobArgs) Kind() string { return "scheduled_jobs" }

// ScheduledJobConfig contains the configuration for the scheduling job worker
type ScheduledJobConfig struct {
	// DatabaseHost for connecting to the postgres database
	DatabaseHost string `koanf:"databaseHost" json:"databaseHost" default:"postgres://postgres:password@0.0.0.0:5432/jobs?sslmode=disable"`
}

type ScheduledJobWorker struct {
	river.WorkerDefaults[ScheduledJobArgs]

	Config ScheduledJobConfig `koanf:"config" json:"config" jsonschema:"description=the scheduled job worker configuration"`

	dbPool *pgxpool.Pool
}

func (s *ScheduledJobWorker) validateConnection() error {

	if s.dbPool == nil {

		ctx, cancelFn := context.WithTimeout(context.Background(), time.Second)
		defer cancelFn()

		dbpool, err := pgxpool.New(ctx, s.Config.DatabaseHost)
		if err != nil {
			return err
		}

		s.dbPool = dbpool
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), time.Second)
	defer cancelFn()

	return s.dbPool.Ping(ctx)
}

type ScheduledJob struct {
	ID            string          `json:"id"`
	DisplayID     string          `json:"display_id"`
	Title         string          `json:"title"`
	Description   string          `json:"description"`
	JobType       string          `json:"job_type"`
	Script        string          `json:"script"`
	Configuration json.RawMessage `json:"configuration"`
	CreatedAt     time.Time       `json:"created_at"`
}

type ControlScheduledJob struct {
	ID            string          `json:"id"`
	JobID         string          `json:"job_id"`
	Configuration json.RawMessage `json:"configuration"`
	Cadence       json.RawMessage `json:"cadence"`
	Cron          string          `json:"cron"`
	Job           *ScheduledJob   `json:"job"`
}

type Run struct {
	ID             string    `json:"id"`
	JobRunnerID    string    `json:"job_runner_id"`
	Status         string    `json:"status"`
	ScheduledJobID string    `json:"scheduled_job_id"`
	CreatedAt      time.Time `json:"created_at"`
}

func (s *ScheduledJobWorker) Work(ctx context.Context, job *river.Job[ScheduledJobArgs]) error {
	if err := s.validateConnection(); err != nil {
		return err
	}

	const batchSize = 5
	var (
		offset  = 0
		allJobs []ControlScheduledJob
		g       = new(errgroup.Group)
	)

	g.SetLimit(10)

	for {
		query := `
			SELECT 
				csj.id, csj.job_id, csj.configuration, csj.cadence, csj.cron,
				sj.id, sj.display_id, sj.title, sj.description, sj.job_type, 
				sj.script, sj.configuration, sj.created_at
			FROM control_scheduled_jobs csj
			JOIN scheduled_jobs sj ON sj.id = csj.job_id
			WHERE csj.deleted_at IS NULL AND sj.deleted_at IS NULL
			ORDER BY csj.id
			LIMIT $1 OFFSET $2
		`

		rows, err := s.dbPool.Query(ctx, query, batchSize, offset)
		if err != nil {
			return err
		}

		jobs, err := scanBatch(rows)
		rows.Close()
		if err != nil {
			return err
		}

		for _, scheduledJob := range jobs {
			job := scheduledJob
			g.Go(func() error {
				if err := s.processJob(ctx, job); err != nil {
					return err
				}
				return nil
			})
		}

		allJobs = append(allJobs, jobs...)

		if len(jobs) < batchSize {
			break
		}
		offset += batchSize
	}

	return g.Wait()
}

func scanBatch(rows pgx.Rows) ([]ControlScheduledJob, error) {
	var jobs []ControlScheduledJob
	for rows.Next() {
		var job ControlScheduledJob
		var scheduledJob ScheduledJob

		err := rows.Scan(
			&job.ID, &job.JobID, &job.Configuration, &job.Cadence, &job.Cron,
			&scheduledJob.ID, &scheduledJob.DisplayID, &scheduledJob.Title,
			&scheduledJob.Description, &scheduledJob.JobType, &scheduledJob.Script,
			&scheduledJob.Configuration, &scheduledJob.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		job.Job = &scheduledJob
		jobs = append(jobs, job)
	}
	return jobs, rows.Err()
}

func (s *ScheduledJobWorker) processJob(ctx context.Context, job ControlScheduledJob) error {
	run := &Run{
		ID:             ulids.New().String(),
		Status:         "PENDING",
		ScheduledJobID: job.ID,
		CreatedAt:      time.Now(),
	}

	query := `
		INSERT INTO scheduled_job_runs (id, status, scheduled_job_id, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := s.dbPool.Exec(ctx, query,
		run.ID,
		run.Status,
		run.ScheduledJobID,
		run.CreatedAt,
	)
	return err
}
