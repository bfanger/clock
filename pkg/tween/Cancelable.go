package tween

import (
	"time"
)

// Cancelable tween
type Cancelable struct {
	canceled bool
	t        Tween
}

// NewCancelable creates cancelable tween from a regular tween
func NewCancelable(t Tween) *Cancelable {
	return &Cancelable{t: t}
}

// Cancel the tween
func (c *Cancelable) Cancel() {
	if c == nil {
		return
	}
	c.canceled = true
}

// Seek to a position when the tween is not canceled
func (c *Cancelable) Seek(d time.Duration) bool {
	if c.canceled {
		return true
	}
	return c.t.Seek(d)
}

// Duration of the tween of 0 when the tween is canceled
func (c *Cancelable) Duration() time.Duration {
	if c.canceled {
		return 0
	}
	return c.t.Duration()
}
