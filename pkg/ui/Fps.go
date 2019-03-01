package ui

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// Fps counter
type Fps struct {
	count  int
	Scene  Composer
	engine *Engine
	done   chan bool
	Text   *Text
}

// NewFps create a new Frames per second counter
func NewFps(e *Engine, f *ttf.Font) *Fps {
	fps := &Fps{Scene: e.Scene, engine: e, Text: NewText("-", f, white), done: make(chan bool)}
	e.Scene = fps
	go fps.tick()
	return fps
}

// Close the fps and stop the tick
func (f *Fps) Close() error {
	close(f.done)
	f.engine.Scene = f.Scene
	return f.Text.Close()
}

// Compose the fps counter
func (f *Fps) Compose(r *sdl.Renderer) error {
	if err := f.Scene.Compose(r); err != nil {
		return err
	}
	f.count++
	return f.Text.Compose(r)
}

func (f *Fps) tick() {
	for {
		select {
		case <-time.After(time.Second):
			f.engine.Go(func() error {
				if err := f.Text.SetText(strconv.Itoa(f.count)); err != nil {
					return errors.Wrap(err, "failed to set text")
				}
				f.count = 0

				return nil
			})
		case <-f.done:
			return
		}
	}
}
