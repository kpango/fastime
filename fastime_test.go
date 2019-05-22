package fastime

import (
	"context"
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
			if f.Now().Unix() != time.Now().Unix() {
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
			if Now().Unix() != time.Now().Unix() {
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
			if now != time.Now().Unix() {
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
			f := New().StartTimerD(context.Background(), 10000)
			now := f.Now().Unix()
			if now != time.Now().Unix() {
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
			if UnixUNanoNow() != uint32(Now().UnixNano()) {
				t.Error("time is not correct")
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
			if f.UnixUNanoNow() != uint32(f.Now().UnixNano()) {
				t.Error("time is not correct")
			}
		})
	}
}

func TestFastime_refresh(t *testing.T) {
	tests := []struct {
		name string
		f    *Fastime
	}{
		{
			name: "refresh",
			f:    New(),
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
			if got := SetFormat(tt.format); !reflect.DeepEqual(got.format.Load().(string), time.RFC3339) {
				t.Errorf("SetFormat() = %v, want %v", got.format.Load().(string), time.RFC3339)
			}
		})
	}
}
func TestFastime_SetFormat(t *testing.T) {
	tests := []struct {
		name   string
		f      *Fastime
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
			if got := tt.f.SetFormat(tt.format); !reflect.DeepEqual(got.format.Load().(string), time.RFC3339) {
				t.Errorf("Fastime.SetFormat() = %v, want %v", got.format.Load().(string), time.RFC3339)
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
		f    *Fastime
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
		f    *Fastime
	}{
		{
			name: "now",
			f:    New(),
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
		f    *Fastime
	}{
		{
			name: "update",
			f:    New(),
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
		f    *Fastime
	}{
		{
			name: "store",
			f:    New(),
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
