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
	done := d > tl.duration
	for _, e := range tl.entries {
		start, duration := e.start, e.t.Duration()
		if d >= start && d <= start+duration {
			// Tween is active
			if e.t.Seek(d-start) == false {
				done = false
			}
		} else if d < start {
			if tl.cursor > start {
				// Rewind tween
				if e.t.Seek(0) == false {
					done = false
				}
			}
		} else {
			if tl.cursor <= start+duration {
				// Forward to ending tween
				if e.t.Seek(e.t.Duration()) == false {
					done = false
				}
			}
		}
	}
	tl.cursor = d
	return done
}
