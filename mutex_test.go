package gsync_test

import (
	"testing"
	"time"

	"github.com/SierraSoftworks/gsync-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMutex(t *testing.T) {
	t.Run("Mutex Behaviour", func(t *testing.T) {
		m := gsync.NewMutex("", false)
		require.NotNil(t, m)

		defer m.Close()

		ch := make(chan struct{})
		ch2 := make(chan struct{})
		go func() {
			assert.NoError(t, m.Wait(200*time.Millisecond))
			ch <- struct{}{}
			<-ch2
			m.Release()
		}()

		select {
		case <-ch:
		case <-time.After(100 * time.Millisecond):
			t.Error("Timed out when waiting for initial wait")
		}

		go func() {
			assert.NoError(t, m.Wait(200*time.Millisecond))
			ch <- struct{}{}
		}()

		select {
		case <-ch:
			t.Error("Failed to wait upon locked mutex")
		case <-time.After(100 * time.Millisecond):
		}

		ch2 <- struct{}{}

		select {
		case <-ch:
		case <-time.After(100 * time.Millisecond):
			t.Error("Timed out when waiting for wait following release")
		}
	})

	t.Run("Shared Mutexes", func(t *testing.T) {
		m1 := gsync.NewMutex("test", false)
		require.NotNil(t, m1)
		defer m1.Close()

		m2 := gsync.NewMutex("test", false)
		require.NotNil(t, m2)
		defer m2.Close()

		ch := make(chan struct{})
		ch2 := make(chan struct{})
		go func() {
			assert.NoError(t, m1.Wait(200*time.Millisecond))
			ch <- struct{}{}
			<-ch2
			m1.Release()
		}()

		select {
		case <-ch:
		case <-time.After(100 * time.Millisecond):
			t.Error("Timed out when waiting for wait on mutex 1")
		}

		go func() {
			assert.NoError(t, m2.Wait(200*time.Millisecond))
			ch <- struct{}{}
		}()

		select {
		case <-ch:
			t.Error("Failed to wait upon locked mutex 1 with mutex 2")
		case <-time.After(100 * time.Millisecond):
		}

		ch2 <- struct{}{}

		select {
		case <-ch:
		case <-time.After(100 * time.Millisecond):
			t.Error("Timed out when waiting for mutex 2 following release")
		}
	})
}
