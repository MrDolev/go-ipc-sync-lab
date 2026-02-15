package prodcons

import "log"

type Consumer struct {
	channel <-chan any
}

func (c *Consumer) consume(done chan<- bool) {
	for value := range c.channel {
		log.Printf("consume data %v", value)
	}
	done <- true
}
