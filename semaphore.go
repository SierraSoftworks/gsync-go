package gsync

import (
	"time"
)

// A Semaphore is a synchronization object which blocks waiting
// tasks when its count reaches zero and allows any thread to increment
// the count up to a pre-configured maximum value.
type Semaphore interface {
	Release(count uint16)
	Wait(timeout time.Duration) error
	Close()
}

// Settable semaphores allow you to adjust the value of the semaphore directly.
type Settable interface {
	Set(count uint16) error
}
