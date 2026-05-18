package compression

import (
	"fmt"
)

type CompressionLevel struct {
	value int
}

var (
	DefaultCompressionLevel = CompressionLevel{6}
)

func (g CompressionLevel) GetLevel() int {
	return g.value
}
func (g CompressionLevel) IsZero() bool {
	return g.value == 0
}

// Увеличивает на 1
func NewCompressionLevel(input int) (CompressionLevel, error) {
	var compressionLevel CompressionLevel

	if input == 0 {
		compressionLevel = DefaultCompressionLevel
	} else {
		compressionLevel.value = input
	}
	if compressionLevel.value < 1 || compressionLevel.value > 10 {
		return CompressionLevel{}, fmt.Errorf("compression: compression quality cannot be lower than 1 or higher than 100")
	} else {
		return compressionLevel, nil
	}
}
