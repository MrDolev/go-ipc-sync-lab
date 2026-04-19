// Package mutex demonstrates the Mutual Exclusion pattern.
// It ensures that only one goroutine can access a shared resource at a time,
// preventing data races in a concurrent environment.
package mutex

import (
	"log"
	"sync"
)

// ServiceMutexI defines the contract for a mutex demonstration runner.
type ServiceMutexI interface {
	// Runner executes the concurrency simulation.
	Runner() CounterRes
}

// Counter is a shared resource that must be accessed safely.
// In Go, the idiomatic way to protect a struct is to include the Mutex
// directly in the struct that requires protection.
type Counter struct {
	mu    sync.Mutex // The lock that protects 'value'
	value int        // The protected shared state
}

// CounterRes holds the final state after all concurrent operations.
type CounterRes struct {
	// FinalIncrement is the value of the counter after all increments.
	FinalIncrement int
}

// Increment safely increments the counter.
// It marks the beginning and end of the "Critical Section" using mu.Lock() and mu.Unlock().
// 1. Lock(): Blocks other goroutines if the lock is already held.
// 2. defer Unlock(): Ensures the lock is released even if a panic occurs.
func (c *Counter) Increment(isDone chan bool) {
	c.mu.Lock()         // Start of Critical Section
	defer c.mu.Unlock() // End of Critical Section (scheduled to run at function exit)

	c.value++
	log.Printf("counter updated: %d", c.value)
	isDone <- true
}

// Value safely returns the current counter value by acquiring the lock.
// This prevents reading the value while another goroutine is mid-increment.
func (c *Counter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// Mutex implements the ServiceMutexI and coordinates the concurrent execution.
type Mutex struct {
	counter *Counter
}

// NewMutex initializes a new Mutex service with a clean counter.
func NewMutex() ServiceMutexI {
	return &Mutex{
		counter: &Counter{},
	}
}

// Runner simulates 100 concurrent goroutines trying to increment the same counter.
// Without a Mutex, this would result in a "Race Condition" where the final value
// would likely be less than 100 because updates would overwrite each other.
func (m *Mutex) Runner() CounterRes {
	chIsDone := make(chan bool)

	// Launch 100 "workers" concurrently.
	for range 100 {
		go m.counter.Increment(chIsDone)
	}

	// Synchronization point: wait for all 100 signals.
	for range 100 {
		<-chIsDone
	}

	return CounterRes{
		FinalIncrement: m.counter.Value(),
	}
}
