package jobs_test

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
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
	"github.com/theopenlane/core/common/models"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/theopenlane/riverboat/pkg/jobs"
	olmocks "github.com/theopenlane/riverboat/pkg/jobs/openlane/mocks"
)

func ptr[T any](v T) *T {
	return &v
}

func newTestEvidenceNode(id, name string, fileIDs []string, controlRefCodes []string) *graphclient.GetEvidences_Evidences_Edges_Node {
	var fileEdges []*graphclient.GetEvidences_Evidences_Edges_Node_Files_Edges
	for _, fid := range fileIDs {
		fileEdges = append(fileEdges, &graphclient.GetEvidences_Evidences_Edges_Node_Files_Edges{
			Node: &graphclient.GetEvidences_Evidences_Edges_Node_Files_Edges_Node{
				ID:           fid,
				PresignedURL: ptr("PLACEHOLDER_URL"),
			},
		})
	}

	var controlEdges []*graphclient.GetEvidences_Evidences_Edges_Node_Controls_Edges
	for _, rc := range controlRefCodes {
		controlEdges = append(controlEdges, &graphclient.GetEvidences_Evidences_Edges_Node_Controls_Edges{
			Node: &graphclient.GetEvidences_Evidences_Edges_Node_Controls_Edges_Node{
				ID:      "ctrl-" + rc,
				RefCode: rc,
			},
		})
	}

	return &graphclient.GetEvidences_Evidences_Edges_Node{
		ID:   id,
		Name: name,
		Files: graphclient.GetEvidences_Evidences_Edges_Node_Files{
			Edges:      fileEdges,
			TotalCount: int64(len(fileEdges)),
		},
		Controls: graphclient.GetEvidences_Evidences_Edges_Node_Controls{
			Edges: controlEdges,
		},
	}
}

func newFileDetail(id, providedName, ext string) *graphclient.GetFileByID {
	return &graphclient.GetFileByID{
		File: graphclient.GetFileByID_File{
			ID:                    id,
			ProvidedFileName:      providedName,
			ProvidedFileExtension: ext,
		},
	}
}

func makeEvidenceEdges(nodes ...*graphclient.GetEvidences_Evidences_Edges_Node) []*graphclient.GetEvidences_Evidences_Edges {
	edges := make([]*graphclient.GetEvidences_Evidences_Edges, len(nodes))
	for i, n := range nodes {
		edges[i] = &graphclient.GetEvidences_Evidences_Edges{Node: n}
	}

	return edges
}

// fileDownloadServer creates a test server that serves file content by path.
func fileDownloadServer(t *testing.T, files map[string][]byte) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if data, ok := files[path]; ok {
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(data) //nolint:errcheck
		} else {
			http.NotFound(w, r)
		}
	}))
}

// extractZipContents reads a zip from bytes and returns a map of path -> content.
func extractZipContents(t *testing.T, data []byte) map[string]string {
	t.Helper()

	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	require.NoError(t, err)

	result := make(map[string]string)
	for _, f := range r.File {
		rc, err := f.Open()
		require.NoError(t, err)

		var buf bytes.Buffer
		_, err = buf.ReadFrom(rc)
		require.NoError(t, err)

		rc.Close() //nolint:errcheck

		result[f.Name] = buf.String()
	}

	return result
}

