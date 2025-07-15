package river

import (
	"github.com/riverqueue/river"
	"github.com/theopenlane/core/pkg/corejobs"
)

func createPeriodicJobs(c Workers) ([]*river.PeriodicJob, error) {
	deleteExportJobs := river.NewPeriodicJob(
		river.PeriodicInterval(c.DeleteExportContentWorker.Config.CutoffDuration),
		func() (river.JobArgs, *river.InsertOpts) {
			return corejobs.DeleteExportContentArgs{}, nil
		},
		&river.PeriodicJobOpts{RunOnStart: true},
	)

	return []*river.PeriodicJob{deleteExportJobs}, nil
}
