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
	// Producer-Consumer Experiment
	start := time.Now()
	printExperimentStart("Producer-Consumer Workflow")

	inputData := []any{10, 10, 10, 31, 41}
	producer := pd.NewProducer(inputData)
	consumer := pd.NewConsumer()
	service := pd.NewProdCons(producer, consumer)
	resProdCons := service.Runner()

	printExperimentEnd("Producer-Consumer Workflow", time.Since(start))

	printSection("Producer-Consumer Workflow Results",
		"Input Data:", inputData,
		"Status:", map[bool]string{true: "✓ Completed", false: "✗ Failed"}[resProdCons.IsDone],
		"Items Consumed:", len(resProdCons.Consumed),
		"Consumed Data:", resProdCons.Consumed,
	)

	// Mutex Experiment
	start = time.Now()
	printExperimentStart("Mutex Synchronization")

	mutexService := mx.NewMutex()
	resMutex := mutexService.Runner()

	printExperimentEnd("Mutex Synchronization", time.Since(start))

	printSection("Mutex Synchronization Results",
		"✓ Completed counter increments:", resMutex.FinalIncrement,
	)

	// Semaphore Experiment
	start = time.Now()
	printExperimentStart("Semaphore Worker Pool")

	var wg sync.WaitGroup
	semaphore := sem.Semaphore{Channel: make(chan struct{}, 3)}
	worker := sem.NewWorker(&wg, semaphore)

	for i := range JOB_COUNT {
		wg.Add(1)
		go worker.Job(i)
	}
	wg.Wait()

	printExperimentEnd("Semaphore Worker Pool", time.Since(start))

	printSection("Semaphore Worker Pool Results",
		"Total Jobs:", JOB_COUNT,
		"Semaphore Limit:", "3 concurrent jobs",
		"Status:", "✓ All jobs completed successfully",
	)
}
