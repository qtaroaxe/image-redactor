package pipeline

import (
	"fmt"
	"math/rand"
	"time"
)

// Task — задача на обработку изображения
type Task struct {
	// Идентификация
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`

	// Входные данные
	InputPath  string `json:"input_path"`
	OutputPath string `json:"output_path"`

	// Параметры обработки
	Format  string `json:"format"`  // jpeg, png, webp
	Width   int    `json:"width"`   // 0 = не менять
	Height  int    `json:"height"`  // 0 = не менять
	Quality int    `json:"quality"` // 1-100 (для JPEG)

	// Обрезка (опционально)
	CropX      int `json:"crop_x"`
	CropY      int `json:"crop_y"`
	CropWidth  int `json:"crop_width"`
	CropHeight int `json:"crop_height"`

	// Переименование
	NewFileName string `json:"new_file_name"`

	// Статус
	Status     string `json:"status"` // pending, processing, done, error
	ErrorMsg   string `json:"error_msg,omitempty"`
	ResultPath string `json:"result_path,omitempty"`
}

// TaskResult — результат для отправки клиенту
type TaskResult struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	OutputURL string `json:"output_url,omitempty"`
	Error     string `json:"error,omitempty"`
}

// NewTask — конструктор задачи
func NewTask(inputPath, outputPath string) *Task {
	return &Task{
		ID:         generateID(),
		CreatedAt:  time.Now(),
		InputPath:  inputPath,
		OutputPath: outputPath,
		Quality:    85,
		Status:     "pending",
	}
}

// generateID — уникальный ID задачи
func generateID() string {
	return fmt.Sprintf("%d_%d", time.Now().UnixNano(), rand.Intn(10000))
}
