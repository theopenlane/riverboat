package jobs_test

import (
	"context"
	"testing"

	"github.com/riverqueue/river"
	"github.com/stretchr/testify/require"

	"github.com/theopenlane/riverboat/pkg/jobs"
)

func (suite *TestSuite) TestDomainWorker() {
	t := suite.T()

	testCases := []struct {
		name          string
		worker        *jobs.DomainWorker
		domain        string
		expectedError string
	}{
		{
			name:   "happy path",
			worker: &jobs.DomainWorker{},
			domain: "google.com",
		},
		{
			name:          "missing domain",
			worker:        &jobs.DomainWorker{},
			expectedError: "domain is required for the domain job",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			err := (tc.worker).Work(ctx, &river.Job[jobs.DomainArgs]{Args: jobs.DomainArgs{
				Domain: tc.domain,
			}})

			if tc.expectedError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError)

				return
			}

			require.NoError(t, err)
		})
	}
}
