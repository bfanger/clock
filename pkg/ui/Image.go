package ui

import (
	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

// Image is the result of a paint and the base ingredient for the Compose()
type Image struct {
	Texture *sdl.Texture
	Frame   sdl.Rect
}

// ImageFromTexture creates an Image from a texture
func ImageFromTexture(t *sdl.Texture, frame sdl.Rect) *Image {
	return &Image{Texture: t, Frame: frame}
}

// Close frees the texture memory used by the image
func (i *Image) Close() error {
	return i.Texture.Destroy()
}

// Image also implements the Imager interface
func (i *Image) Image(_ *sdl.Renderer) (*Image, error) {
	return i, nil
}

// Compose the image
func (i *Image) Compose(r *sdl.Renderer) error {
	return r.Copy(i.Texture, &i.Frame, &i.Frame)
}

// ImageFromSurface creates an Image from a surface
func ImageFromSurface(s *sdl.Surface, r *sdl.Renderer) (*Image, error) {
	t, err := r.CreateTextureFromSurface(s)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &Image{Texture: t, Frame: sdl.Rect{W: s.W, H: s.H}}, nil
}

// ImageFromFile loads an image from disk
func ImageFromFile(filename string, r *sdl.Renderer) (*Image, error) {
	s, err := img.Load(filename)
	if err != nil {
		return nil, errors.Wrap(err, "could not load image")
	}
	defer s.Free()
	s.SetBlendMode(sdl.BLENDMODE_BLEND)
	return ImageFromSurface(s, r)
}

// ImageFromBytes downloads the image from the web
func ImageFromBytes(contents []byte, r *sdl.Renderer) (*Image, error) {
	buffer, err := sdl.RWFromMem(contents)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	surface, err := img.LoadRW(buffer, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer surface.Free()
	return ImageFromSurface(surface, r)
}
