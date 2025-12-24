package river

import (
	"time"

	"github.com/riverqueue/river"
	"github.com/theopenlane/riverboat/pkg/jobs"
)

// minInterval prevents spamming from a bag config
var minInterval = 1 * time.Minute

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

	return j, nil
}
