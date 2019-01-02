package app

import (
	"time"

	"github.com/bfanger/clock/pkg/tween"
	"github.com/bfanger/clock/pkg/ui"
)

// Overlay for notifications
type Overlay struct {
	image  *ui.Image
	sprite *ui.Sprite
	engine *ui.Engine
}

// NewOverlay creates a new overlay
func NewOverlay(engine *ui.Engine) (*Overlay, error) {
	image, err := ui.ImageFromFile(Asset("overlay.png"), engine.Renderer)
	if err != nil {
		return nil, err
	}
	sprite := ui.NewSprite(image)
	sprite.Y = screenHeight
	engine.Append(sprite)

	return &Overlay{
		image:  image,
		sprite: sprite,
		engine: engine}, nil
}

// Close free memory used by the overlay
func (b *Overlay) Close() error {
	return b.image.Close()
}

// Minimize overlay
func (b *Overlay) Minimize() tween.Tween {
	return tween.FromToInt32(196, screenHeight, 650*time.Millisecond, tween.EaseInQuad, func(y int32) {
		b.sprite.Y = y
	})
}

// Maximize overlay
func (b *Overlay) Maximize() tween.Tween {
	return tween.FromToInt32(screenHeight, 170, 800*time.Millisecond, tween.EaseInOutQuad, func(y int32) {
		b.sprite.Y = y
	})
}
