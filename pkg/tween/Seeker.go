package tween

import "time"

// Seeker jump to the position and returns true when the animation has completed
type Seeker interface {
	Seek(time.Duration) bool
}
