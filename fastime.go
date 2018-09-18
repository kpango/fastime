package fastime

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// Fastime is fastime's base struct, it's stores atomic time object
type Fastime struct {
	t      atomic.Value
	cancel context.CancelFunc
}

var (
	once     sync.Once
	instance *Fastime
)

func init() {
	once.Do(func() {
		instance = New(context.Background())
	})
}

// New returns Fastime
func New(ctx context.Context) *Fastime {
	f := new(Fastime)
	f.t.Store(time.Now())
	var ct context.Context
	ct, f.cancel = context.WithCancel(ctx)
	go func() {
		for {
			select {
			case <-ct.Done():
				return
			default:
				f.t.Store(time.Now())
			}
		}
	}()
	return f
}

// Now returns current time
func Now() time.Time {
	return instance.Now()
}

// Stop stops stopping time refresh daemon
func Stop() {
	instance.Stop()
}

// Now returns current time
func (f *Fastime) Now() time.Time {
	return f.t.Load().(time.Time)
}

// Stop stops stopping time refresh daemon
func (f *Fastime) Stop() {
	f.cancel()
}
