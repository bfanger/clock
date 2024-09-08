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
	visible   bool
	value     int
	engine    *ui.Engine
	container *ui.Container
	font      *ttf.Font
	text      *ui.Text
	full      *ui.Image
	empty     *ui.Image
	clip      *ui.Clip
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
	v.empty, err = ui.ImageFromFile(Asset("volume-bar/volume-empty.png"), engine.Renderer)
	if err != nil {
		return nil, err
	}
	sprite := ui.NewSprite(v.empty)
	sprite.X = screenWidth
	sprite.AnchorX = 1
	v.container.Append(sprite)
	v.full, err = ui.ImageFromFile(Asset("volume-bar/volume-full.png"), engine.Renderer)
	if err != nil {
		return nil, err
	}

	sprite = ui.NewSprite(v.full)
	sprite.X = screenWidth
	sprite.AnchorX = 1
	v.clip = &ui.Clip{
		Rect:     &sdl.Rect{X: 0, Y: screenHeight - 1, W: screenWidth, H: screenHeight},
		Composer: sprite,
	}
	v.container.Append(v.clip)

	v.text = ui.NewText("0", v.font, volumeText)
	sprite = ui.NewSprite(v.text)
	sprite.X = screenWidth - 60
	sprite.Y = screenHeight - 80
	sprite.AnchorX = 1
	sprite.AnchorY = 0.5
	v.container.Append(sprite)
	v.SetValue(0)
	return v, nil
}

// Close free resources
func (v *Volume) Close() error {
	v.font.Close()
	if err := v.empty.Close(); err != nil {
		return err
	}
	if err := v.full.Close(); err != nil {
		return err
	}
	if err := v.text.Close(); err != nil {
		return nil
	}
	return nil
}

// Update the volume indicator value
func (v *Volume) SetValue(value int) {
	v.value = value
	v.text.SetText(fmt.Sprintf("%d", value))

	v.visible = true
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
		// Weird bug where the clip rect is not working
		target -= 1
	}

	go func() {
		v.engine.Animate(tween.FromTo(v.clip.Y, target, 300*time.Millisecond, tween.EaseOutQuad, func(y int32) {
			if value == v.value {
				v.clip.Y = y
			}
		}))
		time.Sleep(time.Second)
		if value == v.value {
			v.engine.Go(func() error {
				v.visible = false
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
	return v.container.Compose(r)
}
