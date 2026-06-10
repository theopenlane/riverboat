package render

import (
	"context"
	_ "embed"
	"encoding/base64"
	"fmt"
	"io"
	"time"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/browser_rendering"
	"github.com/cloudflare/cloudflare-go/v6/option"
	"github.com/cloudflare/cloudflare-go/v6/shared"
)

// defaultPDFTimeout is the timeout used when rendering a PDF via the Cloudflare
// browser rendering API
const defaultPDFTimeout = 60 * time.Second

// pdfMargin is the page margin applied to each side of a generated PDF
const pdfMargin = "0.75in"

// pdfFooterMargin is the bottom margin, larger than the others to leave room for the footer
const pdfFooterMargin = "0.9in"

// openlaneLogo is the logo rendered in the PDF footer, replace
// assets/openlane-logo.png to change it
//
//go:embed assets/openlane-logo.png
var openlaneLogo []byte

// footerHTML builds the Chromium footer template containing the Openlane logo and label.
// The logo is embedded as a data URI so it renders without an external request.
func footerHTML() string {
	logoSrc := "data:image/png;base64," + base64.StdEncoding.EncodeToString(openlaneLogo)

	return fmt.Sprintf(
		`<div style="width:100%%;font-size:8px;color:#9aa0a6;text-align:center;padding:0;-webkit-print-color-adjust:exact;">`+
			`<span>Exported from</span>`+
			`<img src="%s" style="height:12px;vertical-align:middle;margin-left:5px;"/></div>`,
		logoSrc,
	)
}

// PDFClient renders HTML documents into PDFs using the Cloudflare browser rendering API
type PDFClient struct {
	// AccountID is the cloudflare account id
	AccountID string
	// APIToken is the cloudflare api token used for authentication
	APIToken string
}

// HTMLToPDF renders a complete HTML document into PDF bytes
func (c *PDFClient) HTMLToPDF(ctx context.Context, html string) ([]byte, error) {
	client := cloudflare.NewClient(
		option.WithAPIToken(c.APIToken),
		option.WithRequestTimeout(defaultPDFTimeout),
	)

	resp, err := client.BrowserRendering.PDF.New(ctx, browser_rendering.PDFNewParams{
		AccountID: cloudflare.F(c.AccountID),
		Body: browser_rendering.PDFNewParamsBodyObject{
			HTML: cloudflare.F(html),
			PDFOptions: cloudflare.F(browser_rendering.PDFNewParamsBodyObjectPDFOptions{
				Format:              cloudflare.F(browser_rendering.PDFNewParamsBodyObjectPDFOptionsFormatLetter),
				PrintBackground:     cloudflare.F(true),
				DisplayHeaderFooter: cloudflare.F(true),
				HeaderTemplate:      cloudflare.F("<span></span>"),
				FooterTemplate:      cloudflare.F(footerHTML()),
				Margin: cloudflare.F(browser_rendering.PDFNewParamsBodyObjectPDFOptionsMargin{
					Top:    cloudflare.F[browser_rendering.PDFNewParamsBodyObjectPDFOptionsMarginTopUnion](shared.UnionString(pdfMargin)),
					Bottom: cloudflare.F[browser_rendering.PDFNewParamsBodyObjectPDFOptionsMarginBottomUnion](shared.UnionString(pdfFooterMargin)),
					Left:   cloudflare.F[browser_rendering.PDFNewParamsBodyObjectPDFOptionsMarginLeftUnion](shared.UnionString(pdfMargin)),
					Right:  cloudflare.F[browser_rendering.PDFNewParamsBodyObjectPDFOptionsMarginRightUnion](shared.UnionString(pdfMargin)),
				}),
			}),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call cloudflare browser rendering: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read pdf response: %w", err)
	}

	return data, nil
}
