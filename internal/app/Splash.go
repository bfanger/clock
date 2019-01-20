package app

import (
	"time"

	"github.com/bfanger/clock/pkg/tween"
	"github.com/bfanger/clock/pkg/ui"
)

// Splash for notifications
type Splash struct {
	image  *ui.Image
	sprite *ui.Sprite
	engine *ui.Engine
}

// NewSplash creates a new Splash
func NewSplash(engine *ui.Engine) (*Splash, error) {
	image, err := ui.ImageFromFile(Asset("splash.jpg"), engine.Renderer)
	if err != nil {
		return nil, err
	}
	sprite := ui.NewSprite(image)
	sprite.SetAlpha(0)
	engine.Scene.Append(sprite)

	return &Splash{
		image:  image,
		sprite: sprite,
		engine: engine}, nil
}

// Close free memory used by the Splash
func (b *Splash) Close() error {
	return b.image.Close()
}

// Splash animation
func (b *Splash) Splash() tween.Tween {
	tl := &tween.Timeline{}
	tl.Add(tween.FromToUint8(0, 255, 300*time.Millisecond, tween.EaseInOutQuad, func(a uint8) {
		b.sprite.SetAlpha(a)
	}))
	tl.AddAt(500*time.Millisecond, tween.FromToUint8(255, 0, 400*time.Millisecond, tween.EaseInOutQuad, func(a uint8) {
		b.sprite.SetAlpha(a)
	}))
	return tl
}
