package river

import (
	"time"

	"github.com/riverqueue/river"
	"github.com/theopenlane/core/pkg/corejobs"
)

// minInterval prevents spamming from a bag config
var minInterval = 1 * time.Minute

func createPeriodicJobs(c Workers) ([]*river.PeriodicJob, error) {
	jobs := []*river.PeriodicJob{}

	if c.DeleteExportContentWorker.Config.Enabled {
		interval := c.DeleteExportContentWorker.Config.Interval
		if interval < minInterval {
			interval = minInterval
		}

		deleteExportJobs := river.NewPeriodicJob(
			river.PeriodicInterval(interval),
			func() (river.JobArgs, *river.InsertOpts) {
				return corejobs.DeleteExportContentArgs{}, nil
			},
			&river.PeriodicJobOpts{RunOnStart: true},
		)
		jobs = append(jobs, deleteExportJobs)
	}

	if c.ValidateCustomDomainWorker.Config.Enabled {
		interval := c.ValidateCustomDomainWorker.Config.ValidateInterval
		if interval < minInterval {
			interval = minInterval
		}

		validateCustomDomainJobs := river.NewPeriodicJob(
			river.PeriodicInterval(interval),
			func() (river.JobArgs, *river.InsertOpts) {
				return corejobs.ValidateCustomDomainArgs{}, nil
			},
			&river.PeriodicJobOpts{RunOnStart: true},
		)
		jobs = append(jobs, validateCustomDomainJobs)
	}

	if c.CacheTrustCenterDataWorker.Config.Enabled {
		interval := c.CacheTrustCenterDataWorker.Config.CacheInterval
		if interval < minInterval {
			interval = minInterval
		}

		cacheTrustCenterDataJobs := river.NewPeriodicJob(
			river.PeriodicInterval(interval),
			func() (river.JobArgs, *river.InsertOpts) {
				return corejobs.CacheTrustCenterDataArgs{}, nil
			},
			&river.PeriodicJobOpts{RunOnStart: true},
		)
		jobs = append(jobs, cacheTrustCenterDataJobs)
	}

	return jobs, nil
}
