package prodcons

import "log"

type Consumer struct {
	results []any
}

func NewConsumer() *Consumer {
	return &Consumer{}
}

func (c *Consumer) Results() []any {
	return c.results
}

// Consume reads data from the provided channel until it is closed.
func (c *Consumer) Consume(ch <-chan any, done chan<- bool) {
	// The 'for range' loop automatically exits when the channel is closed by the producer.
	for value := range ch {
		log.Printf("CONSUMER: received item <- %v", value)
		c.results = append(c.results, value)
	}
	// Signal back to the orchestrator (Runner) that consumption is finished.
	done <- true
}
