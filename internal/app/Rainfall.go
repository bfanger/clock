package app

import (
	"log"
	"math"
	"time"

	"github.com/bfanger/clock/pkg/ui"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	rainfallGap    = int32(4)
	rainfallBar    = int32(19)
	rainfallHeight = int32(64)
)

// Rainfall displays the forecasted rainfall
type Rainfall struct {
	engine  *ui.Engine
	graph   *ui.Image
	texture *sdl.Texture
	start   time.Time
}

// NewRainfall creates a new rainfall widget
func NewRainfall(engine *ui.Engine) (*Rainfall, error) {
	r := &Rainfall{
		engine: engine,
	}
	var err error
	if r.graph, err = ui.ImageFromFile(Asset("rainfall/gradient.png"), engine.Renderer); err != nil {
		return nil, err
	}
	if r.texture, err = engine.Renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_TARGET, rainfallBar*24, rainfallHeight); err != nil {
		return nil, err
	}

	return r, nil
}

// Compose renders the current state of the rainfall widget
func (r *Rainfall) Compose(renderer *sdl.Renderer) error {
	if r.start.IsZero() {
		return nil
	}
	d := int32(time.Since(r.start) / time.Minute)
	quarter := 3 * rainfallBar
	offsetX := screenWidth - 3*rainfallGap - 4*quarter
	renderer.Copy(r.texture, &sdl.Rect{X: 3 * d, Y: 0, W: quarter, H: rainfallHeight}, &sdl.Rect{X: offsetX, Y: screenHeight - rainfallHeight, W: quarter, H: rainfallHeight})
	renderer.Copy(r.texture, &sdl.Rect{X: 3*d + quarter, Y: 0, W: quarter, H: rainfallHeight}, &sdl.Rect{X: offsetX + rainfallGap + quarter, Y: screenHeight - rainfallHeight, W: quarter, H: rainfallHeight})
	renderer.Copy(r.texture, &sdl.Rect{X: 3*d + 2*quarter, Y: 0, W: quarter, H: rainfallHeight}, &sdl.Rect{X: offsetX + 2*rainfallGap + 2*quarter, Y: screenHeight - rainfallHeight, W: quarter, H: rainfallHeight})
	renderer.Copy(r.texture, &sdl.Rect{X: 3*d + 3*quarter, Y: 0, W: quarter, H: rainfallHeight}, &sdl.Rect{X: offsetX + 3*rainfallGap + 3*quarter, Y: screenHeight - rainfallHeight, W: quarter, H: rainfallHeight})
	return nil
}

// Close free resources
func (r *Rainfall) Close() error {
	if err := r.graph.Close(); err != nil {
		return err
	}
	if err := r.texture.Destroy(); err != nil {
		return err
	}
	return nil
}

func (r *Rainfall) SetForecasts(forecasts []RainfallForecast) {
	r.engine.Go(func() error {
		prevTarget := r.engine.Renderer.GetRenderTarget()
		if err := r.engine.Renderer.SetRenderTarget(r.texture); err != nil {
			return err
		}
		defer r.engine.Renderer.SetRenderTarget(prevTarget)
		if err := r.engine.Renderer.Clear(); err != nil {
			return err
		}
		if len(forecasts) < 1 {
			log.Println("no forecasts data (yet)")
			return nil
		}
		if forecasts[len(forecasts)-1].Timestamp.Before(time.Now().Add(time.Hour)) {
			log.Println("no relevant forecasts available")
			return nil

		}

		r.start = forecasts[0].Timestamp
		for i, forecast := range forecasts {
			height := int32(math.Ceil(float64(rainfallHeight) * forecast.Factor))
			if height < 1 {
				height = 1
			}
			src := sdl.Rect{X: 0, Y: rainfallHeight - height, W: rainfallBar, H: height}
			dest := src
			dest.X = int32(i) * rainfallBar
			if err := r.engine.Renderer.Copy(r.graph.Texture, &src, &dest); err != nil {
				return err
			}
		}

		return nil
	})
}

type RainfallForecast struct {
	Timestamp time.Time
	Factor    float64
}
