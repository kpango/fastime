package fastime

import (
	"testing"
	"time"
)

// BenchmarkFastime
func BenchmarkFastime(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Now()
		}
	})
}

// BenchmarkTime
func BenchmarkTime(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			time.Now()
		}
	})
}
