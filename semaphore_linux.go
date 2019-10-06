// +build linux

package gsync

import (
	"syscall"
	"time"
	"unsafe"
)

type semaphoreLinux struct {
	handle int64
}

func NewSemaphore(name string, initial, max uint16) Semaphore {
	handle := createSemaphore(name)

	if handle == -1 {
		return nil
	}

	if initial > 0 {
		semaphoreOperation(handle, int16(initial))
	}

	return &semaphoreLinux{
		handle: handle,
	}
}

func (s *semaphoreLinux) Release(count uint16) {
	semaphoreOperation(s.handle, int16(count))
}

func (s *semaphoreLinux) Wait(timeout time.Duration) error {
	return semaphoreOperationTimed(s.handle, 0, timeout)
}

func (s *semaphoreLinux) Close() {

}

type semop struct {
	sem_num uint16
	sem_op  int16
	sem_flg int16
}

type timespec struct {
	tv_sec  int32
	tv_nsec int32
}

func createSemaphore(name string) int64 {
	namePtr := uintptr(0)
	if name != "" {
		b := append([]byte(name), 0)
		namePtr = uintptr(unsafe.Pointer(&b[0]))
	}

	ret, _, _ := syscall.Syscall(uintptr(syscall.SYS_SEMGET), namePtr, 1, 0)
	if ret == 0 {
		return -1
	}

	return int64(ret)
}

func semaphoreOperation(handle int64, op int16) error {
	ops := []semop{
		semop{
			sem_num: 0,
			sem_op:  op,
			sem_flg: 0,
		},
	}

	_, _, callErr := syscall.Syscall(uintptr(syscall.SYS_SEMOP), uintptr(handle), uintptr(unsafe.Pointer(&ops[0])), uintptr(1))
	if callErr != 0 {
		return callErr
	}

	return nil
}

func semaphoreOperationTimed(handle int64, op int16, timeout time.Duration) error {
	ops := []semop{
		semop{
			sem_num: 0,
			sem_op:  op,
			sem_flg: 0,
		},
	}

	to := &timespec{
		tv_sec:  int32(timeout.Truncate(time.Second).Seconds()),
		tv_nsec: int32((timeout - timeout.Truncate(time.Second)).Nanoseconds()),
	}

	_, _, callErr := syscall.Syscall6(uintptr(syscall.SYS_SEMOP), uintptr(handle), uintptr(unsafe.Pointer(&ops[0])), uintptr(1), uintptr(unsafe.Pointer(&to)), 0, 0)
	if callErr != 0 {
		return callErr
	}

	return nil
}
