// +build linux

package gsync

import (
	"time"
)

type mutexLinux struct {
	handle int64
}

// NewMutex constructs a new mutex with the provided name and initial state. This mutex is initially
// held.
func NewMutex(name string) (Mutex, error) {
	handle, err := createSemaphore(name, 1, iCREAT)

	if err != nil {
		return nil, err
	}

	return &mutexLinux{
		handle: handle,
	}, nil
}

func (m *mutexLinux) Wait(timeout time.Duration) error {
	return semaphoreOperationTimed(m.handle, -1, timeout)
}

func (m *mutexLinux) Release() {
	semaphoreOperation(m.handle, 0, 1, iNOWAIT)
}

func (m *mutexLinux) Close() {

}
