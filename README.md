# Global Sync
**Global implementations of Go's `sync` package**

This library exposes global versions of Go's `sync` package primitives on platforms
which support them. 

## Mutex

```go
import (
    "fmt"

    "github.com/SierraSoftworks/gsync-go"
)

func main() {
    m, err := gsync.NewMutex("/gsync/example")
    if err != nil {
        fmt.Fatalf("failed to create mutex: %s\n", err)
    }

    defer m.Close()

    if m.Wait(50 * time.Millisecond) != nil {
        fmt.Println("You cannot run more than one instance of this application")
        return
    }

    fmt.Println("You are only running one instance of this app!")
}
```

## Semaphore

```go
import (
    "fmt"

    "github.com/SierraSoftworks/gsync-go"
)

func main() {
    s, err := gsync.NewSemaphore("/gsync/example")
    if err != nil {
        fmt.Fatalf("failed to create semaphore: %s\n", err)
    }

    if ss, ok := s.(gsync.Settable); ok {
        ss.Set(10)
    } else {
        ss.Release(10)
    }

    defer s.Release(1)
    defer s.Close()

    if s.Wait(50 * time.Milliseconds) != nil {
        fmt.Println("You cannot run more than 10 instances of this application")
        return
    }

    fmt.Println("You are running fewer than 10 instances of this application")
}
```