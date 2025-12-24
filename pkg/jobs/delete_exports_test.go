package jobs_test

import (
	"context"
	"testing"
	"time"

	"github.com/riverqueue/river"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/theopenlane/core/pkg/enums"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/theopenlane/riverboat/pkg/jobs"

	olmocks "github.com/theopenlane/riverboat/pkg/jobs/openlane/mocks"
)

func TestDeleteExportContentWorker(t *testing.T) {
	t.Parallel()

	exportID1 := "export123"
	exportID2 := "export456"
	cutoffDuration := 30 * time.Minute

	testCases := []struct {
		name                   string
		input                  jobs.DeleteExportContentArgs
		config                 jobs.DeleteExportWorkerConfig
		getExportsResponse     *graphclient.GetExports
		getExportsError        error
		deleteBulkExportError  error
		expectedDeletedIDs     []string
		expectedError          string
		expectDeleteBulkExport bool
	}{
		{
			name:  "happy path - delete exports",
			input: jobs.DeleteExportContentArgs{},
			config: jobs.DeleteExportWorkerConfig{
				OpenlaneConfig: jobs.OpenlaneConfig{
					OpenlaneAPIHost:  "https://api.example.com",
					OpenlaneAPIToken: "tolp_test-token",
				},
				CutoffDuration: cutoffDuration,
			},
			getExportsResponse: &graphclient.GetExports{
				Exports: graphclient.GetExports_Exports{
					Edges: []*graphclient.GetExports_Exports_Edges{
						{
							Node: &graphclient.GetExports_Exports_Edges_Node{
								ID: exportID1,
							},
						},
						{
							Node: &graphclient.GetExports_Exports_Edges_Node{
								ID: exportID2,
							},
						},
					},
				},
			},
			expectedDeletedIDs:     []string{exportID1, exportID2},
			expectDeleteBulkExport: true,
		},
		{
			name:  "happy path - no exports to delete",
			input: jobs.DeleteExportContentArgs{},
			config: jobs.DeleteExportWorkerConfig{
				OpenlaneConfig: jobs.OpenlaneConfig{
					OpenlaneAPIHost:  "https://api.example.com",
					OpenlaneAPIToken: "test-token",
				},
				CutoffDuration: cutoffDuration,
			},
			getExportsResponse: &graphclient.GetExports{
				Exports: graphclient.GetExports_Exports{
					Edges: []*graphclient.GetExports_Exports_Edges{},
				},
			},
			expectDeleteBulkExport: false,
		},
		{
			name:  "error getting exports",
			input: jobs.DeleteExportContentArgs{},
			config: jobs.DeleteExportWorkerConfig{
				OpenlaneConfig: jobs.OpenlaneConfig{
					OpenlaneAPIHost:  "https://api.example.com",
					OpenlaneAPIToken: "test-token",
				},
				CutoffDuration: cutoffDuration,
			},
			getExportsError:        assert.AnError,
			expectedError:          "assert.AnError",
			expectDeleteBulkExport: false,
		},
		{
			name:  "error during bulk delete",
			input: jobs.DeleteExportContentArgs{},
			config: jobs.DeleteExportWorkerConfig{
				OpenlaneConfig: jobs.OpenlaneConfig{
					OpenlaneAPIHost:  "https://api.example.com",
					OpenlaneAPIToken: "test-token",
				},
				CutoffDuration: cutoffDuration,
			},
			getExportsResponse: &graphclient.GetExports{
				Exports: graphclient.GetExports_Exports{
					Edges: []*graphclient.GetExports_Exports_Edges{
						{
							Node: &graphclient.GetExports_Exports_Edges_Node{
								ID: exportID1,
							},
						},
					},
				},
			},
			deleteBulkExportError:  assert.AnError,
			expectedDeletedIDs:     []string{exportID1},
			expectedError:          "assert.AnError",
			expectDeleteBulkExport: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			olMock := olmocks.NewMockGraphClient(t)

			startTime := time.Now()
			pageSize := int64(100) //nolint:mnd

			olMock.EXPECT().GetExports(mock.MatchedBy(func(ctx context.Context) bool {
				return ctx != nil
			}), &pageSize, (*int64)(nil), (*string)(nil), (*string)(nil), mock.MatchedBy(func(input *graphclient.ExportWhereInput) bool {
				if input == nil {
					return false
				}

				if input.CreatedAtLte == nil {
					return false
				}

				expectedCutoff := startTime.Add(-tc.config.CutoffDuration)
				actualCutoff := *input.CreatedAtLte
				timeDiff := expectedCutoff.Sub(actualCutoff)
				if timeDiff < 0 {
					timeDiff = -timeDiff
				}

				if timeDiff > time.Second {
					return false
				}

				if input.StatusIn == nil {
					return false
				}

				expectedStatuses := []enums.ExportStatus{
					enums.ExportStatusNodata,
					enums.ExportStatusReady,
				}

				if len(input.StatusIn) != len(expectedStatuses) {
					return false
				}

				statusMap := make(map[enums.ExportStatus]bool)
				for _, status := range input.StatusIn {
					statusMap[status] = true
				}

				for _, expectedStatus := range expectedStatuses {
					if !statusMap[expectedStatus] {
						return false
					}
				}

				return true
			}), mock.Anything).Return(tc.getExportsResponse, tc.getExportsError)

			if tc.expectDeleteBulkExport {
				olMock.EXPECT().DeleteBulkExport(mock.MatchedBy(func(ctx context.Context) bool {
					return ctx != nil
				}), tc.expectedDeletedIDs).Return(&graphclient.DeleteBulkExport{}, tc.deleteBulkExportError)
			}

			worker := &jobs.DeleteExportContentWorker{
				Config: tc.config,
			}

			worker.WithOpenlaneClient(olMock)

			ctx := context.Background()
			err := worker.Work(ctx, &river.Job[jobs.DeleteExportContentArgs]{Args: tc.input})

			if tc.expectedError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestDeleteExportContentWorker_WithOpenlaneClient(t *testing.T) {
	t.Parallel()

	olMock := olmocks.NewMockGraphClient(t)
	worker := &jobs.DeleteExportContentWorker{}

	result := worker.WithOpenlaneClient(olMock)

	require.Equal(t, worker, result, "WithOpenlaneClient should return the same worker instance")
}

func TestDeleteExportContentArgs_Kind(t *testing.T) {
	t.Parallel()

	args := jobs.DeleteExportContentArgs{}
	require.Equal(t, "delete_export_content", args.Kind())
}
