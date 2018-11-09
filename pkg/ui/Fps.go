package ui

import (
	"fmt"
	"strconv"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// Fps counter
type Fps struct {
	count  int
	engine *Engine
	done   chan bool
	Text   *Text
}

// NewFps create a new Frames per second counter
func NewFps(e *Engine, f *ttf.Font) *Fps {

	fps := &Fps{engine: e, Text: NewText("-", f, white), done: make(chan bool)}
	fps.Text.X = 10
	fps.Text.Y = 7
	e.Append(fps)
	go fps.tick()
	return fps
}

// Close the fps and stop the tick
func (f *Fps) Close() error {
	close(f.done)
	return f.Text.Close()
}

// Compose the fps counter
func (f *Fps) Compose(r *sdl.Renderer) error {
	f.count++
	return f.Text.Compose(r)
}

func (f *Fps) tick() {
	for {
		select {
		case <-time.After(time.Second):
			f.engine.Go(func() error {
				if err := f.Text.SetText(strconv.Itoa(f.count)); err != nil {
					return fmt.Errorf("failed to set text: %v", err)
				}
				f.count = 0
				if len(f.engine.Composers) == 0 || f.engine.Composers[len(f.engine.Composers)-1] != f {
					f.engine.Remove(f)
					f.engine.Append(f)
				}
				return nil
			})
		case <-f.done:
			return
		}
	}
}
