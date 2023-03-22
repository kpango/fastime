package fastime

import (
	"context"
	"math"
	"reflect"
	"sync/atomic"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "is daemon starts?",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := New().StartTimerD(context.Background(), 10000)
			time.Sleep(time.Second * 2)
			if (time.Now().Unix() - f.Now().Unix()) > 1000 {
				t.Error("time is not correct so daemon is not started")
			}
		})
	}
}

func TestNow(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "time equality",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if (time.Now().Unix() - Now().Unix()) > 1000 {
				t.Error("time is not correct")
			}
		})
	}
}

func TestStop(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "check stop",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := Now().Unix()
			if (time.Now().Unix() - now) > 1000 {
				t.Error("time is not correct")
			}
			Stop()
			time.Sleep(time.Second * 3)
			now = Now().Unix()
			if now == time.Now().Unix() {
				t.Error("refresh daemon stopped but time is same")
			}
		})
	}
}

func TestStartStop(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "check start and stop",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			dur := 10 * time.Millisecond
			f := New().StartTimerD(ctx, dur)
			if !f.IsDaemonRunning() {
				t.Error("daemon should be running")
			}
			for i := 0; i < 5; i++ {
				f.StartTimerD(ctx, dur)
				if !f.IsDaemonRunning() {
					t.Error("daemon should be running")
				}
				f.Stop()
				if f.IsDaemonRunning() {
					t.Error("daemon should not be running")
				}
			}
		})
	}
}

func TestFastime_Now(t *testing.T) {
	type fields struct {
		t      atomic.Value
		cancel context.CancelFunc
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		{
			name: "time equality",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := New().StartTimerD(context.Background(), 10000)
			if f.Now().Unix() != time.Now().Unix() {
				t.Error("time is not correct")
			}
		})
	}
}

func TestFastime_Stop(t *testing.T) {
	type fields struct {
		t      atomic.Value
		cancel context.CancelFunc
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "check stop",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := New().StartTimerD(context.Background(), time.Nanosecond*5)
			time.Sleep(time.Second)
			now := f.Now().Unix()
			if (time.Now().Unix() - now) > 1000 {
				t.Error("time is not correct")
			}
			f.Stop()
			time.Sleep(time.Second * 3)
			now = f.Now().Unix()
			if now == time.Now().Unix() {
				t.Error("refresh daemon stopped but time is same")
			}
		})
	}
}

func TestUnixNow(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "time equality",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if UnixNow() != Now().Unix() {
				t.Error("time is not correct")
			}
		})
	}
}

func TestFastime_UnixNow(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "time equality",
		},
	}

	f := New().StartTimerD(context.Background(), time.Millisecond)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if f.UnixNow() != f.Now().Unix() {
				t.Error("time is not correct")
			}
		})
	}
}

func TestUnixUNow(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "time equality",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if UnixUNow() != uint32(Now().Unix()) {
				t.Error("time is not correct")
			}
		})
	}
}

func TestFastime_UnixUNow(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "time equality",
		},
	}

	f := New().StartTimerD(context.Background(), time.Millisecond)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if f.UnixUNow() != uint32(f.Now().Unix()) {
				t.Error("time is not correct")
			}
		})
	}
}

func TestFastime_UnixNanoNow(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "time equality",
		},
	}

	f := New().StartTimerD(context.Background(), time.Nanosecond)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if f.UnixNanoNow() != f.Now().UnixNano() {
				t.Error("time is not correct")
			}
		})
	}
}

func TestUnixUNanoNow(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "time equality",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exp := UnixUNanoNow()
			act := uint32(Now().UnixNano())
			if exp != act {
				t.Errorf("time is not correct, exp: %v, actual: %v", exp, act)
			}
		})
	}
}

func TestFastime_UnixUNanoNow(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "time equality",
		},
	}

	f := New().StartTimerD(context.Background(), time.Nanosecond)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exp := f.UnixUNanoNow()
			act := uint32(f.Now().UnixNano())
			if exp != act {
				t.Errorf("time is not correct, exp: %v, actual: %v", exp, act)
			}
		})
	}
}

func TestFastime_refresh(t *testing.T) {
	tests := []struct {
		name string
		f    *fastime
	}{
		{
			name: "refresh",
			f:    newFastime(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.refresh(); time.Since(got.Now()) > time.Second {
				t.Errorf("time didn't refreshed Fastime.refresh() = %v", got.Now())
			}
		})
	}
}

func TestSetFormat(t *testing.T) {
	tests := []struct {
		name   string
		format string
	}{
		{
			name:   "set RFC3339",
			format: time.RFC3339,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetFormat(tt.format); !reflect.DeepEqual(got.GetFormat(), time.RFC3339) {
				t.Errorf("SetFormat() = %v, want %v", got.GetFormat(), time.RFC3339)
			}
		})
	}
}

func TestFastime_SetFormat(t *testing.T) {
	tests := []struct {
		name   string
		f      Fastime
		format string
	}{
		{
			name:   "set RFC3339",
			f:      New(),
			format: time.RFC3339,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.SetFormat(tt.format); !reflect.DeepEqual(got.GetFormat(), time.RFC3339) {
				t.Errorf("Fastime.SetFormat() = %v, want %v", got.GetFormat(), time.RFC3339)
			}
		})
	}
}

func TestFormattedNow(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "fetch",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(string(FormattedNow()))
		})
	}
}

func TestFastime_FormattedNow(t *testing.T) {
	tests := []struct {
		name string
		f    Fastime
	}{
		{
			name: "fetch",
			f:    New(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(string(tt.f.FormattedNow()))
		})
	}
}

func TestFastime_now(t *testing.T) {
	tests := []struct {
		name string
		f    *fastime
	}{
		{
			name: "now",
			f:    newFastime(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.now(); time.Since(got) > time.Second {
				t.Errorf("time didn't correct Fastime.now() = %v", got)
			}
		})
	}
}

func TestFastime_update(t *testing.T) {
	tests := []struct {
		name string
		f    *fastime
	}{
		{
			name: "update",
			f:    newFastime(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.refresh(); time.Since(got.Now()) > time.Second {
				t.Errorf("time didn't refreshed Fastime.update() = %v", got.Now())
			}
		})
	}
}

func TestFastime_store(t *testing.T) {
	tests := []struct {
		name string
		f    *fastime
	}{
		{
			name: "store",
			f:    newFastime(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := tt.f.now()
			if got := tt.f.store(n); tt.f.Now().UnixNano() != n.UnixNano() {
				t.Errorf("time didn't match Fastime.store() = %v", got.Now())
			}
		})
	}
}

func TestFastime_Since(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "since",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := New().StartTimerD(context.Background(), time.Millisecond*5)
			now := f.Now()
			timeNow := time.Now()
			time.Sleep(time.Second)
			since1 := f.Since(now)
			since2 := time.Since(timeNow)
			if since1 < 50*time.Millisecond {
				t.Errorf("since is not correct.\tfastime.Now: %v,\ttime.Now: %v\tsince1: %d, \tsince2: %d", now.UnixNano(), timeNow.UnixNano(), since1, since2)
			}
			if math.Abs(float64(since1-since2)) > float64(50*time.Millisecond) {
				t.Errorf("since error too large.\tfastime.Now: %v,\ttime.Now: %v\tsince1: %d, \tsince2: %d", now.UnixNano(), timeNow.UnixNano(), since1, since2)
			}
		})
	}
}
