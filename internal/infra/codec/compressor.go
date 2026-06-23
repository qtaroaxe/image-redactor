package compressor

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"golang.org/x/image/draw"

	"github.com/QtaroAXE/image-redactor/internal/domain/compression"
	apperrors "github.com/QtaroAXE/image-redactor/internal/domain/errors"
	"github.com/QtaroAXE/image-redactor/internal/domain/imginfo"
	"github.com/QtaroAXE/image-redactor/internal/infra/fs"
)

type CompressorService struct {
	fs *fs.FileSystem
}

func NewCompressorService(fs *fs.FileSystem) *CompressorService {
	return &CompressorService{fs: fs}
}

func (s *CompressorService) CompressImage(src imginfo.SourceImage, target imginfo.TargetImage) error {
	data, err := s.fs.ReadFile(src.Path())
	if err != nil {
		return err
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return apperrors.WrapWithFile(
			err,
			apperrors.TypeDecode,
			"failed to decode image",
		).WithPath(src.Path())
	}

	if img.Bounds().Dx() == 0 || img.Bounds().Dy() == 0 {
		return apperrors.NewWithFile(
			apperrors.TypeValidate,
			"image has zero dimensions",
		).WithPath(src.Path())
	}

	if err := os.MkdirAll(filepath.Dir(target.Path()), 0755); err != nil {
		return apperrors.WrapWithFile(
			err,
			apperrors.TypeInternal,
			"failed to create output directory",
		).WithPath(target.Path())
	}

	out, err := os.Create(target.Path())
	if err != nil {
		return apperrors.WrapWithFile(
			err,
			apperrors.TypeIO,
			"failed to create output file",
		).WithPath(target.Path())
	}
	defer out.Close()

	switch target.Format().String() {
	case "jpeg":
		return s.encodeJPEG(img, out, target.Quality())
	case "png":
		return s.encodePNG(img, out, target.CompressionLevel())
	case "webp":
		return s.encodeWebP(img, out, target.Quality())
	default:
		return apperrors.NewWithFile(
			apperrors.TypeUnsupported,
			fmt.Sprintf("unsupported format: %s", target.Format()),
		).WithContext("format", target.Format())
	}
}

func (s *CompressorService) encodeJPEG(img image.Image, out *os.File, quality compression.Quality) error {
	opts := &jpeg.Options{
		Quality: quality.Value(),
	}

	if err := jpeg.Encode(out, img, opts); err != nil {
		return apperrors.WrapWithFile(
			err,
			apperrors.TypeEncode,
			"failed to encode JPEG",
		).WithContext("quality", quality.Value())
	}

	return nil
}
func (s *CompressorService) encodePNG(img image.Image, out *os.File, level compression.CompressionLevel) error {
	pngLevel := s.mapCompressionLevel(level.Value())

	encoder := &png.Encoder{
		CompressionLevel: pngLevel,
	}

	if err := encoder.Encode(out, img); err != nil {
		return apperrors.WrapWithFile(
			err,
			apperrors.TypeEncode,
			"failed to encode PNG",
		).WithContext("level", level.Value())
	}

	return nil
}

// заглушка
func (s *CompressorService) encodeWebP(img image.Image, out *os.File, quality compression.Quality) error {
	return apperrors.NewWithFile(
		apperrors.TypeUnsupported,
		"WebP encoding not yet implemented",
	).WithContext("quality", quality.Value())
}

func (s *CompressorService) mapCompressionLevel(level int) png.CompressionLevel {
	switch {
	case level <= 1:
		return png.BestSpeed
	case level >= 9:
		return png.BestCompression
	default:
		return png.CompressionLevel(level - 1)
	}
}

func (s *CompressorService) ResizeImage(img image.Image, width, height int) image.Image {
	srcBounds := img.Bounds()
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()

	if width == 0 && height == 0 {
		return img
	}

	// Сохраняем пропорции
	if width > 0 && height == 0 {
		ratio := float64(width) / float64(srcWidth)
		height = int(float64(srcHeight) * ratio)
	} else if height > 0 && width == 0 {
		ratio := float64(height) / float64(srcHeight)
		width = int(float64(srcWidth) * ratio)
	}

	if width == 0 {
		width = srcWidth
	}
	if height == 0 {
		height = srcHeight
	}

	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), img, srcBounds, draw.Over, nil)

	return dst
}
