package imginfo

import (
	"fmt"
	apperrors "github.com/QtaroAXE/image-redactor/internal/domain/errors"
	"strings"
)

type SourceImage struct {
	path      string
	format    Format
	sizeBytes int64
	width     int // 0 means "unknown"
	height    int // 0 means "unknown"
}

func NewSourceImage(path string, format Format, sizeBytes int64) (SourceImage, error) {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return SourceImage{}, apperrors.New(apperrors.TypeInvalidInput, "path cannot be empty")
	}
	if format.IsZero() {
		return SourceImage{}, apperrors.New(apperrors.TypeInvalidInput, "format cannot be empty")
	}
	if sizeBytes < 0 {
		return SourceImage{}, apperrors.New(apperrors.TypeInvalidInput, "size cannot be lower than 0").WithContext("size_bytes", sizeBytes)
	}
	return SourceImage{
		path:      trimmed,
		format:    format,
		sizeBytes: sizeBytes,
	}, nil
}

func (s SourceImage) WithDimensions(width, height int) (SourceImage, error) {
	if width <= 0 || height <= 0 {
		return SourceImage{}, fmt.Errorf("image: dimensions must be positive, got %dx%d", width, height)
	}
	s.width = width
	s.height = height
	return s, nil
}

func (s SourceImage) Path() string        { return s.path }
func (s SourceImage) Format() Format      { return s.format }
func (s SourceImage) SizeBytes() int64    { return s.sizeBytes }
func (s SourceImage) Width() int          { return s.width }
func (s SourceImage) Height() int         { return s.height }
func (s SourceImage) HasDimensions() bool { return s.width > 0 && s.height > 0 }
