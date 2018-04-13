package tween

import (
	"testing"
	"time"
)

func noop(float32) {}
func TestDuration(t *testing.T) {
	t.Run("empty timeline", func(t *testing.T) {
		tl := Timeline{}
		d := tl.Duration()
		if d != 0 {
			t.Errorf("expecting duration to be 0, got %s", d)
		}
	})

	t.Run("timeline with Add()", func(t *testing.T) {
		tl := Timeline{}
		tween := New(time.Second, Linear, noop)
		tl.Add(tween)
		d := tl.Duration()
		if d != time.Second {
			t.Errorf("expecting duration to be 1s, got %s", d)
		}
		tl.Add(tween)
		d = tl.Duration()
		if d != 2*time.Second {
			t.Errorf("expecting duration to be 2s, got %s", d)
		}
	})

}
