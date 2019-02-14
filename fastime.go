package fastime

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// Fastime is fastime's base struct, it's stores atomic time object
type Fastime struct {
	running bool
	t       atomic.Value
	ut      int64
	uut     uint32
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
	tn := n.UnixNano()
	atomic.StoreInt64(&f.ut, tn)
	atomic.StoreUint32(&f.uut, *(*uint32)(unsafe.Pointer(&tn)))
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

// UnixNanoNow returns current unix nano time
func UnixNanoNow() int64 {
	return instance.UnixNanoNow()
}

// UnixNanoNow returns current unix nano time
func UnixUNanoNow() uint32 {
	return instance.UnixUNanoNow()
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

// UnixNanoNow returns current unix nano time
func (f *Fastime) UnixNanoNow() int64 {
	return atomic.LoadInt64(&f.ut)
}

// UnixNanoNow returns current unix nano time
func (f *Fastime) UnixUNanoNow() uint32 {
	return atomic.LoadUint32(&f.uut)
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
	atomic.StoreUint32(&f.uut, uint32(n.UnixNano()))
	go func() {
		f.running = true
		f.ticker = time.NewTicker(dur)
		for {
			select {
			case <-ct.Done():
				f.ticker.Stop()
				f.running = false
				return
			case <-f.ticker.C:
				n = time.Now()
				f.t.Store(n)
				atomic.StoreInt64(&f.ut, n.UnixNano())
				atomic.StoreUint32(&f.uut, uint32(n.UnixNano()))
			}
		}
	}()

	return f
}
