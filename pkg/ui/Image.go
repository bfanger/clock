package ui

import (
	"github.com/veandco/go-sdl2/sdl"
)

// Imager provide access to underling image/texture
type Imager interface {
	Image(*sdl.Renderer) (*Image, error)
}

// Image is the result of a paint and the base ingredient for the Compose()
type Image struct {
	Texture *sdl.Texture
	Frame   sdl.Rect
}

// Close frees the texture memory used by the image
func (i *Image) Close() error {
	return i.Texture.Destroy()
}

// Compose the image
func (i *Image) Compose(r *sdl.Renderer) error {
	return r.Copy(i.Texture, &i.Frame, &i.Frame)
}

// ImageFromTexture creates an Image from a texture
func ImageFromTexture(t *sdl.Texture, frame sdl.Rect) *Image {
	return &Image{Texture: t, Frame: frame}
}

// ImageFromSurface creates an Image from a surface
func ImageFromSurface(r *sdl.Renderer, s *sdl.Surface) (*Image, error) {
	t, err := r.CreateTextureFromSurface(s)
	if err != nil {
		return nil, err
	}
	return &Image{Texture: t, Frame: sdl.Rect{W: s.W, H: s.H}}, nil
}
