package generator

import "pdfapi/internal/domain/shared"

// DocumentMeta holds optional PDF metadata.
type DocumentMeta struct {
	Title    string
	Author   string
	Subject  string
	Creator  string
	Keywords string
}

// HTMLToPDFParams holds parameters for rendering HTML to PDF.
type HTMLToPDFParams struct {
	HTML       string
	PageConfig shared.PageConfig
	Meta       DocumentMeta
}

// MarkdownToPDFParams holds parameters for rendering Markdown to PDF.
type MarkdownToPDFParams struct {
	Markdown   string
	PageConfig shared.PageConfig
	Meta       DocumentMeta
	Bookmarks  bool
}
