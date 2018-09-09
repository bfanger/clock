package ui

import "github.com/veandco/go-sdl2/sdl"

// Composer is the interface used for composing the final frame
type Composer interface {
	Compose(*sdl.Renderer) error
}
