package operation

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/akkaraponph/presspdf"

	"pdfapi/internal/domain/operation"
	"pdfapi/internal/domain/shared"
)

type service struct{}

// New creates a new operation service backed by presspdf.
func New() operation.Service {
	return &service{}
}

func (s *service) Merge(_ context.Context, params operation.MergeParams) ([]byte, error) {
	tmpDir, err := os.MkdirTemp("", "merge-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	var paths []string
	for i, data := range params.PDFs {
		p := filepath.Join(tmpDir, fmt.Sprintf("%03d.pdf", i))
		if err := os.WriteFile(p, data, 0644); err != nil {
			return nil, err
		}
		paths = append(paths, p)
	}

	outputPath := filepath.Join(tmpDir, "merged.pdf")
	if err := presspdf.MergePDF(outputPath, paths...); err != nil {
		return nil, shared.ErrProcessingFailed{Operation: "merge", Cause: err}
	}

	return os.ReadFile(outputPath)
}

func (s *service) Split(_ context.Context, params operation.SplitParams) (*operation.SplitResult, error) {
	tmpDir, err := os.MkdirTemp("", "split-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	pdfPath := filepath.Join(tmpDir, "input.pdf")
	if err := os.WriteFile(pdfPath, params.PDFData, 0644); err != nil {
		return nil, err
	}

	var opts []presspdf.SplitOption
	if len(params.Ranges) > 0 {
		ranges := make([]presspdf.PageRange, len(params.Ranges))
		for i, r := range params.Ranges {
			ranges[i] = presspdf.PageRange{From: r.From, To: r.To}
		}
		opts = append(opts, presspdf.WithRanges(ranges...))
	}

	outputDir := filepath.Join(tmpDir, "output")
	outputPaths, err := presspdf.SplitPDF(pdfPath, outputDir, opts...)
	if err != nil {
		return nil, shared.ErrProcessingFailed{Operation: "split", Cause: err}
	}

	result := &operation.SplitResult{}
	for _, p := range outputPaths {
		data, err := os.ReadFile(p)
		if err != nil {
			return nil, err
		}
		result.PDFs = append(result.PDFs, operation.PDFOutput{
			Data:     data,
			Filename: filepath.Base(p),
		})
	}

	return result, nil
}

func (s *service) Compress(_ context.Context, params operation.CompressParams) ([]byte, error) {
	tmpDir, err := os.MkdirTemp("", "compress-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	pdfPath := filepath.Join(tmpDir, "input.pdf")
	if err := os.WriteFile(pdfPath, params.PDFData, 0644); err != nil {
		return nil, err
	}

	var opts []presspdf.CompressOption
	if params.ImageQuality > 0 && params.ImageQuality <= 100 {
		opts = append(opts, presspdf.CompressImageQuality(params.ImageQuality))
	}

	outputPath := filepath.Join(tmpDir, "compressed.pdf")
	if err := presspdf.CompressPDF(pdfPath, outputPath, opts...); err != nil {
		return nil, shared.ErrProcessingFailed{Operation: "compress", Cause: err}
	}

	return os.ReadFile(outputPath)
}

func (s *service) Watermark(_ context.Context, params operation.WatermarkParams) ([]byte, error) {
	tmpDir, err := os.MkdirTemp("", "watermark-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	pdfPath := filepath.Join(tmpDir, "input.pdf")
	if err := os.WriteFile(pdfPath, params.PDFData, 0644); err != nil {
		return nil, err
	}

	var opts []presspdf.WatermarkOption

	if params.Template != "" {
		opts = append(opts, presspdf.WatermarkTemplate(params.Template))
	} else if params.Text != "" {
		opts = append(opts, presspdf.WatermarkText(params.Text))
	}

	if params.FontSize > 0 {
		opts = append(opts, presspdf.WatermarkFontSize(params.FontSize))
	}
	if params.Opacity > 0 {
		opts = append(opts, presspdf.WatermarkOpacity(params.Opacity))
	}
	if params.Rotation != 0 {
		opts = append(opts, presspdf.WatermarkRotation(params.Rotation))
	}
	if params.Color != [3]int{} {
		opts = append(opts, presspdf.WatermarkColor(params.Color[0], params.Color[1], params.Color[2]))
	}

	outputPath := filepath.Join(tmpDir, "watermarked.pdf")
	if err := presspdf.WatermarkPDF(pdfPath, outputPath, opts...); err != nil {
		return nil, shared.ErrProcessingFailed{Operation: "watermark", Cause: err}
	}

	return os.ReadFile(outputPath)
}

func (s *service) ExtractText(_ context.Context, pdfData []byte) (*operation.ExtractTextResult, error) {
	pages, err := presspdf.ExtractTextFromBytes(pdfData)
	if err != nil {
		return nil, shared.ErrProcessingFailed{Operation: "extract-text", Cause: err}
	}
	return &operation.ExtractTextResult{Pages: pages}, nil
}

func (s *service) Decrypt(_ context.Context, params operation.DecryptParams) ([]byte, error) {
	tmpDir, err := os.MkdirTemp("", "decrypt-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	pdfPath := filepath.Join(tmpDir, "input.pdf")
	if err := os.WriteFile(pdfPath, params.PDFData, 0644); err != nil {
		return nil, err
	}

	outputPath := filepath.Join(tmpDir, "decrypted.pdf")
	if err := presspdf.DecryptPDF(pdfPath, outputPath, params.Password); err != nil {
		return nil, shared.ErrProcessingFailed{Operation: "decrypt", Cause: err}
	}

	return os.ReadFile(outputPath)
}
