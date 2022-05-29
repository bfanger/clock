package tween

import (
	"time"

	"golang.org/x/exp/constraints"
)

// Tween interface
type Tween interface {
	// Seek to a point on the
	Seek(time.Duration) bool
	// Duration of the tween
	Duration() time.Duration
}

// tween is a transition between 0.0 to 1.0
type tween struct {
	D      time.Duration
	Ease   Ease
	Update func(float32)
}

// Seek to specific position
func (t *tween) Seek(d time.Duration) bool {
	t.Update(t.value(d))
	return d >= t.D
}

// Duration of the tween
func (t *tween) Duration() time.Duration {
	return t.D
}

// value calculates the eased value based on the progress
func (t *tween) value(d time.Duration) float32 {
	return t.Ease(t.progress(d))
}

// progress calculate the progress based on the duration
func (t *tween) progress(d time.Duration) float32 {
	if d < 0 {
		return 0
	}
	if d >= t.D {
		return 1
	}
	return float32(d) / float32(t.D)
}

// FromTo creates a new Tween for number types
func FromTo[T constraints.Float | constraints.Integer](from, to T, d time.Duration, e Ease, update func(T)) Tween {
	fromFloat := float32(from)
	distance := float32(to) - fromFloat
	return &tween{
		D:    d,
		Ease: e,
		Update: func(v float32) {
			update(T(fromFloat + (distance * v)))
		},
	}
}

// Empty creates a tween without a duration
func Empty() Tween {
	return &tween{D: 0, Ease: EaseInOutQuad, Update: func(float32) {}}
}
