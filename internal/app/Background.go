package app

import (
	"time"

	"github.com/bfanger/clock/pkg/tween"
	"github.com/bfanger/clock/pkg/ui"
)

// Background for notifications
type Background struct {
	image  *ui.Image
	sprite *ui.Sprite
	engine *ui.Engine
}

// NewBackground creates a new background
func NewBackground(engine *ui.Engine) (*Background, error) {
	image, err := ui.ImageFromFile(Asset("background.png"), engine.Renderer)
	if err != nil {
		return nil, err
	}
	sprite := ui.NewSprite(image)
	sprite.Y = screenHeight
	engine.Append(sprite)

	return &Background{
		image:  image,
		sprite: sprite,
		engine: engine}, nil
}

// Close free memory used by the background
func (b *Background) Close() error {
	return b.image.Close()
}

// Minimize background
func (b *Background) Minimize() tween.Tween {
	return tween.FromToInt32(196, screenHeight, 650*time.Millisecond, tween.EaseInQuad, func(y int32) {
		b.sprite.Y = y
	})
}

// Maximize background
func (b *Background) Maximize() tween.Tween {
	return tween.FromToInt32(screenHeight, 196, 800*time.Millisecond, tween.EaseInOutQuad, func(y int32) {
		b.sprite.Y = y
	})
}
