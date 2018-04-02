package display

import (
	"fmt"
	"sort"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
)

// Renderer simplifies the render loop
type Renderer struct {
	Mutex    sync.Mutex
	C        chan bool
	renderer *sdl.Renderer
	layers   map[int][]Layer
	zIndexes []int
}

// NewRenderer creates a renderer
func NewRenderer(w *sdl.Window) (*Renderer, error) {
	sdlr, err := sdl.CreateRenderer(w, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		return nil, fmt.Errorf("can't create sdl renderer: %v", err)
	}
	r := &Renderer{
		C:        make(chan bool),
		renderer: sdlr,
		layers:   make(map[int][]Layer)}

	go r.renderLoop()
	return r, nil
}

func (r *Renderer) renderLoop() {
	var err error
	for range r.C {
		if err = r.Render(); err != nil {
			panic(err)
		}
	}
}

// Destroy the render
func (r *Renderer) Destroy() error {
	close(r.C)
	return r.renderer.Destroy()
}

// Render and present to the display
func (r *Renderer) Render() error {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
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
