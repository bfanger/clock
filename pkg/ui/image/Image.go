package image

import (
	"github.com/veandco/go-sdl2/sdl"
)

// Image is the result of a paint and the base ingredient for the Compose()
type Image struct {
	Texture *sdl.Texture
	Frame   sdl.Rect
}

// New creates an Image from a texture
func New(t *sdl.Texture, frame sdl.Rect) *Image {
	return &Image{Texture: t, Frame: frame}
}

// Close frees the texture memory used by the image
func (i *Image) Close() error {
	return i.Texture.Destroy()
}

// Compose the image
func (i *Image) Compose(r *sdl.Renderer) error {
	return r.Copy(i.Texture, &i.Frame, &i.Frame)
}

// FromSurface creates an Image from a surface
func FromSurface(r *sdl.Renderer, s *sdl.Surface) (*Image, error) {
	t, err := r.CreateTextureFromSurface(s)
	if err != nil {
		return nil, err
	}
	return New(t, sdl.Rect{W: s.W, H: s.H}), nil
}
