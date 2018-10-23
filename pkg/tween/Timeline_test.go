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
		tween := FromToFloat32(0, 1, time.Second, Linear, noop)
		tl.Add(tween)
		d := tl.Duration()
		assert.Equal(t, time.Second, d)
		tl.Add(tween)
		d = tl.Duration()
		assert.Equal(t, 2*time.Second, d)
	})

}

func TestLogic(t *testing.T) {
	t.Run("timeline should only seek relevant tweens", func(t *testing.T) {
		tl := Timeline{}
		var x1, x2, x3 int
		tl.Add(FromToInt(1, 10, time.Second, Linear, func(v int) {
			x1 += v
		}))
		tl.Add(FromToInt(20, 40, 2*time.Second, Linear, func(v int) {
			x2 += v
		}))
		tl.Add(FromToInt(1, 100, 1*time.Second, Linear, func(v int) {
			x3 += v
		}))
		assert.Equal(t, 0, x1)
		assert.Equal(t, 0, x2)
		assert.Equal(t, 0, x3)
		tl.Seek(500 * time.Millisecond)
		// Only first tween
		assert.Equal(t, 5, x1)
		assert.Equal(t, 0, x2)
		assert.Equal(t, 0, x3)
		tl.Seek(3500 * time.Millisecond)
		// All tweens (Forward 1 and 2)
		assert.Equal(t, 15, x1)
		assert.Equal(t, 40, x2)
		assert.Equal(t, 50, x3)
		tl.Seek(4000 * time.Millisecond)
		// Only third tween
		assert.Equal(t, 15, x1)
		assert.Equal(t, 40, x2)
		assert.Equal(t, 150, x3)
		tl.Seek(0)
		// Rewind all tweens
		assert.Equal(t, 16, x1)
		assert.Equal(t, 60, x2)
		assert.Equal(t, 151, x3)
	})
}
