package display

import (
	"github.com/veandco/go-sdl2/sdl"
)

// Texture is the result of a paint and the base ingredient for the Compose()
type Texture struct {
	Texture *sdl.Texture
	Frame   *sdl.Rect
}

// NewTexture create a texture
func NewTexture(texture *sdl.Texture, frame *sdl.Rect) *Texture {
	return &Texture{Texture: texture, Frame: frame}
}

// Paint returns the current texture
func (t *Texture) Paint(r *sdl.Renderer) (*sdl.Texture, *sdl.Rect, error) {
	return t.Texture, t.Frame, nil
}

// Destroy the texture
func (t *Texture) Destroy() error {
	return t.Texture.Destroy()
}

// TextureFromSurface create a texture from a surface
func TextureFromSurface(r *sdl.Renderer, s *sdl.Surface) (*Texture, error) {
	t, err := r.CreateTextureFromSurface(s)
	if err != nil {
		return nil, err
	}
	return NewTexture(t, &sdl.Rect{W: s.W, H: s.H}), nil
}
