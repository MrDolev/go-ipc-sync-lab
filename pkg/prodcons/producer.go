package prodcons

import (
	"log"
)

type Producer struct {
	records []any
}

func NewProducer(records []any) *Producer {
	return &Producer{
		records: records,
	}
}

// Produce sends records to the provided channel.
// It iterates through the input data and pushes each item onto the "conveyor belt".
func (p *Producer) Produce(ch chan<- any) {
	for _, r := range p.records {
		log.Printf("PRODUCER: sending item -> %v", r)
		ch <- r
	}
	// CRITICAL: Closing the channel signals the Consumer that no more data is coming.
	// Without this, the Consumer's 'range' loop would wait forever (deadlock).
	close(ch)
}
