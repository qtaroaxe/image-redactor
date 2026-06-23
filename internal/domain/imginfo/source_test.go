package image

import (
	"testing"
)

func TestNewSourceImage(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		format    Format
		sizeBytes int64
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "valid source image",
			path:      "/path/to/image.jpg",
			format:    FormatJPEG,
			sizeBytes: 1024,
			wantErr:   false,
		},
		{
			name:      "valid source image with spaces in path",
			path:      "  /path/to/image.jpg  ",
			format:    FormatPNG,
			sizeBytes: 2048,
			wantErr:   false,
		},
		{
			name:      "empty path",
			path:      "",
			format:    FormatJPEG,
			sizeBytes: 1024,
			wantErr:   true,
			errMsg:    "image: source path is empty",
		},
		{
			name:      "path with only spaces",
			path:      "   ",
			format:    FormatJPEG,
			sizeBytes: 1024,
			wantErr:   true,
			errMsg:    "image: source path is empty",
		},
		{
			name:      "zero format",
			path:      "/path/to/image.jpg",
			format:    Format{""},
			sizeBytes: 1024,
			wantErr:   true,
			errMsg:    "image: source format is not set",
		},
		{
			name:      "negative size",
			path:      "/path/to/image.jpg",
			format:    FormatJPEG,
			sizeBytes: -100,
			wantErr:   true,
			errMsg:    "image: source size must be non-negative, got -100",
		},
		{
			name:      "zero size is valid",
			path:      "/path/to/image.jpg",
			format:    FormatPNG,
			sizeBytes: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSourceImage(tt.path, tt.format, tt.sizeBytes)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewSourceImage() expected error but got none")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("NewSourceImage() error = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("NewSourceImage() unexpected error: %v", err)
				return
			}

			// Проверяем поля
			expectedPath := tt.path
			if tt.path != "" {
				// Путь должен быть обрезан от пробелов
				if tt.path == "  /path/to/image.jpg  " {
					expectedPath = "/path/to/image.jpg"
				}
			}

			if got.Path() != expectedPath {
				t.Errorf("Path() = %v, want %v", got.Path(), expectedPath)
			}
			if got.Format() != tt.format {
				t.Errorf("Format() = %v, want %v", got.Format(), tt.format)
			}
			if got.SizeBytes() != tt.sizeBytes {
				t.Errorf("SizeBytes() = %v, want %v", got.SizeBytes(), tt.sizeBytes)
			}
			if got.Width() != 0 {
				t.Errorf("Width() = %v, want 0", got.Width())
			}
			if got.Height() != 0 {
				t.Errorf("Height() = %v, want 0", got.Height())
			}
			if got.HasDimensions() {
				t.Error("HasDimensions() should be false for new image without dimensions")
			}
		})
	}
}

