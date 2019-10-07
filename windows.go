// +build windows

package gsync

import (
	"syscall"
)

var (
	kernel32, _ = syscall.LoadLibrary("kernel32.dll")

	waitForSingleObjectHandle, _ = syscall.GetProcAddress(kernel32, "WaitForSingleObject")
	closeHandleHandle, _         = syscall.GetProcAddress(kernel32, "CloseHandle")
)

func waitForSingleObject(handle int64, timeoutMilliseconds int64) error {
	_, _, callErr := syscall.Syscall(uintptr(waitForSingleObjectHandle), uintptr(2), uintptr(handle), uintptr(timeoutMilliseconds), 0)
	if callErr != 0 {
		return callErr
	}

	return nil
}

func closeHandle(handle int64) error {
	_, _, callErr := syscall.Syscall(uintptr(closeHandleHandle), uintptr(1), uintptr(handle), 0, 0)
	if callErr != 0 {
		return callErr
	}

	return nil
}
