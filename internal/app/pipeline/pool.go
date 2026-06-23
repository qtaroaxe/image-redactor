package pipeline

// import (
// 	"log"
// 	"sync"
// )

// // Pool — пул воркеров
// type Pool struct {
// 	taskQueue   chan Task       // Очередь задач
// 	workers     int             // Количество воркеров
// 	wg          sync.WaitGroup  // Ожидание завершения
// 	stopChan    chan struct{}   // Сигнал остановки
// 	mu          sync.RWMutex    // Защита статусов
// 	tasksStatus map[string]Task // Хранилище статусов задач
// }

// // NewPool — создаёт новый пул воркеров
// func NewPool(numWorkers, queueSize int) *Pool {
// 	return &Pool{
// 		taskQueue:   make(chan Task, queueSize),
// 		workers:     numWorkers,
// 		stopChan:    make(chan struct{}),
// 		tasksStatus: make(map[string]Task),
// 	}
// }

// // Start — запускает всех воркеров
// func (p *Pool) Start() {
// 	for i := 1; i <= p.workers; i++ {
// 		p.wg.Add(1)
// 		go p.worker(i)
// 	}
// 	log.Printf("✅ Запущено %d воркеров, очередь на %d задач", p.workers, cap(p.taskQueue))
// }

// // worker — горутина-воркер
// func (p *Pool) worker(id int) {
// 	defer p.wg.Done()

// 	for {
// 		select {
// 		case task := <-p.taskQueue:
// 			log.Printf("👷 Воркер %d: начал задачу %s (%s → %s)",
// 				id, task.ID, task.InputPath, task.OutputPath)

// 			// Обновляем статус
// 			p.updateTaskStatus(task.ID, "processing", "", "")

// 			// Обрабатываем изображение
// 			resultPath, err := ProcessImage(task)

// 			if err != nil {
// 				log.Printf("❌ Воркер %d: ошибка задачи %s: %v", id, task.ID, err)
// 				p.updateTaskStatus(task.ID, "error", err.Error(), "")
// 			} else {
// 				log.Printf("✅ Воркер %d: завершил задачу %s → %s", id, task.ID, resultPath)
// 				p.updateTaskStatus(task.ID, "done", "", resultPath)
// 			}

// 		case <-p.stopChan:
// 			log.Printf("🛑 Воркер %d: остановлен", id)
// 			return
// 		}
// 	}
// }

// // AddTask — добавляет задачу в очередь
// func (p *Pool) AddTask(task Task) string {
// 	p.updateTaskStatus(task.ID, "pending", "", "")
// 	p.taskQueue <- task
// 	log.Printf("📥 Задача %s добавлена в очередь (очередь: %d/%d)",
// 		task.ID, len(p.taskQueue), cap(p.taskQueue))
// 	return task.ID
// }

// // GetTaskStatus — возвращает статус задачи
// func (p *Pool) GetTaskStatus(taskID string) (Task, bool) {
// 	p.mu.RLock()
// 	defer p.mu.RUnlock()

// 	task, exists := p.tasksStatus[taskID]
// 	return task, exists
// }

// // updateTaskStatus — обновляет статус задачи
// func (p *Pool) updateTaskStatus(taskID, status, errorMsg, resultPath string) {
// 	p.mu.Lock()
// 	defer p.mu.Unlock()

// 	task, exists := p.tasksStatus[taskID]
// 	if !exists {
// 		task = Task{ID: taskID}
// 	}

// 	task.Status = status
// 	task.ErrorMsg = errorMsg
// 	task.ResultPath = resultPath
// 	p.tasksStatus[taskID] = task
// }

// // Stop — останавливает пул воркеров
// func (p *Pool) Stop() {
// 	log.Println("🛑 Остановка пула воркеров...")
// 	close(p.stopChan)
// 	close(p.taskQueue)
// 	p.wg.Wait()
// 	log.Println("✅ Пул воркеров остановлен")
// }

// // GetQueueLength — длина очереди
// func (p *Pool) GetQueueLength() int {
// 	return len(p.taskQueue)
// }

// // GetStats — статистика пула
// func (p *Pool) GetStats() map[string]interface{} {
// 	p.mu.RLock()
// 	defer p.mu.RUnlock()

// 	var pending, processing, done, failed int
// 	for _, task := range p.tasksStatus {
// 		switch task.Status {
// 		case "pending":
// 			pending++
// 		case "processing":
// 			processing++
// 		case "done":
// 			done++
// 		case "error":
// 			failed++
// 		}
// 	}

// 	return map[string]interface{}{
// 		"total_tasks":    len(p.tasksStatus),
// 		"pending":        pending,
// 		"processing":     processing,
// 		"done":           done,
// 		"failed":         failed,
// 		"queue_length":   len(p.taskQueue),
// 		"queue_capacity": cap(p.taskQueue),
// 		"workers":        p.workers,
// 	}
// }
