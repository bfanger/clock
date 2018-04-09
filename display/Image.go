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
func (i *Image) Paint(r *sdl.Renderer) (*Texture, error) {
	if i.previous == nil {
		i.previous = &Image{}
	}
	if i.texture == nil || i.Filename != i.previous.Filename {
		surface, err := img.Load(i.Filename)
		if err != nil {
			return nil, err
		}
		defer surface.Free()
		i.previous.Filename = i.Filename
		if i.texture != nil {
			if err = i.texture.Destroy(); err != nil {
				return nil, err
			}
		}
		texture, err := TextureFromSurface(r, surface)
		if err != nil {
			return nil, err
		}
		i.texture = texture
	}
	return i.texture, nil
}

// Destroy the loaded texture
func (i *Image) Destroy() error {
	if i.texture != nil {
		return i.texture.Destroy()
	}
	return nil
}
