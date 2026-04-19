package semaphore

import (
	"fmt"
	"sync"
	"time"
)

// Worker represents a process that performs tasks limited by a semaphore.
type Worker struct {
	wg        *sync.WaitGroup
	semaphore Semaphore
}

// NewWorker initializes a new Worker with a sync.WaitGroup and a Semaphore.
func NewWorker(wg *sync.WaitGroup, semaphore Semaphore) *Worker {
	return &Worker{
		wg:        wg,
		semaphore: semaphore,
	}
}

// Job represents a unit of work that requires a semaphore slot to run.
// It uses a semaphore to limit concurrency and a WaitGroup to signal completion.
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
