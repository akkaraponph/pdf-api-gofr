package converter

import "pdfapi/internal/domain/shared"

// ConvertToImagesParams holds parameters for converting a PDF to images.
type ConvertToImagesParams struct {
	PDFData []byte
	Format  shared.ImageFormat
	DPI     int
	Pages   []int // 1-indexed; empty = all pages
}

// ImageData represents a single converted image.
type ImageData struct {
	Data     []byte
	Filename string
}

// ConvertToImagesResult holds the output of a PDF-to-images conversion.
type ConvertToImagesResult struct {
	Images []ImageData
	Format shared.ImageFormat
}

// ImageInput represents an uploaded image for PDF conversion.
type ImageInput struct {
	Data     []byte
	Filename string
}

// ImagesToPDFParams holds parameters for combining images into a PDF.
type ImagesToPDFParams struct {
	Images   []ImageInput
	PageSize shared.PageSize
	FitMode  shared.ImageFitMode
}
