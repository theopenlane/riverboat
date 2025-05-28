package jobs

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"
	"github.com/theopenlane/core/pkg/models"
	"github.com/theopenlane/utils/ulids"
	"golang.org/x/sync/errgroup"
)

type ScheduledJobArgs struct{}

func (ScheduledJobArgs) Kind() string { return "scheduled_jobs" }

// ScheduledJobConfig contains the configuration for the scheduling job worker
type ScheduledJobConfig struct {
	// DatabaseHost for connecting to the postgres database
	// This is for the core server database which can potentially be different from
	// river queue's
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

type Run struct {
	ID                    string    `json:"id"`
	JobRunnerID           string    `json:"job_runner_id"`
	Status                string    `json:"status"`
	ScheduledJobID        string    `json:"scheduled_job_id"`
	CreatedAt             time.Time `json:"created_at"`
	ExpectedExecutionTime time.Time `json:"expected_execution_time"`
	Script                string    `json:"script"`
}

func (s *ScheduledJobWorker) Work(ctx context.Context, job *river.Job[ScheduledJobArgs]) error {
	if err := s.validateConnection(); err != nil {
		return err
	}

	const batchSize = 5
	var (
		offset  = 0
		allJobs []controlScheduledJob
		g       = new(errgroup.Group)
	)

	g.SetLimit(10)

	for {
		query := `
			SELECT 
				csj.id, csj.job_id, csj.configuration, csj.cadence, csj.cron,
				sj.id, sj.display_id, sj.title, sj.description, sj.job_type, 
				sj.script, sj.configuration, sj.created_at, csj.job_runner_id
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

func scanBatch(rows pgx.Rows) ([]controlScheduledJob, error) {
	var jobs []controlScheduledJob
	for rows.Next() {
		var job controlScheduledJob
		var scheduledJob scheduledJob
		var cadence sql.NullString
		var cronStr sql.NullString

		var displayID, title, description, jobType sql.NullString
		var runnerID, script string

		err := rows.Scan(
			&job.ID, &job.JobID, &job.Configuration, &cadence, &cronStr,
			&scheduledJob.ID, &displayID, &title, &description, &jobType,
			&script, &scheduledJob.Configuration, &scheduledJob.CreatedAt,
			&runnerID,
		)
		if err != nil {
			return nil, err
		}

		// we have valid cadence JSON, unmarshal it directly to the struct
		if cadence.Valid && cadence.String != "" {
			if err := json.Unmarshal([]byte(cadence.String), &job.Cadence); err != nil {
				log.Error().Err(err).
					Str("control_job_id", job.ID).
					Str("cadence", cadence.String).
					Msg("failed to unmarshal cadence")
			}
		}

		// we have a valid cron string
		if cronStr.Valid && cronStr.String != "" {
			job.Cron = models.Cron(cronStr.String)
		}

		if jobType.Valid {
			scheduledJob.JobType = jobType.String
		}

		scheduledJob.Script = script
		job.JobRunnerID = runnerID

		job.Job = &scheduledJob
		jobs = append(jobs, job)
	}
	return jobs, rows.Err()
}

func (s *ScheduledJobWorker) processJob(ctx context.Context, job controlScheduledJob) error {
	now := time.Now()
	var nextRun time.Time
	var err error

	if !job.Cadence.IsZero() {
		nextRun, err = job.Cadence.Next(now)
	} else if job.Cron.String() != "" {
		nextRun, err = job.Cron.Next(now)
	} else {
		log.Info().Msg("no cadence skipping")
		return nil
	}

	if err != nil {
		log.Error().Err(err).Msg("could not get the next run date")
		return err
	}

	// if <= 10 mins, we want to schedule the job so the agents can
	// fetch those and run internally on their own
	// This would allow the agents have their own internal cache
	// and won't have to ping home every minute but will do every 10 minutes.
	//
	// The agents will have the list of jobs that will run over the next 10 minutes
	// and execute them if any at the right time
	//
	const scheduleBuffer = 10 * time.Minute
	if nextRun.IsZero() || nextRun.Sub(now) > scheduleBuffer {
		return nil
	}

	var existingCount int
	checkQuery := `
		SELECT COUNT(*) 
		FROM scheduled_job_runs 
		WHERE scheduled_job_id = $1 
		AND expected_execution_time = $2 
		AND status = 'PENDING'
	`
	err = s.dbPool.QueryRow(ctx, checkQuery, job.ID, nextRun).Scan(&existingCount)
	if err != nil {
		return err
	}

	if existingCount > 0 {
		log.Info().
			Str("job_id", job.ID).
			Time("expected_execution_time", nextRun).
			Msg("skipping job creation, run already exists")
		return nil
	}

	run := &Run{
		ID:                    ulids.New().String(),
		Status:                "PENDING",
		ScheduledJobID:        job.ID,
		CreatedAt:             now,
		JobRunnerID:           job.JobRunnerID,
		Script:                job.Job.Script,
		ExpectedExecutionTime: nextRun,
	}

	query := `
		INSERT INTO scheduled_job_runs (
			id, status, scheduled_job_id, created_at, 
			job_runner_id, expected_execution_time, script
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = s.dbPool.Exec(ctx, query,
		run.ID,
		run.Status,
		run.ScheduledJobID,
		run.CreatedAt,
		run.JobRunnerID,
		run.ExpectedExecutionTime,
		run.Script,
	)
	return err
}

type scheduledJob struct {
	ID            string                  `json:"id"`
	DisplayID     string                  `json:"display_id"`
	Title         string                  `json:"title"`
	Description   string                  `json:"description"`
	JobType       string                  `json:"job_type"`
	Script        string                  `json:"script"`
	Configuration models.JobConfiguration `json:"configuration"`
	CreatedAt     time.Time               `json:"created_at"`
}

type controlScheduledJob struct {
	ID            string                  `json:"id,omitempty"`
	JobID         string                  `json:"job_id,omitempty"`
	JobRunnerID   string                  `json:"job_runner_id,omitempty"`
	Configuration models.JobConfiguration `json:"configuration,omitempty"`
	Cadence       models.JobCadence       `json:"cadence,omitempty"`
	Cron          models.Cron             `json:"cron,omitempty"`
	Job           *scheduledJob           `json:"job,omitempty"`
}
