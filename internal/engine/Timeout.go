package engine

import (
	"time"
)

// Timeout creates a timer
func Timeout(callback func(), duration time.Duration) *time.Timer {
	timer := time.NewTimer(duration)
	go func() {
		<-timer.C
		if err := PushEvent(callback); err != nil {
			panic(err)
		}
	}()

	return timer
}
