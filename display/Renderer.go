package display

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

// Renderer simplifies the render loop
type Renderer struct {
	Container
	C         chan bool
	animaters []Animater
	renderer  *sdl.Renderer
	quit      chan bool
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
		quit:     make(chan bool),
		renderer: sdlr,
	}

	go r.renderLoop()
	return r, nil
}

// Destroy the render
func (r *Renderer) Destroy() error {
	r.quit <- true
	close(r.quit)
	close(r.C)
	return r.renderer.Destroy()
}

// Render and present to the display
func (r *Renderer) renderLoop() {
	var err error
	var prevUpdate = time.Now()
	var nextUpdate time.Time
	var dt time.Duration
	var completed []Animater
	for {
		select {
		case <-r.quit:
			return
		case <-r.C:
		default:
			if len(r.animaters) == 0 {
				select {
				case <-r.quit:
					return
				case <-r.C: // Wait for refresh event
				}
			}
		}
		if err = r.renderer.Clear(); err != nil {
			panic(fmt.Errorf("renderer failed to clear: %v", err))
		}
		nextUpdate = time.Now()
		dt = nextUpdate.Sub(prevUpdate)
		r.Mutex.Lock()
		for _, a := range r.animaters {
			if a.Animate(dt) {
				completed = append(completed, a)
			}
		}
		for _, a := range completed {
			for i, existing := range r.animaters {
				if a == existing {
					r.animaters = append(r.animaters[:i], r.animaters[i+1:]...)
				}
			}
		}
		r.Mutex.Unlock()
		completed = nil
		if err = r.Render(r.renderer); err != nil {
			panic(fmt.Errorf("render failed: %v", err))
		}
		r.renderer.Present()
		prevUpdate = nextUpdate
	}
}

// Animate adds the animater to the renderloop
func (r *Renderer) Animate(a Animater) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	r.animaters = append(r.animaters, a)
	if len(r.animaters) == 1 {
		go Refresh()
	}
}

// StopAnimation removes the animater from the renderLoop
func (r *Renderer) StopAnimation(a Animater) error {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	for i, existing := range r.animaters {
		if a == existing {
			r.animaters = append(r.animaters[:i], r.animaters[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("couldn't find animator: %+v", a)
}
