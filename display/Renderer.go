package display

import (
	"fmt"
	"sort"

	"github.com/veandco/go-sdl2/sdl"
)

// Renderer simplifies the render loop
type Renderer struct {
	renderer *sdl.Renderer
	layers   map[int][]Layer
	zIndexes []int
}

// NewRenderer creates a renderer
func NewRenderer(w *sdl.Window) (*Renderer, error) {
	r, err := sdl.CreateRenderer(w, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		return nil, fmt.Errorf("can't create renderer: %v", err)
	}
	layers := make(map[int][]Layer)
	return &Renderer{renderer: r, layers: layers}, nil
}

// Render and present to the display
func (r *Renderer) Render() error {
	if err := r.renderer.Clear(); err != nil {
		return fmt.Errorf("renderer failed to clear: %v", err)
	}
	for _, z := range r.zIndexes {
		for _, layer := range r.layers[z] {
			if err := layer.Render(r.renderer); err != nil {
				return fmt.Errorf("%s failed to render: %v", layer.Name(), err)
			}
		}
	}
	r.renderer.Present()
	return nil
}

// Add a layer
func (r *Renderer) Add(zIndex int, l Layer) {
	newIndex := r.layers[zIndex] == nil
	r.layers[zIndex] = append(r.layers[zIndex], l)
	if newIndex {
		r.zIndexes = append(r.zIndexes, zIndex)
		sort.Ints(r.zIndexes)
	}
}
