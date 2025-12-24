//go:build !trustcenter

package trustcenter

import (
	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"

	"github.com/theopenlane/riverboat/pkg/riverqueue"
)

// AddConditionalWorkers is a no-op when trust center build tag is not present
func AddConditionalWorkers(worker *river.Workers, _ any, _ *riverqueue.Client) (*river.Workers, error) {
	log.Info().Msg("no additional workers to add for non-trustcenter build")

	return worker, nil
}
