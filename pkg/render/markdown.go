package render

import (
	"fmt"
	"strings"

	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/strikethrough"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/table"
	"gopkg.in/yaml.v3"
)

func renderMarkdown(flat map[string]any, details string, excludeMetadata bool) (string, error) {
	if htmlTagRegex.MatchString(details) {
		c := converter.NewConverter(converter.WithPlugins(
			base.NewBasePlugin(),
			commonmark.NewCommonmarkPlugin(),
			strikethrough.NewStrikethroughPlugin(),
			table.NewTablePlugin(),
		))

		parsedMD, err := c.ConvertString(CleanHTML(DetailsToHTML(details)))
		if err != nil {
			return "", err
		}

		details = parsedMD
	}

	var buf strings.Builder

	if !excludeMetadata {
		metadata := make(map[string]any)

		for _, field := range documentMetadataKV {
			if value := flat[field.key]; value != nil {
				metadata[field.key] = value
			}
		}

		frontmatter, err := yaml.Marshal(metadata)
		if err != nil {
			return "", err
		}

		buf.WriteString("---\n")
		buf.Write(frontmatter)
		buf.WriteString("---\n\n")
	}

	if name := flat["name"]; excludeMetadata && name != nil {
		fmt.Fprintf(&buf, "# %v\n\n", name)
	}

	buf.WriteString(details)
	buf.WriteString("\n")

	return buf.String(), nil
}
