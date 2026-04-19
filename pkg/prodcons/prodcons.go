// Package prodcons demonstrates the classic Producer-Consumer synchronization pattern.
// It showcases how two independent processes (goroutines) can communicate safely
// through a shared channel, effectively decoupling the generation of data from its processing.
package prodcons
// ServiceRunnerI defines a generic interface for running services.
type ServiceRunnerI interface {
	// Runner starts the service and returns the results.
	Runner() ProdCondsRes
}

// ProdCons orchestrates the producer-consumer workflow.
// It acts as the "Mediator" that connects a producer to a consumer via a channel.
type ProdCons struct {
	records  []any
	producer ProducerI
	consumer ConsumerI
}

// ProdCondsRes represents the final report of the workflow execution.
type ProdCondsRes struct {
	// IsDone indicates if the workflow finished successfully.
	IsDone bool
	// Consumed contains all items processed by the consumer.
	Consumed []any
}

// ProducerI defines the interface for a data source.
// Producers don't need to know who is consuming the data; they just send to a channel.
type ProducerI interface {
	// Produce generates items and sends them to the channel.
	Produce(chan<- any)
}

// ConsumerI defines the interface for a data sink.
// Consumers don't need to know where the data comes from; they just read from a channel.
type ConsumerI interface {
	// Consume reads from the channel and signals completion on the done channel.
	Consume(<-chan any, chan<- bool)
	// Results returns the items consumed during the process.
	Results() []any
}

// NewProdCons creates a new ProdCons instance with injected dependencies.
// Dependency Injection allows for easier testing and modularity.
func NewProdCons(producer ProducerI, consumer ConsumerI) *ProdCons {
	return &ProdCons{
		producer: producer,
		consumer: consumer,
	}
}

// Runner executes the producer-consumer workflow.
// 1. It creates a 'sharedCh' for data transfer.
// 2. It creates a 'doneCh' to signal when the consumer is finished.
// 3. It launches both as concurrent goroutines.
func (pc *ProdCons) Runner() ProdCondsRes {
...
	sharedCh := make(chan any)
	doneCh := make(chan bool)

	// Start producer: it will send data to sharedCh and close it when done.
	go func() {
		pc.producer.Produce(sharedCh)
	}()

	// Start consumer: it will read from sharedCh until it is closed.
	go func() {
		pc.consumer.Consume(sharedCh, doneCh)
	}()

	// Wait for the completion signal from the consumer.
	isDone := <-doneCh

	return ProdCondsRes{
		IsDone:   isDone,
		Consumed: pc.consumer.Results(),
	}
}
