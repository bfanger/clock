package tween

import (
	"time"
)

// Tween is a transition between 0.0 to 1.0
type Tween struct {
	Duration time.Duration
	Update   func(float32)
	Ease     Easing
	elapsed  time.Duration
}

// Animate the tween
func (t *Tween) Animate(dt time.Duration) bool {
	t.elapsed += dt
	if t.elapsed > t.Duration {
		t.Update(1)
		return true
	}
	d := float32(t.elapsed) / float32(t.Duration)
	t.Update(t.Ease(d))
	return false
}

// WithEase create a new Tween with the specified easing function
func (t Tween) WithEase(e Easing) *Tween {
	t.Ease = e
	return &t
}

// FromToInt32 creates a new Tween for an Int32
func FromToInt32(from, to int32, d time.Duration, update func(int32)) *Tween {
	return &Tween{Duration: d, Ease: EaseInOut, Update: func(v float32) {
		update(from + int32(float32(to-from)*v))
	}}
}
