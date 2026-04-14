package shared

import "strings"

// ImageFormat represents an output image format.
type ImageFormat string

const (
	FormatJPEG ImageFormat = "jpeg"
	FormatPNG  ImageFormat = "png"
)

// ParseImageFormat parses a string into an ImageFormat, defaulting to JPEG.
func ParseImageFormat(s string) ImageFormat {
	if strings.EqualFold(strings.TrimSpace(s), "png") {
		return FormatPNG
	}
	return FormatJPEG
}

// ContentType returns the MIME content type for the format.
func (f ImageFormat) ContentType() string {
	if f == FormatPNG {
		return "image/png"
	}
	return "image/jpeg"
}

// ImageFitMode controls how images are placed on fixed-size pages.
type ImageFitMode string

const (
	FitContain ImageFitMode = "fit"
	FitCover   ImageFitMode = "fill"
	FitStretch ImageFitMode = "stretch"
)

// ParseImageFitMode parses a string into an ImageFitMode, defaulting to Contain.
func ParseImageFitMode(s string) ImageFitMode {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "fill":
		return FitCover
	case "stretch":
		return FitStretch
	default:
		return FitContain
	}
}
