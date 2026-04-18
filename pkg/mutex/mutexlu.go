package mutex

import (
	"log"
	"sync"
)

type ServiceMutexI interface {
	Runner() CounterRes
}

// Counter handles thread-safe counting operations
type Counter struct {
	mu    sync.Mutex
	value int
}

type CounterRes struct {
	FinalIncrement int
}

// Increment safely increments the counter
func (c *Counter) Increment(isDone chan bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
	log.Printf("counter updated: %d", c.value)
	isDone <- true
}

// Value safely returns the current counter value
func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// Mutex coordinates concurrent counter increments
type Mutex struct {
	counter *Counter
}

func NewMutex() ServiceMutexI {
	return &Mutex{
		counter: &Counter{},
	}
}

func (m *Mutex) Runner() CounterRes {
	chIsDone := make(chan bool)
	for i := 0; i < 100; i++ {
		go m.counter.Increment(chIsDone)
	}

	for i := 0; i < 100; i++ {
		<-chIsDone
	}

	return CounterRes{
		FinalIncrement: m.counter.Value(),
	}
}
