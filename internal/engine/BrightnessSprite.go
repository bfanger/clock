package engine

import (
	"errors"

	"github.com/veandco/go-sdl2/sdl"
)

// BrightnessSprite provides software based brightness control
type BrightnessSprite struct {
	Alpha         uint8
	Renderer      *sdl.Renderer
	TextureSprite *TextureSprite
}

// NewBrightnessSprite creates a ready to use BrightnessSprite
func NewBrightnessSprite(renderer *sdl.Renderer, alpha uint8) (*BrightnessSprite, error) {
	brightnessSprite := BrightnessSprite{
		Renderer: renderer,
		Alpha:    alpha}

	if err := brightnessSprite.Update(); err != nil {
		return nil, err
	}
	return &brightnessSprite, nil
}

// Render the text
func (brightnessSprite *BrightnessSprite) Render() error {
	if brightnessSprite.TextureSprite == nil {
		return errors.New("Must call Update() before Render()")
	}
	return brightnessSprite.TextureSprite.Render()
}

// Destroy the sprite and free memory
func (brightnessSprite *BrightnessSprite) Destroy() error {
	if brightnessSprite.TextureSprite != nil {
		return brightnessSprite.TextureSprite.Destroy()
	}
	return nil
}

// Update generates the texture based on the settings
func (brightnessSprite *BrightnessSprite) Update() error {
	rect := brightnessSprite.Renderer.GetViewport()
	surface, err := sdl.CreateRGBSurfaceWithFormat(0, rect.W, rect.H, 32, sdl.PIXELFORMAT_ARGB8888)
	if err != nil {
		return err
	}
	defer surface.Free()
	color := sdl.Color{R: 0, G: 0, B: 0, A: brightnessSprite.Alpha}
	surface.FillRect(nil, color.Uint32())
	textureSprite, err := TextureSpriteFromSurface(brightnessSprite.Renderer, surface)
	if err != nil {
		return err
	}
	if brightnessSprite.TextureSprite != nil {
		textureSprite.Destination.X = brightnessSprite.TextureSprite.Destination.X
		textureSprite.Destination.Y = brightnessSprite.TextureSprite.Destination.Y
		brightnessSprite.TextureSprite.Destroy()
	}
	brightnessSprite.TextureSprite = textureSprite
	return nil
}
