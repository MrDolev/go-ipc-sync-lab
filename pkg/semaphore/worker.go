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

// Job represents a unit of work that requires a semaphore slot to run.
func (w *Worker) Job(id int) {
	// 1. Notify the WaitGroup when this goroutine exits.
	defer w.wg.Done()

	// 2. ACQUIRE: Try to send into the buffered channel.
	// If the channel is full (N slots occupied), this blocks until a slot is free.
	w.semaphore.AcquireLock()
	
	// 3. RELEASE: Ensure the slot is freed when the job is done.
	defer w.semaphore.ReleaseLock()

	start := time.Now()
	fmt.Printf("  [%2d] 🔄 slot acquired, starting work\n", id)
	time.Sleep(1 * time.Second) // Simulate heavy work
	duration := time.Since(start)
	fmt.Printf("  [%2d] ✅ job finished in %v, releasing slot\n", id, duration)
}
