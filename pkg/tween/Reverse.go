package tween

import "time"

func Reverse(t Tween) Tween {
	return &reverse{tween: t}
}

type reverse struct {
	tween Tween
}

func (r *reverse) Seek(d time.Duration) bool {
	total := r.tween.Duration()
	done := d >= total
	if done {
		r.tween.Seek(0)
	} else {
		r.tween.Seek(total - d)
	}
	return done
}

// Duration of the tween
func (r *reverse) Duration() time.Duration {
	return r.tween.Duration()
}
