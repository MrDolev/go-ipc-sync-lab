package main

import (
	"fmt"
	"strings"
	pd "go-ipc/pkg/prodcons"
)

func main() {
	inputData := []any{10, 10, 10, 31, 41}

	// Initialize producer and consumer
	producer := pd.NewProducer(inputData)
	consumer := pd.NewConsumer()

	// Create the service with injected dependencies
	service := pd.NewProdCons(producer, consumer)

	// Run the service
	res := service.Runner()

	// Print formatted results
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("  Producer-Consumer Workflow Results")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Input Data:      %v\n", inputData)
	fmt.Printf("Status:          %v\n", map[bool]string{true: "✓ Completed", false: "✗ Failed"}[res.IsDone])
	fmt.Printf("Items Consumed:  %d\n", len(res.Consumed))
	fmt.Printf("Consumed Data:   %v\n", res.Consumed)
	fmt.Println(strings.Repeat("=", 50) + "\n")
}
