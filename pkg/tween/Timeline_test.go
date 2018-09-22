package tween

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func noop(float32) {}
func TestDuration(t *testing.T) {
	t.Run("empty timeline", func(t *testing.T) {
		tl := Timeline{}
		d := tl.Duration()
		assert.Equal(t, time.Duration(0), d)
	})

	t.Run("timeline with Add()", func(t *testing.T) {
		tl := Timeline{}
		tween := New(time.Second, Linear, noop)
		tl.Add(tween)
		d := tl.Duration()
		assert.Equal(t, time.Second, d)
		tl.Add(tween)
		d = tl.Duration()
		assert.Equal(t, 2*time.Second, d)
	})
}
