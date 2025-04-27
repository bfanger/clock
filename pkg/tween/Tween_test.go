package tween

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFromTo(t *testing.T) {
	t.Run("0 - 100", func(t *testing.T) {
		tween := FromTo(0, 100, 100*time.Second, Linear, func(v int32) {
			assert.Equal(t, int32(25), v)
		})
		tween.Seek(25 * time.Second)
	})
	t.Run("100 - 0", func(t *testing.T) {
		tween := FromTo(100, 0, 100*time.Second, Linear, func(v int32) {
			assert.Equal(t, int32(75), v)
		})
		tween.Seek(25 * time.Second)
	})
}
