// +build windows

package gsync

import (
	"os"
	"runtime"
	"time"
)

type mutexWindows struct {
	handle int64
}

// NewMutex constructs a new mutex with the provided name and initial state.
func NewMutex(name string, initial bool) Mutex {
	handle, err := createMutexW(name, initial)

	if err != nil && !os.IsExist(err) {
		return nil
	}

	return &mutexWindows{
		handle: handle,
	}
}

func (s *mutexWindows) Release() {
	if releaseMutex(s.handle) == nil {
		runtime.UnlockOSThread()
	}
}

func (s *mutexWindows) Wait(timeout time.Duration) error {
	runtime.LockOSThread()

	if err := waitForSingleObject(s.handle, timeout.Milliseconds()); err != nil {
		runtime.UnlockOSThread()
		return err
	}

	return nil
}

func (s *mutexWindows) Close() {
	closeHandle(s.handle)
}
