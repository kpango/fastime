package fastime

import (
	"bytes"
	"context"
	"math"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// Fastime is fastime's base struct, it's stores atomic time object
type Fastime struct {
	uut           uint32
	uunt          uint32
	dur           int64
	ut            int64
	unt           int64
	correctionDur time.Duration
	running       *atomic.Value
	t             *atomic.Value
	ft            *atomic.Value
	format        *atomic.Value
	pool          sync.Pool
	cancel        context.CancelFunc
}

var (
	once     sync.Once
	instance *Fastime
)

func init() {
	once.Do(func() {
		instance = New().StartTimerD(context.Background(), time.Millisecond*5)
	})
}

// New returns Fastime
func New() *Fastime {
	running := new(atomic.Value)
	running.Store(false)
	f := &Fastime{
		t:       new(atomic.Value),
		running: running,
		ut:      math.MaxInt64,
		unt:     math.MaxInt64,
		uut:     math.MaxUint32,
		uunt:    math.MaxUint32,
		format: func() *atomic.Value {
			av := new(atomic.Value)
			av.Store(time.RFC3339)
			return av
		}(),
		correctionDur: time.Millisecond * 100,
	}
	f.pool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, len(f.format.Load().(string))))
		},
	}
	f.ft = func() *atomic.Value {
		av := new(atomic.Value)
		buf := f.pool.Get().(*bytes.Buffer)
		buf.Reset()
		av.Store(buf.Bytes())
		f.pool.Put(buf)
		return av
	}()

	return f.refresh()
}

func (f *Fastime) update() *Fastime {
	return f.store(f.Now().Add(time.Duration(atomic.LoadInt64(&f.dur))))
}

func (f *Fastime) refresh() *Fastime {
	return f.store(f.now())
}

func (f *Fastime) store(t time.Time) *Fastime {
	f.t.Store(t)
	ut := t.Unix()
	unt := t.UnixNano()
	atomic.StoreInt64(&f.ut, ut)
	atomic.StoreInt64(&f.unt, unt)
	atomic.StoreUint32(&f.uut, *(*uint32)(unsafe.Pointer(&ut)))
	atomic.StoreUint32(&f.uunt, *(*uint32)(unsafe.Pointer(&unt)))
	form := f.format.Load().(string)
	buf := f.pool.Get().(*bytes.Buffer)
	if buf == nil || len(form) > buf.Cap() {
		buf.Grow(len(form))
	}
	f.ft.Store(t.AppendFormat(buf.Bytes()[:0], form))
	f.pool.Put(buf)
	return f
}

// SetFormat replaces time format
func SetFormat(format string) *Fastime {
	return instance.SetFormat(format)
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
	if f.running.Load().(bool) {
		f.cancel()
		atomic.StoreInt64(&f.dur, 0)
		return
	}
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
	if f.running.Load().(bool) {
		f.Stop()
	}
	f.refresh()

	var ct context.Context
	ct, f.cancel = context.WithCancel(ctx)
	go func() {
		f.running.Store(true)
		f.dur = math.MaxInt64
		atomic.StoreInt64(&f.dur, dur.Nanoseconds())
		ticker := time.NewTicker(time.Duration(atomic.LoadInt64(&f.dur)))
		ctick := time.NewTicker(f.correctionDur)
		for {
			select {
			case <-ct.Done():
				f.running.Store(false)
				ticker.Stop()
				ctick.Stop()
				return
			case <-ticker.C:
				f.update()
			case <-ctick.C:
				select {
				case <-ct.Done():
					f.running.Store(false)
					ticker.Stop()
					ctick.Stop()
					return
				case <-ticker.C:
					f.refresh()
				}
			}
		}
	}()

	return f
}
