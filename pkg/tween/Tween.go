package tween

import (
	"time"
)

// Tween is a transition between 0.0 to 1.0
type Tween struct {
	Duration time.Duration
	Ease     Ease
	Update   func(float32)
	StartAt  time.Time
}

// New creates a Tween
func New(d time.Duration, e Ease, update func(float32)) *Tween {
	return &Tween{Duration: d, Update: update, Ease: e, StartAt: time.Now()}
}

// Start restart the tween
func (t *Tween) Start() {
	t.StartAt = time.Now()
}
func (t *Tween) Seek(d time.Duration) {
	t.StartAt = time.Now().Add(-d)
	t.Animate()
}

// Animate call the update unction based on the elapsed time
func (t *Tween) Animate() bool {
	now := time.Now()
	if now.Before(t.StartAt) {
		// t.Update(0)
		return false
	}
	dt := now.Sub(t.StartAt)
	if dt > t.Duration {
		t.Update(1)
		return true
	}
	v := float32(dt) / float32(t.Duration)
	t.Update(t.Ease(v))
	return false
}

// FromToInt32 creates a new Tween for an Int32
func FromToInt32(from, to int32, d time.Duration, e Ease, update func(int32)) *Tween {
	return New(d, e, func(v float32) {
		update(from + int32(float32(to-from)*v))
	})
}

// FromToInt32Delta is similar to FromToInt32 but the update method receives the delta
func FromToInt32Delta(from, to int32, d time.Duration, e Ease, update func(int32)) *Tween {
	prev := from
	return FromToInt32(from, to, d, e, func(v int32) {
		d := v - prev
		prev = v
		update(d)
	})
}

// Empty creates a tween without a duration
func Empty() *Tween {
	return New(time.Duration(0), EaseInOutQuad, func(float32) {})
}
