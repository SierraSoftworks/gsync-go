// +build darwin

package gsync

import "fmt"

// NewMutex creates a new named mutex.
func NewMutex(name string) (Mutex, error) {
	return nil, fmt.Errorf("not supported")
}
