package engine

import (
	"errors"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// TextSprite provides text rendering
type TextSprite struct {
	Font          *ttf.Font
	Color         sdl.Color
	Text          string
	Renderer      *sdl.Renderer
	TextureSprite *TextureSprite
}

// Render the text
func (textSprite *TextSprite) Render() error {
	if textSprite.TextureSprite == nil {
		return errors.New("Must call Update() before Render()")
	}
	return textSprite.TextureSprite.Render()
}

// Destroy the sprite and free memory
func (textSprite *TextSprite) Destroy() error {
	if textSprite.TextureSprite != nil {
		return textSprite.TextureSprite.Destroy()
	}
	return nil
}

// Update texture based on the text, font and color settings.
func (textSprite *TextSprite) Update() error {
	if textSprite.Font == nil {
		return errors.New("Font must be set before calling Update()")
	}
	surface, err := textSprite.Font.RenderUTF8Blended(textSprite.Text, textSprite.Color)
	if err != nil {
		return err
	}
	defer surface.Free()
	textureSprite, err := TextureSpriteFromSurface(textSprite.Renderer, surface)
	if err != nil {
		return err
	}
	if textSprite.TextureSprite != nil {
		textSprite.TextureSprite.Destroy()
	}
	textSprite.TextureSprite = textureSprite
	return nil
}
