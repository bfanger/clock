package tween

import (
	"time"
)

// Timeline is a sequence of (overlapping) entries
type Timeline struct {
	// Duration time.Duration
	entries  []entry
	cursor   time.Duration
	duration time.Duration
}

type entry struct {
	start time.Duration
	t     Tween
}

// Add a tween to the timeline
func (tl *Timeline) Add(t Tween) {
	tl.entries = append(tl.entries, entry{start: tl.duration, t: t})
	tl.duration += t.Duration()
}

// AddAt adds tween to a timeline at a specific moment
func (tl *Timeline) AddAt(start time.Duration, t Tween) {
	tl.entries = append(tl.entries, entry{start: start, t: t})
	end := start + t.Duration()
	if end > tl.duration {
		tl.duration = end
	}
}

// Duration returns  the duration of the timeline
func (tl *Timeline) Duration() time.Duration {
	var max time.Duration
	for _, e := range tl.entries {
		if max < e.start+e.t.Duration() {
			max = e.start + e.t.Duration()
		}
	}
	return max
}

// Seek the timeline
func (tl *Timeline) Seek(d time.Duration) bool {
	for _, e := range tl.entries {
		if d >= e.start && d < e.start+e.t.Duration() {
			// Tween is active
			e.t.Seek(d - e.start)
		} else if d < e.start {
			if tl.cursor > e.start {
				// Rewind tween
				e.t.Seek(0)
			}
		} else {
			if tl.cursor < e.start+e.t.Duration() {
				// Forward to ending tween
				e.t.Seek(e.t.Duration())
			}
		}
	}
	tl.cursor = d
	return d > tl.duration
}
