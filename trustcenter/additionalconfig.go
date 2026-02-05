//go:build !trustcenter

package trustcenter

import "github.com/theopenlane/riverboat/pkg/jobs"

// Workers is an empty struct when trustcenter build tag is not present
type Workers struct {
	// OpenlaneConfig configuration for openlane jobs, this is shared across multiple workers
	// if a worker needs specific configuration, it can be set in the worker's config
	OpenlaneConfig jobs.OpenlaneConfig `koanf:"openlaneconfig" json:"openlaneconfig"`
}
