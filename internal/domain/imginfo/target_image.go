package imginfo

import (
	"strings"

	"github.com/QtaroAXE/image-redactor/internal/domain/compression"
	apperrors "github.com/QtaroAXE/image-redactor/internal/domain/errors"
)

type TargetImage struct {
	path             string
	format           Format
	quality          compression.Quality
	compressionLevel compression.CompressionLevel
}

func NewTargetImage(path string, qual compression.Quality, compLevel compression.CompressionLevel, form Format) (TargetImage, error) {
	trimmedPath := strings.TrimSpace(path)
	if trimmedPath == "" {
		return TargetImage{}, apperrors.New(apperrors.TypeInvalidInput, "path cannot be empty")
	}
	if form.IsZero() {
		return TargetImage{}, apperrors.New(apperrors.TypeInvalidInput, "format cannot be empty")
	}
	return TargetImage{path: trimmedPath, quality: qual, compressionLevel: compLevel, format: form}, nil
}
func (t TargetImage) Path() string {
	return t.path
}
func (t TargetImage) Quality() compression.Quality {
	return t.quality
}
func (t TargetImage) CompressionLevel() compression.CompressionLevel {
	return t.compressionLevel
}
func (t TargetImage) Format() Format {
	return t.format
}
