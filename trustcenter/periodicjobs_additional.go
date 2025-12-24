//go:build !trustcenter

package trustcenter

import (
	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"
)

// CreatePeriodicJobs creates periodic jobs for additional modules
// this is a no-op when trust center build tag is not present
func CreatePeriodicJobs(_ Workers) ([]*river.PeriodicJob, error) {
	log.Info().Msg("no additional periodic jobs to add for non-trustcenter build")

	return nil, nil
}
