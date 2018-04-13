package display

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// Text paints text using a font
type Text struct {
	FontFile string
	Size     int
	Color    sdl.Color
	Text     string
	font     *ttf.Font
	texture  *Texture
	previous *Text
}

// NewText creates Text
func NewText(fontFile string, size int, color sdl.Color, text string) *Text {
	return &Text{FontFile: fontFile, Size: size, Color: color, Text: text}
}

// Paint the text
func (t *Text) Paint(r *sdl.Renderer) (*Texture, error) {
	font, err := t.Font()
	if err != nil {
		return nil, err
	}
	if t.texture == nil || t.font != t.previous.font || t.Text != t.previous.Text || t.Color.Uint32() != t.previous.Color.Uint32() {
		surface, err := font.RenderUTF8Blended(t.Text, t.Color)
		// surface, err := font.RenderUTF8Solid(t.Text, t.Color)
		if err != nil {
			return nil, err
		}
		if t.texture != nil {
			err = t.texture.Destroy()
			if err != nil {
				return nil, err
			}
		}
		texture, err := TextureFromSurface(r, surface)
		if err != nil {
			return nil, err
		}
		t.texture = texture
		t.previous.font = t.font
		t.previous.Text = t.Text
		t.previous.Color = t.Color
	}
	return t.texture, nil
}

// Destroy font & texture
func (t *Text) Destroy() error {
	if t.font != nil {
		t.font.Close()
	}
	if t.texture != nil {
		return t.texture.Destroy()
	}
	return nil
}

// Font used to paint the text
func (t *Text) Font() (*ttf.Font, error) {
	if t.previous == nil {
		t.previous = &Text{}
	}
	if t.font == nil || t.FontFile != t.previous.FontFile || t.Size != t.previous.Size {
		if t.font != nil {
			t.font.Close()
		}
		font, err := ttf.OpenFont(t.FontFile, t.Size)
		if err != nil {
			return nil, fmt.Errorf("unable to open font: %v", err)
		}
		t.previous.FontFile = t.FontFile
		t.previous.Size = t.Size
		t.font = font
	}
	return t.font, nil
}
