// +build linux

package gsync

import (
	"fmt"
	"time"
)

type mutexLinux struct {
	handle int64
}

// NewMutex constructs a new mutex with the provided name and initial state.
func NewMutex(name string, initial bool) Mutex {
	initCount := 1
	if initial {
		initCount = 0
	}

	handle, err := createSemaphore(name, 1, iCREAT)

	if err != nil {
		fmt.Println("could not create semaphore", err)
		return nil
	}

	val, err := semaphoreControl(handle, 0, iGETVAL, 0)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// If the semaphore has not yet been initialized then
	// set the initial count.
	if val == 0 {
		semaphoreControl(handle, 0, iSETVAL, int(initCount))
	}

	return &mutexLinux{
		handle: handle,
	}
}

func (m *mutexLinux) Wait(timeout time.Duration) error {
	return semaphoreOperationTimed(m.handle, -1, timeout)
}

func (m *mutexLinux) Release() {
	semaphoreOperation(m.handle, 0, 1, 0)
}

func (m *mutexLinux) Close() {

}
