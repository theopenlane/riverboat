package jobs

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/cloudflare/cloudflare-go/v7"
	"github.com/gertd/go-pluralize"
	"github.com/gocarina/gocsv"
	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"
	"github.com/stoewer/go-strcase"
	"github.com/theopenlane/core/common/enums"
	"github.com/theopenlane/core/common/jobspec"
	goclient "github.com/theopenlane/go-client"
	"github.com/theopenlane/go-client/graphclient"
	"github.com/theopenlane/httpsling"
	"github.com/theopenlane/iam/auth"

	"github.com/theopenlane/riverboat/pkg/jobs/openlane"
	"github.com/theopenlane/riverboat/pkg/render"
)

var defaultPageSize int64 = 100

const (
	defaultExportMaxSnoozes     = 5
	defaultExportSnoozeDuration = 10 * time.Second
)

// pdfFieldsWanted include the fields used in the PDF export we need to make sure we have in the graphql request
var (
	pdfFieldsWanted        = []string{"details", "externalContents", "revision", "name", "createdAt", "updatedAt", "status"}
	additionalPolicyFields = []string{"liveExternalContents"}
)

var (
	// ErrUnexpectedStatus is returned when an HTTP request returns a status code other than 200
	ErrUnexpectedStatus = errors.New("unexpected HTTP status")
	// ErrGraphQLMessage is returned when an error message exists in the response
	ErrGraphQLMessage = errors.New("GraphQL error")
	// ErrUnknownGraphQLError is returned when an GraphQL error occurs but no specific message is available
	ErrUnknownGraphQLError = errors.New("an unknown error occurred")
	// ErrMissingRoot is returned when the GraphQL response is missing the expected root field
	ErrMissingRoot = errors.New("missing root in response")
	// ErrMissingEdges is returned when the response is missing the edges field expected for the data
	ErrMissingEdges = errors.New("missing edges in response")
	// ErrMissingPageInfo is returned when the response is missing pagination data
	ErrMissingPageInfo = errors.New("missing pageInfo in response")
	// ErrMissingHasNextPage is returned when pagination data is missing the hasNextPage field
	ErrMissingHasNextPage = errors.New("missing hasNextPage in pageInfo")
	// ErrMissingEndCursor is returned when pagination data is missing the endCursor field needed for pagination
	ErrMissingEndCursor = errors.New("missing endCursor in pageInfo")
)

// ExportWorkerConfig configuration for the export content worker
type ExportWorkerConfig struct {
	// embed OpenlaneConfig to reuse validation and client creation logic
	OpenlaneConfig `koanf:",squash" jsonschema:"description=the openlane API configuration for exporting"`
	// Enabled indicates if this job is enabled in the server
	Enabled bool `koanf:"enabled" json:"enabled" jsonschema:"required description=whether the export worker is enabled"`
	// MaxZipSize is the maximum allowed size in bytes for a zip archive export, defaults to 500 MB
	MaxZipSize int64 `koanf:"maxzipsize" json:"maxzipsize" jsonschema:"description=the maximum allowed size in bytes for a zip archive export" default:"500000000"`
	// CloudflareAccountID is the cloudflare account id used for browser rendering PDF generation
	CloudflareAccountID string `koanf:"cloudflareaccountid" json:"cloudflareaccountid" jsonschema:"description=the cloudflare account id used for browser rendering pdf generation"`
	// CloudflareAPIKey is the cloudflare api key used for browser rendering PDF generation
	CloudflareAPIKey string `koanf:"cloudflareapikey" json:"cloudflareapikey" jsonschema:"description=the cloudflare api key used for browser rendering pdf generation" sensitive:"true"`
	// MaxSnoozes is the maximum number of times to snooze the job before giving up
	MaxSnoozes int `koanf:"maxsnoozes" json:"maxsnoozes" jsonschema:"default=30 description=maximum number of snoozes before giving up"`
	// SnoozeDuration is the duration to snooze between PDF render retries
	SnoozeDuration time.Duration `koanf:"snoozeduration" json:"snoozeduration" jsonschema:"default=10s description=duration to snooze between pdf render retries"`
}

