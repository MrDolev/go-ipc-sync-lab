package prodcons

// ServiceRunnerI defines a generic interface for running services.
type ServiceRunnerI interface {
	Runner() ProdCondsRes
}

// ProdCons orchestrates the producer-consumer workflow.
type ProdCons struct {
	records  []any
	producer ProducerI
	consumer ConsumerI
}

// ProdCondsRes represents the result of the producer-consumer workflow.
type ProdCondsRes struct {
	IsDone   bool
	Consumed []any
}

// ProducerI defines the interface for a producer.
type ProducerI interface {
	Produce(chan<- any)
}

// ConsumerI defines the interface for a consumer.
type ConsumerI interface {
	Consume(<-chan any, chan<- bool)
	Results() []any
}

// NewProdCons creates a new ProdCons instance with injected dependencies.
func NewProdCons(producer ProducerI, consumer ConsumerI) *ProdCons {
	return &ProdCons{
		producer: producer,
		consumer: consumer,
	}
}

// Runner executes the producer-consumer workflow.
func (pc *ProdCons) Runner() ProdCondsRes {
	sharedCh := make(chan any)
	doneCh := make(chan bool)

	// Start producer in one goroutine.
	go func() {
		pc.producer.Produce(sharedCh)
	}()

	// Start consumer in another goroutine.
	go func() {
		pc.consumer.Consume(sharedCh, doneCh)
	}()

	isDone := <-doneCh

	return ProdCondsRes{
		IsDone:   isDone,
		Consumed: pc.consumer.Results(),
	}
}
