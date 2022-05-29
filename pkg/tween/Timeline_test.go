package tween

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func noop(float32) {}
func TestDuration(t *testing.T) {
	t.Run("empty timeline", func(t *testing.T) {
		assert := assert.New(t)
		tl := Timeline{}
		d := tl.Duration()
		assert.Equal(time.Duration(0), d)
	})

	t.Run("timeline with Add()", func(t *testing.T) {
		assert := assert.New(t)
		tl := Timeline{}
		tween := FromTo(0, 1, time.Second, Linear, noop)
		tl.Add(tween)
		d := tl.Duration()
		assert.Equal(time.Second, d)
		tl.Add(tween)
		d = tl.Duration()
		assert.Equal(2*time.Second, d)
	})

}

func TestLogic(t *testing.T) {
	t.Run("timeline should only seek relevant tweens", func(t *testing.T) {
		assert := assert.New(t)
		tl := Timeline{}
		var x1, x2, x3 int
		tl.Add(FromTo(1, 10, time.Second, Linear, func(v int) {
			x1 += v
		}))
		tl.Add(FromTo(20, 40, 2*time.Second, Linear, func(v int) {
			x2 += v
		}))
		tl.Add(FromTo(1, 100, 1*time.Second, Linear, func(v int) {
			x3 += v
		}))
		assert.Equal(0, x1)
		assert.Equal(0, x2)
		assert.Equal(0, x3)
		tl.Seek(500 * time.Millisecond)
		// Only first tween
		assert.Equal(5, x1)
		assert.Equal(0, x2)
		assert.Equal(0, x3)
		tl.Seek(3500 * time.Millisecond)
		// All tweens (Forward 1 and 2)
		assert.Equal(15, x1)
		assert.Equal(40, x2)
		assert.Equal(50, x3)
		tl.Seek(4000 * time.Millisecond)
		// Only third tween
		assert.Equal(15, x1)
		assert.Equal(40, x2)
		assert.Equal(150, x3)
		tl.Seek(0)
		// Rewind all tweens
		assert.Equal(16, x1)
		assert.Equal(60, x2)
		assert.Equal(151, x3)
	})
}
func TestFunc(t *testing.T) {
	t.Run("timeline should call func", func(t *testing.T) {
		assert := assert.New(t)
		tl := Timeline{}
		var called int
		tl.Add(Func(func() { called++ }))
		assert.Equal(0, called)
		tl.Seek(time.Millisecond)
		assert.Equal(1, called)
		tl.Seek(time.Second)
		assert.Equal(1, called)
		tl.Seek(0)
		assert.Equal(2, called)
		// tl.Seek(0)
		// assert.Equal(3, called)
	})
}
