package jobs

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/gqlgo/gqlgenc/clientv2"
	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/stoewer/go-strcase"
	"github.com/theopenlane/core/common/enums"
	"github.com/theopenlane/core/common/jobspec"
	"github.com/theopenlane/go-client/graphclient"

	goclient "github.com/theopenlane/go-client"
)

const (
	uncategorizedFolder = "_uncategorized"
	metadataFileName    = "metadata.txt"
)

type evidenceFile struct {
	ZipPath      string
	PresignedURL string
	FileID       string
}

type controlInfo struct {
	RefCode            string
	ReferenceFramework string
	ReferenceID        string
	AuditorReferenceID string
}

type evidenceMetadata struct {
	Name  string
	Files []string
}

func (w *ExportContentWorker) exportEvidenceFiles(ctx context.Context,
	job *river.Job[jobspec.ExportContentArgs], export *graphclient.GetExportByID,
) error {
	where, err := parseExportFilters(export)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse export filters")
		return w.updateExportStatus(ctx, job.Args.ExportID, enums.ExportStatusFailed, err)
	}

	impersonationHeader := goclient.WithImpersonationInterceptor(job.Args.UserID, job.Args.OrganizationID)

	evidences, err := w.fetchEvidences(ctx, where, impersonationHeader)
	if err != nil {
		log.Error().Err(err).Msg("failed to fetch evidences for export")
		return w.updateExportStatus(ctx, job.Args.ExportID, enums.ExportStatusFailed, err)
	}

	if len(evidences) == 0 {
		log.Info().Msg("no evidences found for export")
		return w.updateExportStatus(ctx, job.Args.ExportID, enums.ExportStatusNodata, nil)
	}

	// collate unique file IDs across all evidences
	uniqueFileIDs := lo.Uniq(lo.FlatMap(evidences, func(ev *graphclient.GetEvidences_Evidences_Edges_Node, _ int) []string {
		return lo.FilterMap(ev.Files.Edges, func(edge *graphclient.GetEvidences_Evidences_Edges_Node_Files_Edges, _ int) (string, bool) {
			if edge.Node == nil {
				return "", false
			}

			return edge.Node.ID, true
		})
	}))

	fileDetailsMap := make(map[string]*graphclient.GetFileByID_File, len(uniqueFileIDs))

	for _, id := range uniqueFileIDs {
		resp, err := w.olClient.GetFileByID(ctx, id, impersonationHeader)
		if err != nil {
			log.Error().Err(err).Str("file_id", id).Msg("failed to get file details")
			return w.updateExportStatus(ctx, job.Args.ExportID, enums.ExportStatusFailed, err)
		}

		fileDetailsMap[id] = &resp.File
	}

	controls := make(map[string]controlInfo)

	if job.Args.Mode == enums.ExportModeFolder {
		refCodes := collectRefCodes(evidences)

		for _, ref := range refCodes {
			control, err := w.fetchControl(ctx, ref)
			if err != nil {
				log.Warn().Err(err).Str("ref_code", ref).Msg("failed to fetch control info, using defaults")
				controls[ref] = controlInfo{RefCode: ref}

				continue
			}

			controls[ref] = control
		}
	}

	entries, folderEvidences := buildFileEntries(evidences, fileDetailsMap, job.Args.Mode, lo.FromPtr(job.Args.ExportMetadata).KeepFileOriginalName)

	if len(entries) == 0 {
		log.Info().Msg("no evidence files found for export")
		return w.updateExportStatus(ctx, job.Args.ExportID, enums.ExportStatusNodata, nil)
	}

	rootFolder := fmt.Sprintf("evidence_files_export_%s_%s", job.Args.ExportID, time.Now().Format("20060102_150405"))

	data, err := w.createZipArchive(ctx, rootFolder, entries, folderEvidences, controls, job.Args.Mode)
	if err != nil {
		log.Error().Err(err).Msg("failed to build zip archive")
		return w.updateExportStatus(ctx, job.Args.ExportID, enums.ExportStatusFailed, err)
	}

	filename := rootFolder + ".zip"

	reader := bytes.NewReader(data)

	upload := &graphql.Upload{
		File:        reader,
		Filename:    filename,
		Size:        int64(len(data)),
		ContentType: "application/zip",
	}

	updateInput := graphclient.UpdateExportInput{
		Status: &enums.ExportStatusReady,
	}

	_, err = w.olClient.UpdateExport(ctx, job.Args.ExportID, updateInput, []*graphql.Upload{upload}, impersonationHeader)
	if err != nil {
		log.Error().Err(err).Msg("failed to update export with zip file")
		w.updateExportStatus(ctx, job.Args.ExportID, enums.ExportStatusFailed, err) //nolint:errcheck

		return err
	}

	return nil
}

