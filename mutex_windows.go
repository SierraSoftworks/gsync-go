// +build windows

package gsync

import (
	"runtime"
	"time"
)

type mutex struct {
	handle int64
}

// NewMutex constructs a new mutex with the provided name and initial state.
func NewMutex(name string, initial bool) Mutex {
	handle := createMutexW(name, initial)

	if handle == -1 {
		return nil
	}

	return &mutex{
		handle: handle,
	}
}

func (s *mutex) Release() {
	if releaseMutex(s.handle) == nil {
		runtime.UnlockOSThread()
	}
}

func (s *mutex) Wait(timeout time.Duration) error {
	runtime.LockOSThread()

	if err := waitForSingleObject(s.handle, timeout.Milliseconds()); err != nil {
		runtime.UnlockOSThread()
		return err
	}

	return nil
}

func (s *mutex) Close() {
	closeHandle(s.handle)
}
