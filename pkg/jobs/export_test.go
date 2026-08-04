package jobs_test

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/riverqueue/river"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/theopenlane/core/common/enums"
	"github.com/theopenlane/core/common/jobspec"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/theopenlane/riverboat/pkg/jobs"
	olmocks "github.com/theopenlane/riverboat/pkg/jobs/openlane/mocks"
)

func mockGraphQLServer(t *testing.T, responses map[string]interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestBody struct {
			Query string `json:"query"`
		}

		err := json.NewDecoder(r.Body).Decode(&requestBody)
		require.NoError(t, err)

		var response map[string]interface{}
		if strings.Contains(requestBody.Query, "controls") {
			response = map[string]interface{}{
				"controls": responses["controls"],
			}
		} else if strings.Contains(requestBody.Query, "evidences") {
			response = map[string]interface{}{
				"evidences": responses["evidences"],
			}
		} else {
			response = responses
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": response,
		})
	}))
}

func TestExportContentWorker_ControlSubcontrolsExport(t *testing.T) {
	t.Parallel()

	exportID := "export123"
	ownerID := "owner123"

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestBody struct {
			Query string `json:"query"`
		}

		err := json.NewDecoder(r.Body).Decode(&requestBody)
		require.NoError(t, err)
		require.Contains(t, requestBody.Query, "subcontrols")
		require.Contains(t, requestBody.Query, "refCode")
		require.Contains(t, requestBody.Query, "description")
		require.NotContains(t, requestBody.Query, "subcontrolKindName")
		require.NotContains(t, requestBody.Query, "title")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"controls": map[string]interface{}{
					"edges": []interface{}{
						map[string]interface{}{
							"node": map[string]interface{}{
								"description": "Control description",
								"status":      "NOT_IMPLEMENTED",
								"subcontrols": map[string]interface{}{
									"edges": []interface{}{
										map[string]interface{}{
											"node": map[string]interface{}{
												"refCode":     "CC1.2-POF1",
												"description": "Subcontrol 1",
												"status":      "NOT_IMPLEMENTED",
											},
										},
										map[string]interface{}{
											"node": map[string]interface{}{
												"refCode":     "CC1.2-POF2",
												"description": "Subcontrol 2",
												"status":      "APPROVED",
											},
										},
									},
								},
							},
						},
					},
					"pageInfo": map[string]interface{}{
						"hasNextPage": false,
						"endCursor":   nil,
					},
				},
			},
		})
	}))
	defer mockServer.Close()

	olMock := olmocks.NewMockGraphClient(t)
	olMock.EXPECT().GetExportByID(mock.MatchedBy(func(ctx context.Context) bool {
		return ctx != nil
	}), exportID).Return(&graphclient.GetExportByID{
		Export: graphclient.GetExportByID_Export{
			ID:         exportID,
			ExportType: enums.ExportTypeControl,
			Format:     enums.ExportFormatCsv,
			OwnerID:    &ownerID,
			Fields:     []string{"description", "status", "subcontrols"},
		},
	}, nil)

	var uploadedCSVContent []byte

	olMock.EXPECT().UpdateExport(mock.MatchedBy(func(ctx context.Context) bool {
		return ctx != nil
	}), exportID, mock.MatchedBy(func(input graphclient.UpdateExportInput) bool {
		return input.Status != nil && *input.Status == enums.ExportStatusReady
	}), mock.MatchedBy(func(files []*graphql.Upload) bool {
		if len(files) != 1 || files[0].ContentType != "text/csv" {
			return false
		}

		var err error
		uploadedCSVContent, err = io.ReadAll(files[0].File)

		return err == nil
	}), mock.Anything).Return(&graphclient.UpdateExport{}, nil)

	worker := &jobs.ExportContentWorker{
		Config: jobs.ExportWorkerConfig{
			OpenlaneConfig: jobs.OpenlaneConfig{
				OpenlaneAPIHost:  mockServer.URL,
				OpenlaneAPIToken: "tola_test-token",
			},
		},
	}
	worker.WithOpenlaneClient(olMock)

	err := worker.Work(context.Background(), &river.Job[jobspec.ExportContentArgs]{
		Args: jobspec.ExportContentArgs{ExportID: exportID, UserID: "user123", OrganizationID: ownerID},
	})
	require.NoError(t, err)

	records, err := csv.NewReader(bytes.NewReader(uploadedCSVContent)).ReadAll()
	require.NoError(t, err)
	require.Len(t, records, 4)

	headers := records[0]
	descriptionIndex := getHeaderIndex(t, headers, "description")
	statusIndex := getHeaderIndex(t, headers, "status")
	subcontrolsIndex := getHeaderIndex(t, headers, "subcontrols")

	require.NotContains(t, headers, "refCode")
	require.NotContains(t, headers, "subcontrolKindName")
	require.NotContains(t, headers, "title")

	require.Equal(t, "Control description", records[1][descriptionIndex])
	require.Equal(t, "NOT_IMPLEMENTED", records[1][statusIndex])
	require.Equal(t, "CC1.2-POF1, CC1.2-POF2", records[1][subcontrolsIndex])

	require.Equal(t, "Subcontrol 1", records[2][descriptionIndex])
	require.Equal(t, "NOT_IMPLEMENTED", records[2][statusIndex])
	require.Empty(t, records[2][subcontrolsIndex])

	require.Equal(t, "Subcontrol 2", records[3][descriptionIndex])
	require.Equal(t, "APPROVED", records[3][statusIndex])
	require.Empty(t, records[3][subcontrolsIndex])
}

