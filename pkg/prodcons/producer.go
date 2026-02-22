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
func (p *Producer) Produce(ch chan<- any) {
	for _, r := range p.records {
		log.Println("produce data from records", r)
		ch <- r
	}
	close(ch)
}
