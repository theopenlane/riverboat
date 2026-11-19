package render

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/rs/zerolog/log"
	"github.com/theopenlane/core/common/enums"
	"github.com/theopenlane/newman/scrubber"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

// ErrUnsupportedDocumentFormat indicates that a format cannot be rendered as a document.
var ErrUnsupportedDocumentFormat = errors.New("unsupported document format")

var documentMetadataKV = []struct {
	key   string
	label string
}{
	{
		key:   "name",
		label: "Name",
	},
	{
		key:   "status",
		label: "Status",
	},
	{
		key:   "revision",
		label: "Version",
	},
	{
		key:   "createdAt",
		label: "Created At",
	},
	{
		key:   "updatedAt",
		label: "Updated At",
	},
}

// stripHTML removes all HTML tags from a value
var scrub = scrubber.NewPolicyScrubber(
	scrubber.WithStyling(),
	scrubber.WithTables(),
	scrubber.WithImages(),
	scrubber.WithDocumentStructure(),
	scrubber.WithAccessibility(),
	scrubber.WithURLSchemes("http", "https", "mailto", "tel"),
	scrubber.WithNoRelativeURLs(),
	scrubber.WithTargetBlankOnLinks(),
	func(p *bluemonday.Policy) {
		p.AllowAttrs("start").OnElements("ol")
	},
)

// markdownConverter renders markdown content into HTML, supporting GitHub flavored
// markdown features such as tables, strikethrough, and autolinks
var markdownConverter = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
	// render single newlines as <br> so author line breaks are preserved
	goldmark.WithRendererOptions(html.WithHardWraps()),
)

// htmlTagRegex detects the presence of an HTML element tag, used to distinguish
// content that is already HTML from content authored as markdown
var htmlTagRegex = regexp.MustCompile(`<\/?[a-zA-Z][a-zA-Z0-9-]*(\s[^<>]*)?/?>`)

// pdfDocumentTemplate wraps HTML fragments in a styled document so the generated PDF
// uses consistent fonts and spacing instead of the browser defaults
const pdfDocumentTemplate = `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<style>
  html { -webkit-print-color-adjust: exact; }
  body {
    font-family: Arial, "Helvetica Neue", Helvetica, sans-serif;
    font-size: 11pt;
    line-height: 1.5;
    color: #1a1a1a;
  }
  h1, h2, h3, h4, h5, h6 { line-height: 1.25; margin: 1.2em 0 0.4em; font-weight: 600; }
  h1 { font-size: 18pt; }
  h2 { font-size: 15pt; }
  h3 { font-size: 13pt; }
  h4 { font-size: 12pt; }
  h5 { font-size: 11pt; }
  h6 { font-size: 10pt; }
  p { margin: 0 0 0.6em; }
  hr { border: none; border-top: 1px solid #d0d0d0; margin: 1em 0; }
  strong { font-weight: 600; }
  table { border-collapse: collapse; width: 100%%; margin: 1em 0; }
  td, th { border: 1px solid #d0d0d0; padding: 4px 8px; text-align: left; }
</style>
</head>
<body>
%s
</body>
</html>`

// CleanHTML strips all HTML from the value and collapses surrounding whitespace
func CleanHTML(v any) string {
	raw := fmt.Sprint(v)

	return scrub.Scrub(raw)
}

// DetailsToHTML returns the content as HTML. Content that already contains HTML markup
// has its text-level line breaks preserved, otherwise it is treated as markdown and converted
func DetailsToHTML(content string) string {
	if htmlTagRegex.MatchString(content) {
		return preserveHTMLLineBreaks(content)
	}

	var buf bytes.Buffer
	if err := markdownConverter.Convert([]byte(content), &buf); err != nil {
		log.Warn().Err(err).Msg("failed to convert markdown details to HTML, falling back to line breaks")
		return strings.ReplaceAll(content, "\n", "<br/>\n")
	}

	return buf.String()
}

// preserveHTMLLineBreaks converts text-level newlines in HTML content into <br/> so
// author line breaks render, while leaving structural newlines between tags untouched
func preserveHTMLLineBreaks(content string) string {
	lines := strings.Split(content, "\n")
	if len(lines) == 1 {
		return content
	}

	var b strings.Builder

	for i, line := range lines {
		b.WriteString(line)

		if i == len(lines)-1 {
			break
		}

		cur := strings.TrimRight(line, " \t\r")
		next := strings.TrimLeft(lines[i+1], " \t\r")

		// keep the newline as-is when it sits next to a tag or a blank line, otherwise
		// treat it as an author line break within text content
		if cur == "" || next == "" || strings.HasSuffix(cur, ">") || strings.HasPrefix(next, "<") {
			b.WriteString("\n")
		} else {
			b.WriteString("<br/>\n")
		}
	}

	return b.String()
}

// WrapDocument wraps HTML fragments in a complete, styled HTML document suitable for
// rendering to PDF
func WrapDocument(body string) string {
	cleaned := CleanHTML(body)

	return fmt.Sprintf(pdfDocumentTemplate, cleaned)
}

// Flatten flattens a nested map into a flat map with dot notation keys
func Flatten(prefix string, v any, out map[string]any) {
	switch val := v.(type) {
	case map[string]any:
		for k, v2 := range val {
			key := k
			if prefix != "" {
				key = prefix + "." + k
			}

			Flatten(key, v2, out)
		}

	case []any:
		for i, v2 := range val {
			key := fmt.Sprintf("%s.%d", prefix, i)
			Flatten(key, v2, out)
		}

	default:
		// leaf value
		if prefix != "" {
			out[prefix] = val
		}
	}
}

// Document extracts string content from a node and formats it into HTML or Markdown
// for document generation. It includes common headers (name, status,
// timestamps) followed by the details content
func Document(node map[string]any, format enums.ExportFormat, excludeMetadata bool) (string, error) {
	flat := make(map[string]any)
	Flatten("", node, flat)
	details := getDocumentDetail(flat)

	switch format {
	case enums.ExportFormatMarkDown:
		return renderMarkdown(flat, details, excludeMetadata)
	case enums.ExportFormatHTML, enums.ExportFormatPdf:
		return WrapDocument(renderHTML(flat, details, excludeMetadata)), nil
	default:
		return "", fmt.Errorf("%w: %s", ErrUnsupportedDocumentFormat, format)
	}
}

func renderHTML(flat map[string]any, details string, excludeMetadata bool) string {
	var buf strings.Builder

	if name := flat["name"]; excludeMetadata && name != nil {
		fmt.Fprintf(&buf, "<h1>%v</h1>\n", name)
	}

	// include header data unless explicitly set to not include
	if !excludeMetadata {
		addedKeys := make(map[string]bool)
		for _, field := range documentMetadataKV {
			if val, ok := flat[field.key]; ok && val != nil && !addedKeys[field.label] {
				fmt.Fprintf(&buf, "<p><strong>%s:</strong> %v</p>\n", field.label, val)
				addedKeys[field.label] = true
			}
		}

		buf.WriteString("<hr/>\n")
	}

	buf.WriteString(DetailsToHTML(details))

	return buf.String()
}

func getDocumentDetail(flat map[string]any) string {
	for _, key := range []string{"details", "externalContents", "liveExternalContents"} {
		if text := flat[key]; text != nil && text != "" {
			return fmt.Sprint(text)
		}
	}

	return "Empty Document"
}
