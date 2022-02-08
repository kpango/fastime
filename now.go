//go:build aix || darwin || dragonfly || freebsd || (js && wasm) || linux || nacl || netbsd || openbsd || solaris
// +build aix darwin dragonfly freebsd js,wasm linux nacl netbsd openbsd solaris

package fastime

import (
	"syscall"
	"time"
)

func (f *fastime) now() time.Time {
	var tv syscall.Timeval
	err := syscall.Gettimeofday(&tv)
	loc := f.GetLocation()
	if err != nil {
		return time.Now().In(loc)
	}
	return time.Unix(0, syscall.TimevalToNsec(tv)).In(loc)
}
