package river

import (
	"time"

	"github.com/riverqueue/river"
	"github.com/theopenlane/core/common/jobspec"

	"github.com/theopenlane/riverboat/pkg/jobs"
)

// minInterval prevents spamming from a bag config
var minInterval = 1 * time.Minute

const (
	minOrgDeleteReminderInterval  = 1
	twentyFourhours               = time.Hour * 24
	defaultOrgDeletionRunInterval = twentyFourhours
)

// createPeriodicJobs creates periodic jobs for all runtime modules
func createPeriodicJobs(c Workers) ([]*river.PeriodicJob, error) {
	j := []*river.PeriodicJob{}

	if c.DeleteExportContentWorker.Config.Enabled {
		interval := c.DeleteExportContentWorker.Config.Interval
		if interval < minInterval {
			interval = minInterval
		}

		deleteExportJobs := river.NewPeriodicJob(
			river.PeriodicInterval(interval),
			func() (river.JobArgs, *river.InsertOpts) {
				return jobs.DeleteExportContentArgs{}, nil
			},
			&river.PeriodicJobOpts{RunOnStart: true},
		)
		j = append(j, deleteExportJobs)
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
		j = append(j, orgDeletionReminders)
	}

	if c.OrganizationDeletionWorker.Config.Enabled {
		interval := c.OrganizationDeletionWorker.Config.RunInterval
		if interval == 0 {
			interval = defaultOrgDeletionRunInterval
		}

		orgDeletion := river.NewPeriodicJob(
			river.PeriodicInterval(interval),
			func() (river.JobArgs, *river.InsertOpts) {
				return jobspec.OrganizationDeletionArgs{}, nil
			},
			&river.PeriodicJobOpts{},
		)
		j = append(j, orgDeletion)
	}

	return j, nil
}
