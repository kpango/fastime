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
	t       *atomic.Value
	ut      int64
	unt     int64
	uut     uint32
	uunt    uint32
	ft      *atomic.Value
	format  *atomic.Value
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
	return (&Fastime{
		t:       new(atomic.Value),
		running: false,
		format: func() *atomic.Value {
			av := new(atomic.Value)
			av.Store(time.RFC3339)
			return av
		}(),
		ft: func() *atomic.Value {
			av := new(atomic.Value)
			av.Store(make([]byte, 0, len(time.RFC3339))[:0])
			return av
		}(),
	}).refresh()
}

func (f *Fastime) refresh() *Fastime {
	n := time.Now()
	f.t.Store(n)
	ut := n.Unix()
	unt := n.UnixNano()
	atomic.StoreInt64(&f.ut, ut)
	atomic.StoreInt64(&f.unt, unt)
	atomic.StoreUint32(&f.uut, *(*uint32)(unsafe.Pointer(&ut)))
	atomic.StoreUint32(&f.uunt, *(*uint32)(unsafe.Pointer(&unt)))
	form := f.format.Load().(string)
	f.ft.Store(n.AppendFormat(make([]byte, 0, len(form)), form))
	return f
}

// SetFormat replaces time format
func (f *Fastime) SetFormat(format string) *Fastime {
	f.format.Store(format)
	f.refresh()
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

// UnixNow returns current unix time
func UnixNow() int64 {
	return instance.UnixNow()
}

// UnixNow returns current unix time
func UnixUNow() uint32 {
	return instance.UnixUNow()
}

// UnixNanoNow returns current unix nano time
func UnixNanoNow() int64 {
	return instance.UnixNanoNow()
}

// UnixNanoNow returns current unix nano time
func UnixUNanoNow() uint32 {
	return instance.UnixUNanoNow()
}

// FormattedNow returns formatted byte time
func FormattedNow() []byte {
	return instance.FormattedNow()
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

// UnixNow returns current unix time
func (f *Fastime) UnixNow() int64 {
	return atomic.LoadInt64(&f.ut)
}

// UnixNow returns current unix time
func (f *Fastime) UnixUNow() uint32 {
	return atomic.LoadUint32(&f.uut)
}

// UnixNanoNow returns current unix nano time
func (f *Fastime) UnixNanoNow() int64 {
	return atomic.LoadInt64(&f.unt)
}

// UnixNanoNow returns current unix nano time
func (f *Fastime) UnixUNanoNow() uint32 {
	return atomic.LoadUint32(&f.uunt)
}

// FormattedNow returns formatted byte time
func (f *Fastime) FormattedNow() []byte {
	return f.ft.Load().([]byte)
}

// StartTimerD provides time refresh daemon
func (f *Fastime) StartTimerD(ctx context.Context, dur time.Duration) *Fastime {
	if f.running {
		f.Stop()
	}
	f.refresh()

	var ct context.Context
	ct, f.cancel = context.WithCancel(ctx)
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
				f.refresh()
			}
		}
	}()

	return f
}
