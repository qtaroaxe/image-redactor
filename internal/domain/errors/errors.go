// pkg/errors/errors.go
package errors

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"
)

type AppError struct {
	Time time.Time `json:"time"`

	Type string `json:"type"`

	Message string `json:"message"`

	Err error `json:"-"`

	Context map[string]interface{} `json:"context,omitempty"`

	File     string `json:"file,omitempty"`     // имя файла
	Line     int    `json:"line,omitempty"`     // номер строки
	Function string `json:"function,omitempty"` // имя функции
}

const (
	TypeInvalidInput  = "INVALID_INPUT"
	TypeNotFound      = "NOT_FOUND"
	TypeAlreadyExists = "ALREADY_EXISTS"
	TypeInternal      = "INTERNAL_ERROR"
	TypeUnsupported   = "UNSUPPORTED"
	TypePermission    = "PERMISSION_DENIED"
	TypeIO            = "IO_ERROR"
	TypeDecode        = "DECODE_ERROR"
	TypeEncode        = "ENCODE_ERROR"
	TypeCompress      = "COMPRESS_ERROR"
	TypeValidate      = "VALIDATE_ERROR"
	TypeUnknown       = "UNKNOWN_ERROR"
	TypeTimeout       = "TIME_OUT"
)

func New(errType, message string) *AppError {
	err := &AppError{
		Time:    time.Now(),
		Type:    errType,
		Message: message,
		Context: make(map[string]interface{}),
	}
	// Захватываем информацию о файле (пропускаем 2 уровня: New и вызывающий код)
	err.captureFileInfo(2)
	return err
}

// NewWithFile - создает ошибку с явным указанием файла (устаревший, используйте New)
func NewWithFile(errType, message string) *AppError {
	return New(errType, message)
}

// Wrap - оборачивает существующую ошибку с АВТОМАТИЧЕСКИМ захватом файла
func Wrap(err error, errType, message string) *AppError {
	if err == nil {
		return nil
	}
	appErr := &AppError{
		Time:    time.Now(),
		Type:    errType,
		Message: message,
		Err:     err,
		Context: make(map[string]interface{}),
	}
	appErr.captureFileInfo(2)
	return appErr
}

func WrapWithFile(err error, errType, message string) *AppError {
	return Wrap(err, errType, message)
}

func (e *AppError) captureFileInfo(skip int) {
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		e.File = filepath.Base(file)
		e.Line = line

		fn := runtime.FuncForPC(pc)
		if fn != nil {
			e.Function = fn.Name()
		}
	}
}

func (e *AppError) Error() string {
	location := ""
	if e.File != "" {
		location = fmt.Sprintf(" (%s:%d)", e.File, e.Line)
	}

	if e.Err != nil {
		return fmt.Sprintf("[%s]%s %s: %v", e.Type, location, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s]%s %s", e.Type, location, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) WithContext(key string, value interface{}) *AppError {
	e.Context[key] = value
	return e
}

func (e *AppError) WithPath(path string) *AppError {
	e.Context["path"] = path
	return e
}

// WithWorker - добавляет ID воркера (в контекст)
func (e *AppError) WithWorker(id int) *AppError {
	e.Context["worker_id"] = id
	return e
}
