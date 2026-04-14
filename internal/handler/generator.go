package handler

import (
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http/response"

	"pdfapi/internal/domain/generator"
	"pdfapi/internal/domain/shared"
)

// GeneratorHandler handles HTML/Markdown→PDF HTTP requests.
type GeneratorHandler struct {
	svc generator.Service
}

// NewGeneratorHandler creates a new GeneratorHandler.
func NewGeneratorHandler(svc generator.Service) *GeneratorHandler {
	return &GeneratorHandler{svc: svc}
}

type htmlToPDFRequest struct {
	HTML         string  `json:"html"`
	PageSize     string  `json:"page_size"`
	Landscape    bool    `json:"landscape"`
	MarginLeft   float64 `json:"margin_left"`
	MarginTop    float64 `json:"margin_top"`
	MarginRight  float64 `json:"margin_right"`
	MarginBottom float64 `json:"margin_bottom"`
	Title        string  `json:"title"`
	Author       string  `json:"author"`
	Subject      string  `json:"subject"`
}

// HTMLToPDF renders HTML content as a PDF.
func (h *GeneratorHandler) HTMLToPDF(ctx *gofr.Context) (any, error) {
	var req htmlToPDFRequest
	if err := ctx.Bind(&req); err != nil {
		return nil, shared.ErrInvalidInput{Message: "failed to parse request: " + err.Error()}
	}
	if req.HTML == "" {
		return nil, shared.ErrInvalidInput{Message: "html field is required"}
	}

	pdfData, err := h.svc.HTMLToPDF(ctx, generator.HTMLToPDFParams{
		HTML:       req.HTML,
		PageConfig: buildPageConfig(req.PageSize, req.Landscape, req.MarginLeft, req.MarginTop, req.MarginRight, req.MarginBottom),
		Meta: generator.DocumentMeta{
			Title:   req.Title,
			Author:  req.Author,
			Subject: req.Subject,
		},
	})
	if err != nil {
		return nil, err
	}

	return response.File{
		Content:     pdfData,
		ContentType: "application/pdf",
	}, nil
}

type markdownToPDFRequest struct {
	Markdown     string  `json:"markdown"`
	PageSize     string  `json:"page_size"`
	Landscape    bool    `json:"landscape"`
	Bookmarks    bool    `json:"bookmarks"`
	MarginLeft   float64 `json:"margin_left"`
	MarginTop    float64 `json:"margin_top"`
	MarginRight  float64 `json:"margin_right"`
	MarginBottom float64 `json:"margin_bottom"`
	Title        string  `json:"title"`
	Author       string  `json:"author"`
	Subject      string  `json:"subject"`
}

// MarkdownToPDF renders Markdown content as a PDF.
func (h *GeneratorHandler) MarkdownToPDF(ctx *gofr.Context) (any, error) {
	var req markdownToPDFRequest
	if err := ctx.Bind(&req); err != nil {
		return nil, shared.ErrInvalidInput{Message: "failed to parse request: " + err.Error()}
	}
	if req.Markdown == "" {
		return nil, shared.ErrInvalidInput{Message: "markdown field is required"}
	}

	pdfData, err := h.svc.MarkdownToPDF(ctx, generator.MarkdownToPDFParams{
		Markdown:   req.Markdown,
		PageConfig: buildPageConfig(req.PageSize, req.Landscape, req.MarginLeft, req.MarginTop, req.MarginRight, req.MarginBottom),
		Meta: generator.DocumentMeta{
			Title:   req.Title,
			Author:  req.Author,
			Subject: req.Subject,
		},
		Bookmarks: req.Bookmarks,
	})
	if err != nil {
		return nil, err
	}

	return response.File{
		Content:     pdfData,
		ContentType: "application/pdf",
	}, nil
}

func buildPageConfig(pageSize string, landscape bool, ml, mt, mr, mb float64) shared.PageConfig {
	margins := shared.DefaultMargins()
	if ml > 0 {
		margins.Left = ml
	}
	if mt > 0 {
		margins.Top = mt
	}
	if mr > 0 {
		margins.Right = mr
	}
	if mb > 0 {
		margins.Bottom = mb
	}
	return shared.PageConfig{
		Size:        shared.ParsePageSize(pageSize),
		Orientation: shared.ParseOrientation(landscape),
		Margins:     margins,
	}
}