// ExportContentWorker exports the content into csv and makes it downloadable
type ExportContentWorker struct {
	river.WorkerDefaults[jobspec.ExportContentArgs]

	Config ExportWorkerConfig `koanf:"config" json:"config" jsonschema:"description=the configuration for exporting"`

	olClient    openlane.GraphClient
	requester   *httpsling.Requester
	pdfRenderer PDFRenderer
}

// SetDefaultsIfUnset sets default values for export worker config fields if they are not already set
func (c *ExportWorkerConfig) SetDefaultsIfUnset(input OpenlaneConfig) error {
	if c.MaxSnoozes == 0 {
		c.MaxSnoozes = defaultExportMaxSnoozes
	}

	if c.SnoozeDuration == 0 {
		c.SnoozeDuration = defaultExportSnoozeDuration
	}

	return c.OpenlaneConfig.SetDefaultsIfUnset(input)
}

// PDFRenderer renders a complete HTML document into PDF bytes
type PDFRenderer interface {
	HTMLToPDF(ctx context.Context, html string) ([]byte, error)
}

// WithOpenlaneClient sets the Openlane client for the worker
// and returns the worker for method chaining
func (w *ExportContentWorker) WithOpenlaneClient(cl openlane.GraphClient) *ExportContentWorker {
	w.olClient = cl
	return w
}

// WithRequester sets the httpsling requester to use for HTTP requests
func (w *ExportContentWorker) WithRequester(requester *httpsling.Requester) *ExportContentWorker {
	w.requester = requester
	return w
}

// WithPDFRenderer sets the renderer used to convert HTML documents into PDFs
func (w *ExportContentWorker) WithPDFRenderer(renderer PDFRenderer) *ExportContentWorker {
	w.pdfRenderer = renderer
	return w
}

