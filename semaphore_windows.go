// +build windows

package gsync

import (
	"time"
)

type semaphore struct {
	handle int64
}

// NewSemaphore creates a new named semaphore with the specified initial
// and maximum counts.
func NewSemaphore(name string, initial, max int64) Semaphore {
	handle := createSemaphoreW(name, initial, max)

	if handle == -1 {
		return nil
	}

	return &semaphore{
		handle: handle,
	}
}

func (s *semaphore) Release(count int) {
	releaseSemaphore(s.handle, count)
}

func (s *semaphore) Wait(timeout time.Duration) error {
	return waitForSingleObject(s.handle, timeout.Milliseconds())
}

func (s *semaphore) Close() {
	closeHandle(s.handle)
}
