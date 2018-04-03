package display

import (
	"fmt"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
)

// Renderer simplifies the render loop
type Renderer struct {
	Container
	Mutex    sync.Mutex
	C        chan bool
	renderer *sdl.Renderer
}

// NewRenderer creates a renderer
func NewRenderer(w *sdl.Window) (*Renderer, error) {
	sdlr, err := sdl.CreateRenderer(w, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		return nil, fmt.Errorf("can't create sdl renderer: %v", err)
	}
	r := &Renderer{
		Container: Container{
			layers: make(map[int][]Layer),
			depths: []int{0},
		},
		C:        make(chan bool),
		renderer: sdlr,
	}

	go r.renderLoop()
	return r, nil
}

// Render and present to the display
func (r *Renderer) renderLoop() {
	var err error
	for range r.C {
		r.Mutex.Lock()
		if err = r.renderer.Clear(); err != nil {
			panic(fmt.Errorf("renderer failed to clear: %v", err))
		}
		if err = r.Render(r.renderer); err != nil {
			panic(fmt.Errorf("render failed: %v", err))
		}
		r.renderer.Present()
		r.Mutex.Unlock()
	}
}

// Destroy the render
func (r *Renderer) Destroy() error {
	close(r.C)
	return r.renderer.Destroy()
}

// // Add a layer
// func (r *Renderer) Add(zIndex int, l Layer) {
// 	newIndex := r.layers[zIndex] == nil
// 	r.layers[zIndex] = append(r.layers[zIndex], l)
// 	if newIndex {
// 		r.zIndexes = append(r.zIndexes, zIndex)
// 		sort.Ints(r.zIndexes)
// 	}
// }
