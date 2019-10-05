package gsync_test

import (
	"testing"
	"time"

	"github.com/SierraSoftworks/gsync-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSemaphore(t *testing.T) {
	t.Run("Semaphore Behaviour", func(t *testing.T) {
		s := gsync.NewSemaphore("", 1, 1)
		require.NotNil(t, s)

		defer s.Close()

		ch := make(chan struct{})
		go func() {
			assert.NoError(t, s.Wait(200*time.Millisecond))
			ch <- struct{}{}
		}()

		select {
		case <-ch:
		case <-time.After(100 * time.Millisecond):
			t.Error("Timed out when waiting for initial wait")
		}

		go func() {
			assert.NoError(t, s.Wait(200*time.Millisecond))
			ch <- struct{}{}
		}()

		select {
		case <-ch:
			t.Error("Failed to wait upon locked semaphore")
		case <-time.After(100 * time.Millisecond):
		}

		s.Release(1)

		select {
		case <-ch:
		case <-time.After(100 * time.Millisecond):
			t.Error("Timed out when waiting for wait following release")
		}
	})

	t.Run("Shared Semaphores", func(t *testing.T) {
		s1 := gsync.NewSemaphore("/test", 1, 1)
		require.NotNil(t, s1)
		defer s1.Close()

		s2 := gsync.NewSemaphore("/test", 0, 1)
		require.NotNil(t, s2)
		defer s2.Close()

		ch := make(chan struct{})
		go func() {
			assert.NoError(t, s1.Wait(200*time.Millisecond))
			ch <- struct{}{}
		}()

		select {
		case <-ch:
		case <-time.After(100 * time.Millisecond):
			t.Error("Timed out when waiting for initial wait on semaphore 1")
		}

		go func() {
			assert.NoError(t, s2.Wait(200*time.Millisecond))
			ch <- struct{}{}
		}()

		select {
		case <-ch:
			t.Error("Failed to wait upon locked semaphore 2")
		case <-time.After(100 * time.Millisecond):
		}

		s1.Release(1)

		select {
		case <-ch:
		case <-time.After(100 * time.Millisecond):
			t.Error("Timed out when waiting for wait following release of semaphore 1")
		}
	})
}
