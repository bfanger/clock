package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTile(t *testing.T) {
	t.Run("minmax", func(t *testing.T) {
		assert.Equal(t, 512, tileSize) // test is based around 512x512 tiles
		xmin, xmax := minmax(400, 256)
		assert.Equal(t, 0, xmin)
		assert.Equal(t, 0, xmax)

		xmin, xmax = minmax(700, 200)
		assert.Equal(t, -1, xmin)
		assert.Equal(t, 1, xmax)
	})
}
