package display

import (
	"math"
	"strconv"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

// Fps displays the average framerate over the last 10 frames (when added as an animation)
type Fps struct {
	text *Text
	avg  []time.Duration
}

// NewFps create a new Fps and updates every minute
func NewFps(r *Renderer, font string, fontSize int) *Fps {
	white := sdl.Color{R: 255, G: 255, B: 255, A: 255}
	text := NewText(font, fontSize, white, "-")
	return &Fps{
		text: text,
	}
}

// Destroy the Fps
func (f *Fps) Destroy() error {
	return f.text.Destroy()
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
	v := math.Round(ps / total)
	f.text.Text = strconv.Itoa(int(v))

	return false
}

// Paint the FPS
func (f *Fps) Paint(r *sdl.Renderer) (*Texture, error) {
	return f.text.Paint(r)
}
