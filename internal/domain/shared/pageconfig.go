package shared

import "strings"

// PageSize represents a standard page size.
type PageSize string

const (
	PageSizeA3     PageSize = "a3"
	PageSizeA4     PageSize = "a4"
	PageSizeA5     PageSize = "a5"
	PageSizeLetter PageSize = "letter"
	PageSizeLegal  PageSize = "legal"
)

// ParsePageSize parses a string into a PageSize, defaulting to A4.
func ParsePageSize(s string) PageSize {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "a3":
		return PageSizeA3
	case "a5":
		return PageSizeA5
	case "letter":
		return PageSizeLetter
	case "legal":
		return PageSizeLegal
	default:
		return PageSizeA4
	}
}

// Orientation represents page orientation.
type Orientation string

const (
	Portrait  Orientation = "portrait"
	Landscape Orientation = "landscape"
)

// ParseOrientation returns Landscape if true, Portrait otherwise.
func ParseOrientation(landscape bool) Orientation {
	if landscape {
		return Landscape
	}
	return Portrait
}

// Margins defines page margins in millimeters.
type Margins struct {
	Left   float64
	Top    float64
	Right  float64
	Bottom float64
}

// DefaultMargins returns 10mm margins on all sides.
func DefaultMargins() Margins {
	return Margins{Left: 10, Top: 10, Right: 10, Bottom: 10}
}

// PageConfig combines size, orientation, and margins.
type PageConfig struct {
	Size        PageSize
	Orientation Orientation
	Margins     Margins
}

// DefaultPageConfig returns an A4 portrait page with 10mm margins.
func DefaultPageConfig() PageConfig {
	return PageConfig{
		Size:        PageSizeA4,
		Orientation: Portrait,
		Margins:     DefaultMargins(),
	}
}
