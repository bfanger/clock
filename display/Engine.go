package display

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

// Engine handles animation, syncronisation and // simplifies the render loop
type Engine struct {
	scene          Layer
	queueMutex     sync.Mutex
	queue          []*sync.WaitGroup
	animatersMutex sync.Mutex
	animaters      []Animater
	rendererMutex  sync.RWMutex
	renderer       *sdl.Renderer
	lastUpdate     time.Time
	err            error
	destroyed      bool
}

// NewEngine creates a new Engine
func NewEngine(r *sdl.Renderer, scene Layer) *Engine {
	return &Engine{
		scene:    scene,
		renderer: r,
	}
}

// Destroy the engine
func (e *Engine) Destroy() error {
	e.rendererMutex.Lock()
	defer e.rendererMutex.Unlock()
	e.destroyed = true
	return nil
}

const (
	slowdown       = 100 * time.Millisecond // Slowdown animation if below 10fps
	maxRefreshRate = 8 * time.Millisecond   // Limit FPS at ~125fps
)

// Do calls the function
func (e *Engine) Do(tasks ...interface{}) error {
	e.rendererMutex.RLock()
	if e.destroyed {
		e.rendererMutex.RUnlock()
		return errors.New("engine was destroyed")
	}
	for _, t := range tasks {
		switch cb := t.(type) {
		case func():
			cb()
		case func() error:
			if err := cb(); err != nil {
				e.rendererMutex.RUnlock()
				return err
			}
		default:
			e.rendererMutex.RUnlock()
			return fmt.Errorf("Invalid task %T", cb)
		}
	}
	e.rendererMutex.RUnlock()
	return e.Refresh()
}

// Refresh updates the display and waits until the frame is rendered.
func (e *Engine) Refresh() error {
	var wg sync.WaitGroup
	wg.Add(1)
	e.queueMutex.Lock()
	e.queue = append(e.queue, &wg)
	refresh := len(e.queue) == 1
	e.queueMutex.Unlock()

	if refresh {
		e.lastUpdate = time.Now().Add(-maxRefreshRate)
		more, err := e.render()
		if err != nil {
			return fmt.Errorf("count not render: %v", err)
		}
		if more {
			go func() {
				// keep rendering until all refreshes are flushed
				for more {
					more, err = e.render()
					if err != nil {
						e.err = fmt.Errorf("count not render: %v", err)
						return
					}
				}
			}()
		}
	}
	wg.Wait()
	return nil
}

func (e *Engine) render() (bool, error) {
	e.rendererMutex.Lock()
	defer e.rendererMutex.Unlock()
	if e.destroyed {
		return false, errors.New("engine was destroyed")
	}
	if err := e.err; e.err != nil {
		e.err = nil
		return false, err
	}
	e.queueMutex.Lock()
	l := len(e.queue)
	e.queueMutex.Unlock()

	e.prepareFrame()
	if err := e.drawFrame(); err != nil {
		return false, fmt.Errorf("could not render: %v", err)
	}

	e.queueMutex.Lock()
	queue := e.queue[:l]
	e.queue = e.queue[l:]
	rerender := len(e.queue) > 0
	if !rerender {
		e.animatersMutex.Lock()
		rerender = len(e.animaters) > 0
		e.animatersMutex.Unlock()
		if rerender {
			var wg sync.WaitGroup
			wg.Add(1)
			e.queue = append(e.queue, &wg)
		}
	}
	e.queueMutex.Unlock()
	for _, wg := range queue {
		wg.Done()
	}
	return rerender, nil
}

func (e *Engine) prepareFrame() {
	e.animatersMutex.Lock()
	defer e.animatersMutex.Unlock()

	now := time.Now()
	dt := now.Sub(e.lastUpdate)
	e.lastUpdate = now
	if dt > slowdown {
		dt = slowdown
	} else if dt < maxRefreshRate {
		time.Sleep(maxRefreshRate - dt)
	}

	var completed []Animater
	for _, a := range e.animaters {
		if a.Animate(dt) {
			completed = append(completed, a)
		}
	}
	// Remove completed animaters
	for _, a := range completed {
		for i, existing := range e.animaters {
			if a == existing {
				e.animaters = append(e.animaters[:i], e.animaters[i+1:]...)
			}
		}
	}
}

func (e *Engine) drawFrame() error {
	if err := e.renderer.Clear(); err != nil {
		return fmt.Errorf("could not clear: %v", err)
	}
	if err := e.scene.Render(e.renderer); err != nil {
		return err
	}
	e.renderer.Present()
	return nil
}

// Animate adds the animater to the renderloop
func (e *Engine) Animate(a Animater) {
	e.animatersMutex.Lock()
	defer e.animatersMutex.Unlock()
	e.animaters = append(e.animaters, a)
	if len(e.animaters) == 1 {
		go func() {
			if err := e.Refresh(); err != nil {
				e.err = err
			}
		}()
	}
}

// StopAnimation removes the animater from the renderLoop
func (e *Engine) StopAnimation(a Animater) bool {
	e.animatersMutex.Lock()
	defer e.animatersMutex.Unlock()
	for i, existing := range e.animaters {
		if a == existing {
			e.animaters = append(e.animaters[:i], e.animaters[i+1:]...)
			return true
		}
	}
	return false
}
