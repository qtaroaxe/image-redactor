package compression

import (
	"fmt"

	apperrors "github.com/QtaroAXE/image-redactor/internal/domain/errors"
)

type CompressionLevel struct {
	value int
}

var (
	DefaultCompressionLevel = CompressionLevel{6}
)

func (g CompressionLevel) Value() int {
	return g.value
}
func (g CompressionLevel) IsZero() bool {
	return g.value == 0
}

// Увеличивает на 1
func NewCompressionLevel(input int) (CompressionLevel, error) {
	if input == 0 {
		return DefaultCompressionLevel, nil
	}
	if input < 1 || input > 10 {
		err := fmt.Errorf("compression level must be between 1 and 10, got %d", input)
		return CompressionLevel{}, apperrors.Wrap(
			err,
			apperrors.TypeInvalidInput,
			"invalid compression level",
		).WithContext("level", input)
	}
	return CompressionLevel{value: input}, nil
}
