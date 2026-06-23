package imginfo

import (
	"strings"

	apperrors "github.com/QtaroAXE/image-redactor/internal/domain/errors"
)

// struct to create supported formats
type Format struct {
	name string
}

// supported formats
var (
	FormatJPEG = Format{name: "jpeg"}
	FormatPNG  = Format{name: "png"}
	FormatWebP = Format{name: "webp"}
)

// String() returns image format real name
func (f Format) String() string {
	return f.name
}

func (f Format) IsZero() bool {
	return f.name == ""
}

// true or false if format matching with other
func (f Format) Equals(otherF string) bool {
	return f.name == otherF
}

func NewFormat(input string) (Format, error) {
	inputForm := strings.TrimSpace(input)
	inputForm = strings.ToLower(inputForm)

	if inputForm == "" {
		return Format{}, apperrors.New(apperrors.TypeInvalidInput, "image format cannot be empty")
	}
	switch inputForm {
	case "jpeg", "jpg":
		return FormatJPEG, nil
	case "png":
		return FormatPNG, nil
	case "webp":
		return FormatWebP, nil
	default:
		return Format{}, apperrors.New(apperrors.TypeUnsupported, "unsupported image format: "+inputForm).WithContext("format", inputForm)
	}
}
