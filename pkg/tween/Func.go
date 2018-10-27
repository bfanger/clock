package tween

import (
	"time"
)

type funcTween struct {
	fn func()
}

// Func add a func as a tween to a timeline.
func Func(fn func()) Tween {
	return &funcTween{fn: fn}
}

func (*funcTween) Duration() time.Duration {
	return 0
}

func (f *funcTween) Seek(_ time.Duration) bool {
	f.fn()
	return true
}
