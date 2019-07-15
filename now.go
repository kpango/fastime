// +build aix darwin dragonfly freebsd js,wasm linux nacl netbsd openbsd solaris

package fastime

import (
	"syscall"
	"time"
)

func (f *Fastime) now() time.Time {
	var tv syscall.Timeval
	err := syscall.Gettimeofday(&tv)
	if err != nil {
		return time.Now()
	}
	return time.Unix(0, syscall.TimevalToNsec(tv))
}
