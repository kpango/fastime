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
	ticker *time.Ticker
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
	f.ticker = time.NewTicker(time.Millisecond * 100)
	go func() {
		for {
			select {
			case <-ct.Done():
				f.ticker.Stop()
				return
			case <-f.ticker.C:
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

// SetDuration changes time refresh duration
func SetDuration(dur time.Duration) *Fastime {
	return instance.SetDuration(dur)
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
	f.ticker.Stop()
	f.ticker = time.NewTicker(dur)
	return f
}
