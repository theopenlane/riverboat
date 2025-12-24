//go:build !trustcenter

package river

import (
	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"
	"github.com/theopenlane/riverboat/pkg/riverqueue"
)

// addConditionalWorkers is a no-op when trust center build tag is not present
func addConditionalWorkers(worker *river.Workers, w any, insertOnlyClient *riverqueue.Client) (*river.Workers, error) {
	log.Info().Msg("no additional workers to add for non-trustcenter build")

	return worker, nil
}