func parseExportFilters(export *graphclient.GetExportByID) (*graphclient.EvidenceWhereInput, error) {
	if export.Export.Filters == nil {
		return nil, nil
	}

	filters := *export.Export.Filters
	if filters == "" {
		return nil, nil
	}

	var where graphclient.EvidenceWhereInput
	if err := json.Unmarshal([]byte(filters), &where); err != nil {
		return nil, err
	}

	return &where, nil
}

func (w *ExportContentWorker) fetchEvidences(ctx context.Context, where *graphclient.EvidenceWhereInput,
	impersonation clientv2.RequestInterceptor,
) ([]*graphclient.GetEvidences_Evidences_Edges_Node, error) {
	var (
		allEvidences []*graphclient.GetEvidences_Evidences_Edges_Node
		after        *string
		pageSize     = defaultPageSize
	)

	for {
		resp, err := w.olClient.GetEvidences(ctx, &pageSize, nil, after, nil, where, nil, impersonation)
		if err != nil {
			return nil, err
		}

		nodes := lo.FilterMap(resp.Evidences.Edges, func(edge *graphclient.GetEvidences_Evidences_Edges, _ int) (*graphclient.GetEvidences_Evidences_Edges_Node, bool) {
			if edge.Node == nil {
				return nil, false
			}

			return edge.Node, true
		})

		allEvidences = append(allEvidences, nodes...)

		if !resp.Evidences.PageInfo.HasNextPage {
			break
		}

		after = resp.Evidences.PageInfo.EndCursor
	}

	return allEvidences, nil
}

func collectRefCodes(evidences []*graphclient.GetEvidences_Evidences_Edges_Node) []string {
	return lo.Uniq(lo.FlatMap(evidences, func(ev *graphclient.GetEvidences_Evidences_Edges_Node, _ int) []string {
		return lo.FilterMap(ev.Controls.Edges, func(edge *graphclient.GetEvidences_Evidences_Edges_Node_Controls_Edges, _ int) (string, bool) {
			if edge.Node == nil || edge.Node.RefCode == "" {
				return "", false
			}

			return edge.Node.RefCode, true
		})
	}))
}

func (w *ExportContentWorker) fetchControl(ctx context.Context, refCode string) (controlInfo, error) {
	var first int64 = 1

	where := &graphclient.ControlWhereInput{
		RefCode: &refCode,
	}

	resp, err := w.olClient.GetControls(ctx, &first, nil, nil, nil, where, nil)
	if err != nil {
		return controlInfo{}, err
	}

	if len(resp.Controls.Edges) == 0 || resp.Controls.Edges[0].Node == nil {
		return controlInfo{RefCode: refCode}, nil
	}

	node := resp.Controls.Edges[0].Node

	return controlInfo{
		RefCode:            node.RefCode,
		ReferenceFramework: lo.FromPtrOr(node.ReferenceFramework, ""),
		ReferenceID:        lo.FromPtrOr(node.ReferenceID, ""),
		AuditorReferenceID: lo.FromPtrOr(node.AuditorReferenceID, ""),
	}, nil
}

