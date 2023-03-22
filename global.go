package fastime

import (
	"context"
	"sync"
	"time"
)

var (
	once     sync.Once
	instance Fastime
)

func init() {
	once.Do(func() {
		instance = New().StartTimerD(context.Background(), time.Millisecond*5)
	})
}

func IsDaemonRunning() (running bool) {
	return instance.IsDaemonRunning()
}

func GetLocation() (loc *time.Location) {
	return instance.GetLocation()
}

func GetFormat() (form string) {
	return instance.GetFormat()
}

// SetLocation replaces time location
func SetLocation(location *time.Location) (ft Fastime) {
	return instance.SetLocation(location)
}

// SetFormat replaces time format
func SetFormat(format string) (ft Fastime) {
	return instance.SetFormat(format)
}

// Now returns current time
func Now() (now time.Time) {
	return instance.Now()
}

// Since returns the time elapsed since t.
// It is shorthand for fastime.Now().Sub(t).
func Since(t time.Time) (dur time.Duration) {
	return instance.Since(t)
}

// Stop stops stopping time refresh daemon
func Stop() {
	instance.Stop()
}

// UnixNow returns current unix time
func UnixNow() (now int64) {
	return instance.UnixNow()
}

// UnixUNow returns current unix time
func UnixUNow() (now uint32) {
	return instance.UnixUNow()
}

// UnixNanoNow returns current unix nano time
func UnixNanoNow() (now int64) {
	return instance.UnixNanoNow()
}

// UnixUNanoNow returns current unix nano time
func UnixUNanoNow() (now uint32) {
	return instance.UnixUNanoNow()
}

// FormattedNow returns formatted byte time
func FormattedNow() (now []byte) {
	return instance.FormattedNow()
}

// StartTimerD provides time refresh daemon
func StartTimerD(ctx context.Context, dur time.Duration) (ft Fastime) {
	return instance.StartTimerD(ctx, dur)
}
