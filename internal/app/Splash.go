package app

import (
	"time"

	"github.com/bfanger/clock/pkg/tween"
	"github.com/bfanger/clock/pkg/ui"
	"github.com/veandco/go-sdl2/sdl"
)

// Splash for notifications
type Splash struct {
	image  *ui.Image
	sprite *ui.Sprite
}

// NewSplash creates a new Splash
func NewSplash(r *sdl.Renderer) (*Splash, error) {
	image, err := ui.ImageFromFile(Asset("splash.jpg"), r)
	if err != nil {
		return nil, err
	}
	sprite := ui.NewSprite(image)
	sprite.SetAlpha(0)

	return &Splash{
		image:  image,
		sprite: sprite,
	}, nil
}

// Close free memory used by the Splash
func (s *Splash) Close() error {
	return s.image.Close()
}

// Compose the splash
func (s *Splash) Compose(r *sdl.Renderer) error {
	return s.sprite.Compose(r)
}

// Splash animation
func (s *Splash) Splash() tween.Tween {
	tl := &tween.Timeline{}
	tl.Add(tween.FromTo(0, 255, 300*time.Millisecond, tween.EaseInOutQuad, func(a uint8) {
		s.sprite.SetAlpha(a)
	}))
	tl.AddAt(500*time.Millisecond, tween.FromTo(255, 0, 400*time.Millisecond, tween.EaseInOutQuad, func(a uint8) {
		s.sprite.SetAlpha(a)
	}))
	return tl
}
