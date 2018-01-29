package engine

import (
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

// Texture is a renderable from a sdl.Texture
type Texture struct {
	Renderer    *sdl.Renderer
	Texture     *sdl.Texture
	Frame       *sdl.Rect
	Destination *sdl.Rect
}

// Render the texture
func (texture *Texture) Render() error {
	return texture.Renderer.Copy(texture.Texture, texture.Frame, texture.Destination)
}

// Dispose and free resources
func (texture *Texture) Dispose() error {
	return texture.Texture.Destroy()
}

// TextureFromImage creates a Texture from an image
func TextureFromImage(renderer *sdl.Renderer, path string) (*Texture, error) {
	image, err := img.Load(path)
	if err != nil {
		return nil, err
	}
	defer image.Free()

	texture, err := renderer.CreateTextureFromSurface(image)
	if err != nil {
		return nil, err
	}
	frame := sdl.Rect{X: 0, Y: 0, W: image.W, H: image.H}
	destination := frame
	return &Texture{
		Renderer:    renderer,
		Texture:     texture,
		Frame:       &frame,
		Destination: &destination}, nil
}

// TextureFromSurface creates a Texture from a surface
func TextureFromSurface(renderer *sdl.Renderer, surface *sdl.Surface) (*Texture, error) {
	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return nil, err
	}
	source := sdl.Rect{X: 0, Y: 0, W: surface.W, H: surface.H}
	destination := sdl.Rect{X: 0, Y: 0, W: surface.W, H: surface.H}
	return &Texture{
		Renderer:    renderer,
		Texture:     texture,
		Frame:       &source,
		Destination: &destination}, nil
}
