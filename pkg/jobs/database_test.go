package jobs_test

import (
	"context"
	"testing"

	"github.com/riverqueue/river"
	"github.com/stretchr/testify/require"
	"github.com/theopenlane/dbx/pkg/dbxclient"
	"github.com/theopenlane/utils/ulids"

	"github.com/theopenlane/riverboat/pkg/jobs"
)

// TODO :this currently does not test the actual database creation because the dbx client is not mocked
// so the only thing we can test is that the worker returns early when disabled
func (suite *TestSuite) TestDatabaseWorker() {
	t := suite.T()

	testCases := []struct {
		name          string
		worker        *jobs.DatabaseWorker
		args          jobs.DatabaseArgs
		expectedError string
	}{
		{
			name: "happy path, skip while disabled",
			worker: &jobs.DatabaseWorker{
				Config: dbxclient.Config{
					Enabled: false,
				},
			},
			args: jobs.DatabaseArgs{
				Location:       "AMER",
				OrganizationID: ulids.New().String(),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			err := (tc.worker).Work(ctx, &river.Job[jobs.DatabaseArgs]{Args: tc.args})

			if tc.expectedError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError)

				return
			}

			require.NoError(t, err)
		})
	}
}
