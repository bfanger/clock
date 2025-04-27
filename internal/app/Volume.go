package app

import (
	"fmt"
	"time"

	"github.com/bfanger/clock/pkg/tween"
	"github.com/bfanger/clock/pkg/ui"
	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var volumeText = sdl.Color{R: 200, G: 200, B: 200}

// Volume displays the volume
type Volume struct {
	value     int
	engine    *ui.Engine
	container *ui.Container
	font      *ttf.Font
	text      *ui.Text
	clip      *ui.Clip
	images    struct {
		empty *ui.Image
		full  *ui.Image
	}
	sprites struct {
		empty *ui.Sprite
		full  *ui.Sprite
		text  *ui.Sprite
	}
}

// NewVolume creates a new time widget
func NewVolume(engine *ui.Engine) (*Volume, error) {
	v := &Volume{
		engine:    engine,
		container: &ui.Container{},
	}
	var err error
	v.font, err = ttf.OpenFont(Asset("Roboto-Regular.ttf"), 80)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open font")
	}
	v.images.empty, err = ui.ImageFromFile(Asset("volume-bar/volume-empty.png"), engine.Renderer)
	if err != nil {
		return nil, err
	}
	empty := ui.NewSprite(v.images.empty)
	empty.X = screenWidth
	empty.AnchorX = 1
	v.sprites.empty = empty
	v.container.Append(empty)

	v.images.full, err = ui.ImageFromFile(Asset("volume-bar/volume-full.png"), engine.Renderer)
	if err != nil {
		return nil, err
	}

	full := ui.NewSprite(v.images.full)
	full.X = screenWidth
	full.AnchorX = 1
	v.sprites.full = full
	v.clip = &ui.Clip{
		Rect:     &sdl.Rect{X: 0, Y: screenHeight - 1, W: screenWidth, H: screenHeight},
		Composer: full,
	}
	v.container.Append(v.clip)

	v.text = ui.NewText("0", v.font, volumeText)
	text := ui.NewSprite(v.text)
	text.X = screenWidth - 60
	text.Y = screenHeight - 80
	text.AnchorX = 1
	text.AnchorY = 0.5
	v.sprites.text = text
	v.container.Append(text)
	return v, nil
}

// Close free resources
func (v *Volume) Close() error {
	v.font.Close()
	if err := v.images.empty.Close(); err != nil {
		return err
	}
	if err := v.images.full.Close(); err != nil {
		return err
	}
	if err := v.text.Close(); err != nil {
		return nil
	}
	return nil
}

// Update the volume indicator value
func (v *Volume) SetValue(value int) {
	if v.value == value {
		return
	}
	v.value = value
	v.text.SetText(fmt.Sprintf("%d", value))
	v.sprites.text.SetAlpha(255)

	height := 0

	if value < 20 {
		height += value * 10
	} else if value < 40 {
		height = 200 + (value-20)*4
	} else {
		height = 280 + (value-40)*2
	}
	padding := (height) / 10
	height += padding * 2
	target := screenHeight - int32(height)
	if target == screenHeight {
		target -= 1 // Fixes a weird bug where the clip is not working
	}
	go v.engine.Animate(tween.FromTo(v.clip.Y, target, 300*time.Millisecond, tween.EaseOutQuad, func(y int32) {
		if v.value == value {
			v.clip.Y = y
		}
	}))
	go v.engine.Animate(tween.FromTo(v.sprites.empty.GetAlpha(), 255, 300*time.Millisecond, tween.Linear, func(alpha uint8) {
		if v.value == value {
			v.sprites.empty.SetAlpha(alpha)
		}
	}))
	go v.engine.Animate(tween.FromTo(v.sprites.full.GetAlpha(), 255, 100*time.Millisecond, tween.Linear, func(alpha uint8) {
		if v.value == value {
			v.sprites.full.SetAlpha(alpha)
		}
	}))
	go func() {
		time.Sleep(1750 * time.Millisecond)
		if value == v.value {
			v.engine.Animate(tween.FromTo(v.sprites.empty.GetAlpha(), 0, 1500*time.Millisecond, tween.Linear, func(alpha uint8) {
				if v.value == value {
					v.sprites.empty.SetAlpha(alpha)
					v.sprites.full.SetAlpha(alpha)
					v.sprites.text.SetAlpha(alpha)
				}
			}))
		}
	}()
}

// Compose renders the volume indicator
func (v *Volume) Compose(r *sdl.Renderer) error {
	if v.sprites.text.GetAlpha() == 0 {
		return nil
	}
	return v.container.Compose(r)
}
