// internal/service/compressor/pool.go
package pipeline

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	apperrors "github.com/QtaroAXE/image-redactor/internal/domain/errors"
	"github.com/QtaroAXE/image-redactor/internal/domain/imginfo"
	compressor "github.com/QtaroAXE/image-redactor/internal/infra/codec"
	"github.com/QtaroAXE/image-redactor/internal/infra/fs"
)

// WorkerPool - пул воркеров для многопоточного сжатия
type WorkerPool struct {
	// Конфигурация
	workerCount int
	batchSize   int

	// Сервисы
	compressor *compressor.CompressorService
	fs         *fs.FileSystem

	// Каналы
	jobQueue    chan Job
	resultQueue chan Result

	// Управление
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// Статистика
	stats   *PoolStats
	statsMu sync.RWMutex

	// Обработка ошибок
	errorHandler func(error)
}

// Job - задача на сжатие
type Job struct {
	ID         string
	SourcePath string
	Source     imginfo.SourceImage
	Target     imginfo.TargetImage
	Priority   int
	CreatedAt  time.Time
	RetryCount int
}

// Result - результат сжатия
type Result struct {
	JobID       string
	SourcePath  string
	OutputPath  string
	Error       error
	Duration    time.Duration
	SizeBefore  int64
	SizeAfter   int64
	Success     bool
	ProcessedAt time.Time
}

// PoolStats - статистика пула
type PoolStats struct {
	TotalJobs     int64
	CompletedJobs int64
	FailedJobs    int64
	RetriedJobs   int64
	ActiveWorkers int32
	QueueLength   int32
	TotalDuration time.Duration
	StartTime     time.Time
}

// PoolConfig - конфигурация пула
type PoolConfig struct {
	WorkerCount  int
	BatchSize    int
	QueueSize    int
	ErrorHandler func(error)
}

// NewWorkerPool - создает новый пул воркеров
func NewWorkerPool(
	compressor *compressor.CompressorService,
	fs *fs.FileSystem,
	cfg PoolConfig,
) *WorkerPool {
	if cfg.WorkerCount <= 0 {
		cfg.WorkerCount = 4
	}
	if cfg.BatchSize <= 0 {
		cfg.BatchSize = 10
	}
	if cfg.QueueSize <= 0 {
		cfg.QueueSize = cfg.WorkerCount * 2
	}

	ctx, cancel := context.WithCancel(context.Background())

	pool := &WorkerPool{
		workerCount:  cfg.WorkerCount,
		batchSize:    cfg.BatchSize,
		compressor:   compressor,
		fs:           fs,
		jobQueue:     make(chan Job, cfg.QueueSize),
		resultQueue:  make(chan Result, cfg.QueueSize),
		ctx:          ctx,
		cancel:       cancel,
		stats:        &PoolStats{StartTime: time.Now()},
		errorHandler: cfg.ErrorHandler,
	}

	return pool
}

// Start - запускает пул воркеров
func (p *WorkerPool) Start() {
	log.Printf("Starting worker pool with %d workers", p.workerCount)

	// Запускаем воркеров
	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}

	// Запускаем обработчик результатов
	p.wg.Add(1)
	go p.resultProcessor()

	// Запускаем мониторинг
	p.wg.Add(1)
	go p.monitor()
}

// Stop - останавливает пул воркеров
func (p *WorkerPool) Stop() {
	log.Println("Stopping worker pool...")
	p.cancel()
	close(p.jobQueue)
	p.wg.Wait()
	close(p.resultQueue)
	log.Println("Worker pool stopped")
}

// AddJob - добавляет задачу в очередь
func (p *WorkerPool) AddJob(source imginfo.SourceImage, target imginfo.TargetImage) error {
	select {
	case <-p.ctx.Done():
		return apperrors.New(
			apperrors.TypeInternal,
			"pool is stopped",
		)
	default:
		job := Job{
			ID:         fmt.Sprintf("job_%d", time.Now().UnixNano()),
			SourcePath: source.Path(),
			Source:     source,
			Target:     target,
			CreatedAt:  time.Now(),
			RetryCount: 0,
		}

		atomic.AddInt64(&p.stats.TotalJobs, 1)
		atomic.AddInt32(&p.stats.QueueLength, 1)

		p.jobQueue <- job
		return nil
	}
}

// AddJobs - добавляет несколько задач в очередь
func (p *WorkerPool) AddJobs(jobs []Job) error {
	for _, job := range jobs {
		if err := p.AddJob(job.Source, job.Target); err != nil {
			return err
		}
	}
	return nil
}

// GetStats - возвращает статистику пула
func (p *WorkerPool) GetStats() PoolStats {
	p.statsMu.RLock()
	defer p.statsMu.RUnlock()

	stats := *p.stats
	stats.QueueLength = atomic.LoadInt32(&p.stats.QueueLength)
	stats.ActiveWorkers = atomic.LoadInt32(&p.stats.ActiveWorkers)
	return stats
}

