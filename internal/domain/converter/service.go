package converter

import "context"

// Service defines the converter domain operations.
type Service interface {
	ConvertToImages(ctx context.Context, params ConvertToImagesParams) (*ConvertToImagesResult, error)
	ImagesToPDF(ctx context.Context, params ImagesToPDFParams) ([]byte, error)
}
