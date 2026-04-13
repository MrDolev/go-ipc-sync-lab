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

	fmt.Printf("processing job %d\n", id)
	time.Sleep(1 * time.Second)
	fmt.Printf("job done %d\n", id)

	defer w.semaphore.ReleaseLock()
}