// Wait - ожидает завершения всех задач
func (p *WorkerPool) Wait() {
	// Ждем пока очередь опустеет
	for {
		if len(p.jobQueue) == 0 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	// Даем время на завершение текущих задач
	time.Sleep(time.Second)
}

// worker - горутина воркера
func (p *WorkerPool) worker(id int) {
	defer p.wg.Done()

	log.Printf("Worker %d started", id)

	for {
		select {
		case <-p.ctx.Done():
			log.Printf("Worker %d stopping", id)
			return

		case job, ok := <-p.jobQueue:
			if !ok {
				log.Printf("Worker %d: job queue closed", id)
				return
			}

			atomic.AddInt32(&p.stats.ActiveWorkers, 1)
			atomic.AddInt32(&p.stats.QueueLength, -1)

			result := p.processJob(id, job)
			p.resultQueue <- result

			atomic.AddInt32(&p.stats.ActiveWorkers, -1)
		}
	}
}

// processJob - обрабатывает одну задачу
func (p *WorkerPool) processJob(workerID int, job Job) Result {
	startTime := time.Now()

	result := Result{
		JobID:       job.ID,
		SourcePath:  job.Source.Path(),
		ProcessedAt: time.Now(),
	}

	// Получаем размер исходного файла
	sizeBefore, err := p.fs.GetFileSize(job.Source.Path())
	if err != nil {
		result.Error = apperrors.Wrap(
			err,
			apperrors.TypeIO,
			"failed to get source file size",
		).WithWorker(workerID)
		result.Success = false
		atomic.AddInt64(&p.stats.FailedJobs, 1)
		return result
	}
	result.SizeBefore = sizeBefore

	// Проверяем существование файла
	if !p.fs.FileExists(job.Source.Path()) {
		err := apperrors.New(
			apperrors.TypeNotFound,
			"source file does not exist",
		).WithPath(job.Source.Path()).WithWorker(workerID)

		result.Error = err
		result.Success = false
		atomic.AddInt64(&p.stats.FailedJobs, 1)

		// Перемещаем в error директорию
		p.fs.MoveToError(job.Source.Path())
		return result
	}

	// Выполняем сжатие
	err = p.compressor.CompressImage(job.Source, job.Target)
	if err != nil {
		// Добавляем информацию о воркере
		if appErr, ok := err.(*apperrors.AppError); ok {
			appErr.WithWorker(workerID)
		}

		result.Error = err
		result.Success = false
		atomic.AddInt64(&p.stats.FailedJobs, 1)

		// Если ошибка временная - пробуем повторить
		if job.RetryCount < 3 && p.isRetryableError(err) {
			job.RetryCount++
			atomic.AddInt64(&p.stats.RetriedJobs, 1)

			// Повторная попытка с задержкой
			backoff := time.Duration(job.RetryCount) * time.Second
			time.Sleep(backoff)

			log.Printf("Retrying job %s (attempt %d)", job.ID, job.RetryCount)
			p.jobQueue <- job
			return result
		}

		// Перемещаем в error директорию
		p.fs.MoveToError(job.Source.Path())
		return result
	}

	// Получаем размер выходного файла
	sizeAfter, err := p.fs.GetFileSize(job.Target.Path())
	if err == nil {
		result.SizeAfter = sizeAfter
	}

	result.OutputPath = job.Target.Path()
	result.Success = true
	result.Duration = time.Since(startTime)

	atomic.AddInt64(&p.stats.CompletedJobs, 1)

	// Обновляем общую длительность
	p.statsMu.Lock()
	p.stats.TotalDuration += result.Duration
	p.statsMu.Unlock()

	// Перемещаем исходный файл в processed
	if err := p.fs.MoveToProcessed(job.Source.Path()); err != nil {
		log.Printf("Failed to move processed file: %v", err)
	}

	return result
}

// resultProcessor - обрабатывает результаты
func (p *WorkerPool) resultProcessor() {
	defer p.wg.Done()

	for result := range p.resultQueue {
		if result.Error != nil {
			if p.errorHandler != nil {
				p.errorHandler(result.Error)
			}
			log.Printf("Job %s failed: %v", result.JobID, result.Error)
		} else {
			log.Printf("Job %s completed in %v (%.2f%% reduction)",
				result.JobID,
				result.Duration,
				p.calculateReduction(result.SizeBefore, result.SizeAfter),
			)
		}
	}
}

// monitor - мониторинг пула
func (p *WorkerPool) monitor() {
	defer p.wg.Done()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker.C:
			stats := p.GetStats()
			log.Printf("Pool stats: Jobs=%d, Completed=%d, Failed=%d, Active=%d, Queue=%d",
				stats.TotalJobs,
				stats.CompletedJobs,
				stats.FailedJobs,
				stats.ActiveWorkers,
				stats.QueueLength,
			)
		}
	}
}

// isRetryableError - проверяет, можно ли повторить задачу
func (p *WorkerPool) isRetryableError(err error) bool {
	if appErr, ok := err.(*apperrors.AppError); ok {
		switch appErr.Type {
		case apperrors.TypeIO:
			return true
		case apperrors.TypeTimeout:
			return true
		case apperrors.TypeInternal:
			return true
		default:
			return false
		}
	}
	return false
}

// calculateReduction - вычисляет процент сжатия
func (p *WorkerPool) calculateReduction(before, after int64) float64 {
	if before == 0 {
		return 0
	}
	return (1 - float64(after)/float64(before)) * 100
}
