package engine

import (
	"errors"

	"github.com/veandco/go-sdl2/sdl"
)

// Brightness provides software based brightness control
type Brightness struct {
	Alpha    uint8
	Renderer *sdl.Renderer
	Texture  *Texture
}

// NewBrightness creates a ready to use Brightness
func NewBrightness(renderer *sdl.Renderer, alpha uint8) (*Brightness, error) {
	brightness := Brightness{
		Renderer: renderer,
		Alpha:    alpha}

	if err := brightness.Update(); err != nil {
		return nil, err
	}
	return &brightness, nil
}

// Render the overlay
func (brightness *Brightness) Render() error {
	if brightness.Texture == nil {
		return errors.New("Must call Update() before Render()")
	}
	return brightness.Texture.Render()
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
	rect := brightness.Renderer.GetViewport()
	surface, err := sdl.CreateRGBSurfaceWithFormat(0, rect.W, rect.H, 32, sdl.PIXELFORMAT_ARGB8888)
	if err != nil {
		return err
	}
	defer surface.Free()
	color := sdl.Color{R: 0, G: 0, B: 0, A: brightness.Alpha}
	surface.FillRect(nil, color.Uint32())
	texture, err := TextureFromSurface(brightness.Renderer, surface)
	if err != nil {
		return err
	}
	if brightness.Texture != nil {
		texture.Destination.X = brightness.Texture.Destination.X
		texture.Destination.Y = brightness.Texture.Destination.Y
		brightness.Texture.Dispose()
	}
	brightness.Texture = texture
	return nil
}
