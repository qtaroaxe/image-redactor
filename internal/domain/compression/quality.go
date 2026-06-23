package compression

import (
	"fmt"

	apperrors "github.com/QtaroAXE/image-redactor/internal/domain/errors"
)

type Quality struct {
	value int
}

var (
	DefaultQuality = Quality{85}
)

func (q Quality) Value() int {
	return q.value
}
func (q Quality) IsZero() bool {
	return q.value == 0
}

func NewQuality(input int) (Quality, error) {
	if input == 0 {
		return DefaultQuality, nil
	}
	if input < 1 || input > 100 {
		err := fmt.Errorf("quality must be between 1 and 100, got %d", input)
		return Quality{}, apperrors.Wrap(
			err,
			apperrors.TypeInvalidInput,
			"invalid compression quality",
		).WithContext("quality_value", input)
	}
	return Quality{value: input}, nil
}
