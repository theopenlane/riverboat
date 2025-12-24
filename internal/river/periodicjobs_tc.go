//go:build trustcenter

package river

import (
	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"
	"github.com/theopenlane/core/pkg/jobspec"
)

// createAdditionalPeriodicJobs creates periodic jobs for the trust center module
func createAdditionalPeriodicJobs(c AdditionalWorkers) ([]*river.PeriodicJob, error) {
	log.Info().Msg("adding additional trust censter periodic jobs")

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
	}

	return jobs, nil
}