// Work satisfies the river.Worker interface for the export content worker
// it creates a csv, uploads it and associates it with the export
func (w *ExportContentWorker) Work(ctx context.Context, job *river.Job[jobspec.ExportContentArgs]) error {
	if job.Args.ExportID == "" {
		return newMissingRequiredArg("export_id", jobspec.ExportContentArgs{}.Kind())
	}

	// Exports must be done on behalf of a user in an organization
	if job.Args.OrganizationID == "" {
		return newMissingRequiredArg("organization_id", jobspec.ExportContentArgs{}.Kind())
	}

	if job.Args.UserID == "" {
		return newMissingRequiredArg("user_id", jobspec.ExportContentArgs{}.Kind())
	}

	if w.olClient == nil {
		cl, err := w.Config.getOpenlaneClient()
		if err != nil {
			return err
		}

		w.olClient = cl
	}

	if w.requester == nil {
		var err error

		w.requester, err = httpsling.New(
			httpsling.URL(w.Config.OpenlaneAPIHost),
			httpsling.BearerAuth(w.Config.OpenlaneAPIToken),
			httpsling.Header(httpsling.HeaderContentType, httpsling.ContentTypeJSONUTF8),
			httpsling.Header(httpsling.HeaderAccept, "application/graphql-response+json"),
		)
		if err != nil {
			return err
		}
	}

	export, err := w.olClient.GetExportByID(ctx, job.Args.ExportID)
	if err != nil {
		log.Error().Err(err).Str("export_id", job.Args.ExportID).Msg("failed to get export")
		return w.updateExportStatus(ctx, job.Args.ExportID, enums.ExportStatusFailed, err)
	}

	if job.Args.Mode == enums.ExportModeFlat || job.Args.Mode == enums.ExportModeFolder {
		return w.exportEvidenceFiles(ctx, job, export)
	}

	var filterMap map[string]any

	filtersPtr := export.Export.Filters
	if filtersPtr != nil {
		filters := *filtersPtr
		if filters != "" {
			if err := json.Unmarshal([]byte(filters), &filterMap); err != nil {
				log.Error().Err(err).Msg("failed to parse filters")
				return w.updateExportStatus(ctx, job.Args.ExportID, enums.ExportStatusFailed, err)
			}
		}
	}

	where := make(map[string]any)

	maps.Copy(where, filterMap)

	hasWhere := len(where) > 0

	exportType := strcase.LowerCamelCase(export.Export.ExportType.String())
	rootQuery := pluralize.NewClient().Plural(exportType)

	fields := export.Export.Fields

	// Specific fields wanted for non-csv export
	if export.Export.Format != enums.ExportFormatCsv {
		fields = pdfFieldsWanted
	}

	if export.Export.ExportType == enums.ExportTypeInternalPolicy {
		fields = append(fields, additionalPolicyFields...)
	}

	query := w.buildGraphQLQuery(rootQuery, exportType, fields, hasWhere)

	var (
		allNodes []map[string]any
		after    *string
	)

	for {
		nodes, hasNext, nextCursor, err := w.fetchPage(ctx, query, rootQuery, after, where, job.Args)
		if err != nil {
			return w.updateExportStatus(ctx, job.Args.ExportID, enums.ExportStatusFailed, err)
		}

		allNodes = append(allNodes, nodes...)

		if !hasNext {
			break
		}

		after = &nextCursor
	}

	if len(allNodes) == 0 {
		log.Info().Msg("no data found for export")
		return w.updateExportStatus(ctx, job.Args.ExportID, enums.ExportStatusNodata, nil)
	}

	// Determine the output format from the export record
	exportFormat := export.Export.Format
	timestamp := time.Now().Format("20060102_150405")

	var (
		fileData    []byte
		filename    string
		contentType string
	)

	switch exportFormat {
	// TODO: determine support for implementation, for now leaving off implementation
	case enums.ExportFormatDocx, enums.ExportFormatMarkDown:
		return ErrUnsupportedExportType
	case enums.ExportFormatPdf:
		// PDF generation relies on the Cloudflare browser rendering API, when it is
		// not configured the PDF format is effectively unsupported
		if w.Config.CloudflareAccountID == "" || w.Config.CloudflareAPIKey == "" {
			return ErrUnsupportedExportType
		}

		pdfs, err := w.generatePolicyPDFs(ctx, allNodes, rootQuery)
		if err != nil {
			if ok, newErr := w.isRetryable(job, err); ok {
				return newErr
			}

			log.Error().Err(err).Msg("failed to generate PDFs")

			return w.updateExportStatus(ctx, job.Args.ExportID, enums.ExportStatusFailed, err)
		}

		// a single document is uploaded as a standalone PDF, multiple documents are
		// bundled into a zip archive so each document keeps its own file
		if len(pdfs) == 1 {
			fileData = pdfs[0].data
			filename = pdfs[0].name
			contentType = "application/pdf"
		} else {
			fileData, err = w.buildPDFZip(pdfs)
			if err != nil {
				log.Error().Err(err).Msg("failed to build PDF zip archive")

				return w.updateExportStatus(ctx, job.Args.ExportID, enums.ExportStatusFailed, err)
			}

			filename = fmt.Sprintf("%s_export_%s.zip", strcase.SnakeCase(rootQuery), timestamp)

			contentType = "application/zip"
		}
	default:
		// Default to CSV (includes ExportFormatCsv and any unknown format)
		fileData, err = w.marshalToCSV(allNodes)
		if err != nil {
			log.Error().Err(err).Msg("failed to marshal to CSV")
			return w.updateExportStatus(ctx, job.Args.ExportID, enums.ExportStatusFailed, err)
		}

		filename = fmt.Sprintf("%s_export_%s.csv", rootQuery, timestamp)
		contentType = "text/csv"
	}

	reader := bytes.NewReader(fileData)

	upload := &graphql.Upload{
		File:        reader,
		Filename:    filename,
		Size:        int64(len(fileData)),
		ContentType: contentType,
	}

	updateInput := graphclient.UpdateExportInput{
		Status: &enums.ExportStatusReady,
	}

	// impersonate the requesting user so the uploaded file is created in their
	// organization and is linked to the export
	impersonation := goclient.WithImpersonationInterceptor(job.Args.UserID, job.Args.OrganizationID)

	_, err = w.olClient.UpdateExport(ctx, job.Args.ExportID, updateInput, []*graphql.Upload{upload}, impersonation)
	if err != nil {
		log.Error().Err(err).Msg("failed to update export with file")
		return w.updateExportStatus(ctx, job.Args.ExportID, enums.ExportStatusFailed, err)
	}

	return nil
}

