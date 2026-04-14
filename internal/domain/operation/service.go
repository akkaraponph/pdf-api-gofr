package operation

import "context"

// Service defines the operation domain operations.
type Service interface {
	Merge(ctx context.Context, params MergeParams) ([]byte, error)
	Split(ctx context.Context, params SplitParams) (*SplitResult, error)
	Compress(ctx context.Context, params CompressParams) ([]byte, error)
	Watermark(ctx context.Context, params WatermarkParams) ([]byte, error)
	ExtractText(ctx context.Context, pdfData []byte) (*ExtractTextResult, error)
	Decrypt(ctx context.Context, params DecryptParams) ([]byte, error)
}
