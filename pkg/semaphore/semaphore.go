package semaphore

type Semaphore struct {
	Channel chan struct{}
}

func (s *Semaphore) AcquireLock() {
	s.Channel <- struct{}{}
}

func (s *Semaphore) ReleaseLock() {
	<-s.Channel
}
