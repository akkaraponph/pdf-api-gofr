package handler

import (
	"mime/multipart"
	"strconv"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http/response"

	"pdfapi/internal/domain/operation"
	"pdfapi/internal/domain/shared"
)

// OperationHandler handles PDF manipulation HTTP requests.
type OperationHandler struct {
	svc operation.Service
}

// NewOperationHandler creates a new OperationHandler.
func NewOperationHandler(svc operation.Service) *OperationHandler {
	return &OperationHandler{svc: svc}
}

// --- Merge ---

type mergePDFRequest struct {
	File1  *multipart.FileHeader `file:"file1"`
	File2  *multipart.FileHeader `file:"file2"`
	File3  *multipart.FileHeader `file:"file3"`
	File4  *multipart.FileHeader `file:"file4"`
	File5  *multipart.FileHeader `file:"file5"`
	File6  *multipart.FileHeader `file:"file6"`
	File7  *multipart.FileHeader `file:"file7"`
	File8  *multipart.FileHeader `file:"file8"`
	File9  *multipart.FileHeader `file:"file9"`
	File10 *multipart.FileHeader `file:"file10"`
}

// MergePDF merges multiple PDF files into one.
func (h *OperationHandler) MergePDF(ctx *gofr.Context) (any, error) {
	var req mergePDFRequest
	if err := ctx.Bind(&req); err != nil {
		return nil, shared.ErrInvalidInput{Message: "failed to parse request: " + err.Error()}
	}

	headers := collectFileHeaders(
		req.File1, req.File2, req.File3, req.File4, req.File5,
		req.File6, req.File7, req.File8, req.File9, req.File10,
	)
	if len(headers) < 2 {
		return nil, shared.ErrInvalidInput{Message: "at least 2 PDF files are required (upload as file1, file2, ...)"}
	}

	var pdfs [][]byte
	for _, fh := range headers {
		data, err := saveUpload(fh)
		if err != nil {
			return nil, err
		}
		pdfs = append(pdfs, data)
	}

	merged, err := h.svc.Merge(ctx, operation.MergeParams{PDFs: pdfs})
	if err != nil {
		return nil, err
	}

	return response.File{
		Content:     merged,
		ContentType: "application/pdf",
	}, nil
}

// --- Split ---

type splitPDFRequest struct {
	PDF    *multipart.FileHeader `file:"pdf"`
	Ranges string                `form:"ranges"` // "1-3,5-7"
}

