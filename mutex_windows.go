// +build windows

package gsync

import (
	"fmt"
	"os"
	"runtime"
	"syscall"
	"time"
	"unsafe"
)

var (
	createMutexWHandle, _ = syscall.GetProcAddress(kernel32, "CreateMutexW")
	releaseMutexHandle, _ = syscall.GetProcAddress(kernel32, "ReleaseMutex")
)

type mutexWindows struct {
	handle int64
}

// NewMutex constructs a new mutex with the provided name. This mutex is initially
// held by the creating thread.
func NewMutex(name string) (Mutex, error) {
	handle, err := createMutexW(name, true)

	if err != nil && !os.IsExist(err) {
		return nil, err
	}

	return &mutexWindows{
		handle: handle,
	}, err
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

func createMutexW(name string, initial bool) (int64, error) {
	namePtr := uintptr(0)
	if name != "" {
		namePtr = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(name)))
	}

	initialInt := 0
	if initial {
		initialInt = 1
	}

	ret, _, callErr := syscall.Syscall(uintptr(createMutexWHandle), uintptr(3), 0, uintptr(initialInt), namePtr)
	if callErr != 0 && !os.IsExist(callErr) {
		return -1, callErr
	}

	if ret == 0 {
		return -1, fmt.Errorf("unable to create mutex")
	}

	return int64(ret), nil
}

func releaseMutex(handle int64) error {
	_, _, callErr := syscall.Syscall(uintptr(releaseMutexHandle), uintptr(1), uintptr(handle), 0, 0)
	if callErr != 0 {
		return callErr
	}

	return nil
}
