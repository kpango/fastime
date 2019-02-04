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
	ut      int64
	cancel  context.CancelFunc
	ticker  *time.Ticker
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
	n := time.Now()
	f.t.Store(n)
	atomic.StoreInt64(&f.ut, n.UnixNano())
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

// SetDuration changes time refresh duration
func SetDuration(dur time.Duration) *Fastime {
	return instance.SetDuration(dur)
}

// UnixNanoNow returns current unix nano time
func UnixNanoNow() int64 {
	return instance.UnixNanoNow()
}

// StartTimerD provides time refresh daemon
func StartTimerD(ctx context.Context, dur time.Duration) *Fastime {
	return instance.StartTimerD(ctx, dur)
}

// Now returns current time
func (f *Fastime) Now() time.Time {
	return f.t.Load().(time.Time)
}

// Stop stops stopping time refresh daemon
func (f *Fastime) Stop() {
	f.cancel()
}

// SetDuration changes time refresh duration
func (f *Fastime) SetDuration(dur time.Duration) *Fastime {
	if f.running && f.ticker != nil {
		f.ticker.Stop()
	}
	f.ticker = time.NewTicker(dur)
	return f
}

// UnixNanoNow returns current unix nano time
func (f *Fastime) UnixNanoNow() int64 {
	return atomic.LoadInt64(&f.ut)
}

// StartTimerD provides time refresh daemon
func (f *Fastime) StartTimerD(ctx context.Context, dur time.Duration) *Fastime {
	if f.running {
		f.Stop()
	}

	var ct context.Context
	ct, f.cancel = context.WithCancel(ctx)
	n := time.Now()
	f.t.Store(n)
	atomic.StoreInt64(&f.ut, n.UnixNano())
	go func() {
		f.running = true
		f.ticker = time.NewTicker(dur)
		for {
			select {
			case <-ct.Done():
				f.ticker.Stop()
				return
			case <-f.ticker.C:
				n = time.Now()
				f.t.Store(n)
				atomic.StoreInt64(&f.ut, n.UnixNano())
			}
		}
	}()

	return f
}
