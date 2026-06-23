// pkg/errors/errors.go
package errors

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"
)

// AppError - единый формат ошибок приложения
type AppError struct {
	// Время возникновения
	Time time.Time `json:"time"`

	// Тип ошибки
	Type string `json:"type"`

	// Сообщение об ошибке
	Message string `json:"message"`

	// Внутренняя ошибка
	Err error `json:"-"`

	// Дополнительный контекст
	Context map[string]interface{} `json:"context,omitempty"`

	// 👇 ИНФОРМАЦИЯ О ФАЙЛЕ - ДОБАВЛЯЕТСЯ АВТОМАТИЧЕСКИ
	File     string `json:"file,omitempty"`     // имя файла
	Line     int    `json:"line,omitempty"`     // номер строки
	Function string `json:"function,omitempty"` // имя функции
}

// Константы типов ошибок
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

// New - создает новую ошибку с АВТОМАТИЧЕСКИМ захватом файла
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
	// Захватываем информацию о файле (пропускаем 2 уровня: Wrap и вызывающий код)
	appErr.captureFileInfo(2)
	return appErr
}

// WrapWithFile - оборачивает ошибку с информацией о файле (устаревший, используйте Wrap)
func WrapWithFile(err error, errType, message string) *AppError {
	return Wrap(err, errType, message)
}

// captureFileInfo - ЗАХВАТЫВАЕТ информацию о файле, где произошла ошибка
func (e *AppError) captureFileInfo(skip int) {
	// Получаем информацию о caller
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		// Сохраняем только имя файла (без полного пути)
		e.File = filepath.Base(file)
		e.Line = line

		// Получаем имя функции
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			e.Function = fn.Name()
		}
	}
}

// Error - реализует интерфейс error с ОТОБРАЖЕНИЕМ информации о файле
func (e *AppError) Error() string {
	// Формируем информацию о местоположении
	location := ""
	if e.File != "" {
		location = fmt.Sprintf(" (%s:%d)", e.File, e.Line)
	}

	// Формируем сообщение
	if e.Err != nil {
		return fmt.Sprintf("[%s]%s %s: %v", e.Type, location, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s]%s %s", e.Type, location, e.Message)
}

// Unwrap - позволяет использовать errors.Is и errors.As
func (e *AppError) Unwrap() error {
	return e.Err
}

// WithContext - добавляет контекст к ошибке
func (e *AppError) WithContext(key string, value interface{}) *AppError {
	e.Context[key] = value
	return e
}

// WithPath - добавляет путь к файлу (в контекст)
func (e *AppError) WithPath(path string) *AppError {
	e.Context["path"] = path
	return e
}

// WithWorker - добавляет ID воркера (в контекст)
func (e *AppError) WithWorker(id int) *AppError {
	e.Context["worker_id"] = id
	return e
}
