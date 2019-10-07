// +build darwin

package gsync

import "fmt"

// NewSemaphore creates a new named semaphore.
func NewSemaphore(name string) (Semaphore, error) {
	return nil, fmt.Errorf("not supported")
}