func TestExportEvidenceFiles_HappyPathFlat(t *testing.T) {
	t.Parallel()

	exportID := "export-flat-001"
	userID := "user-001"
	orgID := "org-001"
	ownerID := orgID

	fileContent := map[string][]byte{
		"file1": []byte("pdf-content-here"),
		"file2": []byte("csv-content-here"),
	}
	dlServer := fileDownloadServer(t, fileContent)
	defer dlServer.Close()

	ev1 := newTestEvidenceNode("ev1", "Quarterly Access Review", []string{"f1", "f2"}, []string{"SOC2::CC1.1"})
	ev1.Files.Edges[0].Node.PresignedURL = ptr(dlServer.URL + "/file1")
	ev1.Files.Edges[1].Node.PresignedURL = ptr(dlServer.URL + "/file2")

	olMock := olmocks.NewMockGraphClient(t)

	olMock.EXPECT().GetExportByID(mock.Anything, exportID).Return(&graphclient.GetExportByID{
		Export: graphclient.GetExportByID_Export{
			ID:         exportID,
			ExportType: enums.ExportTypeEvidence,
			OwnerID:    &ownerID,
		},
	}, nil)

	olMock.EXPECT().GetEvidences(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&graphclient.GetEvidences{
			Evidences: graphclient.GetEvidences_Evidences{
				Edges:      makeEvidenceEdges(ev1),
				TotalCount: 1,
				PageInfo:   graphclient.GetEvidences_Evidences_PageInfo{HasNextPage: false},
			},
		}, nil)

	olMock.EXPECT().GetFileByID(mock.Anything, "f1", mock.Anything).Return(newFileDetail("f1", "access-review", ".pdf"), nil)
	olMock.EXPECT().GetFileByID(mock.Anything, "f2", mock.Anything).Return(newFileDetail("f2", "scan-results", ".csv"), nil)

	var uploadedZip []byte

	olMock.EXPECT().UpdateExport(mock.Anything, exportID,
		mock.MatchedBy(func(input graphclient.UpdateExportInput) bool {
			return input.Status != nil && *input.Status == enums.ExportStatusReady
		}),
		mock.MatchedBy(func(files []*graphql.Upload) bool {
			if len(files) != 1 || files[0].ContentType != "application/zip" {
				return false
			}

			buf := new(bytes.Buffer)
			buf.ReadFrom(files[0].File) //nolint:errcheck
			uploadedZip = buf.Bytes()

			return true
		}),
		mock.Anything,
	).Return(&graphclient.UpdateExport{}, nil)

	worker := &jobs.ExportContentWorker{
		Config: jobs.ExportWorkerConfig{
			OpenlaneConfig: jobs.OpenlaneConfig{
				OpenlaneAPIHost:  "https://api.example.com",
				OpenlaneAPIToken: "tola_test-token",
			},
			MaxZipSize: 50 * 1024 * 1024,
		},
	}
	worker.WithOpenlaneClient(olMock)

	err := worker.Work(context.Background(), &river.Job[jobspec.ExportContentArgs]{
		Args: jobspec.ExportContentArgs{
			ExportID:       exportID,
			UserID:         userID,
			OrganizationID: orgID,
			Mode:           enums.ExportModeFlat,
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, uploadedZip)

	contents := extractZipContents(t, uploadedZip)
	require.Len(t, contents, 2)

	var foundPDF, foundCSV bool
	for path, content := range contents {
		require.NotContains(t, path, "/", "FLAT mode should have no directories")

		if strings.HasSuffix(path, ".pdf") {
			require.Equal(t, "pdf-content-here", content)
			foundPDF = true
		}

		if strings.HasSuffix(path, ".csv") {
			require.Equal(t, "csv-content-here", content)
			foundCSV = true
		}
	}

	require.True(t, foundPDF, "should contain PDF file")
	require.True(t, foundCSV, "should contain CSV file")
}

func TestExportEvidenceFiles_HappyPathFolder(t *testing.T) {
	t.Parallel()

	exportID := "export-folder-001"
	userID := "user-001"
	orgID := "org-001"
	ownerID := orgID

	fileContent := map[string][]byte{
		"file1": []byte("pdf-content"),
	}
	dlServer := fileDownloadServer(t, fileContent)
	defer dlServer.Close()

	ev1 := newTestEvidenceNode("ev1", "Quarterly Access Review", []string{"f1"}, []string{"SOC2::CC1.1"})
	ev1.Files.Edges[0].Node.PresignedURL = ptr(dlServer.URL + "/file1")

	olMock := olmocks.NewMockGraphClient(t)

	olMock.EXPECT().GetExportByID(mock.Anything, exportID).Return(&graphclient.GetExportByID{
		Export: graphclient.GetExportByID_Export{
			ID:         exportID,
			ExportType: enums.ExportTypeEvidence,
			OwnerID:    &ownerID,
		},
	}, nil)

	olMock.EXPECT().GetEvidences(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&graphclient.GetEvidences{
			Evidences: graphclient.GetEvidences_Evidences{
				Edges:      makeEvidenceEdges(ev1),
				TotalCount: 1,
				PageInfo:   graphclient.GetEvidences_Evidences_PageInfo{HasNextPage: false},
			},
		}, nil)

	olMock.EXPECT().GetFileByID(mock.Anything, "f1", mock.Anything).Return(newFileDetail("f1", "access-review", ".pdf"), nil)

	var uploadedZip []byte

	olMock.EXPECT().UpdateExport(mock.Anything, exportID,
		mock.MatchedBy(func(input graphclient.UpdateExportInput) bool {
			return input.Status != nil && *input.Status == enums.ExportStatusReady
		}),
		mock.MatchedBy(func(files []*graphql.Upload) bool {
			if len(files) != 1 || files[0].ContentType != "application/zip" {
				return false
			}

			buf := new(bytes.Buffer)
			buf.ReadFrom(files[0].File) //nolint:errcheck
			uploadedZip = buf.Bytes()

			return true
		}),
		mock.Anything,
	).Return(&graphclient.UpdateExport{}, nil)

	olMock.EXPECT().GetControls(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.MatchedBy(func(where *graphclient.ControlWhereInput) bool {
			return where != nil && where.RefCode != nil && *where.RefCode == "SOC2::CC1.1"
		}),
		mock.Anything,
	).Return(&graphclient.GetControls{
		Controls: graphclient.GetControls_Controls{
			Edges: []*graphclient.GetControls_Controls_Edges{
				{
					Node: &graphclient.GetControls_Controls_Edges_Node{
						RefCode:            "SOC2::CC1.1",
						ReferenceFramework: ptr("SOC 2"),
						ReferenceID:        ptr("CC1.1"),
						AuditorReferenceID: ptr("AUD-001"),
					},
				},
			},
		},
	}, nil)

	worker := &jobs.ExportContentWorker{
		Config: jobs.ExportWorkerConfig{
			OpenlaneConfig: jobs.OpenlaneConfig{
				OpenlaneAPIHost:  "https://api.example.com",
				OpenlaneAPIToken: "tola_test-token",
			},
			MaxZipSize: 50 * 1024 * 1024,
		},
	}
	worker.WithOpenlaneClient(olMock)

	err := worker.Work(context.Background(), &river.Job[jobspec.ExportContentArgs]{
		Args: jobspec.ExportContentArgs{
			ExportID:       exportID,
			UserID:         userID,
			OrganizationID: orgID,
			Mode:           enums.ExportModeFolder,
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, uploadedZip)

	contents := extractZipContents(t, uploadedZip)

	require.Contains(t, contents, "SOC2--CC1.1/metadata.txt")

	metadata := contents["SOC2--CC1.1/metadata.txt"]
	require.Contains(t, metadata, "RefCode: SOC2::CC1.1")
	require.Contains(t, metadata, "ReferenceFramework: SOC 2")
	require.Contains(t, metadata, "Quarterly Access Review")

	var foundPDF bool
	for path, content := range contents {
		if strings.HasPrefix(path, "SOC2--CC1.1/") && strings.HasSuffix(path, ".pdf") {
			require.Equal(t, "pdf-content", content)
			foundPDF = true
		}
	}

	require.True(t, foundPDF, "should contain PDF file in folder")
}

func TestExportEvidenceFiles_FolderWithUncategorized(t *testing.T) {
	t.Parallel()

	exportID := "export-uncat-001"
	userID := "user-001"
	orgID := "org-001"
	ownerID := orgID

	fileContent := map[string][]byte{
		"file1": []byte("orphan-content"),
	}
	dlServer := fileDownloadServer(t, fileContent)
	defer dlServer.Close()

	ev1 := newTestEvidenceNode("ev1", "Orphan Evidence", []string{"f1"}, nil)
	ev1.Files.Edges[0].Node.PresignedURL = ptr(dlServer.URL + "/file1")

	olMock := olmocks.NewMockGraphClient(t)

	olMock.EXPECT().GetExportByID(mock.Anything, exportID).Return(&graphclient.GetExportByID{
		Export: graphclient.GetExportByID_Export{
			ID:         exportID,
			ExportType: enums.ExportTypeEvidence,
			OwnerID:    &ownerID,
		},
	}, nil)

	olMock.EXPECT().GetEvidences(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&graphclient.GetEvidences{
			Evidences: graphclient.GetEvidences_Evidences{
				Edges:      makeEvidenceEdges(ev1),
				TotalCount: 1,
				PageInfo:   graphclient.GetEvidences_Evidences_PageInfo{HasNextPage: false},
			},
		}, nil)

	olMock.EXPECT().GetFileByID(mock.Anything, "f1", mock.Anything).Return(newFileDetail("f1", "orphan", ".docx"), nil)

	var uploadedZip []byte

	olMock.EXPECT().UpdateExport(mock.Anything, exportID,
		mock.MatchedBy(func(input graphclient.UpdateExportInput) bool {
			return input.Status != nil && *input.Status == enums.ExportStatusReady
		}),
		mock.MatchedBy(func(files []*graphql.Upload) bool {
			if len(files) != 1 {
				return false
			}

			buf := new(bytes.Buffer)
			buf.ReadFrom(files[0].File) //nolint:errcheck
			uploadedZip = buf.Bytes()

			return true
		}),
		mock.Anything,
	).Return(&graphclient.UpdateExport{}, nil)

	worker := &jobs.ExportContentWorker{
		Config: jobs.ExportWorkerConfig{
			OpenlaneConfig: jobs.OpenlaneConfig{
				OpenlaneAPIHost:  "https://api.example.com",
				OpenlaneAPIToken: "tola_test-token",
			},
			MaxZipSize: 50 * 1024 * 1024,
		},
	}
	worker.WithOpenlaneClient(olMock)

	err := worker.Work(context.Background(), &river.Job[jobspec.ExportContentArgs]{
		Args: jobspec.ExportContentArgs{
			ExportID:       exportID,
			UserID:         userID,
			OrganizationID: orgID,
			Mode:           enums.ExportModeFolder,
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, uploadedZip)

	contents := extractZipContents(t, uploadedZip)
	require.Contains(t, contents, "_uncategorized/metadata.txt")

	var foundDocx bool
	for path := range contents {
		if strings.HasPrefix(path, "_uncategorized/") && strings.HasSuffix(path, ".docx") {
			foundDocx = true
		}
	}

	require.True(t, foundDocx, "should contain docx file in _uncategorized folder")
}

func TestExportEvidenceFiles_FolderMultipleControls(t *testing.T) {
	t.Parallel()

	exportID := "export-multi-001"
	userID := "user-001"
	orgID := "org-001"
	ownerID := orgID

	fileContent := map[string][]byte{
		"file1": []byte("shared-file-content"),
	}
	dlServer := fileDownloadServer(t, fileContent)
	defer dlServer.Close()

	ev1 := newTestEvidenceNode("ev1", "Quarterly Access Review", []string{"f1"}, []string{"SOC2::CC1.1", "ISO27001::A.5.1"})
	ev1.Files.Edges[0].Node.PresignedURL = ptr(dlServer.URL + "/file1")

	olMock := olmocks.NewMockGraphClient(t)

	olMock.EXPECT().GetExportByID(mock.Anything, exportID).Return(&graphclient.GetExportByID{
		Export: graphclient.GetExportByID_Export{
			ID:         exportID,
			ExportType: enums.ExportTypeEvidence,
			OwnerID:    &ownerID,
		},
	}, nil)

	olMock.EXPECT().GetEvidences(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&graphclient.GetEvidences{
			Evidences: graphclient.GetEvidences_Evidences{
				Edges:      makeEvidenceEdges(ev1),
				TotalCount: 1,
				PageInfo:   graphclient.GetEvidences_Evidences_PageInfo{HasNextPage: false},
			},
		}, nil)

	olMock.EXPECT().GetFileByID(mock.Anything, "f1", mock.Anything).Return(newFileDetail("f1", "access-review", ".pdf"), nil)

	var uploadedZip []byte

	olMock.EXPECT().UpdateExport(mock.Anything, exportID,
		mock.MatchedBy(func(input graphclient.UpdateExportInput) bool {
			return input.Status != nil && *input.Status == enums.ExportStatusReady
		}),
		mock.MatchedBy(func(files []*graphql.Upload) bool {
			if len(files) != 1 {
				return false
			}

			buf := new(bytes.Buffer)
			buf.ReadFrom(files[0].File) //nolint:errcheck
			uploadedZip = buf.Bytes()

			return true
		}),
		mock.Anything,
	).Return(&graphclient.UpdateExport{}, nil)

	olMock.EXPECT().GetControls(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.MatchedBy(func(where *graphclient.ControlWhereInput) bool {
			return where != nil && where.RefCode != nil && *where.RefCode == "SOC2::CC1.1"
		}),
		mock.Anything,
	).Return(&graphclient.GetControls{
		Controls: graphclient.GetControls_Controls{
			Edges: []*graphclient.GetControls_Controls_Edges{
				{
					Node: &graphclient.GetControls_Controls_Edges_Node{
						RefCode:            "SOC2::CC1.1",
						ReferenceFramework: ptr("SOC 2"),
					},
				},
			},
		},
	}, nil)

	olMock.EXPECT().GetControls(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.MatchedBy(func(where *graphclient.ControlWhereInput) bool {
			return where != nil && where.RefCode != nil && *where.RefCode == "ISO27001::A.5.1"
		}),
		mock.Anything,
	).Return(&graphclient.GetControls{
		Controls: graphclient.GetControls_Controls{
			Edges: []*graphclient.GetControls_Controls_Edges{
				{
					Node: &graphclient.GetControls_Controls_Edges_Node{
						RefCode:            "ISO27001::A.5.1",
						ReferenceFramework: ptr("ISO 27001"),
					},
				},
			},
		},
	}, nil)

	worker := &jobs.ExportContentWorker{
		Config: jobs.ExportWorkerConfig{
			OpenlaneConfig: jobs.OpenlaneConfig{
				OpenlaneAPIHost:  "https://api.example.com",
				OpenlaneAPIToken: "tola_test-token",
			},
			MaxZipSize: 50 * 1024 * 1024,
		},
	}
	worker.WithOpenlaneClient(olMock)

	err := worker.Work(context.Background(), &river.Job[jobspec.ExportContentArgs]{
		Args: jobspec.ExportContentArgs{
			ExportID:       exportID,
			UserID:         userID,
			OrganizationID: orgID,
			Mode:           enums.ExportModeFolder,
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, uploadedZip)

	contents := extractZipContents(t, uploadedZip)

	var foundInSOC2, foundInISO bool
	for path, content := range contents {
		if strings.HasPrefix(path, "SOC2--CC1.1/") && strings.HasSuffix(path, ".pdf") {
			require.Equal(t, "shared-file-content", content)
			foundInSOC2 = true
		}

		if strings.HasPrefix(path, "ISO27001--A.5.1/") && strings.HasSuffix(path, ".pdf") {
			require.Equal(t, "shared-file-content", content)
			foundInISO = true
		}
	}

	require.True(t, foundInSOC2, "should contain file in SOC2 folder")
	require.True(t, foundInISO, "should contain file in ISO folder")
}

func TestExportEvidenceFiles_NoEvidencesFound(t *testing.T) {
	t.Parallel()

	exportID := "export-nodata-001"
	userID := "user-001"
	orgID := "org-001"
	ownerID := orgID

	olMock := olmocks.NewMockGraphClient(t)

	olMock.EXPECT().GetExportByID(mock.Anything, exportID).Return(&graphclient.GetExportByID{
		Export: graphclient.GetExportByID_Export{
			ID:         exportID,
			ExportType: enums.ExportTypeEvidence,
			OwnerID:    &ownerID,
		},
	}, nil)

	olMock.EXPECT().GetEvidences(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&graphclient.GetEvidences{
			Evidences: graphclient.GetEvidences_Evidences{
				Edges:      nil,
				TotalCount: 0,
				PageInfo:   graphclient.GetEvidences_Evidences_PageInfo{HasNextPage: false},
			},
		}, nil)

	olMock.EXPECT().UpdateExport(mock.Anything, exportID,
		mock.MatchedBy(func(input graphclient.UpdateExportInput) bool {
			return input.Status != nil && *input.Status == enums.ExportStatusNodata
		}),
		([]*graphql.Upload)(nil),
	).Return(&graphclient.UpdateExport{}, nil)

	worker := &jobs.ExportContentWorker{
		Config: jobs.ExportWorkerConfig{
			OpenlaneConfig: jobs.OpenlaneConfig{
				OpenlaneAPIHost:  "https://api.example.com",
				OpenlaneAPIToken: "tola_test-token",
			},
		},
	}
	worker.WithOpenlaneClient(olMock)

	err := worker.Work(context.Background(), &river.Job[jobspec.ExportContentArgs]{
		Args: jobspec.ExportContentArgs{
			ExportID:       exportID,
			UserID:         userID,
			OrganizationID: orgID,
			Mode:           enums.ExportModeFlat,
		},
	})

	require.NoError(t, err)
}

func TestExportEvidenceFiles_ErrorFetchingEvidences(t *testing.T) {
	t.Parallel()

	exportID := "export-err-001"
	userID := "user-001"
	orgID := "org-001"
	ownerID := orgID

	olMock := olmocks.NewMockGraphClient(t)

	olMock.EXPECT().GetExportByID(mock.Anything, exportID).Return(&graphclient.GetExportByID{
		Export: graphclient.GetExportByID_Export{
			ID:         exportID,
			ExportType: enums.ExportTypeEvidence,
			OwnerID:    &ownerID,
		},
	}, nil)

	olMock.EXPECT().GetEvidences(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, assert.AnError)

	olMock.EXPECT().UpdateExport(mock.Anything, exportID,
		mock.MatchedBy(func(input graphclient.UpdateExportInput) bool {
			return input.Status != nil && *input.Status == enums.ExportStatusFailed
		}),
		([]*graphql.Upload)(nil),
	).Return(&graphclient.UpdateExport{}, nil)

	worker := &jobs.ExportContentWorker{
		Config: jobs.ExportWorkerConfig{
			OpenlaneConfig: jobs.OpenlaneConfig{
				OpenlaneAPIHost:  "https://api.example.com",
				OpenlaneAPIToken: "tola_test-token",
			},
		},
	}
	worker.WithOpenlaneClient(olMock)

	err := worker.Work(context.Background(), &river.Job[jobspec.ExportContentArgs]{
		Args: jobspec.ExportContentArgs{
			ExportID:       exportID,
			UserID:         userID,
			OrganizationID: orgID,
			Mode:           enums.ExportModeFlat,
		},
	})

	require.NoError(t, err)
}

func TestExportEvidenceFiles_ErrorDownloadingFile(t *testing.T) {
	t.Parallel()

	exportID := "export-dl-err-001"
	userID := "user-001"
	orgID := "org-001"
	ownerID := orgID

	// server that returns 500 for all requests
	dlServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer dlServer.Close()

	ev1 := newTestEvidenceNode("ev1", "Bad Evidence", []string{"f1"}, nil)
	ev1.Files.Edges[0].Node.PresignedURL = ptr(dlServer.URL + "/file1")

	olMock := olmocks.NewMockGraphClient(t)

	olMock.EXPECT().GetExportByID(mock.Anything, exportID).Return(&graphclient.GetExportByID{
		Export: graphclient.GetExportByID_Export{
			ID:         exportID,
			ExportType: enums.ExportTypeEvidence,
			OwnerID:    &ownerID,
		},
	}, nil)

	olMock.EXPECT().GetEvidences(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&graphclient.GetEvidences{
			Evidences: graphclient.GetEvidences_Evidences{
				Edges:      makeEvidenceEdges(ev1),
				TotalCount: 1,
				PageInfo:   graphclient.GetEvidences_Evidences_PageInfo{HasNextPage: false},
			},
		}, nil)

	olMock.EXPECT().GetFileByID(mock.Anything, "f1", mock.Anything).Return(newFileDetail("f1", "bad-file", ".pdf"), nil)

	// should update status to failed
	olMock.EXPECT().UpdateExport(mock.Anything, exportID,
		mock.MatchedBy(func(input graphclient.UpdateExportInput) bool {
			return input.Status != nil && *input.Status == enums.ExportStatusFailed
		}),
		([]*graphql.Upload)(nil),
	).Return(&graphclient.UpdateExport{}, nil)

	worker := &jobs.ExportContentWorker{
		Config: jobs.ExportWorkerConfig{
			OpenlaneConfig: jobs.OpenlaneConfig{
				OpenlaneAPIHost:  "https://api.example.com",
				OpenlaneAPIToken: "tola_test-token",
			},
			MaxZipSize: 50 * 1024 * 1024,
		},
	}
	worker.WithOpenlaneClient(olMock)

	err := worker.Work(context.Background(), &river.Job[jobspec.ExportContentArgs]{
		Args: jobspec.ExportContentArgs{
			ExportID:       exportID,
			UserID:         userID,
			OrganizationID: orgID,
			Mode:           enums.ExportModeFlat,
		},
	})

	require.NoError(t, err)
}

func TestExportEvidenceFiles_KeepFileOriginalName(t *testing.T) {
	t.Parallel()

	exportID := "export-orig-001"
	userID := "user-001"
	orgID := "org-001"
	ownerID := orgID

	fileContent := map[string][]byte{
		"file1": []byte("original-name-content"),
	}
	dlServer := fileDownloadServer(t, fileContent)
	defer dlServer.Close()

	ev1 := newTestEvidenceNode("ev1", "Quarterly Access Review", []string{"f1"}, nil)
	ev1.Files.Edges[0].Node.PresignedURL = ptr(dlServer.URL + "/file1")

	olMock := olmocks.NewMockGraphClient(t)

	olMock.EXPECT().GetExportByID(mock.Anything, exportID).Return(&graphclient.GetExportByID{
		Export: graphclient.GetExportByID_Export{
			ID:         exportID,
			ExportType: enums.ExportTypeEvidence,
			OwnerID:    &ownerID,
		},
	}, nil)

	olMock.EXPECT().GetEvidences(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&graphclient.GetEvidences{
			Evidences: graphclient.GetEvidences_Evidences{
				Edges:      makeEvidenceEdges(ev1),
				TotalCount: 1,
				PageInfo:   graphclient.GetEvidences_Evidences_PageInfo{HasNextPage: false},
			},
		}, nil)

	olMock.EXPECT().GetFileByID(mock.Anything, "f1", mock.Anything).Return(newFileDetail("f1", "My Original Document", ".pdf"), nil)

	var uploadedZip []byte

	olMock.EXPECT().UpdateExport(mock.Anything, exportID,
		mock.MatchedBy(func(input graphclient.UpdateExportInput) bool {
			return input.Status != nil && *input.Status == enums.ExportStatusReady
		}),
		mock.MatchedBy(func(files []*graphql.Upload) bool {
			if len(files) != 1 {
				return false
			}

			buf := new(bytes.Buffer)
			buf.ReadFrom(files[0].File) //nolint:errcheck
			uploadedZip = buf.Bytes()

			return true
		}),
		mock.Anything,
	).Return(&graphclient.UpdateExport{}, nil)

	worker := &jobs.ExportContentWorker{
		Config: jobs.ExportWorkerConfig{
			OpenlaneConfig: jobs.OpenlaneConfig{
				OpenlaneAPIHost:  "https://api.example.com",
				OpenlaneAPIToken: "tola_test-token",
			},
			MaxZipSize: 50 * 1024 * 1024,
		},
	}
	worker.WithOpenlaneClient(olMock)

	err := worker.Work(context.Background(), &river.Job[jobspec.ExportContentArgs]{
		Args: jobspec.ExportContentArgs{
			ExportID:       exportID,
			UserID:         userID,
			OrganizationID: orgID,
			Mode:           enums.ExportModeFlat,
			ExportMetadata: &models.ExportMetadata{
				KeepFileOriginalName: true,
			},
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, uploadedZip)

	contents := extractZipContents(t, uploadedZip)
	require.Contains(t, contents, "My Original Document.pdf")
	require.Equal(t, "original-name-content", contents["My Original Document.pdf"])
}

func TestExportEvidenceFiles_DefaultModeCSVExport(t *testing.T) {
	t.Parallel()

	// when mode is empty, should fall through to CSV export
	exportID := "export-csv-001"
	ownerID := "org-001"

	olMock := olmocks.NewMockGraphClient(t)

	olMock.EXPECT().GetExportByID(mock.Anything, exportID).Return(&graphclient.GetExportByID{
		Export: graphclient.GetExportByID_Export{
			ID:         exportID,
			ExportType: enums.ExportTypeControl,
			OwnerID:    &ownerID,
			Fields:     []string{"id", "name"},
		},
	}, nil)

	// mock graphQL server for CSV export
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{ //nolint:errcheck
			"data": map[string]interface{}{
				"controls": map[string]interface{}{
					"edges": []interface{}{
						map[string]interface{}{
							"node": map[string]interface{}{
								"id":   "ctrl-1",
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
		})
	}))
	defer mockServer.Close()

	// the existing CSV export path should be taken - expect CSV upload
	olMock.EXPECT().UpdateExport(mock.Anything, exportID,
		mock.MatchedBy(func(input graphclient.UpdateExportInput) bool {
			return input.Status != nil && *input.Status == enums.ExportStatusReady
		}),
		mock.MatchedBy(func(files []*graphql.Upload) bool {
			return len(files) == 1 && files[0].ContentType == "text/csv"
		}),
		mock.Anything,
	).Return(&graphclient.UpdateExport{}, nil)

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
		Args: jobspec.ExportContentArgs{
			ExportID:       exportID,
			UserID:         "user-001",
			OrganizationID: ownerID,
			// mode is empty - should use CSV path
		},
	})

	require.NoError(t, err)
}

func TestExportEvidenceFiles_MissingRequiredArgs(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		args          jobspec.ExportContentArgs
		expectedError string
	}{
		{
			name: "missing export ID",
			args: jobspec.ExportContentArgs{
				UserID:         "user-001",
				OrganizationID: "org-001",
				Mode:           enums.ExportModeFlat,
			},
			expectedError: "export_id is required",
		},
		{
			name: "missing user ID",
			args: jobspec.ExportContentArgs{
				ExportID:       "export-001",
				OrganizationID: "org-001",
				Mode:           enums.ExportModeFlat,
			},
			expectedError: "user_id is required",
		},
		{
			name: "missing organization ID",
			args: jobspec.ExportContentArgs{
				ExportID: "export-001",
				UserID:   "user-001",
				Mode:     enums.ExportModeFlat,
			},
			expectedError: "organization_id is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			worker := &jobs.ExportContentWorker{
				Config: jobs.ExportWorkerConfig{
					OpenlaneConfig: jobs.OpenlaneConfig{
						OpenlaneAPIHost:  "https://api.example.com",
						OpenlaneAPIToken: "tola_test-token",
					},
				},
			}

			olMock := olmocks.NewMockGraphClient(t)
			worker.WithOpenlaneClient(olMock)

			err := worker.Work(context.Background(), &river.Job[jobspec.ExportContentArgs]{Args: tc.args})
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.expectedError)
		})
	}
}

