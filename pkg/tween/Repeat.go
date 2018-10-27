package tween

import (
	"time"
)

type repeat struct {
	duration time.Duration
	infinite bool
	t        Tween
	cursor   time.Duration
}

// Repeat a Tween x times, 0 for "infinite"
func Repeat(times int, t Tween) Tween {
	return &repeat{
		duration: time.Duration(times) * t.Duration(),
		infinite: times == 0,
		t:        t}
}

func (r *repeat) Duration() time.Duration {
	return r.duration
}

func (r *repeat) Seek(d time.Duration) bool {
	tweenDuration := r.t.Duration()
	if int(d/tweenDuration) != int(r.cursor/tweenDuration) {
		r.t.Seek(0) // new repeat cycle, rewind tween
	}
	r.cursor = d
	r.t.Seek(d % tweenDuration)
	if r.infinite {
		r.duration = d + time.Minute // the duration grows dynamicly
		return false
	}
	return d > r.Duration()
}
