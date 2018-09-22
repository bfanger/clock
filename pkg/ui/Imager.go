package ui

import (
	"github.com/veandco/go-sdl2/sdl"
)

// Imager provide access to underling image/texture
type Imager interface {
	Image(*sdl.Renderer) (*Image, error)
}
