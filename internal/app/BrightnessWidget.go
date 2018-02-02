package app

import (
	"time"

	"../../internal/engine"
	"github.com/veandco/go-sdl2/sdl"
)

// BrightnessWidget adapts the brighness to time of day
type BrightnessWidget struct {
	RequestUpdate chan Widget
	World         *engine.Container
	Brightness    *engine.Brightness
	Disposed      chan bool
}

// NewBrightnessWidget creates an active BrightnessWidget
func NewBrightnessWidget(world *engine.Container, requestUpdate chan Widget) (*BrightnessWidget, error) {

	brightnessWidget := &BrightnessWidget{
		RequestUpdate: requestUpdate,
		World:         world}

	var displayMode sdl.DisplayMode
	if err := sdl.GetCurrentDisplayMode(0, &displayMode); err != nil {
		return nil, err
	}

	if displayMode.W <= 320 {
		brightness, err := engine.NewBrightness(world.Renderer, brightnessAlphaForTime(time.Now().Local()))
		if err != nil {
			return nil, err
		}
		world.Add(brightness)
		brightnessWidget.Brightness = brightness
		brightnessWidget.Disposed = make(chan bool)
		go brightnessWidgetLifecycle(brightnessWidget)
	}

	return brightnessWidget, nil
}

// Dispose resources
func (brightnessWidget *BrightnessWidget) Dispose() error {
	if brightnessWidget.Brightness != nil {
		brightnessWidget.Disposed <- true
		return brightnessWidget.Brightness.Dispose()
	}
	return nil
}

// Update the brightness texture
func (brightnessWidget *BrightnessWidget) Update() error {
	return brightnessWidget.Brightness.Update()
}

func brightnessWidgetLifecycle(brightnessWidget *BrightnessWidget) {
	for {
		// Calculate the delay to the start of the next hour
		started := time.Now().Local()
		delay := (time.Duration(1) * time.Hour)
		delay -= (time.Duration(started.Minute()) * time.Minute)
		delay -= (time.Duration(started.Second()) * time.Second)

		select {
		case <-brightnessWidget.Disposed:
			return
		case <-time.After(delay):
			now := time.Now().Local()
			brightnessWidget.Brightness.Alpha = brightnessAlphaForTime(now)
			brightnessWidget.RequestUpdate <- brightnessWidget
		}
	}
}

func brightnessAlphaForTime(time time.Time) uint8 {
	levels := map[int]uint8{
		0:  180,
		1:  180,
		2:  180,
		3:  180,
		4:  160,
		5:  140,
		6:  120,
		7:  100,
		8:  90,
		9:  80,
		10: 50,
		11: 30,
		12: 10,
		13: 10,
		14: 20,
		15: 0,
		16: 80,
		17: 100,
		18: 110,
		19: 120,
		20: 130,
		21: 150,
		22: 160,
		23: 180,
	}
	return levels[time.Hour()]
}