func (w *ExportContentWorker) isRetryable(job *river.Job[jobspec.ExportContentArgs], err error) (bool, error) {
	var errCheck *cloudflare.Error

	// 422 could be for a plethora of reasons from timeouts to browser/page crashes
	if !errors.As(err, &errCheck) || errCheck.StatusCode != http.StatusUnprocessableEntity {
		return false, nil
	}

	snoozes := struct{ Snoozes int }{}
	if err := json.Unmarshal(job.Metadata, &snoozes); err != nil {
		snoozes.Snoozes = 0
	}

	if snoozes.Snoozes >= w.Config.MaxSnoozes {
		log.Warn().Msg("max snoozes reached, giving up")
		return true, ErrMaxSnoozesReached
	}

	log.Info().
		Int("snoozes", snoozes.Snoozes).
		Msg("cloudflare browser rendering returned 422, snoozing")

	return true, river.JobSnooze(w.Config.SnoozeDuration)
}

// buildGraphQLQuery generates a query that can be used to paginate and fetch all data
//
// e.g
//
// query GetControls(
//
//	$first: Int
//	$last: Int
//	$after: Cursor
//	$before: Cursor
//	$where: ControlWhereInput
//
//	) {
//	  controls(
//	    first: $first
//	    last: $last
//	    after: $after
//	    before: $before
//	    where: $where
//	    orderBy: $orderBy
//	  ) {
//	    totalCount
//	    pageInfo {
//	      startCursor
//	      endCursor
//	      hasPreviousPage
//	      hasNextPage
//	    }
//	    edges {
//	      node {
//	        id
//	      }
//	    }
//	  }
//	}
func (w *ExportContentWorker) buildGraphQLQuery(root, singular string, fields []string, hasWhere bool) string {
	fieldStr := CreateFieldsStr(fields)

	var (
		varStr string
		argStr string
	)

	if hasWhere {
		whereInputType := strcase.UpperCamelCase(singular) + "WhereInput"
		varStr = fmt.Sprintf(", $where: %s!", whereInputType)
		argStr = ", where: $where"
	}

	return fmt.Sprintf(`query ($first: Int, $after: Cursor%s) {
  %s(first: $first, after: $after%s) {
    totalCount
    pageInfo {
      hasNextPage
      endCursor
    }
    edges {
      node {
        %s
      }
    }
  }
}`, varStr, root, argStr, fieldStr)
}

// CreateFieldsStr creates a graphql fields string from a list of fields
// supporting nested fields using dot notation
//
// e.g:
//
//	["id", "name", "owner.name", "tasks.title"]
//
// becomes:
//
//	id
//	name
//	owner { name }
//	tasks { edges { node { title } }
func CreateFieldsStr(fields []string) string {
	if len(fields) == 0 {
		fields = []string{"id"}
	}

	fieldStr := ""

	for _, f := range fields {
		if strings.Contains(f, ".") {
			// split and create nested fields
			parts := strings.Split(f, ".")

			fieldStr += parts[0] + "\n        "

			numClosingBraces := 0

			for i, p := range parts[1:] {
				// check the parent to see if it is plural, which will be the same index as the loop
				// because we are looping over parts[1:]
				isParentPlural := pluralize.NewClient().IsPlural(parts[i])

				if isParentPlural {
					fieldStr += "{ edges { node { " + p + " "
					numClosingBraces += 3
				} else {
					fieldStr += "{  " + p + " "
					numClosingBraces++
				}

				if i == len(parts[1:])-1 {
					for range numClosingBraces {
						fieldStr += " } "
					}
				}
			}

			// Add a newline after the closing braces
			fieldStr += "\n        "
		} else {
			fieldStr += f + "\n        "
		}
	}

	return fieldStr
}

