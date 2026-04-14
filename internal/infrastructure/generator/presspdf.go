package generator

import (
	"context"

	"github.com/akkaraponph/presspdf"

	"pdfapi/internal/domain/generator"
	"pdfapi/internal/domain/shared"
)

type service struct{}

// New creates a new generator service backed by presspdf.
func New() generator.Service {
	return &service{}
}

func (s *service) HTMLToPDF(_ context.Context, params generator.HTMLToPDFParams) ([]byte, error) {
	doc := presspdf.New()
	applyDocConfig(doc, params.PageConfig, params.Meta)

	page := doc.AddPage(resolvePageSize(params.PageConfig))
	page.HTML(params.HTML)

	data, err := doc.Bytes()
	if err != nil {
		return nil, shared.ErrProcessingFailed{Operation: "html-to-pdf", Cause: err}
	}
	return data, nil
}

func (s *service) MarkdownToPDF(_ context.Context, params generator.MarkdownToPDFParams) ([]byte, error) {
	doc := presspdf.New()
	applyDocConfig(doc, params.PageConfig, params.Meta)

	page := doc.AddPage(resolvePageSize(params.PageConfig))

	var opts []presspdf.MarkdownOption
	if params.Bookmarks {
		opts = append(opts, presspdf.WithBookmarks())
	}
	page.Markdown(params.Markdown, opts...)

	data, err := doc.Bytes()
	if err != nil {
		return nil, shared.ErrProcessingFailed{Operation: "markdown-to-pdf", Cause: err}
	}
	return data, nil
}

func applyDocConfig(doc *presspdf.Document, pc shared.PageConfig, meta generator.DocumentMeta) {
	m := pc.Margins
	if m.Left == 0 && m.Top == 0 && m.Right == 0 {
		m = shared.DefaultMargins()
	}
	doc.SetMargins(m.Left, m.Top, m.Right)
	doc.SetAutoPageBreak(true, m.Bottom)
	doc.SetFont("helvetica", "", 12)

	if meta.Title != "" {
		doc.SetTitle(meta.Title)
	}
	if meta.Author != "" {
		doc.SetAuthor(meta.Author)
	}
	if meta.Subject != "" {
		doc.SetSubject(meta.Subject)
	}
	if meta.Creator != "" {
		doc.SetCreator(meta.Creator)
	}
	if meta.Keywords != "" {
		doc.SetKeywords(meta.Keywords)
	}
}

func resolvePageSize(pc shared.PageConfig) presspdf.PageSize {
	var size presspdf.PageSize
	switch pc.Size {
	case shared.PageSizeA3:
		size = presspdf.A3
	case shared.PageSizeA5:
		size = presspdf.A5
	case shared.PageSizeLetter:
		size = presspdf.Letter
	case shared.PageSizeLegal:
		size = presspdf.Legal
	default:
		size = presspdf.A4
	}
	if pc.Orientation == shared.Landscape {
		size = size.Landscape()
	}
	return size
}
