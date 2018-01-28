package engine

import (
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

// TextureSprite /container
type TextureSprite struct {
	Renderer    *sdl.Renderer
	Texture     *sdl.Texture
	Frame       *sdl.Rect
	Destination *sdl.Rect
}

// Render the sprite
func (textureSprite *TextureSprite) Render() error {
	return textureSprite.Renderer.Copy(textureSprite.Texture, textureSprite.Frame, textureSprite.Destination)
}

// Destroy the sprite and free memory
func (textureSprite *TextureSprite) Destroy() error {
	return textureSprite.Texture.Destroy()
}

// TextureSpriteFromImage creates a sprite from an image
func TextureSpriteFromImage(renderer *sdl.Renderer, path string) (*TextureSprite, error) {
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
	return &TextureSprite{
		Renderer:    renderer,
		Texture:     texture,
		Frame:       &frame,
		Destination: &destination}, nil
}

// TextureSpriteFromSurface creates a sprite from a surface
func TextureSpriteFromSurface(renderer *sdl.Renderer, surface *sdl.Surface) (*TextureSprite, error) {
	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return nil, err
	}
	source := sdl.Rect{X: 0, Y: 0, W: surface.W, H: surface.H}
	destination := sdl.Rect{X: 0, Y: 0, W: surface.W, H: surface.H}
	return &TextureSprite{
		Renderer:    renderer,
		Texture:     texture,
		Frame:       &source,
		Destination: &destination}, nil
}
