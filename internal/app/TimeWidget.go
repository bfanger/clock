package app

import (
	"strconv"
	"time"

	"../../internal/engine"
	"github.com/veandco/go-sdl2/ttf"
)

// TimeWidget displays the current time
type TimeWidget struct {
	World       *engine.Container
	Font        *ttf.Font
	Hours       *engine.Text
	Dots        *engine.Text
	Minutes     *engine.Text
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

	timeWidget := &TimeWidget{
		NeedsRedraw: needsRedraw,
		World:       world,
		Font:        font,
		Hours:       hours,
		Dots:        dots,
		Minutes:     minutes,
		Disposed:    make(chan bool)}

	world.Add(hours)
	world.Add(dots)
	world.Add(minutes)

	timeWidget.Update()

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
	if err := timeWidget.Hours.Dispose(); err != nil {
		return err
	}
	if err := timeWidget.Minutes.Dispose(); err != nil {
		return err
	}
	return timeWidget.Dots.Dispose()
}

const left int32 = 178

// Update based on current time and center-align the elements.
func (timeWidget *TimeWidget) Update() {
	now := time.Now().Local()
	timeWidget.Hours.Content = strconv.Itoa(now.Hour())
	if err := timeWidget.Hours.Update(); err != nil {
		panic(err)
	}
	timeWidget.Minutes.Content = now.Format("04")
	if err := timeWidget.Minutes.Update(); err != nil {
		panic(err)
	}
	offset := -44 + (timeWidget.Hours.Texture.Destination.W / 2)
	timeWidget.Hours.Texture.Destination.X = left - timeWidget.Hours.Texture.Destination.W + offset
	timeWidget.Dots.Texture.Destination.X = left + 4 + offset
	timeWidget.Minutes.Texture.Destination.X = left + 24 + offset
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
			timeWidget.Update()
			timeWidget.NeedsRedraw <- true
		}
	}
}
