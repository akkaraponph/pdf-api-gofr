package converter

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/akkaraponph/presspdf"

	"pdfapi/internal/domain/converter"
	"pdfapi/internal/domain/shared"
)

type service struct{}

// New creates a new converter service backed by presspdf.
func New() converter.Service {
	return &service{}
}

func (s *service) ConvertToImages(_ context.Context, params converter.ConvertToImagesParams) (*converter.ConvertToImagesResult, error) {
	tmpDir, err := os.MkdirTemp("", "conv-to-img-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	pdfPath := filepath.Join(tmpDir, "input.pdf")
	if err := os.WriteFile(pdfPath, params.PDFData, 0644); err != nil {
		return nil, err
	}

	var opts []presspdf.ConvertOption
	opts = append(opts, presspdf.WithFormat(mapImageFormat(params.Format)))

	dpi := params.DPI
	if dpi <= 0 {
		dpi = 300
	}
	opts = append(opts, presspdf.WithDPI(dpi))

	if len(params.Pages) > 0 {
		opts = append(opts, presspdf.WithPages(params.Pages...))
	}

	outputDir := filepath.Join(tmpDir, "output")
	imagePaths, err := presspdf.ConvertToImages(pdfPath, outputDir, opts...)
	if err != nil {
		if strings.Contains(err.Error(), "no PDF renderer found") {
			return nil, shared.ErrToolNotFound{Tool: "pdftoppm/mutool/gs"}
		}
		return nil, shared.ErrProcessingFailed{Operation: "pdf-to-images", Cause: err}
	}

	result := &converter.ConvertToImagesResult{Format: params.Format}
	for _, p := range imagePaths {
		data, err := os.ReadFile(p)
		if err != nil {
			return nil, err
		}
		result.Images = append(result.Images, converter.ImageData{
			Data:     data,
			Filename: filepath.Base(p),
		})
	}

	return result, nil
}

func (s *service) ImagesToPDF(_ context.Context, params converter.ImagesToPDFParams) ([]byte, error) {
	tmpDir, err := os.MkdirTemp("", "img-to-pdf-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	var imagePaths []string
	for i, img := range params.Images {
		p := filepath.Join(tmpDir, fmt.Sprintf("%03d_%s", i, img.Filename))
		if err := os.WriteFile(p, img.Data, 0644); err != nil {
			return nil, err
		}
		imagePaths = append(imagePaths, p)
	}

	var opts []presspdf.ImageToPDFOption
	if ps, ok := mapPageSize(params.PageSize); ok {
		opts = append(opts, presspdf.ImagePageSize(ps))
	}
	if params.FitMode != "" {
		opts = append(opts, presspdf.ImageFit(string(params.FitMode)))
	}

	outputPath := filepath.Join(tmpDir, "output.pdf")
	if err := presspdf.ImagesToPDF(outputPath, imagePaths, opts...); err != nil {
		return nil, shared.ErrProcessingFailed{Operation: "images-to-pdf", Cause: err}
	}

	return os.ReadFile(outputPath)
}

func mapImageFormat(f shared.ImageFormat) presspdf.ImageFormat {
	if f == shared.FormatPNG {
		return presspdf.PNG
	}
	return presspdf.JPEG
}

func mapPageSize(s shared.PageSize) (presspdf.PageSize, bool) {
	switch s {
	case shared.PageSizeA3:
		return presspdf.A3, true
	case shared.PageSizeA4:
		return presspdf.A4, true
	case shared.PageSizeA5:
		return presspdf.A5, true
	case shared.PageSizeLetter:
		return presspdf.Letter, true
	case shared.PageSizeLegal:
		return presspdf.Legal, true
	default:
		return presspdf.PageSize{}, false
	}
}
