package display

import "github.com/veandco/go-sdl2/sdl"

// Layer is the interface used for compositing the final frame
type Layer interface {
	Name() string
	Render(*sdl.Renderer) error
}