// extractErrors converts a slice of errors from the request into one
func extractErrors(errs []any) error {
	if len(errs) == 0 {
		return nil
	}

	var errMsgs []error

	for _, e := range errs {
		if msg, ok := e.(map[string]any); ok {
			if m, ok := msg["message"].(string); ok {
				errMsgs = append(errMsgs, fmt.Errorf("%w: %s", ErrGraphQLMessage, m))
			}
		}
	}

	if len(errMsgs) > 0 {
		return errors.Join(errMsgs...)
	}

	return ErrUnknownGraphQLError
}

// executeGraphQLQuery performs a GraphQL query against the Openlane API
func (w *ExportContentWorker) executeGraphQLQuery(ctx context.Context, query string, variables map[string]any, jobArgs jobspec.ExportContentArgs) (map[string]any, error) {
	body := map[string]any{"query": query}
	if len(variables) > 0 {
		body["variables"] = variables
	}

	// Prepare request options
	opts := []httpsling.Option{
		httpsling.Post("/query"),
		httpsling.Body(body),
	}

	// Add user context headers if provided (for system admin operations)
	if jobArgs.UserID != "" {
		opts = append(opts, httpsling.Header(auth.UserIDHeader, jobArgs.UserID))
	}

	if jobArgs.OrganizationID != "" {
		opts = append(opts, httpsling.Header(auth.OrganizationIDHeader, jobArgs.OrganizationID))
	}

	var result struct {
		Data   map[string]any `json:"data"`
		Errors []any          `json:"errors"`
	}

	resp, err := w.requester.ReceiveWithContext(ctx, &result, opts...)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck

	if len(result.Errors) > 0 {
		return nil, extractErrors(result.Errors)
	}

	return result.Data, nil
}

