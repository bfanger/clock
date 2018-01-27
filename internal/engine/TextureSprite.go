package engine

import (
	"fmt"

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
func (sprite *TextureSprite) Render() error {
	fmt.Println("Render Texture")
	return sprite.Renderer.Copy(sprite.Texture, sprite.Frame, sprite.Destination)
}

// Destroy the sprite and free memory
func (sprite *TextureSprite) Destroy() error {
	return sprite.Texture.Destroy()
}

// TextureSpriteFromImage creates a sprite from an image
func TextureSpriteFromImage(path string, renderer *sdl.Renderer) (*TextureSprite, error) {
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
