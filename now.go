//go:build aix || darwin || dragonfly || freebsd || (js && wasm) || linux || nacl || netbsd || openbsd || solaris
// +build aix darwin dragonfly freebsd js,wasm linux nacl netbsd openbsd solaris

package fastime

import (
	"syscall"
	"time"
)

func (f *fastime) now() (now time.Time) {
	var tv syscall.Timeval
	err := syscall.Gettimeofday(&tv)
	loc := f.GetLocation()
	if err != nil {
		now = time.Now()
		if loc != nil{
		    return now.In(loc)
		}
		return now
	}
	now = time.Unix(0, syscall.TimevalToNsec(tv))
	if loc != nil{
	    return now.In(loc)
	}
	return now
}
