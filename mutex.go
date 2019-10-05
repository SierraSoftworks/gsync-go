package gsync

import (
	"time"
)

// A Mutex represents a global mutual exclusion lock
type Mutex interface {
	Wait(timeout time.Duration) error
	Release()
	Close()
}
