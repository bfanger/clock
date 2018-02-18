package app

import (
	"time"

	"../../internal/engine"
	"github.com/veandco/go-sdl2/ttf"
)

// ClockWidget displays the current time
type ClockWidget struct {
	Parent     engine.ContainerInterface
	Font       *ttf.Font
	Container  *engine.Container
	Background *engine.Texture
	Hours      *engine.Text
	Dots       *engine.Text
	Minutes    *engine.Text
	Timer      *time.Timer
}

// Mount activates the clock
func (clockWidget *ClockWidget) Mount(parent engine.ContainerInterface) error {
	// Background
	background, err := engine.TextureFromImage(ResourcePath("time_background.png"))
	if err != nil {
		return err
	}
	background.Destination.Y = 84

	// Text
	font, err := ttf.OpenFont(ResourcePath("Teko-Light.ttf"), 135)
	if err != nil {
		return err
	}

	hours, err := engine.NewText(
		font,
		White(),
		"--")
	if err != nil {
		return err
	}
	hours.Texture.Destination.Y = 80

	dotFont, err := ttf.OpenFont(ResourcePath("Teko-Light.ttf"), 110)
	if err != nil {
		return err
	}
	defer dotFont.Close()
	dots, err := engine.NewText(
		dotFont,
		White(),
		":")
	if err != nil {
		return err
	}
	dots.Texture.Destination.Y = 90

	minutes, err := engine.NewText(
		font,
		White(),
		"--")
	if err != nil {
		return err
	}
	minutes.Texture.Destination.Y = 80

	container := engine.Container{}
	container.Add(background)
	container.Add(hours)
	container.Add(dots)
	container.Add(minutes)

	clockWidget.Parent = parent
	clockWidget.Font = font
	clockWidget.Container = &container
	clockWidget.Background = background
	clockWidget.Hours = hours
	clockWidget.Dots = dots
	clockWidget.Minutes = minutes
	clockWidget.tick()

	parent.Add(&container)
	return nil
}

// Unmount Dispose resources
func (clockWidget *ClockWidget) Unmount() error {
	if err := clockWidget.Parent.Remove(clockWidget.Container); err != nil {
		return err
	}
	clockWidget.Timer.Stop()
	clockWidget.Font.Close()
	return clockWidget.Container.DisposeItems()
}

const left int32 = 178

// Redraw based on current time and center-align the elements.
func (clockWidget *ClockWidget) Redraw() error {
	now := time.Now().Local()
	clockWidget.Hours.Content = now.Format("3")
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

func (clockWidget *ClockWidget) tick() {
	clockWidget.Redraw()

	// Calculate the delay to the start of the next minute
	started := time.Now().Local()
	delay := (time.Duration(1) * time.Minute)
	delay -= (time.Duration(started.Second()) * time.Second)
	delay -= (time.Duration(started.Nanosecond()) * time.Nanosecond)

	clockWidget.Timer = engine.Timeout(clockWidget.tick, delay)
}
