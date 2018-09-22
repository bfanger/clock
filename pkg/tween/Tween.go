package tween

import (
	"time"
)

// Tween is a transition between 0.0 to 1.0
type Tween struct {
	Duration  time.Duration
	Ease      Ease
	Update    func(float32)
	StartedAt time.Time
}

// New creates a Tween
func New(d time.Duration, e Ease, update func(float32)) *Tween {
	return &Tween{Duration: d, Update: update, Ease: e, StartedAt: time.Now()}
}

// Start the tween
func (t *Tween) Start() {
	t.StartedAt = time.Now()
	t.Update(0)
}

// Seek to specific
func (t *Tween) Seek(d time.Duration) {
	now := time.Now()
	t.StartedAt = now.Add(-d)
	t.Update(t.Value(now))
}

// Value calculated the eased value based on the current time
func (t *Tween) Value(now time.Time) float32 {
	return t.Ease(t.Progress(now))
}

// Progress calculate the progress based on the current time
func (t *Tween) Progress(now time.Time) float32 {
	if now.Before(t.StartedAt) {
		return 0
	}
	dt := now.Sub(t.StartedAt)
	if dt > t.Duration {
		return 1
	}
	return float32(dt) / float32(t.Duration)
}

// Animate call the update unction based on the elapsed time
func (t *Tween) Animate(now time.Time) bool {
	v := t.Value(now)
	t.Update(v)
	return v == 1
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
