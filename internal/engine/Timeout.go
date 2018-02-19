package engine

import (
	"sync"
	"time"
)

// Timeout represents a running timer
type Timeout struct {
	timer     *time.Timer
	cancelled bool
}

var timeoutMutex sync.Mutex

// SetTimeout creates a timer which can be cancelled
func SetTimeout(callback func(), duration time.Duration) *Timeout {
	timer := time.NewTimer(duration)
	timeout := Timeout{cancelled: false, timer: timer}
	go func() {
		<-timer.C
		timeoutMutex.Lock()
		defer timeoutMutex.Unlock()
		if timeout.cancelled {
			return
		}
		if err := PushEvent(callback); err != nil {
			panic(err)
		}
	}()
	return &timeout

}

// Cancel the timeout, the scheduled callback wont be called
func (timeout *Timeout) Cancel() {
	timeoutMutex.Lock()
	timeout.cancelled = true
	timeoutMutex.Unlock()
	timeout.timer.Stop()
}
