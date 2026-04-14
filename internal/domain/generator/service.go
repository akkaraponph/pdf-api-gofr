package generator

import "context"

// Service defines the generator domain operations.
type Service interface {
	HTMLToPDF(ctx context.Context, params HTMLToPDFParams) ([]byte, error)
	MarkdownToPDF(ctx context.Context, params MarkdownToPDFParams) ([]byte, error)
}
