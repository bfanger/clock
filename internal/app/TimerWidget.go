package app

import (
	"fmt"
	"time"

	"../../internal/engine"
	"github.com/veandco/go-sdl2/ttf"
)

// TimerWidget displays a count-down timer.
type TimerWidget struct {
	World             *engine.Container
	Container         *engine.Container
	HideableContainer *engine.Hideable
	Font              *ttf.Font
	Hour              int
	Minute            int
	Countdown         int64
	Blink             int64
	Timer             *engine.Text
	BlinkingTimer     *engine.Hideable
	Background        *engine.Texture
	Visible           bool
	RequestUpdate     chan Widget
	Disposed          chan bool
}

// NewTimerWidget creates an active TimerWidget
func NewTimerWidget(backgroundPath string, hour int, minute int, world *engine.Container, requestUpdate chan Widget) (*TimerWidget, error) {

	// Background
	background, err := engine.TextureFromImage(world.Renderer, ResourcePath(backgroundPath))
	if err != nil {
		return nil, err
	}

	// Text
	font, err := ttf.OpenFont(ResourcePath("Teko-Light.ttf"), 48)
	if err != nil {
		return nil, err
	}

	timer, err := engine.NewText(
		font,
		Black(),
		"--:--",
		world.Renderer)
	if err != nil {
		return nil, err
	}
	timer.Texture.Destination.X = 20
	timer.Texture.Destination.Y = 10

	container := engine.NewContainer(world.Renderer)

	timerWidget := &TimerWidget{
		RequestUpdate:     requestUpdate,
		World:             world,
		Font:              font,
		Hour:              hour,
		Minute:            minute,
		Countdown:         900,  // 15 min
		Blink:             -120, // blink for 2 min
		Timer:             timer,
		BlinkingTimer:     engine.NewHideable(timer),
		Background:        background,
		Container:         container,
		HideableContainer: engine.NewHideable(container),
		Visible:           false,
		Disposed:          make(chan bool)}

	timerWidget.Update()

	container.Add(timerWidget.Background)
	container.Add(timerWidget.BlinkingTimer)
	world.Add(timerWidget.HideableContainer)

	go timerWidgetLifecycle(timerWidget)

	return timerWidget, nil
}

// Dispose resources
func (timerWidget *TimerWidget) Dispose() error {
	timerWidget.Disposed <- true
	timerWidget.Font.Close()
	if err := timerWidget.Background.Dispose(); err != nil {
		return err
	}
	if err := timerWidget.Timer.Dispose(); err != nil {
		return err
	}

	return timerWidget.Container.Dispose()
}

// Update based on current time and center-align the elements.
func (timerWidget *TimerWidget) Update() error {
	remaining := timerWidget.secondsRemaining()
	if !timerWidget.HideableContainer.Visible {
		return nil
	}

	if remaining < 0 {
		if remaining < -120 {
			timerWidget.HideableContainer.Visible = false
		} else {
			timerWidget.Timer.Color = Red()
			timerWidget.Timer.Content = "00:00" //fmt.Sprintf("-%d:%02d", int(remaining*-1)/60, int(remaining*-1)%60)
			timerWidget.BlinkingTimer.Visible = time.Now().Second()%2 == 0
			if err := timerWidget.Timer.Update(); err != nil {
				return err
			}
		}
	} else {
		timerWidget.Timer.Color = Black()
		timerWidget.Timer.Content = fmt.Sprintf("%d:%02d", int(remaining)/60, int(remaining)%60)
		if err := timerWidget.Timer.Update(); err != nil {
			return err
		}
	}
	return nil
}

func timerWidgetLifecycle(timerWidget *TimerWidget) {
	var delay time.Duration

	for {
		remaining := timerWidget.secondsRemaining()
		currentlyVisible := timerWidget.HideableContainer.Visible
		timerWidget.HideableContainer.Visible = remaining > timerWidget.Blink && remaining < timerWidget.Countdown
		if timerWidget.HideableContainer.Visible {
			// Update every second when visible
			started := time.Now().Local()
			delay = (time.Duration(1) * time.Second)
			delay -= (time.Duration(started.Nanosecond()) * time.Nanosecond)
		} else if remaining > 0 {
			//delay  count down to 15 minutes before
			delay = (time.Duration(remaining-timerWidget.Countdown) * time.Second)
		} else {
			if currentlyVisible {
				timerWidget.RequestUpdate <- timerWidget
			}
			// @todo calculate the delay until the next day
			delay = (time.Duration(1) * time.Hour)
		}
		select {
		case <-timerWidget.Disposed:
			return
		case <-time.After(delay):
			timerWidget.RequestUpdate <- timerWidget

		}
	}
}
func (timerWidget *TimerWidget) secondsRemaining() int64 {
	now := time.Now().Local()
	target := time.Date(now.Year(), now.Month(), now.Day(), timerWidget.Hour, timerWidget.Minute, 0, 0, time.Local)
	return target.Unix() - now.Unix()
}
