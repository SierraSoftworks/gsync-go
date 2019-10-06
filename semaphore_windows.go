// +build windows

package gsync

import (
	"math"
	"time"
)

type semaphoreWindows struct {
	handle int64
}

// NewSemaphore creates a new named semaphore with the specified initial
// and maximum counts.
func NewSemaphoreWithValue(name string, initial int) (Semaphore, error) {
	handle, err := createSemaphoreW(name, int64(initial), math.MaxInt32)

	if err != nil {
		return nil, err
	}

	return &semaphoreWindows{
		handle: handle,
	}, nil
}

// NewSemaphore creates a new named semaphore with the specified initial
// and maximum counts.
func NewSemaphore(name string) (Semaphore, error) {
	return NewSemaphoreWithValue(name, 0)
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
