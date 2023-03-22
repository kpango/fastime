package fastime

import (
	"context"
	"math"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"
)

type Fastime interface {
	IsDaemonRunning() bool
	GetFormat() string
	SetFormat(format string) Fastime
	GetLocation() *time.Location
	SetLocation(location *time.Location) Fastime
	Now() time.Time
	Stop()
	UnixNow() int64
	UnixUNow() uint32
	UnixNanoNow() int64
	UnixUNanoNow() uint32
	FormattedNow() []byte
	Since(t time.Time) time.Duration
	StartTimerD(ctx context.Context, dur time.Duration) Fastime
}

// Fastime is fastime's base struct, it's stores atomic time object
type fastime struct {
	uut           uint32
	uunt          uint32
	dur           int64
	ut            int64
	unt           int64
	correctionDur time.Duration
	mu            sync.Mutex
	wg            sync.WaitGroup
	running       atomic.Bool
	t             atomic.Pointer[time.Time]
	ft            atomic.Pointer[[]byte]
	format        atomic.Pointer[string]
	formatValid   atomic.Bool
	location      atomic.Pointer[time.Location]
}

const bufSize = 64

// New returns Fastime
func New() Fastime {
	return newFastime()
}

func newFastime() *fastime {
	f := &fastime{
		ut:   math.MaxInt64,
		unt:  math.MaxInt64,
		uut:  math.MaxUint32,
		uunt: math.MaxUint32,
		correctionDur: time.Millisecond * 100,
	}

	form := time.RFC3339
	f.format.Store(&form)
	loc := func() (loc *time.Location) {
			tz, ok := syscall.Getenv("TZ")
			if ok && tz != "" {
				var err error
				loc, err = time.LoadLocation(tz)
				if err == nil {
					return loc
				}
			}
			return new(time.Location)
		}()

	f.location.Store(loc)


	buf := f.newBuffer(len(f.GetFormat()) + 10)
	f.ft.Store(&buf)

	return f.refresh()
}

func (f *fastime) update() *fastime {
	return f.store(f.Now().Add(time.Duration(atomic.LoadInt64(&f.dur))))
}

func (f *fastime) refresh() *fastime {
	return f.store(f.now())
}

func (f *fastime) newBuffer(max int) (b []byte) {
	if max < bufSize {
		var buf [bufSize]byte
		b = buf[:0]
	} else {
		b = make([]byte, 0, max)
	}
	return b
}

func (f *fastime) store(t time.Time) *fastime {
	f.t.Store(&t)
	f.formatValid.Store(false)
	ut := t.Unix()
	unt := t.UnixNano()
	atomic.StoreInt64(&f.ut, ut)
	atomic.StoreInt64(&f.unt, unt)
	atomic.StoreUint32(&f.uut, *(*uint32)(unsafe.Pointer(&ut)))
	atomic.StoreUint32(&f.uunt, *(*uint32)(unsafe.Pointer(&unt)))
	return f
}

func (f *fastime) IsDaemonRunning() bool {
	return f.running.Load()
}

func (f *fastime) GetLocation() *time.Location {
	loc := f.location.Load()
	if loc == nil {
		return nil
	}
	return loc
}

func (f *fastime) GetFormat() string {
	return *f.format.Load()
}

// SetLocation replaces time location
func (f *fastime) SetLocation(location *time.Location) Fastime {
	if location == nil {
		return f
	}
	f.location.Store(location)
	f.refresh()
	return f
}

// SetFormat replaces time format
func (f *fastime) SetFormat(format string) Fastime {
	f.format.Store(&format)
	f.formatValid.Store(false)
	f.refresh()
	return f
}

// Now returns current time
func (f *fastime) Now() time.Time {
	return *f.t.Load()
}

// Stop stops stopping time refresh daemon
func (f *fastime) Stop() {
	f.mu.Lock()
	f.stop()
	f.mu.Unlock()
}

func (f *fastime) stop() {
	if f.IsDaemonRunning() {
		atomic.StoreInt64(&f.dur, 0)
	}
	f.wg.Wait()
}

func (f *fastime) Since(t time.Time) time.Duration {
	return f.Now().Sub(t)
}

// UnixNow returns current unix time
func (f *fastime) UnixNow() int64 {
	return atomic.LoadInt64(&f.ut)
}

// UnixNow returns current unix time
func (f *fastime) UnixUNow() uint32 {
	return atomic.LoadUint32(&f.uut)
}

// UnixNanoNow returns current unix nano time
func (f *fastime) UnixNanoNow() int64 {
	return atomic.LoadInt64(&f.unt)
}

// UnixNanoNow returns current unix nano time
func (f *fastime) UnixUNanoNow() uint32 {
	return atomic.LoadUint32(&f.uunt)
}

// FormattedNow returns formatted byte time
func (f *fastime) FormattedNow() []byte {
	// only update formatted value on swap
	if f.formatValid.CompareAndSwap(false, true) {
		form := f.GetFormat()
		buf := f.Now().AppendFormat(f.newBuffer(len(form)+10), form)
		f.ft.Store(&buf)
	}
	return *f.ft.Load()
}

// StartTimerD provides time refresh daemon
func (f *fastime) StartTimerD(ctx context.Context, dur time.Duration) Fastime {
	f.mu.Lock()
	defer f.mu.Unlock()
	// if the daemon was already running, restart
	if f.IsDaemonRunning() {
		f.stop()
	}
	f.running.Store(true)
	f.dur = math.MaxInt64
	atomic.StoreInt64(&f.dur, dur.Nanoseconds())
	ticker := time.NewTicker(time.Duration(atomic.LoadInt64(&f.dur)))
	lastCorrection := f.now()
	f.wg.Add(1)
	f.refresh()

	go func() {
		// daemon cleanup
		defer func() {
			f.running.Store(false)
			ticker.Stop()
			f.wg.Done()
		}()
		for atomic.LoadInt64(&f.dur) > 0 {
			t := <-ticker.C
			// rely on ticker for approximation
			if t.Sub(lastCorrection) < f.correctionDur {
				f.update()
			} else { // correct the system time at a fixed interval
				select {
				case <-ctx.Done():
					return
				default:
				}
				f.refresh()
				lastCorrection = t
			}
		}
	}()
	return f
}
