package semaphore

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestWorkerJob(t *testing.T) {
	tests := []struct {
		name          string
		semaphoreSize int
		jobCount      int
	}{
		{
			name:          "semaphore limits concurrent access",
			semaphoreSize: 3,
			jobCount:      10,
		},
		{
			name:          "single semaphore permit",
			semaphoreSize: 1,
			jobCount:      5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup = sync.WaitGroup{}
			w := &Worker{
				wg: &wg,
				semaphore: Semaphore{
					Channel: make(chan struct{}, tt.semaphoreSize),
				},
			}
			cj := NewCompletedJobs()
			for i := 0; i < tt.jobCount; i++ {
				w.wg.Add(1)
				cj.jobDone(i, w)
			}
			w.wg.Wait()
			if tt.jobCount != cj.Count() {
				t.Errorf("Runner() = %v, want %v", cj.Count(), tt.jobCount)
			}
		})
	}
}

type CompletedJobs struct {
	completed atomic.Int64
}

func NewCompletedJobs() *CompletedJobs {
	return &CompletedJobs{}
}

func (cj *CompletedJobs) jobDone(id int, worker *Worker) {
	go func() {
		worker.Job(id)
		cj.completed.Add(1)
	}()
}

func (cj *CompletedJobs) Count() int {
	return int(cj.completed.Load())
}