func foldersForEvidence(ev *graphclient.GetEvidences_Evidences_Edges_Node, mode enums.ExportMode) []string {
	if mode != enums.ExportModeFolder {
		return []string{""}
	}

	folders := lo.FilterMap(ev.Controls.Edges, func(ctrl *graphclient.GetEvidences_Evidences_Edges_Node_Controls_Edges, _ int) (string, bool) {
		if ctrl.Node == nil || ctrl.Node.RefCode == "" {
			return "", false
		}

		return sanitizeRefCode(ctrl.Node.RefCode), true
	})

	if len(folders) == 0 {
		return []string{uncategorizedFolder}
	}

	return folders
}

func buildFileEntries(evidences []*graphclient.GetEvidences_Evidences_Edges_Node, fileDetails map[string]*graphclient.GetFileByID_File,
	mode enums.ExportMode, retainOriginalFileName bool,
) ([]evidenceFile, map[string][]evidenceMetadata) {
	usedNames := make(map[string]map[string]int)

	// we need metadata.txt for folder exports
	folderEvidences := make(map[string][]evidenceMetadata)

	seenInFlat := make(map[string]struct{})

	// filter out evidences without files
	evidencesWithFiles := lo.Filter(evidences, func(ev *graphclient.GetEvidences_Evidences_Edges_Node, _ int) bool {
		return len(ev.Files.Edges) > 0
	})

	entries := lo.FlatMap(evidencesWithFiles, func(ev *graphclient.GetEvidences_Evidences_Edges_Node, _ int) []evidenceFile {
		folders := foldersForEvidence(ev, mode)
		padWidth := padFileNameWidth(len(ev.Files.Edges))

		var names []string

		evidenceEntries := lo.FlatMap(ev.Files.Edges, func(file *graphclient.GetEvidences_Evidences_Edges_Node_Files_Edges, idx int) []evidenceFile {
			if file.Node == nil {
				return nil
			}

			fileID := file.Node.ID
			presignedURL := lo.FromPtrOr(file.Node.PresignedURL, "")

			detail, ok := fileDetails[fileID]
			if !ok || detail == nil {
				return nil
			}

			ext := getFileExtension(detail)

			baseName := detail.ProvidedFileName
			if !retainOriginalFileName {
				baseName = fmt.Sprintf("%s-%0*d", strcase.KebabCase(ev.Name), padWidth, idx+1)
			}

			fileName := strings.TrimSuffix(baseName, ext) + ext

			return lo.FilterMap(folders, func(folder string, _ int) (evidenceFile, bool) {
				if usedNames[folder] == nil {
					usedNames[folder] = make(map[string]int)
				}

				finalName := resolveCollision(usedNames[folder], fileName)
				usedNames[folder][finalName]++

				zipPath := finalName
				if folder != "" {
					zipPath = folder + "/" + finalName
				}

				// if in FLAT mode, we need to deduplicate by the fileid
				if mode == enums.ExportModeFlat {
					if _, seen := seenInFlat[fileID]; seen {
						return evidenceFile{}, false
					}

					seenInFlat[fileID] = struct{}{}
				}

				names = append(names, finalName)

				return evidenceFile{
					ZipPath:      zipPath,
					PresignedURL: presignedURL,
					FileID:       fileID,
				}, true
			})
		})

		// create metadata.txt for each folder
		if mode == enums.ExportModeFolder && len(names) > 0 {
			meta := evidenceMetadata{Name: ev.Name, Files: names}

			lo.ForEach(folders, func(folder string, _ int) {
				folderEvidences[folder] = append(folderEvidences[folder], meta)
			})
		}

		return evidenceEntries
	})

	return entries, folderEvidences
}

