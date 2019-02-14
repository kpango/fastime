package main

import (
	"fmt"
	"time"

	"github.com/kpango/fastime"
)

func main() {
	s1 := fastime.Now()
	s2 := fastime.Now()
	s3 := fastime.Now()
	time.Sleep(time.Second * 2)
	s4 := fastime.Now()

	time.Sleep(time.Second * 5)

	fmt.Printf("s1=%v\ns2=%v\ns3=%v\ns4=%v\n", s1, s2, s3, s4)

	fmt.Printf("nanonow %v\nnow unixnano %v\nnow add unixnano%v\nnanonow + dur %v\n", fastime.UnixNanoNow(), fastime.Now().Unix(), fastime.Now().Add(time.Second), fastime.UnixNanoNow()+int64(time.Second))
}
