package display

import "time"

// Animater is called in the renderloop with the time since the last frame.
// Animate() returns true when animation is complete.
type Animater interface {
	Animate(dt time.Duration) bool
}
