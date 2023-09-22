package ui

import (
	"github.com/veandco/go-sdl2/sdl"
)

// Imager provide access to underling image/texture
type Imager interface {
	Image(*sdl.Renderer) (*Image, error)
}

// ImageWidth get the width of the image
func ImageWidth(i Imager, r *sdl.Renderer) (int32, error) {
	image, err := i.Image(r)
	if err != nil {
		return 0, err
	}
	return image.Frame.W, nil
}

// ImageHeight get the height of the image
func ImageHeight(i Imager, r *sdl.Renderer) (int32, error) {
	image, err := i.Image(r)
	if err != nil {
		return 0, err
	}
	return image.Frame.H, nil
}
