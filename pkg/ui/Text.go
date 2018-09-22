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

var white = sdl.Color{R: 255, G: 255, B: 255, A: 255}

// NewText creates new Text layer
func NewText(text string, f *ttf.Font, c sdl.Color) *Text {
	return &Text{text: text, font: f, color: c}
}

// Close free the texture memory
func (t *Text) Close() error {
	return t.needsUpdate()
}

// SetText update the contents
func (t *Text) SetText(text string) error {
	if text == t.text {
		return nil
	}
	t.text = text
	return t.needsUpdate()
}

// SetColor changes the color
func (t *Text) SetColor(c sdl.Color) error {
	if c == t.color {
		return nil
	}
	t.color = c
	return t.needsUpdate()
}

// SetFont changes the font
func (t *Text) SetFont(f *ttf.Font) error {
	t.font = f
	return t.needsUpdate()
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

// Width of the text
func (t *Text) Width(r *sdl.Renderer) (int32, error) {
	image, err := t.Image(r)
	if err != nil {
		return 0, err
	}
	return image.Frame.W, nil
}

// Height of the text
func (t *Text) Height(r *sdl.Renderer) (int32, error) {
	image, err := t.Image(r)
	if err != nil {
		return 0, err
	}
	return image.Frame.H, nil
}

// needsUpdate destroys the texture so the next call to Image() will generate a new image.
func (t *Text) needsUpdate() error {
	if t.image != nil {
		if err := t.image.Close(); err != nil {
			return err
		}
		t.image = nil
	}
	return nil
}