func TestSourceImage_WithDimensions(t *testing.T) {
	baseImage, err := NewSourceImage("/test/image.jpg", FormatJPEG, 1024)
	if err != nil {
		t.Fatalf("Failed to create base image: %v", err)
	}

	tests := []struct {
		name    string
		width   int
		height  int
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid dimensions",
			width:   1920,
			height:  1080,
			wantErr: false,
		},
		{
			name:    "valid square dimensions",
			width:   800,
			height:  800,
			wantErr: false,
		},
		{
			name:    "zero width",
			width:   0,
			height:  1080,
			wantErr: true,
			errMsg:  "image: dimensions must be positive, got 0x1080",
		},
		{
			name:    "zero height",
			width:   1920,
			height:  0,
			wantErr: true,
			errMsg:  "image: dimensions must be positive, got 1920x0",
		},
		{
			name:    "negative width",
			width:   -100,
			height:  1080,
			wantErr: true,
			errMsg:  "image: dimensions must be positive, got -100x1080",
		},
		{
			name:    "negative height",
			width:   1920,
			height:  -100,
			wantErr: true,
			errMsg:  "image: dimensions must be positive, got 1920x-100",
		},
		{
			name:    "both negative",
			width:   -1920,
			height:  -1080,
			wantErr: true,
			errMsg:  "image: dimensions must be positive, got -1920x-1080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img, err := baseImage.WithDimensions(tt.width, tt.height)

			if tt.wantErr {
				if err == nil {
					t.Errorf("WithDimensions() expected error but got none")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("WithDimensions() error = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("WithDimensions() unexpected error: %v", err)
				return
			}

			// Проверяем, что оригинальное изображение не изменилось
			if baseImage.Width() != 0 || baseImage.Height() != 0 {
				t.Error("Original image should not be modified")
			}

			// Проверяем новое изображение
			if img.Width() != tt.width {
				t.Errorf("Width() = %v, want %v", img.Width(), tt.width)
			}
			if img.Height() != tt.height {
				t.Errorf("Height() = %v, want %v", img.Height(), tt.height)
			}
			if !img.HasDimensions() {
				t.Error("HasDimensions() should be true for image with dimensions")
			}

			// Проверяем, что остальные поля сохранились
			if img.Path() != baseImage.Path() {
				t.Errorf("Path() = %v, want %v", img.Path(), baseImage.Path())
			}
			if img.Format() != baseImage.Format() {
				t.Errorf("Format() = %v, want %v", img.Format(), baseImage.Format())
			}
			if img.SizeBytes() != baseImage.SizeBytes() {
				t.Errorf("SizeBytes() = %v, want %v", img.SizeBytes(), baseImage.SizeBytes())
			}
		})
	}
}

func TestSourceImage_HasDimensions(t *testing.T) {
	img, _ := NewSourceImage("/test.jpg", FormatJPEG, 100)

	if img.HasDimensions() {
		t.Error("HasDimensions() should return false for image without dimensions")
	}

	imgWithDim, _ := img.WithDimensions(100, 200)
	if !imgWithDim.HasDimensions() {
		t.Error("HasDimensions() should return true for image with dimensions")
	}
}

func TestSourceImage_Getters(t *testing.T) {
	path := "/test/unique/path.jpg"
	format := FormatPNG
	size := int64(54321)

	img, err := NewSourceImage(path, format, size)
	if err != nil {
		t.Fatalf("Failed to create image: %v", err)
	}

	if img.Path() != path {
		t.Errorf("Path() = %v, want %v", img.Path(), path)
	}

	if img.Format() != format {
		t.Errorf("Format() = %v, want %v", img.Format(), format)
	}

	if img.SizeBytes() != size {
		t.Errorf("SizeBytes() = %v, want %v", img.SizeBytes(), size)
	}

	if img.Width() != 0 {
		t.Errorf("Width() = %v, want 0", img.Width())
	}

	if img.Height() != 0 {
		t.Errorf("Height() = %v, want 0", img.Height())
	}
}

func TestSourceImage_Immutability(t *testing.T) {
	// Тест проверяет иммутабельность методов WithDimensions
	original, err := NewSourceImage("/test.jpg", FormatJPEG, 1000)
	if err != nil {
		t.Fatalf("Failed to create image: %v", err)
	}

	modified, err := original.WithDimensions(1920, 1080)
	if err != nil {
		t.Fatalf("Failed to set dimensions: %v", err)
	}

	// Оригинал не должен измениться
	if original.Width() != 0 || original.Height() != 0 {
		t.Error("Original image should not have dimensions")
	}

	// Модифицированная копия должна иметь размеры
	if modified.Width() != 1920 || modified.Height() != 1080 {
		t.Error("Modified image should have dimensions")
	}

	// Другие поля должны совпадать
	if original.Path() != modified.Path() {
		t.Error("Path should be the same")
	}
	if original.Format() != modified.Format() {
		t.Error("Format should be the same")
	}
	if original.SizeBytes() != modified.SizeBytes() {
		t.Error("SizeBytes should be the same")
	}

}
