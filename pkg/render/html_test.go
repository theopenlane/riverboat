package render_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/theopenlane/riverboat/pkg/render"
)

func TestDetailsToHTML(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		content     string
		contains    []string
		notContains []string
		exact       string
	}{
		{
			name:    "markdown is converted to html",
			content: "# Policy Title\n\nSome **bold** text\n\n- first\n- second",
			contains: []string{
				"<h1>Policy Title</h1>",
				"<strong>bold</strong>",
				"<li>first</li>",
				"<li>second</li>",
			},
		},
		{
			name:    "single newlines become line breaks",
			content: "line one\nline two",
			contains: []string{
				"line one<br>",
				"line two",
			},
		},
		{
			name:    "github flavored markdown tables are supported",
			content: "| a | b |\n| - | - |\n| 1 | 2 |",
			contains: []string{
				"<table>",
				"<td>1</td>",
			},
		},
		{
			name:    "existing html is returned unchanged",
			content: "<p>already <strong>html</strong></p>",
			exact:   "<p>already <strong>html</strong></p>",
		},
		{
			name:    "html fragment with attributes is treated as html",
			content: `<div class="note">hello</div>`,
			exact:   `<div class="note">hello</div>`,
		},
		{
			name:    "text level newlines in html become line breaks",
			content: "<p>line one\nline two</p>",
			contains: []string{
				"line one<br/>",
				"line two",
			},
		},
		{
			name:    "structural newlines between html tags are preserved",
			content: "<p>one</p>\n<p>two</p>",
			exact:   "<p>one</p>\n<p>two</p>",
		},
		{
			name:    "plain text with angle brackets is treated as markdown",
			content: "if a < b and c > d then ok",
			contains: []string{
				"<p>if a &lt; b and c &gt; d then ok</p>",
			},
			notContains: []string{
				"if a < b",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := render.DetailsToHTML(tc.content)

			if tc.exact != "" {
				assert.Equal(t, tc.exact, got)
			}

			for _, want := range tc.contains {
				assert.Contains(t, got, want)
			}

			for _, notWant := range tc.notContains {
				assert.NotContains(t, got, notWant)
			}
		})
	}
}
