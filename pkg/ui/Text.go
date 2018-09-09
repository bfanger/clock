package ui

import (
	"errors"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// Text manages text
type Text struct {
	X, Y  int32
	text  string
	color sdl.Color
	font  *ttf.Font
	image *Image
}

// NewText creates new Text layer
func NewText(text string, c sdl.Color, f *ttf.Font, x, y int32) *Text {
	return &Text{text: text, font: f, color: c, X: x, Y: y}
}

// Close free the texture memory
func (t *Text) Close() error {
	if t.image != nil {
		if err := t.image.Close(); err != nil {
			return err
		}
		t.image = nil
	}
	return nil
}

// SetText update the contents
func (t *Text) SetText(text string) error {
	if text == t.text {
		return nil
	}
	t.text = text
	return t.Close()
}

// SetColor changes the color
func (t *Text) SetColor(c sdl.Color) error {
	if c == t.color {
		return nil
	}
	t.color = c
	return t.Close()
}

// SetFont changes the font
func (t *Text) SetFont(f *ttf.Font) error {
	t.font = f
	return t.Close()
}

// Image convert the text into an image (and caches the result)
func (t *Text) Image(r *sdl.Renderer) (*Image, error) {
	if t.text == "" {
		return nil, nil
	}
	if t.image == nil {
		if t.font == nil {
			return nil, errors.New("(*ui.Text).Font was nil")
		}
		surface, err := t.font.RenderUTF8Blended(t.text, t.color)
		if err != nil {
			return nil, err
		}
		defer surface.Free()
		t.image, err = ImageFromSurface(r, surface)
		if err != nil {
			return nil, err
		}
	}
	return t.image, nil
}

// Compose the Text
func (t *Text) Compose(r *sdl.Renderer) error {
	image, err := t.Image(r)
	if err != nil {
		return err
	}
	if image == nil {
		return nil
	}
	frame := &image.Frame
	return r.Copy(image.Texture, frame, &sdl.Rect{X: t.X, Y: t.Y, W: frame.W, H: frame.H})
}
