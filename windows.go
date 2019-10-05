// +build windows

package gsync

import (
	"syscall"
	"unsafe"
)

var (
	kernel32, _               = syscall.LoadLibrary("kernel32.dll")
	createSemaphoreWHandle, _ = syscall.GetProcAddress(kernel32, "CreateSemaphoreW")
	releaseSemaphoreHandle, _ = syscall.GetProcAddress(kernel32, "ReleaseSemaphore")

	createMutexWHandle, _ = syscall.GetProcAddress(kernel32, "CreateMutexW")
	releaseMutexHandle, _ = syscall.GetProcAddress(kernel32, "ReleaseMutex")

	waitForSingleObjectHandle, _ = syscall.GetProcAddress(kernel32, "WaitForSingleObject")
	closeHandleHandle, _         = syscall.GetProcAddress(kernel32, "CloseHandle")
)

func createSemaphoreW(name string, initial, max int64) int64 {
	namePtr := uintptr(0)
	if name != "" {
		namePtr = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(name)))
	}

	ret, _, _ := syscall.Syscall6(uintptr(createSemaphoreWHandle), uintptr(4), 0, uintptr(initial), uintptr(max), namePtr, 0, 0)
	if ret == 0 {
		return -1
	}

	return int64(ret)
}

func releaseSemaphore(handle int64, count int) error {
	_, _, callErr := syscall.Syscall(uintptr(releaseSemaphoreHandle), uintptr(3), uintptr(handle), uintptr(count), 0)
	if callErr != 0 {
		return callErr
	}

	return nil
}

func createMutexW(name string, initial bool) int64 {
	namePtr := uintptr(0)
	if name != "" {
		namePtr = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(name)))
	}

	initialInt := 0
	if initial {
		initialInt = 1
	}

	ret, _, _ := syscall.Syscall(uintptr(createMutexWHandle), uintptr(3), 0, uintptr(initialInt), namePtr)
	if ret == 0 {
		return -1
	}

	return int64(ret)
}

func releaseMutex(handle int64) error {
	_, _, callErr := syscall.Syscall(uintptr(releaseMutexHandle), uintptr(1), uintptr(handle), 0, 0)
	if callErr != 0 {
		return callErr
	}

	return nil
}

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
