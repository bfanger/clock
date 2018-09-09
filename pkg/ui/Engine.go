package ui

import (
	"sync"
	"sync/atomic"

	"github.com/bfanger/clock/pkg/tween"
	"github.com/veandco/go-sdl2/sdl"
)

// Engine handles the event- & renderloop
type Engine struct {
	Renderer  *sdl.Renderer
	Composers []Composer
	updates   []func() error
	mutex     sync.Mutex
	waiting   atomic.Value
}

// NewEngine create a new engine
func NewEngine(r *sdl.Renderer) *Engine {
	e := &Engine{Renderer: r}
	e.waiting.Store(false)
	return e
}

// Go shedules work that needs to be done in the ui thread
// When calling Go inside a update, that work will be shedules for the next frame.
func (e *Engine) Go(update func() error) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	// @todo Set a max for queued updates?
	e.updates = append(e.updates, update)
	if e.waiting.Load().(bool) {
		event := &sdl.UserEvent{Type: sdl.USEREVENT, Timestamp: sdl.GetTicks(), Code: 808}
		if _, err := sdl.PushEvent(event); err != nil {
			return err
		}
	}
	return nil
}

// EventLoop runs the evenloop
func (e *Engine) EventLoop(handle func(sdl.Event)) error {
	for {
		if e.needsUpdate() {
			if err := e.update(); err != nil {
				return err
			}
			if err := e.render(); err != nil {
				return err
			}
		}
		event := sdl.PollEvent()
		if event == nil {
			if e.needsUpdate() {
				continue
			}
			e.waiting.Store(true)
			event = sdl.WaitEvent()
			e.waiting.Store(false)
		}
		switch typedEvent := event.(type) {
		case *sdl.QuitEvent:
			return nil
		case *sdl.WindowEvent:
			if typedEvent.Event == sdl.WINDOWEVENT_EXPOSED {
				if err := e.render(); err != nil {
					return err
				}
			}
			handle(event)
		default:
			handle(event)
		}
	}
}

// Animate the tween
func (e *Engine) Animate(t *tween.Tween) {
	t.Start()
	wg := sync.WaitGroup{}
	done := false
	update := func() error {
		done = t.Animate()
		wg.Done()
		return nil
	}
	for {
		wg.Add(1)
		e.Go(update)
		wg.Wait()
		if done {
			break
		}
	}
}

// Append a composer
func (e *Engine) Append(c Composer) {
	e.Composers = append(e.Composers, c)
	e.Go(noop)
}

// Remove a composer
func (e *Engine) Remove(layer Composer) {
	for i, l := range e.Composers {
		if l == layer {
			e.Composers = append(e.Composers[:i], e.Composers[i+1:]...)
			e.Go(noop)
			return
		}
	}
}
func noop() error {
	return nil
}

func (e *Engine) needsUpdate() bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	return len(e.updates) > 0
}

func (e *Engine) update() error {
	e.mutex.Lock()
	tmp := make([]func() error, len(e.updates))
	copy(tmp, e.updates)
	e.updates = nil
	e.mutex.Unlock()
	// log.Printf("updates: %d", len(tmp))
	for _, update := range tmp {
		if err := update(); err != nil {
			return err
		}
	}
	return nil
}

func (e *Engine) render() error {
	if err := e.Renderer.Clear(); err != nil {
		return err
	}
	for _, layer := range e.Composers {
		if err := layer.Compose(e.Renderer); err != nil {
			return err
		}
	}
	e.Renderer.Present()
	return nil
}
