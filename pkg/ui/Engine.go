package ui

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/bfanger/clock/pkg/tween"
	"github.com/veandco/go-sdl2/sdl"
)

// Engine handles the event- & render-loop
type Engine struct {
	Renderer   *sdl.Renderer
	Scene      Composer
	Wait       time.Duration // Limit framerate, 30 FPS = time.Second / 30
	updates    []func() error
	mutex      sync.Mutex
	waiting    atomic.Value
	lastRender time.Time
}

// NewEngine create a new engine
func NewEngine(scene Composer, r *sdl.Renderer) *Engine {
	e := &Engine{Renderer: r, Scene: scene}
	e.waiting.Store(false)
	return e
}

// Go schedules work that needs to be done in the ui thread
// When calling Go inside a update, that work will be schedules for the next frame.
// The error wil be propagated to the EventLoop()
func (e *Engine) Go(update func() error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	// @todo Set a max for queued updates?
	e.updates = append(e.updates, update)
	if e.waiting.Load().(bool) {
		event := &sdl.UserEvent{Type: sdl.USEREVENT, Timestamp: sdl.GetTicks(), Code: 808}
		if _, err := sdl.PushEvent(event); err != nil {
			panic(err) // @todo propagate to EventLoop
		}
	}
}

// Do run the code on the ui thread. Engine.Do() is the synchronisch variant of Engine.Go()
// Is also returns the error that occurred in the callback instead of ending the EventLoop.
func (e *Engine) Do(fn func() error) error {
	wg := sync.WaitGroup{}
	wg.Add(1)
	var err error
	e.Go(func() error {
		err = fn()
		wg.Done()
		return nil
	})
	wg.Wait()
	return err
}

// EventLoop runs the event-loop
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
func (e *Engine) Animate(a tween.Seeker) {
	start := time.Now()
	wg := sync.WaitGroup{}
	done := false
	update := func() error {
		done = a.Seek(time.Since(start))
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
	for _, update := range tmp {
		if err := update(); err != nil {
			return err
		}
	}
	return nil
}

func (e *Engine) render() error {
	started := time.Now()
	if err := e.Renderer.Clear(); err != nil {
		return err
	}
	if err := e.Scene.Compose(e.Renderer); err != nil {
		return err
	}
	e.Renderer.Present()
	completed := time.Now()
	dt := completed.Sub(e.lastRender)
	if dt < e.Wait {
		render := completed.Sub(started)
		update := started.Sub(e.lastRender)
		time.Sleep(e.Wait - dt - render - update)
	}
	e.lastRender = time.Now()
	return nil
}

// func (e *Engine) forceUpdate() {
// 	e.Go(noop)
// }

// func noop() error {
// 	return nil
// }