func (w *ExportContentWorker) createZipArchive(
	ctx context.Context, rootFolder string, files []evidenceFile, folders map[string][]evidenceMetadata,
	controls map[string]controlInfo, mode enums.ExportMode,
) ([]byte, error) {
	var buf bytes.Buffer

	zw := zip.NewWriter(&buf)

	currentTime := time.Now()

	if mode == enums.ExportModeFolder {
		for folder, evidences := range folders {
			detail, found := lo.Find(lo.Values(controls), func(cd controlInfo) bool {
				return sanitizeRefCode(cd.RefCode) == folder
			})

			if !found && folder == uncategorizedFolder {
				detail = controlInfo{RefCode: uncategorizedFolder}
			}

			metadataContent := createMetadataContent(detail.RefCode, detail, evidences)

			f, err := zw.CreateHeader(&zip.FileHeader{
				Name:     rootFolder + "/" + folder + "/" + metadataFileName,
				Method:   zip.Deflate,
				Modified: currentTime,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create metadata.txt in zip: %w", err)
			}

			if _, err := f.Write(metadataContent); err != nil {
				return nil, fmt.Errorf("failed to write metadata.txt: %w", err)
			}
		}
	}

	// now we need to download each file so we can add them to the zip
	for _, file := range files {
		if file.PresignedURL == "" {
			continue
		}

		data, err := downloadFile(ctx, file.PresignedURL)
		if err != nil {
			return nil, fmt.Errorf("failed to download file %s: %w", file.FileID, err)
		}

		// we need to make sure the zip file cannot be larger than what is expected
		if int64(buf.Len()+len(data)) > w.Config.MaxZipSize {
			return nil, ErrZipTooLarge
		}

		f, err := zw.CreateHeader(&zip.FileHeader{
			Name:     rootFolder + "/" + file.ZipPath,
			Method:   zip.Deflate,
			Modified: currentTime,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create file in zip: %w", err)
		}

		if _, err := f.Write(data); err != nil {
			return nil, fmt.Errorf("failed to write file to zip: %w", err)
		}
	}

	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("failed to close zip writer: %w", err)
	}

	return buf.Bytes(), nil
}

func getFileExtension(detail *graphclient.GetFileByID_File) string {
	if detail.ProvidedFileExtension == "" {
		return filepath.Ext(detail.ProvidedFileName)
	}

	knownExtension := detail.ProvidedFileExtension
	if !strings.HasPrefix(knownExtension, ".") {
		knownExtension = "." + knownExtension
	}

	return knownExtension
}

// sanitizeRefCode makes sure we can safely use a refCode for a folder name
// e.g. "SOC2::CC1.1" will be converted into "SOC2--CC1.1"
func sanitizeRefCode(refCode string) string {
	return strings.ReplaceAll(refCode, "::", "--")
}

func downloadFile(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create download request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req) //nolint:gosec // URL comes from presigned file URLs
	if err != nil {
		return nil, fmt.Errorf("could not download file: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", ErrUnexpectedStatus, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file body: %w", err)
	}

	return data, nil
}

func createMetadataContent(refCode string, control controlInfo, evidences []evidenceMetadata) []byte {
	var b strings.Builder

	fmt.Fprintf(&b, "RefCode: %s\n", refCode)

	if control.ReferenceFramework != "" {
		fmt.Fprintf(&b, "ReferenceFramework: %s\n", control.ReferenceFramework)
	}

	if control.ReferenceID != "" {
		fmt.Fprintf(&b, "ReferenceID: %s\n", control.ReferenceID)
	}

	if control.AuditorReferenceID != "" {
		fmt.Fprintf(&b, "AuditorReferenceID: %s\n", control.AuditorReferenceID)
	}

	b.WriteString("\nEvidence:\n")

	lo.ForEach(evidences, func(ev evidenceMetadata, _ int) {
		fmt.Fprintf(&b, "  - Name: %s\n", ev.Name)
		fmt.Fprintf(&b, "    Files: %s\n", strings.Join(ev.Files, ", "))
	})

	return []byte(b.String())
}

//nolint:mnd // we want to have files like file-001, file-002 e.g
func padFileNameWidth(count int) int {
	switch {
	case count <= 99:
		return 2
	case count <= 999:
		return 3
	default:
		return 4
	}
}

func resolveCollision(used map[string]int, filename string) string {
	if _, exists := used[filename]; !exists {
		return filename
	}

	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)

	for i := 2; ; i++ {
		if candidate := fmt.Sprintf("%s-%d%s", base, i, ext); used[candidate] == 0 {
			return candidate
		}
	}
}
