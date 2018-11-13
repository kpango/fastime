package main

import (
	"fmt"
	"time"

	"github.com/kpango/fastime"
)

func main() {
	s1 := fastime.Now()
	s2 := fastime.Now()
	fastime.SetDuration(time.Millisecond * 500)
	s3 := fastime.Now()
	time.Sleep(time.Second * 2)
	s4 := fastime.Now()

	fmt.Printf("s1=%v\ns2=%v\ns3=%v\ns4=%v\n", s1, s2, s3, s4)
}
