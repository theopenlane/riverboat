//go:build !trustcenter

package river

import (
	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"
)

// createAdditionalPeriodicJobs creates periodic jobs for additional modules
// this is a no-op when trust center build tag is not present
func createAdditionalPeriodicJobs(c any) ([]*river.PeriodicJob, error) {
	log.Info().Msg("no additional periodic jobs to add for non-trustcenter build")

	return nil, nil
}
