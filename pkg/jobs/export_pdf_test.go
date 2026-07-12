package jobs_test

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/cloudflare/cloudflare-go/v7"
	"github.com/cloudflare/cloudflare-go/v7/shared"
	"github.com/gqlgo/gqlgenc/clientv2"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/theopenlane/core/common/enums"
	"github.com/theopenlane/core/common/jobspec"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/theopenlane/riverboat/pkg/jobs"
	pdfmocks "github.com/theopenlane/riverboat/pkg/jobs/mocks"
	olmocks "github.com/theopenlane/riverboat/pkg/jobs/openlane/mocks"
)

func TestExportContentWorkerPDF(t *testing.T) {
	t.Parallel()

	const (
		exportID = "export123"
		ownerID  = "owner123"
	)

	input := jobspec.ExportContentArgs{
		ExportID:       exportID,
		UserID:         "user123",
		OrganizationID: ownerID,
	}

	t.Run("single document uploads a standalone pdf", func(t *testing.T) {
		graphServer := mockGraphQLServer(t, controlsResponse(
			map[string]interface{}{"id": "c1", "name": "Access Control Policy", "details": "<p>body</p>"},
		))
		defer graphServer.Close()

		olMock := olmocks.NewMockGraphClient(t)
		olMock.EXPECT().GetExportByID(mock.Anything, exportID).Return(newPDFExport(exportID, ownerID), nil)

		var captured []*graphql.Upload

		olMock.EXPECT().UpdateExport(mock.Anything, exportID, mock.Anything, mock.Anything, mock.Anything).
			Run(func(_ context.Context, _ string, _ graphclient.UpdateExportInput, exportFiles []*graphql.Upload, _ ...clientv2.RequestInterceptor) {
				captured = exportFiles
			}).Return(&graphclient.UpdateExport{}, nil)

		worker := newPDFWorker(graphServer.URL, fakePDFRenderer(t)).WithOpenlaneClient(olMock)

		err := worker.Work(context.Background(), &river.Job[jobspec.ExportContentArgs]{Args: input})
		require.NoError(t, err)

		require.Len(t, captured, 1)
		assert.Equal(t, "application/pdf", captured[0].ContentType)
		assert.Equal(t, "access_control_policy.pdf", captured[0].Filename)
	})

	t.Run("multiple documents are bundled into a zip", func(t *testing.T) {
		graphServer := mockGraphQLServer(t, controlsResponse(
			map[string]interface{}{"id": "c1", "name": "Access Control Policy", "details": "<p>a</p>"},
			map[string]interface{}{"id": "c2", "name": "Data Retention Policy", "details": "<p>b</p>"},
		))
		defer graphServer.Close()

		olMock := olmocks.NewMockGraphClient(t)
		olMock.EXPECT().GetExportByID(mock.Anything, exportID).Return(newPDFExport(exportID, ownerID), nil)

		var captured []*graphql.Upload

		olMock.EXPECT().UpdateExport(mock.Anything, exportID, mock.Anything, mock.Anything, mock.Anything).
			Run(func(_ context.Context, _ string, _ graphclient.UpdateExportInput, exportFiles []*graphql.Upload, _ ...clientv2.RequestInterceptor) {
				captured = exportFiles
			}).Return(&graphclient.UpdateExport{}, nil)

		worker := newPDFWorker(graphServer.URL, fakePDFRenderer(t)).WithOpenlaneClient(olMock)

		err := worker.Work(context.Background(), &river.Job[jobspec.ExportContentArgs]{Args: input})
		require.NoError(t, err)

		require.Len(t, captured, 1)
		assert.Equal(t, "application/zip", captured[0].ContentType)

		names := zipEntryNames(t, captured[0])
		assert.ElementsMatch(t, []string{"access_control_policy.pdf", "data_retention_policy.pdf"}, names)
	})

	t.Run("unconfigured cloudflare reports unsupported export type", func(t *testing.T) {
		graphServer := mockGraphQLServer(t, controlsResponse(
			map[string]interface{}{"id": "c1", "name": "Access Control Policy", "details": "<p>a</p>"},
		))
		defer graphServer.Close()

		olMock := olmocks.NewMockGraphClient(t)
		olMock.EXPECT().GetExportByID(mock.Anything, exportID).Return(newPDFExport(exportID, ownerID), nil)

		worker := &jobs.ExportContentWorker{
			Config: jobs.ExportWorkerConfig{
				OpenlaneConfig: jobs.OpenlaneConfig{
					OpenlaneAPIHost:  graphServer.URL,
					OpenlaneAPIToken: "tola_test-token",
				},
			},
		}
		worker.WithOpenlaneClient(olMock)

		err := worker.Work(context.Background(), &river.Job[jobspec.ExportContentArgs]{Args: input})
		require.ErrorIs(t, err, jobs.ErrUnsupportedExportType)
	})

	t.Run("cloudflare 422 error snoozes for retries", func(t *testing.T) {
		srv := mockGraphQLServer(t, controlsResponse(
			map[string]interface{}{"id": "c1", "name": "Access Control Policy", "details": "<p>a</p>"},
		))
		defer srv.Close()

		olMock := olmocks.NewMockGraphClient(t)
		olMock.EXPECT().GetExportByID(mock.Anything, exportID).Return(newPDFExport(exportID, ownerID), nil)

		cfErr := cloudflareRenderUnprocessableError(t)
		renderer := pdfmocks.NewMockPDFRenderer(t)
		renderer.EXPECT().HTMLToPDF(mock.Anything, mock.Anything).Return(nil, cfErr)

		worker := newPDFWorker(srv.URL, renderer).WithOpenlaneClient(olMock)

		err := worker.Work(context.Background(), &river.Job[jobspec.ExportContentArgs]{
			JobRow: &rivertype.JobRow{
				Metadata: []byte(`{"Snoozes":0}`),
			},
			Args: input,
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "JobSnoozeError")
	})

	t.Run("cloudflare 422 returns max snoozes error when retries are exhausted", func(t *testing.T) {
		srv := mockGraphQLServer(t, controlsResponse(
			map[string]interface{}{"id": "c1", "name": "Access Control Policy", "details": "<p>a</p>"},
		))
		defer srv.Close()

		olMock := olmocks.NewMockGraphClient(t)
		olMock.EXPECT().GetExportByID(mock.Anything, exportID).Return(newPDFExport(exportID, ownerID), nil)

		cfErr := cloudflareRenderUnprocessableError(t)
		renderer := pdfmocks.NewMockPDFRenderer(t)
		renderer.EXPECT().HTMLToPDF(mock.Anything, mock.Anything).Return(nil, cfErr)

		worker := newPDFWorker(srv.URL, renderer).WithOpenlaneClient(olMock)

		err := worker.Work(context.Background(), &river.Job[jobspec.ExportContentArgs]{
			JobRow: &rivertype.JobRow{
				Metadata: []byte(`{"Snoozes":30}`),
			},
			Args: input,
		})
		require.ErrorIs(t, err, jobs.ErrMaxSnoozesReached)
	})
}

// fakePDFRenderer returns a mock renderer that produces deterministic fake PDF bytes
func fakePDFRenderer(t *testing.T) *pdfmocks.MockPDFRenderer {
	m := pdfmocks.NewMockPDFRenderer(t)
	m.EXPECT().HTMLToPDF(mock.Anything, mock.Anything).Return([]byte("%PDF-1.4 fake"), nil)

	return m
}

func newPDFExport(exportID, ownerID string) *graphclient.GetExportByID {
	return &graphclient.GetExportByID{
		Export: graphclient.GetExportByID_Export{
			ID:         exportID,
			ExportType: enums.ExportTypeControl,
			Format:     enums.ExportFormatPdf,
			OwnerID:    &ownerID,
			Fields:     []string{"id", "name", "details"},
		},
	}
}

func cloudflareRenderUnprocessableError(t *testing.T) error {
	t.Helper()

	reqURL, err := url.Parse("https://api.cloudflare.com/client/v4/accounts/cf-account/browser-rendering/pdf")
	require.NoError(t, err)

	return fmt.Errorf("failed to call cloudflare browser rendering: %w", &cloudflare.Error{
		StatusCode: http.StatusUnprocessableEntity,
		Errors: []shared.ErrorData{
			{
				Code:    1000,
				Message: "unprocessable browser rendering request",
			},
		},
		Request: &http.Request{
			Method: http.MethodPost,
			URL:    reqURL,
		},
		Response: &http.Response{
			StatusCode: http.StatusUnprocessableEntity,
		},
	})
}

func controlsResponse(nodes ...map[string]interface{}) map[string]interface{} {
	edges := make([]interface{}, 0, len(nodes))
	for _, n := range nodes {
		edges = append(edges, map[string]interface{}{"node": n})
	}

	return map[string]interface{}{
		"controls": map[string]interface{}{
			"edges": edges,
			"pageInfo": map[string]interface{}{
				"hasNextPage": false,
				"endCursor":   nil,
			},
		},
	}
}

func newPDFWorker(graphURL string, renderer jobs.PDFRenderer) *jobs.ExportContentWorker {
	worker := &jobs.ExportContentWorker{
		Config: jobs.ExportWorkerConfig{
			OpenlaneConfig: jobs.OpenlaneConfig{
				OpenlaneAPIHost:  graphURL,
				OpenlaneAPIToken: "tola_test-token",
			},
			MaxZipSize:          52428800,
			CloudflareAccountID: "cf-account",
			CloudflareAPIKey:    "cf-test-key",
			MaxSnoozes:          30,
			SnoozeDuration:      10 * time.Second,
		},
	}

	return worker.WithPDFRenderer(renderer)
}

func zipEntryNames(t *testing.T, upload *graphql.Upload) []string {
	data, err := io.ReadAll(upload.File)
	require.NoError(t, err)

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	require.NoError(t, err)

	names := make([]string, 0, len(zr.File))
	for _, f := range zr.File {
		names = append(names, f.Name)
	}

	return names
}
