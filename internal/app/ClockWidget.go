package app

import (
	"strconv"
	"time"

	"../../internal/engine"
	"github.com/veandco/go-sdl2/ttf"
)

// ClockWidget displays the current time
type ClockWidget struct {
	World         *engine.Container
	Font          *ttf.Font
	Hours         *engine.Text
	Dots          *engine.Text
	Minutes       *engine.Text
	Background    *engine.Texture
	RequestUpdate chan Widget
	Disposed      chan bool
}

// NewClockWidget creates an active ClockWidget
func NewClockWidget(world *engine.Container, requestUpdate chan Widget) (*ClockWidget, error) {

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

	hours, err := engine.NewText(
		font,
		White(),
		"--",
		world.Renderer)
	if err != nil {
		return nil, err
	}
	hours.Texture.Destination.Y = 80

	dotFont, err := ttf.OpenFont(ResourcePath("Teko-Light.ttf"), 110)
	if err != nil {
		return nil, err
	}
	defer dotFont.Close()
	dots, err := engine.NewText(
		dotFont,
		White(),
		":",
		world.Renderer)
	if err != nil {
		return nil, err
	}
	dots.Texture.Destination.Y = 90

	minutes, err := engine.NewText(
		font,
		White(),
		"--",
		world.Renderer)
	if err != nil {
		return nil, err
	}
	minutes.Texture.Destination.Y = 80

	clockWidget := &ClockWidget{
		RequestUpdate: requestUpdate,
		World:         world,
		Font:          font,
		Hours:         hours,
		Dots:          dots,
		Minutes:       minutes,
		Background:    background,
		Disposed:      make(chan bool)}

	world.Add(hours)
	world.Add(dots)
	world.Add(minutes)

	clockWidget.Update()

	go clockWidgetLifecycle(clockWidget)

	return clockWidget, nil
}

// Dispose resources
func (clockWidget *ClockWidget) Dispose() error {
	clockWidget.Disposed <- true
	clockWidget.Font.Close()
	if err := clockWidget.Background.Dispose(); err != nil {
		return err
	}
	if err := clockWidget.Hours.Dispose(); err != nil {
		return err
	}
	if err := clockWidget.Minutes.Dispose(); err != nil {
		return err
	}
	return clockWidget.Dots.Dispose()
}

const left int32 = 178

// Update based on current time and center-align the elements.
func (clockWidget *ClockWidget) Update() error {
	now := time.Now().Local()
	clockWidget.Hours.Content = strconv.Itoa(now.Hour())
	if err := clockWidget.Hours.Update(); err != nil {
		return err
	}
	clockWidget.Minutes.Content = now.Format("04")
	if err := clockWidget.Minutes.Update(); err != nil {
		return err
	}
	offset := -44 + (clockWidget.Hours.Texture.Destination.W / 2)
	clockWidget.Hours.Texture.Destination.X = left - clockWidget.Hours.Texture.Destination.W + offset
	clockWidget.Dots.Texture.Destination.X = left + 4 + offset
	clockWidget.Minutes.Texture.Destination.X = left + 24 + offset
	return nil
}

func clockWidgetLifecycle(clockWidget *ClockWidget) {
	for {
		// Calculate the delay to the start of the next minute
		started := time.Now().Local()
		delay := (time.Duration(1) * time.Minute)
		delay -= (time.Duration(started.Second()) * time.Second)
		delay -= (time.Duration(started.Nanosecond()) * time.Nanosecond)

		select {
		case <-clockWidget.Disposed:
			return
		case <-time.After(delay):
			clockWidget.RequestUpdate <- clockWidget
		}
	}
}
