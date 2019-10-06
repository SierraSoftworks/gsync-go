// +build linux

package gsync

import (
	"hash/fnv"
	"syscall"
	"time"
	"unsafe"
)

const (
	iCREAT  = 00001000
	iEXCL   = 00002000
	iNOWAIT = 00004000
	iUNDO   = 0x1000

	iGETVAL = 12
	iGETALL = 13
	iSETVAL = 16
	iSETALL = 17
)

type semaphoreLinux struct {
	handle int64
}

// NewSemaphoreWithValue creates a new named semaphore with the specified initial count.
func NewSemaphoreWithValue(name string, initial int) (Semaphore, error) {
	handle, err := createSemaphore(name, 1, iCREAT)

	if err != nil {
		return nil, err
	}

	semaphoreControl(handle, 0, iSETVAL, int(initial))

	return &semaphoreLinux{
		handle: handle,
	}, nil
}

// NewSemaphore creates a new named semaphore.
func NewSemaphore(name string) (Semaphore, error) {
	handle, err := createSemaphore(name, 1, iCREAT)

	if err != nil {
		return nil, err
	}

	return &semaphoreLinux{
		handle: handle,
	}, nil
}

func (s *semaphoreLinux) Release(count uint16) {
	if count <= 0 {
		return
	}

	semaphoreOperation(s.handle, 0, int16(count), 0)
}

func (s *semaphoreLinux) Wait(timeout time.Duration) error {
	return semaphoreOperationTimed(s.handle, -1, timeout)
}

func (s *semaphoreLinux) Close() {

}

type semop struct {
	semNum  uint16
	semOp   int16
	semFlag int16
}

type timespec struct {
	tvSec  int32
	tvNsec int32
}

func ftok(name string) int32 {
	h := fnv.New32a()
	h.Write([]byte(name))
	return int32(h.Sum32())
}

func createSemaphore(name string, num, flags int) (int64, error) {
	namePtr := uintptr(0)
	if name != "" {
		namePtr = uintptr(ftok(name))
	}

	ret, _, callErr := syscall.Syscall(uintptr(syscall.SYS_SEMGET), namePtr, uintptr(num), uintptr(flags|0660))
	if callErr != 0 {
		return -1, callErr
	}

	return int64(ret), nil
}

func semaphoreControl(handle int64, semNo, cmd, value int) (int64, error) {
	ret, _, callErr := syscall.Syscall6(uintptr(syscall.SYS_SEMCTL), uintptr(handle), uintptr(semNo), uintptr(cmd), uintptr(value), 0, 0)
	if callErr != 0 {
		return int64(ret), callErr
	}

	return int64(ret), nil
}

func semaphoreOperation(handle int64, semNum uint16, op, flag int16) error {
	ops := []semop{
		semop{
			semNum:  semNum,
			semOp:   op,
			semFlag: flag,
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
			semNum:  0,
			semOp:   op,
			semFlag: 0,
		},
	}

	to := &timespec{
		tvSec:  int32(timeout.Truncate(time.Second).Seconds()),
		tvNsec: int32((timeout - timeout.Truncate(time.Second)).Nanoseconds()),
	}

	_, _, callErr := syscall.Syscall6(uintptr(syscall.SYS_SEMOP), uintptr(handle), uintptr(unsafe.Pointer(&ops[0])), uintptr(1), uintptr(unsafe.Pointer(&to)), 0, 0)
	if callErr != 0 {
		return callErr
	}

	return nil
}
