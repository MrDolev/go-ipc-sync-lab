package main

import (
	"fmt"
	mx "go-ipc/pkg/mutex"
	pd "go-ipc/pkg/prodcons"
	sem "go-ipc/pkg/semaphore"
	"strings"
	"sync"
	"time"
)

const JOB_COUNT int = 10
const REPEATS = 50

// printSection formats and prints a section with title and key-value pairs
func printSection(title string, pairs ...any) {
	fmt.Println("\n" + strings.Repeat("=", REPEATS))
	fmt.Println("  " + title)
	fmt.Println(strings.Repeat("=", REPEATS))

	for i := 0; i < len(pairs); i += 2 {
		key := pairs[i].(string)
		value := pairs[i+1]
		fmt.Printf("%s  %v\n", key, value)
	}

	fmt.Println(strings.Repeat("=", REPEATS) + "\n")
}

// printExperimentStart prints the beginning of an experiment
func printExperimentStart(name string) {
	fmt.Printf("\n▶ Starting: %s\n", name)
}

// printExperimentEnd prints the end of an experiment with duration
func printExperimentEnd(name string, duration time.Duration) {
	fmt.Printf("✓ Completed: %s (Duration: %v)\n\n", name, duration)
}

func main() {
	// -------------------------------------------------------------------------
	// 1. Producer-Consumer Workflow
	// -------------------------------------------------------------------------
	// Goal: Decouple data source (Producer) from processing (Consumer).
	start := time.Now()
	printExperimentStart("PATTERN: PRODUCER-CONSUMER (Pipeline)")

	inputData := []any{10, 20, 30, 41, 51}
	producer := pd.NewProducer(inputData)
	consumer := pd.NewConsumer()
	service := pd.NewProdCons(producer, consumer)
	
	// service.Runner() launches two goroutines connected by a channel.
	resProdCons := service.Runner()

	printExperimentEnd("PRODUCER-CONSUMER", time.Since(start))
	printSection("Results",
		"Input Items:", len(inputData),
		"Status:", map[bool]string{true: "✓ Processed all items", false: "✗ Failed"}[resProdCons.IsDone],
		"Items Received:", len(resProdCons.Consumed),
	)

	// -------------------------------------------------------------------------
	// 2. Mutex Synchronization
	// -------------------------------------------------------------------------
	// Goal: Prevent 'Race Conditions' when multiple goroutines write to the same var.
	start = time.Now()
	printExperimentStart("PATTERN: MUTEX (Mutual Exclusion)")

	mutexService := mx.NewMutex()
	// mutexService.Runner() launches 100 concurrent increments.
	resMutex := mutexService.Runner()

	printExperimentEnd("MUTEX", time.Since(start))
	printSection("Results",
		"Total Concurrent Increments:", 100,
		"Final Counter Value:", resMutex.FinalIncrement,
		"Status:", "✓ No data lost",
	)

	// -------------------------------------------------------------------------
	// 3. Semaphore Worker Pool
	// -------------------------------------------------------------------------
	// Goal: Limit the number of concurrent goroutines (e.g. rate limiting).
	start = time.Now()
	printExperimentStart("PATTERN: SEMAPHORE (Bounded Concurrency)")

	var wg sync.WaitGroup
	// Create a semaphore with capacity 3 (only 3 workers can run at once).
	semaphore := sem.Semaphore{Channel: make(chan struct{}, 3)}
	worker := sem.NewWorker(&wg, semaphore)

	for i := range JOB_COUNT {
		wg.Add(1)
		go worker.Job(i) // Try to run a job (may block if no slots in semaphore)
	}
	wg.Wait()

	printExperimentEnd("SEMAPHORE", time.Since(start))
	printSection("Results",
		"Total Jobs Dispatched:", JOB_COUNT,
		"Concurrency Limit (N):", "3 simultaneous workers",
		"Status:", "✓ Controlled execution finished",
	)
}
