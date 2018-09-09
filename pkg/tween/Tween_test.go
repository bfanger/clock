package tween

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFromToInt32(t *testing.T) {
	t.Run("0 - 100", func(t *testing.T) {
		tween := FromToInt32(0, 100, 100*time.Second, Linear, func(v int32) {
			assert.Equal(t, int32(25), v)
		})
		tween.Seek(25 * time.Second)
	})
	t.Run("100 - 0", func(t *testing.T) {
		tween := FromToInt32(100, 0, 100*time.Second, Linear, func(v int32) {
			assert.Equal(t, int32(75), v)
		})
		tween.Seek(25 * time.Second)
	})

}
