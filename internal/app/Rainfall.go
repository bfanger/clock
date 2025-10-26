package app

import (
	"log"
	"time"

	"github.com/bfanger/clock/pkg/ui"
	"github.com/veandco/go-sdl2/sdl"
)

// Rainfall displays the forecasted rainfall
type Rainfall struct {
	forecasts  []RainfallForecast
	engine     *ui.Engine
	background *ui.Image
	foreground *ui.Image
	graph      *ui.Image
	display    *ui.Image
	sprite     *ui.Sprite
	done       chan bool
}

// NewRainfall creates a new rainfall widget
func NewRainfall(engine *ui.Engine) (*Rainfall, error) {
	r := &Rainfall{
		engine: engine,
	}
	var err error
	if r.background, err = ui.ImageFromFile(Asset("rainfall/gradient.png"), engine.Renderer); err != nil {
		return nil, err
	}
	if r.foreground, err = ui.ImageFromFile(Asset("rainfall/timelines.png"), engine.Renderer); err != nil {
		return nil, err
	}
	if r.graph, err = ui.ImageFromFile(Asset("rainfall/forecast.png"), engine.Renderer); err != nil {
		return nil, err
	}
	texture, err := engine.Renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_TARGET, 240, 64)
	if err != nil {
		return nil, err
	}
	r.display = &ui.Image{Texture: texture, Frame: sdl.Rect{W: 240, H: 64}}

	r.sprite = ui.NewSprite(r.display)
	r.sprite.X = screenWidth
	r.sprite.Y = screenHeight
	r.sprite.AnchorX = 1
	r.sprite.AnchorY = 1

	go r.tick()
	return r, nil
}

// Compose renders the current state of the rainfall widget
func (r *Rainfall) Compose(renderer *sdl.Renderer) error {
	return r.sprite.Compose(renderer)
}

// Close free resources
func (r *Rainfall) Close() error {
	if err := r.background.Close(); err != nil {
		return err
	}
	if err := r.graph.Close(); err != nil {
		return err
	}
	if err := r.display.Close(); err != nil {
		return err
	}
	close(r.done)
	return nil
}

func (r *Rainfall) SetForecasts(forecasts []RainfallForecast) {
	r.engine.Go(func() error {
		r.forecasts = forecasts
		return r.update()
	})
}

func (r *Rainfall) tick() {
	for {
		select {
		case <-r.done:
			return
		case <-time.After(30 * time.Second):
			r.engine.Go(r.update)
		}
	}
}

func (r *Rainfall) update() error {
	prevTarget := r.engine.Renderer.GetRenderTarget()
	if err := r.engine.Renderer.SetRenderTarget(r.display.Texture); err != nil {
		return err
	}
	defer r.engine.Renderer.SetRenderTarget(prevTarget)
	if err := r.engine.Renderer.Clear(); err != nil {
		return err
	}
	if len(r.forecasts) < 1 {
		log.Println("no forecasts data (yet)")
		return nil
	}
	if r.forecasts[len(r.forecasts)-1].Timestamp.Before(time.Now().Add(time.Hour)) {
		log.Println("no relevant forecasts available")
		return nil
	}
	if err := r.engine.Renderer.Copy(r.background.Texture, &r.background.Frame, &r.background.Frame); err != nil {
		return err
	}
	start := time.Now().UTC()

	for _, forecast := range r.forecasts {
		d := forecast.Timestamp.Sub(start)
		if d < -5*time.Minute || d > time.Hour {
			continue
		}
		dest := sdl.Rect{W: 20, H: 64}
		dest.X = int32((float64(d) / float64(time.Hour)) * 240.0)
		dest.Y = int32(-64.0 * (forecast.Percentage / 100.0))

		if err := r.engine.Renderer.Copy(r.graph.Texture, &r.graph.Frame, &dest); err != nil {
			return err
		}
	}
	if err := r.engine.Renderer.Copy(r.foreground.Texture, &r.foreground.Frame, &r.foreground.Frame); err != nil {
		return err
	}

	return nil
}

type RainfallForecast struct {
	Timestamp  time.Time
	Percentage float64
}
