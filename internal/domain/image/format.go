package image

import(
	"fmt"
	"strings"
)
//struct to create supported formats
type Format struct{
	name string
}
//supported formats
var(
	FormatJPEG = Format{name:"jpeg"}
	FormatPNG = Format{name:"png"}
	FormatWebP = Format{name:"webp"}
)
//String() returns image format real name
func (f Format) String() string{
	return f.name
}

func (f Format) IsZero()bool{
	return f.name == ""
}
//true or false if format matching with other
func (f Format) Equals(otherF Format)bool{
	return f.name == otherF
}

func NewFormat(input string) (Format, error){
	inputForm := strings.TrimSpace(input)
	inputForm = strings.ToLower(inputForm)

	if inputForm == ""{
		return Format{}, fmt.Errorf("image format cannot be empty")
	}
	switch inputForm {
	case "jpeg", "jpg":
		return FormatJPEG, nil
	case "png":
		 return FormatPNG, nil
	case "webp":
		return FormatWebP, nil
	default:
		return Format{}, fmt.Errorf("unsupported image format: %s", input)
	}
}