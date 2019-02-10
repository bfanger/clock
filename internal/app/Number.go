package app

import (
	"strconv"

	"github.com/bfanger/clock/pkg/ui"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// Number manages Text layer per character
type Number struct {
	number int
	color  sdl.Color
	font   *ttf.Font
	image  *ui.Image
}

// NewNumber creates a new Number layer
func NewNumber(number int, f *ttf.Font, c sdl.Color) *Number {
	return &Number{number: number, font: f, color: c}
}

// Close free the texture memory
func (n *Number) Close() error {
	return n.needsUpdate()
}

// Image generate a image of the number
func (n *Number) Image(r *sdl.Renderer) (*ui.Image, error) {
	if n.image == nil {
		var err error
		characters := []*ui.Text{}
		var offset int32
		var letterSpacing int32 = -6
		frame := sdl.Rect{W: letterSpacing * -1}

		for _, c := range strconv.Itoa(n.number) {
			text := ui.NewText(string(c), n.font, n.color)
			defer text.Close()
			image, err := text.Image(r)
			if err != nil {
				return nil, err
			}
			text.X = offset
			offset += image.Frame.W + letterSpacing
			frame.W += image.Frame.W + letterSpacing
			frame.H = image.Frame.H
			characters = append(characters, text)
		}
		prevTarget := r.GetRenderTarget()
		defer r.SetRenderTarget(prevTarget)
		n.image = &ui.Image{Frame: frame}
		if n.image.Texture, err = r.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_TARGET, frame.W, frame.H); err != nil {
			return nil, err
		}
		if err := n.image.Texture.SetBlendMode(sdl.BLENDMODE_BLEND); err != nil {
			return nil, err
		}
		if err := r.SetRenderTarget(n.image.Texture); err != nil {
			return nil, err
		}
		if err := r.Clear(); err != nil {
			return nil, err
		}
		for _, glyp := range characters {
			glyp.Compose(r)
		}
	}
	return n.image, nil
}

// SetColor of the number
func (n *Number) SetColor(c sdl.Color) error {
	if c == n.color {
		return nil
	}
	n.color = c
	return n.needsUpdate()
}

func (n *Number) needsUpdate() error {
	if n.image != nil {
		if err := n.image.Close(); err != nil {
			return err
		}
		n.image = nil
	}
	return nil
}
