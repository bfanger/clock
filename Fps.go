package main

import (
	"math"
	"strconv"
	"time"

	"github.com/bfanger/clock/display"
	"github.com/bfanger/clock/sprite"
	"github.com/veandco/go-sdl2/sdl"
)

// Fps displays the average framerate over the last 10 frames
type Fps struct {
	Layer display.Layer
	text  *display.Text
	avg   []time.Duration
	r     *display.Renderer
}

// NewFps create a new Fps and updates every minute
func NewFps(r *display.Renderer, font string, opts ...sprite.Option) *Fps {
	white := sdl.Color{R: 255, G: 255, B: 255, A: 255}
	text := display.NewText(font, 16, white, "-")
	l := sprite.New("Fps", text, opts...)
	f := &Fps{
		Layer: l,
		text:  text,
		r:     r,
	}
	r.Animate(f)
	return f
}

// Destroy the Fps
func (f *Fps) Destroy() error {
	f.r.Mutex.Lock()
	defer f.r.Mutex.Unlock()
	err := f.text.Destroy()
	return err
}

// Animate the fps
func (f *Fps) Animate(dt time.Duration) bool {
	const weighted = 10
	if len(f.avg) < weighted {
		f.avg = append(f.avg, dt)
	} else {
		f.avg = append(f.avg[1:], dt)
	}
	var total float64
	for _, t := range f.avg {
		total += float64(t)
	}
	ps := float64(len(f.avg)) * float64(time.Second)

	// f.text.Text = strconv.Itoa(int((n * time.Second) / total))

	v := math.Round(ps / total)
	f.text.Text = strconv.Itoa(int(v))

	return false
}
