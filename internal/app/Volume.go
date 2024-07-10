package app

import (
	"fmt"
	"time"

	"github.com/bfanger/clock/pkg/ui"
	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// var orange = sdl.Color{R: 203, G: 87, B: 0, A: 255}
var volumeText = sdl.Color{R: 255, G: 255, B: 255}

// Volume displays the volume
type Volume struct {
	visible bool
	value   int
	engine  *ui.Engine
	font    *ttf.Font
	text    *ui.Text
	sprite  *ui.Sprite
}

// NewVolume creates a new time widget
func NewVolume(engine *ui.Engine) (*Volume, error) {
	font, err := ttf.OpenFont(Asset("Roboto-Regular.ttf"), 90)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open font")
	}
	text := ui.NewText("", font, volumeText)
	sprite := ui.NewSprite(text)
	sprite.X = screenWidth - 0
	sprite.Y = screenHeight / 2
	sprite.AnchorX = 1
	sprite.AnchorY = 0.5

	// sprite.SetScale(0.2)
	// text.SetText("0")

	return &Volume{
		engine: engine,
		text:   text,
		font:   font,
		sprite: sprite,
	}, nil
}

// Close free resources
func (v *Volume) Close() error {
	if err := v.text.Close(); err != nil {
		return err
	}
	v.font.Close()
	return nil
}

// Update the volume indicator value
func (v *Volume) SetValue(value int) {
	v.value = value
	v.visible = true
	v.sprite.SetAlpha(200)
	v.text.SetText(fmt.Sprintf("%d", value))
	go func() {
		time.Sleep(time.Second)
		if value == v.value {
			v.visible = false
			v.engine.Go(func() error {
				v.sprite.SetAlpha(0)
				return nil
			})
		}
	}()
}

// Compose renders the volume indicator
func (v *Volume) Compose(r *sdl.Renderer) error {
	if !v.visible {
		return nil
	}
	return v.sprite.Compose(r)
}
