package ui

import (
	"image"

	"github.com/veandco/go-sdl2/sdl"
)

// Imager provide access to underling image/texture
type Imager interface {
	Image(*sdl.Renderer) (*image.Image, error)
}
