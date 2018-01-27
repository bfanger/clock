package engine

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// TextSprite provides text rendering
type TextSprite struct {
	dirty         bool
	font          *ttf.Font
	color         sdl.Color
	text          string
	textureSprite TextureSprite
}

// Render the text
func (sprite *TextSprite) Render() error {
	fmt.Println("Render Text")
	return sprite.textureSprite.Render()
}

// Destroy the sprite and free memory
func (sprite *TextSprite) Destroy() error {
	return sprite.textureSprite.Destroy()
}

// TextSpriteFromText create a sprite from text
func TextSpriteFromText(text string, renderer *sdl.Renderer) (*TextSprite, error) {
	font, err := DefaultFont()
	if err != nil {
		return nil, err
	}
	color := DefaultColor()
	surface, err := font.RenderUTF8Blended(text, color)
	defer surface.Free()

	texture, err := renderer.CreateTextureFromSurface(surface)

	source := sdl.Rect{X: 0, Y: 0, W: surface.W, H: surface.H}
	destination := sdl.Rect{X: 95, Y: 90, W: surface.W, H: surface.H}
	textureSprite := TextureSprite{
		Renderer:    renderer,
		Texture:     texture,
		Frame:       &source,
		Destination: &destination}
	return &TextSprite{
		dirty:         true,
		text:          text,
		font:          font,
		color:         DefaultColor(),
		textureSprite: textureSprite}, nil
}

var defaultFont *ttf.Font

// DefaultFont default font
func DefaultFont() (*ttf.Font, error) {
	if defaultFont == nil {
		var err error
		defaultFont, err = ttf.OpenFont("./assets/Teko-Light.ttf", 120)
		if err != nil {
			return nil, err
		}
		// 	defer defaultFont.Close()
	}
	return defaultFont, nil
}

var defaultColor sdl.Color

// DefaultColor determines the text color
func DefaultColor() sdl.Color {
	return sdl.Color{R: 255, G: 255, B: 255, A: 0}
}
