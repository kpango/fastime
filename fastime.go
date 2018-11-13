package fastime

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// Fastime is fastime's base struct, it's stores atomic time object
type Fastime struct {
	running bool
	t       atomic.Value
	cancel  context.CancelFunc
}

var (
	once     sync.Once
	instance *Fastime
)

func init() {
	once.Do(func() {
		instance = New().StartTimerD(context.Background(), time.Millisecond*100)
	})
}

// New returns Fastime
func New() *Fastime {
	f := new(Fastime)
	f.t.Store(time.Now())
	return f
}

func StartTimerD(ctx context.Context, dur time.Duration) *Fastime {
	return instance.StartTimerD(ctx, dur)
}

func (f *Fastime) StartTimerD(ctx context.Context, dur time.Duration) *Fastime {
	if f.running {
		f.Stop()
	}

	var ct context.Context
	ct, f.cancel = context.WithCancel(ctx)

	f.t.Store(time.Now())

	go func() {
		f.running = true
		ticker := time.NewTicker(dur)
		for {
			select {
			case <-ct.Done():
				return
			case <-ticker.C:
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
