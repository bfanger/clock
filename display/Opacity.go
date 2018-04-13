package display

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

// Opacity adds opacity to any Painter
type Opacity struct {
	Painter Painter
	Alpha   uint8
}

// NewOpacity creates a new opacity wrapper
func NewOpacity(painter Painter, alpha uint8) *Opacity {
	return &Opacity{Painter: painter, Alpha: alpha}
}

// Destroy the Opacity
func (o *Opacity) Destroy() error {
	return nil
}

// Paint the painter and add the opacity
func (o *Opacity) Paint(r *sdl.Renderer) (*Texture, error) {
	t, err := o.Painter.Paint(r)
	if err != nil {
		return nil, fmt.Errorf("failed to paint: %v", err)
	}
	t.Texture.SetAlphaMod(o.Alpha)
	return t, nil
}
