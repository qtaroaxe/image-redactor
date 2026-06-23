// internal/service/fs/filesystem.go
package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	apperrors "github.com/QtaroAXE/image-redactor/internal/domain/errors"
)

type FileSystem struct {
	inputDir     string
	outputDir    string
	processedDir string
	errorDir     string
	mu           sync.RWMutex
}

func NewFileSystem(inputDir, outputDir, processedDir, errorDir string) (*FileSystem, error) {
	dirs := []string{inputDir, outputDir, processedDir, errorDir}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, apperrors.WrapWithFile(
				err,
				apperrors.TypeInternal,
				fmt.Sprintf("failed to create directory: %s", dir),
			).WithContext("directory", dir)
		}
	}

	return &FileSystem{
		inputDir:     inputDir,
		outputDir:    outputDir,
		processedDir: processedDir,
		errorDir:     errorDir,
	}, nil
}

func (fs *FileSystem) GetInputFiles() ([]string, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	var files []string

	err := filepath.Walk(fs.inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return apperrors.WrapWithFile(
				err,
				apperrors.TypeInternal,
				"failed to walk directory",
			).WithContext("path", path)
		}

		if info.IsDir() {
			return nil
		}

		if fs.isImageFile(path) {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func (fs *FileSystem) isImageFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp":
		return true
	default:
		return false
	}
}

func (fs *FileSystem) GetOutputPath(inputPath string, format string) string {
	fileName := filepath.Base(inputPath)
	ext := filepath.Ext(fileName)
	nameWithoutExt := fileName[:len(fileName)-len(ext)]

	return filepath.Join(fs.outputDir, nameWithoutExt+"."+format)
}

func (fs *FileSystem) SaveOutput(data []byte, outputPath string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return apperrors.WrapWithFile(
			err,
			apperrors.TypeInternal,
			"failed to create output directory",
		).WithContext("path", outputPath)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return apperrors.WrapWithFile(
			err,
			apperrors.TypeIO,
			"failed to write output file",
		).WithContext("path", outputPath)
	}

	return nil
}

func (fs *FileSystem) MoveToProcessed(path string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	dest := filepath.Join(fs.processedDir, filepath.Base(path))

	if err := os.Rename(path, dest); err != nil {
		return apperrors.WrapWithFile(
			err,
			apperrors.TypeIO,
			"failed to move file to processed directory",
		).WithContext("source", path).
			WithContext("destination", dest)
	}

	return nil
}

func (fs *FileSystem) MoveToError(path string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	dest := filepath.Join(fs.errorDir, filepath.Base(path))

	if err := os.Rename(path, dest); err != nil {
		return apperrors.WrapWithFile(
			err,
			apperrors.TypeIO,
			"failed to move file to error directory",
		).WithContext("source", path).
			WithContext("destination", dest)
	}

	return nil
}

func (fs *FileSystem) ReadFile(path string) ([]byte, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, apperrors.WrapWithFile(
			err,
			apperrors.TypeIO,
			"failed to read file",
		).WithContext("path", path)
	}

	return data, nil
}

func (fs *FileSystem) FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
func (fs *FileSystem) GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, apperrors.WrapWithFile(
			err,
			apperrors.TypeNotFound,
			"failed to get file info",
		).WithContext("path", path)
	}
	return info.Size(), nil
}

func (fs *FileSystem) GetInputDir() string {
	return fs.inputDir
}
func (fs *FileSystem) GetOutputDir() string {
	return fs.outputDir
}
