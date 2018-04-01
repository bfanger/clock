package display

import "github.com/veandco/go-sdl2/sdl"

type Layer interface {
	Name() string
	Render(*sdl.Renderer) error
}
