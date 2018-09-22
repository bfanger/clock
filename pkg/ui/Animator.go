package ui

import "time"

// Animator get the duration since it started and return true then the animation is complete
type Animator interface {
	Animate(time.Duration) bool
}
