package main

import (
	"fmt"
	mx "go-ipc/pkg/mutex"
	pd "go-ipc/pkg/prodcons"
	"strings"
)

func main() {
	inputData := []any{10, 10, 10, 31, 41}

	// Initialize producer and consumer
	producer := pd.NewProducer(inputData)
	consumer := pd.NewConsumer()

	// Create the service with injected dependencies
	service := pd.NewProdCons(producer, consumer)

	// Run the service
	resProdCons := service.Runner()

	// Print formatted results
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("  Producer-Consumer Workflow Results")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Input Data:      %v\n", inputData)
	fmt.Printf("Status:          %v\n", map[bool]string{true: "✓ Completed", false: "✗ Failed"}[resProdCons.IsDone])
	fmt.Printf("Items Consumed:  %d\n", len(resProdCons.Consumed))
	fmt.Printf("Consumed Data:   %v\n", resProdCons.Consumed)
	fmt.Println(strings.Repeat("=", 50) + "\n")

	// Initialize and run mutex service
	mutexService := mx.NewMutex()
	resMutex := mutexService.Runner()

	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("  Mutex Synchronization Results")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("✓ Completed %d concurrent counter increments\n", resMutex.FinalIncrement)
	fmt.Println(strings.Repeat("=", 50) + "\n")
}
