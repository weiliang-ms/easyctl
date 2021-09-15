package util

import (
	"fmt"
	"time"
)

// Timer allows for injecting fake or real timers into code that
// needs to do arbitrary things based on time.
type Timer interface {
	C() <-chan time.Time
	Stop() bool
	Reset(d time.Duration) bool
}

// count秒后执行msg
func CountDown(count int, msg string) {
	tick := time.Tick(1 * time.Second)
	for countdown := count; countdown > 0; countdown-- {
		fmt.Printf("\r%2d秒后%s...", countdown, msg)
		<-tick
	}

	fmt.Println()
}
