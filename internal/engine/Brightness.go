package engine

import (
	"errors"

	"github.com/veandco/go-sdl2/sdl"
)

// Brightness provides software based brightness control
type Brightness struct {
	Alpha   uint8
	Texture *Texture
}

// NewBrightness creates a ready to use Brightness
func NewBrightness(alpha uint8) (*Brightness, error) {
	brightness := Brightness{Alpha: alpha}

	if err := brightness.Update(); err != nil {
		return nil, err
	}
	return &brightness, nil
}

// Draw the overlay
func (brightness *Brightness) Draw() error {
	if brightness.Texture == nil {
		return errors.New("Must call Update() before Draw()")
	}
	return brightness.Texture.Draw()
}

// Dispose and free resources
func (brightness *Brightness) Dispose() error {
	if brightness.Texture != nil {
		return brightness.Texture.Dispose()
	}
	return nil
}

// Update generates the texture based on the settings
func (brightness *Brightness) Update() error {
	rect := Renderer().GetViewport()
	color := sdl.Color{R: 0, G: 0, B: 0, A: brightness.Alpha}
	texture, err := TextureFromColor(rect.W, rect.H, color)
	if err != nil {
		return err
	}
	if brightness.Texture != nil {
		brightness.Texture.Dispose()
	}
	brightness.Texture = texture
	return nil
}
