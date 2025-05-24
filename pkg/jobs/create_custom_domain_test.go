package jobs_test

import (
	"context"

	"github.com/riverqueue/river"
	"github.com/stretchr/testify/require"
	"github.com/theopenlane/riverboat/pkg/jobs"
)

func (suite *TestSuite) TestCreateCustomDomainWorker() {
	t := suite.T()
	worker := &jobs.CreateCustomDomainWorker{
		Config: jobs.CreateCustomDomainConfig{
			CloudflareAPIKey: "test",
		},
	}
	ctx := context.Background()

	err := worker.Work(ctx, &river.Job[jobs.CreateCustomDomainArgs]{Args: jobs.CreateCustomDomainArgs{
		CustomDomainID: "test",
	}})

	require.NoError(t, err)
}
