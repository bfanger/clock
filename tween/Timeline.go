package tween

import (
	"time"
)

// Timeline is a sequence of (overlapping) entries
type Timeline struct {
	// Duration time.Duration
	entries []entry
	cursor  time.Duration
	end     time.Duration
}

type entry struct {
	start time.Duration
	t     *Tween
}

// Add a tween to the timeline
func (tl *Timeline) Add(t *Tween) {
	tl.entries = append(tl.entries, entry{start: tl.end, t: t})
	tl.end += t.Duration
}

// AddAt adds tween to a timeline at a specific momeny
func (tl *Timeline) AddAt(start time.Duration, t *Tween) {
	tl.entries = append(tl.entries, entry{start: start, t: t})
}

// Duration returns  the duration of the timeline
func (tl *Timeline) Duration() time.Duration {
	var max time.Duration
	for _, e := range tl.entries {
		if max < e.start+e.t.Duration {
			max = e.start + e.t.Duration
		}
	}
	return max
}

// Animate the timeline
func (tl *Timeline) Animate(dt time.Duration) bool {
	tl.cursor += dt
	var total time.Duration
	for _, e := range tl.entries {
		if tl.cursor < e.start {
			// e.t.Animate(0)
			// } else if tl.cursor > e.start+e.t.Duration {
			// 	e.t.Animate(e.t.Duration + time.Millisecond) // @todo check done?
		} else {
			e.t.Animate(dt) //tl.cursor - e.start
		}
		if total < e.start+e.t.Duration {
			total = e.start + e.t.Duration
		}
		// if max < e.start+e.t.Duration {
		// 	max = e.start + e.t.Duration
		// }
	}
	return tl.cursor > total
}
