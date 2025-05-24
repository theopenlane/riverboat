package jobs_test

import (
	"context"

	"github.com/riverqueue/river"
	"github.com/stretchr/testify/require"
	"github.com/theopenlane/riverboat/pkg/jobs"
)

func (suite *TestSuite) TestDeleteCustomDomainWorker() {
	t := suite.T()
	worker := &jobs.DeleteCustomDomainWorker{
		Config: jobs.DeleteCustomDomainConfig{
			CloudflareAPIKey: "test",
		},
	}
	ctx := context.Background()

	err := worker.Work(ctx, &river.Job[jobs.DeleteCustomDomainArgs]{Args: jobs.DeleteCustomDomainArgs{
		CustomDomainID: "test",
	}})

	require.NoError(t, err)
}
