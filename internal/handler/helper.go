package handler

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"strconv"
	"strings"

	"pdfapi/internal/domain/operation"
)

// saveUpload reads an uploaded file into memory.
func saveUpload(fh *multipart.FileHeader) ([]byte, error) {
	f, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

// collectFileHeaders filters nil entries from a list of file header pointers.
func collectFileHeaders(headers ...*multipart.FileHeader) []*multipart.FileHeader {
	var result []*multipart.FileHeader
	for _, h := range headers {
		if h != nil {
			result = append(result, h)
		}
	}
	return result
}

// parsePages parses a comma-separated list of 1-indexed page numbers.
func parsePages(s string) ([]int, error) {
	if s == "" {
		return nil, nil
	}
	parts := strings.Split(s, ",")
	pages := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("invalid page number: %q", p)
		}
		if n < 1 {
			return nil, fmt.Errorf("page number must be >= 1, got %d", n)
		}
		pages = append(pages, n)
	}
	return pages, nil
}

// parsePageRanges parses ranges like "1-3,5-7" or "1-3" into PageRange slices.
func parsePageRanges(s string) ([]operation.PageRange, error) {
	if s == "" {
		return nil, nil
	}
	parts := strings.Split(s, ",")
	var ranges []operation.PageRange
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		dashIdx := strings.Index(p, "-")
		if dashIdx < 0 {
			n, err := strconv.Atoi(p)
			if err != nil {
				return nil, fmt.Errorf("invalid range: %q", p)
			}
			ranges = append(ranges, operation.PageRange{From: n, To: n})
			continue
		}
		from, err := strconv.Atoi(strings.TrimSpace(p[:dashIdx]))
		if err != nil {
			return nil, fmt.Errorf("invalid range: %q", p)
		}
		to, err := strconv.Atoi(strings.TrimSpace(p[dashIdx+1:]))
		if err != nil {
			return nil, fmt.Errorf("invalid range: %q", p)
		}
		if from < 1 || to < from {
			return nil, fmt.Errorf("invalid range: %q", p)
		}
		ranges = append(ranges, operation.PageRange{From: from, To: to})
	}
	return ranges, nil
}

type zipEntry struct {
	Data []byte
	Name string
}

// zipBytes creates a zip archive from named byte slices.
func zipBytes(entries []zipEntry) ([]byte, error) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for _, e := range entries {
		w, err := zw.Create(e.Name)
		if err != nil {
			return nil, err
		}
		if _, err := w.Write(e.Data); err != nil {
			return nil, err
		}
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