// SplitPDF splits a PDF into parts by page ranges.
func (h *OperationHandler) SplitPDF(ctx *gofr.Context) (any, error) {
	var req splitPDFRequest
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

	ranges, err := parsePageRanges(req.Ranges)
	if err != nil {
		return nil, shared.ErrInvalidInput{Message: err.Error()}
	}

	result, err := h.svc.Split(ctx, operation.SplitParams{
		PDFData: pdfData,
		Ranges:  ranges,
	})
	if err != nil {
		return nil, err
	}

	// Single output: return directly
	if len(result.PDFs) == 1 {
		return response.File{
			Content:     result.PDFs[0].Data,
			ContentType: "application/pdf",
		}, nil
	}

	// Multiple outputs: return as zip
	entries := make([]zipEntry, len(result.PDFs))
	for i, p := range result.PDFs {
		entries[i] = zipEntry{Data: p.Data, Name: p.Filename}
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

// --- Compress ---

type compressPDFRequest struct {
	PDF          *multipart.FileHeader `file:"pdf"`
	ImageQuality string                `form:"image_quality"` // 1-100
}

// CompressPDF compresses a PDF.
func (h *OperationHandler) CompressPDF(ctx *gofr.Context) (any, error) {
	var req compressPDFRequest
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

	var quality int
	if req.ImageQuality != "" {
		quality, _ = strconv.Atoi(req.ImageQuality)
	}

	compressed, err := h.svc.Compress(ctx, operation.CompressParams{
		PDFData:      pdfData,
		ImageQuality: quality,
	})
	if err != nil {
		return nil, err
	}

	return response.File{
		Content:     compressed,
		ContentType: "application/pdf",
	}, nil
}

// --- Watermark ---

type watermarkPDFRequest struct {
	PDF      *multipart.FileHeader `file:"pdf"`
	Text     string                `form:"text"`
	Template string                `form:"template"` // "draft", "confidential", "copy", "sample", "do-not-copy"
	FontSize string                `form:"font_size"`
	Opacity  string                `form:"opacity"`
	Rotation string                `form:"rotation"`
	ColorR   string                `form:"color_r"`
	ColorG   string                `form:"color_g"`
	ColorB   string                `form:"color_b"`
}

// WatermarkPDF adds a watermark to a PDF.
func (h *OperationHandler) WatermarkPDF(ctx *gofr.Context) (any, error) {
	var req watermarkPDFRequest
	if err := ctx.Bind(&req); err != nil {
		return nil, shared.ErrInvalidInput{Message: "failed to parse request: " + err.Error()}
	}
	if req.PDF == nil {
		return nil, shared.ErrInvalidInput{Message: "pdf file is required"}
	}
	if req.Text == "" && req.Template == "" {
		return nil, shared.ErrInvalidInput{Message: "text or template is required"}
	}

	pdfData, err := saveUpload(req.PDF)
	if err != nil {
		return nil, err
	}

	params := operation.WatermarkParams{
		PDFData:  pdfData,
		Text:     req.Text,
		Template: req.Template,
	}
	if req.FontSize != "" {
		params.FontSize, _ = strconv.ParseFloat(req.FontSize, 64)
	}
	if req.Opacity != "" {
		params.Opacity, _ = strconv.ParseFloat(req.Opacity, 64)
	}
	if req.Rotation != "" {
		params.Rotation, _ = strconv.ParseFloat(req.Rotation, 64)
	}
	if req.ColorR != "" || req.ColorG != "" || req.ColorB != "" {
		r, _ := strconv.Atoi(req.ColorR)
		g, _ := strconv.Atoi(req.ColorG)
		b, _ := strconv.Atoi(req.ColorB)
		params.Color = [3]int{r, g, b}
	}

	watermarked, err := h.svc.Watermark(ctx, params)
	if err != nil {
		return nil, err
	}

	return response.File{
		Content:     watermarked,
		ContentType: "application/pdf",
	}, nil
}

// --- Extract Text ---

type extractTextRequest struct {
	PDF *multipart.FileHeader `file:"pdf"`
}

// ExtractText extracts text content from a PDF, returning JSON.
func (h *OperationHandler) ExtractText(ctx *gofr.Context) (any, error) {
	var req extractTextRequest
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

	result, err := h.svc.ExtractText(ctx, pdfData)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// --- Decrypt ---

type decryptPDFRequest struct {
	PDF      *multipart.FileHeader `file:"pdf"`
	Password string                `form:"password"`
}

// DecryptPDF removes password protection from a PDF.
func (h *OperationHandler) DecryptPDF(ctx *gofr.Context) (any, error) {
	var req decryptPDFRequest
	if err := ctx.Bind(&req); err != nil {
		return nil, shared.ErrInvalidInput{Message: "failed to parse request: " + err.Error()}
	}
	if req.PDF == nil {
		return nil, shared.ErrInvalidInput{Message: "pdf file is required"}
	}
	if req.Password == "" {
		return nil, shared.ErrInvalidInput{Message: "password is required"}
	}

	pdfData, err := saveUpload(req.PDF)
	if err != nil {
		return nil, err
	}

	decrypted, err := h.svc.Decrypt(ctx, operation.DecryptParams{
		PDFData:  pdfData,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return response.File{
		Content:     decrypted,
		ContentType: "application/pdf",
	}, nil
}
