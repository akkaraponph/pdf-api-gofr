package handler

import (
	"mime/multipart"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http/response"

	"pdfapi/internal/domain/converter"
	"pdfapi/internal/domain/shared"
)

// ConverterHandler handles PDF↔image conversion HTTP requests.
type ConverterHandler struct {
	svc converter.Service
}

// NewConverterHandler creates a new ConverterHandler.
func NewConverterHandler(svc converter.Service) *ConverterHandler {
	return &ConverterHandler{svc: svc}
}

type pdfToImagesRequest struct {
	PDF    *multipart.FileHeader `file:"pdf"`
	Format string                `form:"format"` // "jpeg" or "png" (default: jpeg)
	DPI    string                `form:"dpi"`    // resolution (default: 300)
	Pages  string                `form:"pages"`  // "1,3,5"
}

// PDFToImages converts a PDF file into images.
func (h *ConverterHandler) PDFToImages(ctx *gofr.Context) (any, error) {
	var req pdfToImagesRequest
	if err := ctx.Bind(&req); err != nil {
		return nil, shared.ErrInvalidInput{Message: "failed to parse request: " + err.Error()}
	}
	if req.PDF == nil {
		return nil, shared.ErrInvalidInput{Message: "pdf file is required"}
	}

	pdfData, err := saveUpload(req.PDF)
	if err != nil {
		return nil, err
	}

	pages, err := parsePages(req.Pages)
	if err != nil {
		return nil, shared.ErrInvalidInput{Message: err.Error()}
	}

	dpi := 300
	if req.DPI != "" {
		if v, err := parseInt(req.DPI); err == nil && v > 0 {
			dpi = v
		}
	}

	format := shared.ParseImageFormat(req.Format)

	result, err := h.svc.ConvertToImages(ctx, converter.ConvertToImagesParams{
		PDFData: pdfData,
		Format:  format,
		DPI:     dpi,
		Pages:   pages,
	})
	if err != nil {
		return nil, err
	}

	// Single image: return directly
	if len(result.Images) == 1 {
		return response.File{
			Content:     result.Images[0].Data,
			ContentType: format.ContentType(),
		}, nil
	}

	// Multiple images: return as zip
	entries := make([]zipEntry, len(result.Images))
	for i, img := range result.Images {
		entries[i] = zipEntry{Data: img.Data, Name: img.Filename}
	}
	zipData, err := zipBytes(entries)
	if err != nil {
		return nil, err
	}

	return response.File{
		Content:     zipData,
		ContentType: "application/zip",
	}, nil
}

type imagesToPDFRequest struct {
	Image1   *multipart.FileHeader `file:"image1"`
	Image2   *multipart.FileHeader `file:"image2"`
	Image3   *multipart.FileHeader `file:"image3"`
	Image4   *multipart.FileHeader `file:"image4"`
	Image5   *multipart.FileHeader `file:"image5"`
	Image6   *multipart.FileHeader `file:"image6"`
	Image7   *multipart.FileHeader `file:"image7"`
	Image8   *multipart.FileHeader `file:"image8"`
	Image9   *multipart.FileHeader `file:"image9"`
	Image10  *multipart.FileHeader `file:"image10"`
	PageSize string                `form:"page_size"` // "a4", "a3", "letter", "auto"
	Fit      string                `form:"fit"`       // "fit", "fill", "stretch"
}

// ImagesToPDF converts uploaded images into a single PDF.
func (h *ConverterHandler) ImagesToPDF(ctx *gofr.Context) (any, error) {
	var req imagesToPDFRequest
	if err := ctx.Bind(&req); err != nil {
		return nil, shared.ErrInvalidInput{Message: "failed to parse request: " + err.Error()}
	}

	headers := collectFileHeaders(
		req.Image1, req.Image2, req.Image3, req.Image4, req.Image5,
		req.Image6, req.Image7, req.Image8, req.Image9, req.Image10,
	)
	if len(headers) == 0 {
		return nil, shared.ErrInvalidInput{Message: "at least one image is required (upload as image1, image2, ...)"}
	}

	var images []converter.ImageInput
	for _, fh := range headers {
		data, err := saveUpload(fh)
		if err != nil {
			return nil, err
		}
		images = append(images, converter.ImageInput{Data: data, Filename: fh.Filename})
	}

	pdfData, err := h.svc.ImagesToPDF(ctx, converter.ImagesToPDFParams{
		Images:   images,
		PageSize: shared.ParsePageSize(req.PageSize),
		FitMode:  shared.ParseImageFitMode(req.Fit),
	})
	if err != nil {
		return nil, err
	}

	return response.File{
		Content:     pdfData,
		ContentType: "application/pdf",
	}, nil
}

func parseInt(s string) (int, error) {
	var n int
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, shared.ErrInvalidInput{Message: "invalid number: " + s}
		}
		n = n*10 + int(c-'0')
	}
	return n, nil
}
