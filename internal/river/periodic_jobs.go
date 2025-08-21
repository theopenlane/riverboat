package river

import (
	"github.com/riverqueue/river"
	"github.com/theopenlane/core/pkg/corejobs"
)

func createPeriodicJobs(c Workers) ([]*river.PeriodicJob, error) {
	jobs := []*river.PeriodicJob{}
	if c.DeleteExportContentWorker.Config.Enabled {
		deleteExportJobs := river.NewPeriodicJob(
			river.PeriodicInterval(c.DeleteExportContentWorker.Config.Interval),
			func() (river.JobArgs, *river.InsertOpts) {
				return corejobs.DeleteExportContentArgs{}, nil
			},
			&river.PeriodicJobOpts{RunOnStart: true},
		)
		jobs = append(jobs, deleteExportJobs)
	}

	return jobs, nil
}
