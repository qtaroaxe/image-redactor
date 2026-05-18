package compression

import (
	"fmt"
)

type Quality struct {
	value int
}

var (
	DefaultQuality = Quality{85}
)

func (q Quality) GetPercent() int {
	return q.value
}
func (q Quality) IsZero() bool {
	return q.value == 0
}

// Увеличивает на 1
func NewQuality(input int) (Quality, error) {
	quality := Quality{input}
	if quality.value < 1 || quality.value > 100 {
		return Quality{}, fmt.Errorf("compression: compression quality cannot be lower than 1 or higher than 100")
	} else {
		return quality, nil
	}
}
