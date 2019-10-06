// +build windows

package gsync

import (
	"time"
)

type semaphoreWindows struct {
	handle int64
}

// NewSemaphore creates a new named semaphore with the specified initial
// and maximum counts.
func NewSemaphore(name string, initial, max uint16) Semaphore {
	handle := createSemaphoreW(name, int64(initial), int64(max))

	if handle == -1 {
		return nil
	}

	return &semaphoreWindows{
		handle: handle,
	}
}

func (s *semaphoreWindows) Release(count uint16) {
	releaseSemaphore(s.handle, int(count))
}

func (s *semaphoreWindows) Wait(timeout time.Duration) error {
	return waitForSingleObject(s.handle, timeout.Milliseconds())
}

func (s *semaphoreWindows) Close() {
	closeHandle(s.handle)
}
