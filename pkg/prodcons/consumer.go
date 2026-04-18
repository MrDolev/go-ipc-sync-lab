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

// Consume reads data from the provided channel and appends it to results.
func (c *Consumer) Consume(ch <-chan any, done chan<- bool) {
	for value := range ch {
		log.Printf("consume data %v", value)
		c.results = append(c.results, value)
	}
	done <- true
}
