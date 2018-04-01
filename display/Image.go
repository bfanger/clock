package display

import (
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

// Image load the image and converts it to a texture
type Image struct {
	Filename string
	texture  *Texture
	previous *Image
}

// NewImage creates an Image
func NewImage(filename string) *Image {
	return &Image{Filename: filename}
}

// Paint the image
func (im *Image) Paint(r *sdl.Renderer) (*Texture, error) {
	if im.previous == nil {
		im.previous = &Image{}
	}
	if im.texture == nil || im.Filename != im.previous.Filename {
		surface, err := img.Load(im.Filename)
		if err != nil {
			return nil, err
		}
		defer surface.Free()
		im.previous.Filename = im.Filename
		if im.texture != nil {
			if err = im.texture.Destroy(); err != nil {
				return nil, err
			}
		}
		texture, err := TextureFromSurface(r, surface)
		if err != nil {
			return nil, err
		}
		im.texture = texture
	}
	return im.texture, nil
}

// Destroy the loaded texture
func (im *Image) Destroy() error {
	if im.texture != nil {
		return im.texture.Destroy()
	}
	return nil
}
