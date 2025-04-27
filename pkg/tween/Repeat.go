package tween

import (
	"time"
)

type repeat struct {
	duration time.Duration
	infinite bool
	t        Tween
}

// Repeat a Tween x times, 0 for "infinite"
func Repeat(times int, t Tween) Tween {
	return &repeat{
		duration: time.Duration(times) * t.Duration(),
		infinite: times == 0,
		t:        t,
	}
}

func (r *repeat) Duration() time.Duration {
	return r.duration
}

func (r *repeat) Seek(d time.Duration) bool {
	r.t.Seek(d % r.t.Duration())
	if r.infinite {
		r.duration = d + time.Minute // the duration grows dynamically
		return false
	}
	return d > r.Duration()
}
