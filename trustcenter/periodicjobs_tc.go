//go:build trustcenter

package trustcenter

import (
	"time"

	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"
	"github.com/theopenlane/core/common/jobspec"
)

const (
	// minInterval prevents spamming from a bag config
	minInterval = 1 * time.Minute

	// atleast one day
	minOrgDeleteReminderInterval = 1

	twentyFourhours                     = time.Hour * 24
	defaultOrgDeletionWorkerRunInterval = twentyFourhours
)

// CreatePeriodicJobs creates periodic jobs for the trust center module
func CreatePeriodicJobs(c Workers) ([]*river.PeriodicJob, error) {
	log.Info().Msg("adding additional trust center periodic jobs")

	jobs := []*river.PeriodicJob{}

	if c.ValidateCustomDomainWorker.Config.Enabled {
		interval := c.ValidateCustomDomainWorker.Config.ValidateInterval
		if interval < minInterval {
			interval = minInterval
		}

		validateCustomDomainJobs := river.NewPeriodicJob(
			river.PeriodicInterval(interval),
			func() (river.JobArgs, *river.InsertOpts) {
				return jobspec.ValidateCustomDomainArgs{}, nil
			},
			&river.PeriodicJobOpts{RunOnStart: true},
		)
		jobs = append(jobs, validateCustomDomainJobs)

		log.Info().Msg("periodic worker enabled: validate custom domain")
	}

	if c.OrganizationDeletionReminderWorker.Config.Enabled {
		days := c.OrganizationDeletionReminderWorker.Config.PaymentMethodInterval
		if days < minOrgDeleteReminderInterval {
			days = minOrgDeleteReminderInterval
		}

		orgDeletionReminders := river.NewPeriodicJob(
			river.PeriodicInterval(time.Duration(days)*twentyFourhours),
			func() (river.JobArgs, *river.InsertOpts) {
				return jobspec.OrganizationDeletionReminderArgs{}, nil
			},
			&river.PeriodicJobOpts{},
		)
		jobs = append(jobs, orgDeletionReminders)

		log.Info().Msg("periodic worker enabled: organization deletion reminder")
	}

	if c.OrganizationDeletionWorker.Config.Enabled {
		interval := c.OrganizationDeletionWorker.Config.RunInterval
		if interval == 0 {
			interval = defaultOrgDeletionWorkerRunInterval
		}

		orgDeletion := river.NewPeriodicJob(
			river.PeriodicInterval(interval),
			func() (river.JobArgs, *river.InsertOpts) {
				return jobspec.OrganizationDeletionArgs{}, nil
			},
			&river.PeriodicJobOpts{},
		)
		jobs = append(jobs, orgDeletion)

		log.Info().Msg("periodic worker enabled: organization deletion")
	}

	return jobs, nil
}
