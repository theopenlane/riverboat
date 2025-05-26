package jobs_test

import (
	"context"

	"github.com/riverqueue/river"
	"github.com/stretchr/testify/require"

	"github.com/theopenlane/riverboat/pkg/jobs"
)

func (suite *TestSuite) TestValidateCustomDomainWorker() {
	t := suite.T()
	worker := &jobs.ValidateCustomDomainWorker{
		Config: jobs.CustomDomainConfig{
			CloudflareAPIKey: "test",
		},
	}
	ctx := context.Background()

	err := worker.Work(ctx, &river.Job[jobs.ValidateCustomDomainArgs]{Args: jobs.ValidateCustomDomainArgs{
		CustomDomainID: "test",
	}})

	require.NoError(t, err)
}
