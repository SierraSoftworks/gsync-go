// +build windows

package gsync

import (
	"fmt"
	"math"
	"os"
	"syscall"
	"time"
	"unsafe"
)

var (
	createSemaphoreWHandle, _ = syscall.GetProcAddress(kernel32, "CreateSemaphoreW")
	releaseSemaphoreHandle, _ = syscall.GetProcAddress(kernel32, "ReleaseSemaphore")
)

type semaphoreWindows struct {
	handle int64
}

// NewSemaphore creates a new named semaphore.
func NewSemaphore(name string) (Semaphore, error) {
	handle, err := createSemaphoreW(name, 0, math.MaxInt16)

	if err != nil {
		return nil, err
	}

	return &semaphoreWindows{
		handle: handle,
	}, nil
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

func createSemaphoreW(name string, initial, max int64) (int64, error) {
	namePtr := uintptr(0)
	if name != "" {
		namePtr = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(name)))
	}

	ret, _, callErr := syscall.Syscall6(uintptr(createSemaphoreWHandle), uintptr(4), 0, uintptr(initial), uintptr(max), namePtr, 0, 0)
	if callErr != 0 && !os.IsExist(callErr) {
		return -1, callErr
	}

	if ret == 0 {
		return -1, fmt.Errorf("unable to create semaphore")
	}

	return int64(ret), nil
}

func releaseSemaphore(handle int64, count int) error {
	_, _, callErr := syscall.Syscall(uintptr(releaseSemaphoreHandle), uintptr(3), uintptr(handle), uintptr(count), 0)
	if callErr != 0 {
		return callErr
	}

	return nil
}
