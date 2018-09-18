package fastime

import (
	"context"
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
			f := New(context.Background())
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
			f := New(context.Background())
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
			f := New(context.Background())
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
