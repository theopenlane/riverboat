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

	return jobs, nil
}