// marshalToCSV converts a list of nodes (maps) into CSV format
// and stripes any HTML from the values
func (w *ExportContentWorker) marshalToCSV(nodes []map[string]any) ([]byte, error) {
	if len(nodes) == 0 {
		return nil, nil
	}

	// 1) Flatten all nodes
	flatNodes := make([]map[string]any, len(nodes))
	headerSet := make(map[string]struct{})

	for i, n := range nodes {
		flat := make(map[string]any)
		render.Flatten("", n, flat)
		flatNodes[i] = flat

		for k := range flat {
			headerSet[k] = struct{}{}
		}
	}

	if len(headerSet) == 0 {
		return nil, nil
	}

	// 2) Build sorted, stable headers
	headers := make([]string, 0, len(headerSet))
	for k := range headerSet {
		headers = append(headers, k)
	}

	sort.Strings(headers)

	// 3) Write CSV
	var buf bytes.Buffer

	wr := csv.NewWriter(&buf)
	writer := gocsv.NewSafeCSVWriter(wr)

	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	for _, node := range flatNodes {
		row := make([]string, len(headers))
		for i, h := range headers {
			val, ok := node[h]
			if !ok || val == nil {
				row[i] = ""
				continue
			}

			row[i] = render.CleanHTML(val)
		}

		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (w *ExportContentWorker) updateExportStatus(ctx context.Context, exportID string, status enums.ExportStatus, err error) error {
	updateInput := graphclient.UpdateExportInput{
		Status: &status,
	}

	if status == enums.ExportStatusFailed && err != nil {
		msg := err.Error()
		updateInput.ErrorMessage = &msg
	}

	_, err = w.olClient.UpdateExport(ctx, exportID, updateInput, nil)
	if err != nil {
		log.Error().Err(err).
			Str("export_id", exportID).
			Str("status", string(status)).
			Msg("failed to update export status")

		return err
	}

	log.Info().Str("export_id", exportID).Msg("export status updated")

	return nil
}

func (w *ExportContentWorker) fetchPage(ctx context.Context, query, rootQuery string, after *string, where map[string]any, jobArgs jobspec.ExportContentArgs) ([]map[string]any, bool, string, error) {
	vars := map[string]any{"first": defaultPageSize}
	if after != nil {
		vars["after"] = *after
	}

	if len(where) > 0 {
		vars["where"] = where
	}

	data, err := w.executeGraphQLQuery(ctx, query, vars, jobArgs)
	if err != nil {
		log.Error().Err(err).Msg("failed to execute GraphQL query")
		return nil, false, "", err
	}

	rootData, ok := data[rootQuery].(map[string]any)
	if !ok {
		log.Error().Msg("missing root in response")
		return nil, false, "", ErrMissingRoot
	}

	edges, ok := rootData["edges"].([]any)
	if !ok {
		log.Error().Msg("missing edges in response")
		return nil, false, "", ErrMissingEdges
	}

	nodes := make([]map[string]any, 0, len(edges))
	for _, edge := range edges {
		edgeMap, ok := edge.(map[string]any)
		if !ok {
			continue
		}

		node, ok := edgeMap["node"].(map[string]any)
		if ok {
			nodes = append(nodes, node)
		}
	}

	pageInfo, ok := rootData["pageInfo"].(map[string]any)
	if !ok {
		log.Error().Msg("missing pageInfo in response")
		return nil, false, "", ErrMissingPageInfo
	}

	hasNext, ok := pageInfo["hasNextPage"].(bool)
	if !ok {
		log.Error().Msg("missing hasNextPage in pageInfo")
		return nil, false, "", ErrMissingHasNextPage
	}

	var endCursor string
	if hasNext {
		endCursor, ok = pageInfo["endCursor"].(string)
		if !ok {
			log.Error().Msg("missing endCursor in pageInfo")
			return nil, false, "", ErrMissingEndCursor
		}
	}

	return nodes, hasNext, endCursor, nil
}

// pdfFile holds a single rendered PDF document along with the filename it should
// take inside the export (or the zip archive when multiple documents are exported)
type pdfFile struct {
	name string
	data []byte
}

// generatePolicyPDFs renders each node into its own PDF document, named after the document's name field
func (w *ExportContentWorker) generatePolicyPDFs(ctx context.Context, nodes []map[string]any, rootQuery string) ([]pdfFile, error) {
	renderer := w.pdfRenderer
	if renderer == nil {
		renderer = &render.PDFClient{
			AccountID: w.Config.CloudflareAccountID,
			APIToken:  w.Config.CloudflareAPIKey,
		}
	}

	usedNames := make(map[string]int)
	pdfs := make([]pdfFile, 0, len(nodes))

	for _, node := range nodes {
		doc := render.WrapDocument(strings.Join(render.ExtractDetailsStrings([]map[string]any{node}), "\n"))

		data, err := renderer.HTMLToPDF(ctx, doc)
		if err != nil {
			return nil, err
		}

		name := resolveCollision(usedNames, pdfFileName(node, rootQuery))
		usedNames[name]++

		pdfs = append(pdfs, pdfFile{name: name, data: data})
	}

	return pdfs, nil
}

// pdfFileName derives a safe pdf filename from a node's name field, falling back to
// the rootQuery when the node has no usable name
func pdfFileName(node map[string]any, fallback string) string {
	name := fallback

	if raw, ok := node["name"]; ok && raw != nil {
		if str := strings.TrimSpace(fmt.Sprint(raw)); str != "" {
			name = str
		}
	}

	return strcase.SnakeCase(name) + ".pdf"
}

// buildPDFZip bundles the rendered PDFs into a single zip archive, enforcing the
// configured maximum archive size
func (w *ExportContentWorker) buildPDFZip(pdfs []pdfFile) ([]byte, error) {
	var buf bytes.Buffer

	zw := zip.NewWriter(&buf)

	currentTime := time.Now()

	for _, pdf := range pdfs {
		if w.Config.MaxZipSize > 0 && int64(buf.Len()+len(pdf.data)) > w.Config.MaxZipSize {
			return nil, ErrZipTooLarge
		}

		f, err := zw.CreateHeader(&zip.FileHeader{
			Name:     pdf.name,
			Method:   zip.Deflate,
			Modified: currentTime,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create file in zip: %w", err)
		}

		if _, err := f.Write(pdf.data); err != nil {
			return nil, fmt.Errorf("failed to write file to zip: %w", err)
		}
	}

	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("failed to close zip writer: %w", err)
	}

	return buf.Bytes(), nil
}
