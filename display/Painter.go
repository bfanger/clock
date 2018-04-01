package display

import "github.com/veandco/go-sdl2/sdl"

// Painter paints it contents onto a Texture
// The painted texture is cached and only repainted if the properties change.
// When the painter is no longer used call Destroy() to free it's resources.
type Painter interface {
	Paint(*sdl.Renderer) (*Texture, error)
	Destroy() error
}
