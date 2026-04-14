package operation

// MergeParams holds parameters for merging multiple PDFs.
type MergeParams struct {
	PDFs [][]byte
}

// PageRange defines a range of pages (1-indexed, inclusive).
type PageRange struct {
	From int
	To   int
}

// SplitParams holds parameters for splitting a PDF.
type SplitParams struct {
	PDFData []byte
	Ranges  []PageRange // empty = one file per page
}

// PDFOutput represents a single output PDF from a split operation.
type PDFOutput struct {
	Data     []byte
	Filename string
}

// SplitResult holds the output of a split operation.
type SplitResult struct {
	PDFs []PDFOutput
}

// CompressParams holds parameters for compressing a PDF.
type CompressParams struct {
	PDFData      []byte
	ImageQuality int // 1-100; 0 = skip image re-encoding
}

// WatermarkParams holds parameters for adding a watermark to a PDF.
type WatermarkParams struct {
	PDFData  []byte
	Text     string
	FontSize float64
	Opacity  float64
	Rotation float64
	Color    [3]int
	Template string // "draft", "confidential", "copy", "sample", "do-not-copy"
}

// ExtractTextResult holds the extracted text per page.
type ExtractTextResult struct {
	Pages []string
}

// DecryptParams holds parameters for decrypting a PDF.
type DecryptParams struct {
	PDFData  []byte
	Password string
}
