package app

import (
	"time"

	"../../internal/engine"
	"github.com/veandco/go-sdl2/ttf"
)

// TimeWidget displays the current time
type TimeWidget struct {
	World       *engine.Container
	Font        *ttf.Font
	Time        *engine.Text
	Background  *engine.Texture
	NeedsRedraw chan bool
	Disposed    chan bool
}

// NewTimeWidget creates an active TimeWidget
func NewTimeWidget(world *engine.Container, needsRedraw chan bool) (*TimeWidget, error) {

	// Background
	background, err := engine.TextureFromImage(world.Renderer, ResourcePath("time_background.png"))
	if err != nil {
		return nil, err
	}
	background.Destination.Y = 84
	world.Add(background)

	// Text
	font, err := ttf.OpenFont(ResourcePath("Teko-Light.ttf"), 135)
	if err != nil {
		return nil, err
	}

	date := time.Now()
	text, err := engine.NewText(
		font,
		engine.White(),
		date.Format("15:04"),
		world.Renderer)
	if err != nil {
		return nil, err
	}
	text.Texture.Destination.X = 95
	text.Texture.Destination.Y = 80
	world.Add(text)

	timeWidget := &TimeWidget{
		NeedsRedraw: needsRedraw,
		World:       world,
		Font:        font,
		Time:        text,
		Disposed:    make(chan bool)}

	go timeWidgetLifecycle(timeWidget)

	return timeWidget, nil
}

// Dispose resources
func (timeWidget *TimeWidget) Dispose() error {
	timeWidget.Disposed <- true
	timeWidget.Font.Close()
	if err := timeWidget.Background.Dispose(); err != nil {
		return err
	}
	return timeWidget.Time.Dispose()
}

func timeWidgetLifecycle(timeWidget *TimeWidget) {
	for {
		// Calculate the delay to the start of the next minute
		started := time.Now().Local()
		delay := (time.Duration(1) * time.Minute)
		delay -= (time.Duration(started.Second()) * time.Second)
		delay -= (time.Duration(started.Nanosecond()) * time.Nanosecond)

		select {
		case <-timeWidget.Disposed:
			return
		case <-time.After(delay):
			now := time.Now().Local()
			timeWidget.Time.Content = now.Format("15:04")
			if err := timeWidget.Time.Update(); err != nil {
				panic(err)
			}
			timeWidget.NeedsRedraw <- true
		}
	}
}
