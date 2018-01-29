package engine

import (
	"errors"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// Text provides text rendering
type Text struct {
	Font     *ttf.Font
	Color    sdl.Color
	Content  string
	Renderer *sdl.Renderer
	Texture  *Texture
}

// NewText creates a ready to use
func NewText(font *ttf.Font, color sdl.Color, content string, renderer *sdl.Renderer) (*Text, error) {
	text := Text{
		Font:     font,
		Color:    color,
		Content:  content,
		Renderer: renderer}
	if err := text.Update(); err != nil {
		return nil, err
	}
	return &text, nil

}

// Render the text
func (text *Text) Render() error {
	if text.Texture == nil {
		return errors.New("Must call Update() before Render()")
	}
	return text.Texture.Render()
}

// Dispose the sprite and free memory
func (text *Text) Dispose() error {
	if text.Texture != nil {
		return text.Texture.Dispose()
	}
	return nil
}

// Update texture based on the text, font and color settings.
func (text *Text) Update() error {
	if text.Font == nil {
		return errors.New("Font must be set before calling Update()")
	}
	surface, err := text.Font.RenderUTF8Blended(text.Content, text.Color)
	if err != nil {
		return err
	}
	defer surface.Free()
	Texture, err := TextureFromSurface(text.Renderer, surface)
	if err != nil {
		return err
	}
	if text.Texture != nil {
		Texture.Destination.X = text.Texture.Destination.X
		Texture.Destination.Y = text.Texture.Destination.Y
		text.Texture.Dispose()
	}
	text.Texture = Texture
	return nil
}