func getHeaderIndex(t *testing.T, headers []string, header string) int {
	t.Helper()

	for i, h := range headers {
		if h == header {
			return i
		}
	}

	require.Failf(t, "missing CSV header", "header %q not found in %v", header, headers)

	return -1
}

func TestExportContentWorker(t *testing.T) {
	t.Parallel()

	exportID := "export123"
	controlID1 := "control123"
	controlID2 := "control456"
	ownerID := "owner123"

	testCases := []struct {
		name                     string
		input                    jobspec.ExportContentArgs
		getExportByIDResponse    *graphclient.GetExportByID
		getExportByIDError       error
		graphQLResponses         map[string]interface{}
		updateExportError        error
		updateExportErrorSecond  error
		expectedError            string
		expectGetExportByID      bool
		expectUpdateExport       bool
		expectUpdateExportFiles  bool
		expectUpdateExportSecond bool
		expectedExportStatus     *enums.ExportStatus
	}{
		{
			name: "happy path - export controls",
			input: jobspec.ExportContentArgs{
				ExportID:       exportID,
				UserID:         "user123",
				OrganizationID: ownerID,
			},
			getExportByIDResponse: &graphclient.GetExportByID{
				Export: graphclient.GetExportByID_Export{
					ID:         exportID,
					ExportType: enums.ExportTypeControl,
					OwnerID:    &ownerID,
					Fields:     []string{"id", "name", "description"},
				},
			},
			graphQLResponses: map[string]interface{}{
				"controls": map[string]interface{}{
					"edges": []interface{}{
						map[string]interface{}{
							"node": map[string]interface{}{
								"id":          controlID1,
								"name":        "Control 1",
								"description": "First control",
							},
						},
						map[string]interface{}{
							"node": map[string]interface{}{
								"id":          controlID2,
								"name":        "Control 2",
								"description": "Second control",
							},
						},
					},
					"pageInfo": map[string]interface{}{
						"hasNextPage": false,
						"endCursor":   nil,
					},
				},
			},
			expectGetExportByID:     true,
			expectUpdateExport:      true,
			expectUpdateExportFiles: true,
			expectedExportStatus:    &enums.ExportStatusReady,
		},
		{
			name: "missing export ID",
			input: jobspec.ExportContentArgs{
				ExportID:       "",
				UserID:         "user123",
				OrganizationID: ownerID,
			},
			expectedError:       "export_id is required for the export_content job",
			expectGetExportByID: false,
		},
		{
			name: "missing organization ID",
			input: jobspec.ExportContentArgs{
				ExportID:       exportID,
				UserID:         "user123",
				OrganizationID: "",
			},
			expectedError:       "organization_id is required for the export_content job",
			expectGetExportByID: false,
		},
		{
			name: "missing user ID",
			input: jobspec.ExportContentArgs{
				ExportID:       exportID,
				UserID:         "",
				OrganizationID: ownerID,
			},
			expectedError:       "user_id is required for the export_content job",
			expectGetExportByID: false,
		},
		{
			name: "error getting export",
			input: jobspec.ExportContentArgs{
				ExportID:       exportID,
				UserID:         "user123",
				OrganizationID: ownerID,
			},
			getExportByIDError:   assert.AnError,
			expectGetExportByID:  true,
			expectUpdateExport:   true,
			expectedExportStatus: &enums.ExportStatusFailed,
		},
		{
			name: "no data found for export",
			input: jobspec.ExportContentArgs{
				ExportID:       exportID,
				UserID:         "user123",
				OrganizationID: ownerID,
			},
			getExportByIDResponse: &graphclient.GetExportByID{
				Export: graphclient.GetExportByID_Export{
					ID:         exportID,
					ExportType: enums.ExportTypeControl,
					OwnerID:    &ownerID,
					Fields:     []string{"id", "name"},
				},
			},
			graphQLResponses: map[string]interface{}{
				"controls": map[string]interface{}{
					"edges": []interface{}{},
					"pageInfo": map[string]interface{}{
						"hasNextPage": false,
						"endCursor":   nil,
					},
				},
			},
			expectGetExportByID:  true,
			expectUpdateExport:   true,
			expectedExportStatus: &enums.ExportStatusNodata,
		},
		{
			name: "error updating export with file",
			input: jobspec.ExportContentArgs{
				ExportID:       exportID,
				UserID:         "user123",
				OrganizationID: ownerID,
			},
			getExportByIDResponse: &graphclient.GetExportByID{
				Export: graphclient.GetExportByID_Export{
					ID:         exportID,
					ExportType: enums.ExportTypeControl,
					OwnerID:    &ownerID,
					Fields:     []string{"id", "name"},
				},
			},
			graphQLResponses: map[string]interface{}{
				"controls": map[string]interface{}{
					"edges": []interface{}{
						map[string]interface{}{
							"node": map[string]interface{}{
								"id":   controlID1,
								"name": "Control 1",
							},
						},
					},
					"pageInfo": map[string]interface{}{
						"hasNextPage": false,
						"endCursor":   nil,
					},
				},
			},
			updateExportError:        assert.AnError,
			expectGetExportByID:      true,
			expectUpdateExport:       true,
			expectUpdateExportFiles:  true,
			expectUpdateExportSecond: true,
			expectedExportStatus:     &enums.ExportStatusFailed,
		},
		{
			name: "export evidence type",
			input: jobspec.ExportContentArgs{
				ExportID:       exportID,
				UserID:         "user123",
				OrganizationID: ownerID,
			},
			getExportByIDResponse: &graphclient.GetExportByID{
				Export: graphclient.GetExportByID_Export{
					ID:         exportID,
					ExportType: enums.ExportTypeEvidence,
					OwnerID:    &ownerID,
					Fields:     []string{"id", "name", "type"},
				},
			},
			graphQLResponses: map[string]interface{}{
				"evidences": map[string]interface{}{
					"edges": []interface{}{
						map[string]interface{}{
							"node": map[string]interface{}{
								"id":   "evidence123",
								"name": "Evidence 1",
								"type": "Document",
							},
						},
					},
					"pageInfo": map[string]interface{}{
						"hasNextPage": false,
						"endCursor":   nil,
					},
				},
			},
			expectGetExportByID:     true,
			expectUpdateExport:      true,
			expectUpdateExportFiles: true,
			expectedExportStatus:    &enums.ExportStatusReady,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var mockServer *httptest.Server
			if tc.graphQLResponses != nil {
				mockServer = mockGraphQLServer(t, tc.graphQLResponses)
				defer mockServer.Close()
			}

			olMock := olmocks.NewMockGraphClient(t)

			if tc.expectGetExportByID {
				olMock.EXPECT().GetExportByID(mock.MatchedBy(func(ctx context.Context) bool {
					return ctx != nil
				}), tc.input.ExportID).Return(tc.getExportByIDResponse, tc.getExportByIDError)
			}

			if tc.expectUpdateExport {
				if tc.expectUpdateExportFiles {
					olMock.EXPECT().UpdateExport(mock.MatchedBy(func(ctx context.Context) bool {
						return ctx != nil
					}), tc.input.ExportID, mock.MatchedBy(func(input graphclient.UpdateExportInput) bool {
						return input.Status != nil && *input.Status == enums.ExportStatusReady
					}), mock.MatchedBy(func(files []*graphql.Upload) bool {
						return len(files) == 1 && files[0].ContentType == "text/csv"
					}), mock.Anything,
					).Return(&graphclient.UpdateExport{}, tc.updateExportError)

					if tc.expectUpdateExportSecond {
						olMock.EXPECT().UpdateExport(mock.MatchedBy(func(ctx context.Context) bool {
							return ctx != nil
						}), tc.input.ExportID, mock.MatchedBy(func(input graphclient.UpdateExportInput) bool {
							return input.Status != nil && *input.Status == *tc.expectedExportStatus
						}), ([]*graphql.Upload)(nil)).Return(&graphclient.UpdateExport{}, tc.updateExportErrorSecond)
					}
				} else {
					olMock.EXPECT().UpdateExport(mock.MatchedBy(func(ctx context.Context) bool {
						return ctx != nil
					}), tc.input.ExportID, mock.MatchedBy(func(input graphclient.UpdateExportInput) bool {
						if input.Status == nil {
							return false
						}
						expectedStatus := *tc.expectedExportStatus
						actualStatus := *input.Status
						return actualStatus == expectedStatus
					}), ([]*graphql.Upload)(nil)).Return(&graphclient.UpdateExport{}, tc.updateExportError)
				}
			}

			config := jobs.ExportWorkerConfig{
				OpenlaneConfig: jobs.OpenlaneConfig{
					OpenlaneAPIHost:  "https://api.example.com",
					OpenlaneAPIToken: "tola_test-token",
				},
			}
			if mockServer != nil {
				config.OpenlaneAPIHost = mockServer.URL
			}

			worker := &jobs.ExportContentWorker{
				Config: config,
			}

			worker.WithOpenlaneClient(olMock)

			ctx := context.Background()
			err := worker.Work(ctx, &river.Job[jobspec.ExportContentArgs]{Args: tc.input})

			if tc.expectedError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestExportContentWorker_GraphQLErrorResponse(t *testing.T) {
	t.Parallel()

	exportID := "export123"
	ownerID := "owner123"

	// mock server that returns GraphQL errors
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": []interface{}{
				map[string]interface{}{
					"message": "Field 'controls' not found",
				},
			},
		})
	}))
	defer mockServer.Close()

	olMock := olmocks.NewMockGraphClient(t)

	// mock successful GetExportByID
	olMock.EXPECT().GetExportByID(mock.MatchedBy(func(ctx context.Context) bool {
		return ctx != nil
	}), exportID).Return(&graphclient.GetExportByID{
		Export: graphclient.GetExportByID_Export{
			ID:         exportID,
			ExportType: enums.ExportTypeControl,
			OwnerID:    &ownerID,
			Fields:     []string{"id", "name"},
		},
	}, nil)

	// mock failed UpdateExport call
	olMock.EXPECT().UpdateExport(mock.MatchedBy(func(ctx context.Context) bool {
		return ctx != nil
	}), exportID, mock.MatchedBy(func(input graphclient.UpdateExportInput) bool {
		return input.Status != nil && *input.Status == enums.ExportStatusFailed
	}), ([]*graphql.Upload)(nil)).Return(&graphclient.UpdateExport{}, nil)

	worker := &jobs.ExportContentWorker{
		Config: jobs.ExportWorkerConfig{
			OpenlaneConfig: jobs.OpenlaneConfig{
				OpenlaneAPIHost:  mockServer.URL,
				OpenlaneAPIToken: "tola_test-token",
			},
		},
	}

	worker.WithOpenlaneClient(olMock)

	ctx := context.Background()
	err := worker.Work(ctx, &river.Job[jobspec.ExportContentArgs]{
		Args: jobspec.ExportContentArgs{ExportID: exportID, UserID: "user123", OrganizationID: ownerID},
	})

	require.NoError(t, err)
}

func TestExportContentWorker_WithOpenlaneClient(t *testing.T) {
	t.Parallel()

	olMock := olmocks.NewMockGraphClient(t)
	worker := &jobs.ExportContentWorker{}

	result := worker.WithOpenlaneClient(olMock)

	require.Equal(t, worker, result, "WithOpenlaneClient should return the same worker instance")
}

func TestExportContentArgs_Kind(t *testing.T) {
	t.Parallel()

	args := jobspec.ExportContentArgs{}
	require.Equal(t, "export_content", args.Kind())
}
