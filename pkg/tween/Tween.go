package tween

import (
	"time"
)

// Tween is a transition between 0.0 to 1.0
type Tween struct {
	Duration time.Duration
	Ease     Ease
	Update   func(float32)
}

// New creates a Tween
func New(d time.Duration, e Ease, update func(float32)) *Tween {
	return &Tween{Duration: d, Update: update, Ease: e}
}

// Seek to specific
func (t *Tween) Seek(d time.Duration) {
	t.Update(t.Value(d))
}

// Value calculates the eased value based on the progress
func (t *Tween) Value(d time.Duration) float32 {
	return t.Ease(t.Progress(d))
}

// Progress calculate the progress based on the duration
func (t *Tween) Progress(d time.Duration) float32 {
	if d < 0 {
		return 0
	}
	if d > t.Duration {
		return 1
	}
	return float32(d) / float32(t.Duration)
}

// Animate returns true when the tween completed
func (t *Tween) Animate(d time.Duration) bool {
	v := t.Value(d)
	t.Update(v)
	return v == 1
}

// FromToFloat32 creates a new Tween for an float32
func FromToFloat32(from, to float32, d time.Duration, e Ease, update func(float32)) *Tween {
	distance := to - from
	return New(d, e, func(v float32) {
		update(from + (distance * v))
	})
}

// FromToInt32 creates a new Tween for an Int32
func FromToInt32(from, to int32, d time.Duration, e Ease, update func(int32)) *Tween {
	f := float32(from)
	distance := float32(to) - f
	return New(d, e, func(v float32) {
		update(int32(f + (distance * v)))
	})
}

// FromToUint8 creates a new Tween for an Int8
func FromToUint8(from, to uint8, d time.Duration, e Ease, update func(uint8)) *Tween {
	f := float32(from)
	distance := float32(to) - f
	return New(d, e, func(v float32) {
		update(uint8(f + (distance * v)))
	})
}

// Empty creates a tween without a duration
func Empty() *Tween {
	return New(time.Duration(0), EaseInOutQuad, func(float32) {})
}
