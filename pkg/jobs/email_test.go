package jobs_test

import (
	"context"
	"testing"

	"github.com/riverqueue/river"
	"github.com/stretchr/testify/require"
	"github.com/theopenlane/newman"

	"github.com/theopenlane/riverboat/pkg/jobs"
)

func (suite *TestSuite) TestEmailWorker() {
	t := suite.T()

	emailWithFrom := newman.NewEmailMessageWithOptions(
		newman.WithTo([]string{"ted@mosby.com"}),
		newman.WithFrom("robin@scherbatsky.com"),
	)

	emailWithoutFrom := newman.NewEmailMessageWithOptions(
		newman.WithTo([]string{"ted@mosby.com"}),
	)

	testCases := []struct {
		name          string
		worker        *jobs.EmailWorker
		msg           *newman.EmailMessage
		expectedError string
	}{
		{
			name: "happy path, dev mode",
			worker: &jobs.EmailWorker{
				EmailConfig: jobs.EmailConfig{
					DevMode:   true,
					TestDir:   "test",
					FromEmail: "robin@scherbatsky.net",
				},
			},
			msg: emailWithoutFrom,
		},
		{
			name: "missing test directory",
			worker: &jobs.EmailWorker{
				EmailConfig: jobs.EmailConfig{
					DevMode:   true,
					FromEmail: "robin@scherbatsky.net",
				},
			},
			msg:           emailWithoutFrom,
			expectedError: jobs.ErrMissingTestDir.Error(),
		},
		{
			name: "missing from email",
			worker: &jobs.EmailWorker{
				EmailConfig: jobs.EmailConfig{
					DevMode: true,
					TestDir: "test",
				},
			},
			msg:           emailWithoutFrom,
			expectedError: "from is required",
		},
		{
			name: "happy path, missing from email but in message",
			worker: &jobs.EmailWorker{
				EmailConfig: jobs.EmailConfig{
					DevMode: true,
					TestDir: "test",
				},
			},
			msg: emailWithFrom,
		},
		{
			name: "missing token",
			worker: &jobs.EmailWorker{
				EmailConfig: jobs.EmailConfig{
					DevMode:   false,
					FromEmail: "robin@scherbatsky.net",
				},
			},
			msg:           emailWithoutFrom,
			expectedError: jobs.ErrMissingToken.Error(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			err := (tc.worker).Work(ctx, &river.Job[jobs.EmailArgs]{Args: jobs.EmailArgs{
				Message: *tc.msg,
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
