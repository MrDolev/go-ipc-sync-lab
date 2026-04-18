package semaphore

import (
	"fmt"
	"sync"
	"time"
)

type Worker struct {
	wg        *sync.WaitGroup
	semaphore Semaphore
}

func NewWorker(wg *sync.WaitGroup, semaphore Semaphore) *Worker {
	return &Worker{
		wg:        wg,
		semaphore: semaphore,
	}
}

func (w *Worker) Job(id int) {
	defer w.wg.Done()

	w.semaphore.AcquireLock()
	defer w.semaphore.ReleaseLock()

	start := time.Now()
	fmt.Printf("  [%2d] 🔄 acquiring semaphore lock\n", id)
	fmt.Printf("  [%2d] ▶️  processing job\n", id)
	time.Sleep(1 * time.Second)
	duration := time.Since(start)
	fmt.Printf("  [%2d] ✅ completed in %v\n", id, duration)
}
