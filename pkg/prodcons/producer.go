package prodcons

import (
	"log"
)

type Producer struct {
	records []any
}

func (p *Producer) produce(ch chan<- any) {
	for _, r := range p.records {
		log.Println("produce data from records", r)
		ch <- r
	}
	close(ch)
}
