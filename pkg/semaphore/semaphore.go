// Package semaphore implements a counting semaphore using Go's buffered channels.
// A counting semaphore is used to control access to a pool of resources or
// to limit the number of concurrent operations.
package semaphore

// Semaphore provides a mechanism to limit concurrent access.
// It uses an empty struct channel 'chan struct{}' because it occupies zero memory.
type Semaphore struct {
	// Channel capacity represents the maximum number of concurrent workers allowed.
	Channel chan struct{}
}

// AcquireLock occupies a slot in the semaphore.
// If the buffered channel is full, this call blocks until another goroutine calls ReleaseLock().
// This is equivalent to the 'P' (proberen/test) or 'wait' operation in classic semaphore theory.
func (s *Semaphore) AcquireLock() {
	s.Channel <- struct{}{}
}

// ReleaseLock frees a slot in the semaphore.
// By receiving from the channel, it allows one more blocked AcquireLock() to proceed.
// This is equivalent to the 'V' (verhogen/increment) or 'signal' operation in classic semaphore theory.
func (s *Semaphore) ReleaseLock() {
	<-s.Channel
}
