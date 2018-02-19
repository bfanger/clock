package app

import (
	"fmt"
	"strconv"
	"time"

	"../../internal/engine"
	"github.com/veandco/go-sdl2/ttf"
)

// TimerWidget displays a count-down timer.
type TimerWidget struct {
	Parent            engine.ContainerInterface
	HideableContainer *engine.Hideable
	Container         *engine.Container
	Font              *ttf.Font
	Hour              int
	Minute            int
	Second            int
	Repeat            bool
	Countdown         int64
	Blink             int64
	Alarm             string
	TimeLeft          *engine.Text
	BlinkingTimer     *engine.Hideable
	Background        *engine.Texture
	Dial              *engine.Texture
	Timeout           *engine.Timeout
}

// NewTimerWidget creates an active TimerWidget
func NewTimerWidget(backgroundPath string, hour int, minute int) (*TimerWidget, error) {

	// Background
	background, err := engine.TextureFromImage(ResourcePath(backgroundPath))
	if err != nil {
		return nil, err
	}

	// Text
	font, err := ttf.OpenFont(ResourcePath("Roboto-Medium.ttf"), 50)
	if err != nil {
		return nil, err
	}

	timeLeft, err := engine.NewText(
		font,
		Black(),
		"-")
	if err != nil {
		return nil, err
	}
	timeLeft.Texture.Destination.Y = 10

	dial, err := engine.TextureFromImage(ResourcePath("dial.png"))
	if err != nil {
		return nil, err
	}
	dial.Destination.X = 248
	dial.Destination.Y = 8
	dial.Destination.W = 64
	dial.Destination.H = 64
	dial.Frame.W = 64
	dial.Frame.H = 64

	container := engine.Container{}

	timerWidget := &TimerWidget{
		// RequestUpdate:     requestUpdate,
		Font:              font,
		Hour:              hour,
		Minute:            minute,
		Countdown:         900,  // 15 min
		Blink:             -120, // blink for 2 min
		Alarm:             "0",
		TimeLeft:          timeLeft,
		Dial:              dial,
		BlinkingTimer:     engine.NewHideable(timeLeft),
		Background:        background,
		Container:         &container,
		HideableContainer: engine.NewHideable(&container)}

	container.Add(timerWidget.Background)
	container.Add(timerWidget.BlinkingTimer)
	container.Add(timerWidget.Dial)

	return timerWidget, nil
}

// Mount timer and update timer & visibity
func (timerWidget *TimerWidget) Mount(parent engine.ContainerInterface) error {
	parent.Add(timerWidget.HideableContainer)
	timerWidget.Parent = parent
	timerWidget.tick()
	return nil
}

// Unmount and dispose resources
func (timerWidget *TimerWidget) Unmount() error {
	timerWidget.Timeout.Cancel()
	timerWidget.Font.Close()
	if err := timerWidget.Parent.Remove(timerWidget.HideableContainer); err != nil {
		return err
	}
	return timerWidget.Container.DisposeItems()
}

// Redraw based on current time and center-align the elements.
func (timerWidget *TimerWidget) Redraw() error {
	remaining := timerWidget.secondsRemaining()
	if !timerWidget.HideableContainer.Visible {
		return nil
	}
	if remaining <= 0 {
		timerWidget.Dial.Frame.Y = 64 * 7
		if remaining < -120 {
			timerWidget.HideableContainer.Visible = false
		} else {
			timerWidget.TimeLeft.Color = Red()
			timerWidget.TimeLeft.Content = timerWidget.Alarm //fmt.Sprintf("-%d:%02d", int(remaining*-1)/60, int(remaining*-1)%60)
			timerWidget.BlinkingTimer.Visible = time.Now().Second()%2 == 0
			if err := timerWidget.TimeLeft.Update(); err != nil {
				return err
			}
		}
	} else {
		timerWidget.TimeLeft.Color = Black()
		timerWidget.BlinkingTimer.Visible = true
		minutes := int(remaining) / 60
		seconds := int(remaining) % 60
		if minutes == 0 {
			timerWidget.TimeLeft.Content = strconv.Itoa(seconds)
		} else {
			timerWidget.TimeLeft.Content = fmt.Sprintf("%d:%02d", minutes, seconds)
		}
		if err := timerWidget.TimeLeft.Update(); err != nil {
			return err
		}
		timerWidget.TimeLeft.Texture.Destination.X = (320 / 2) - (timerWidget.TimeLeft.Texture.Frame.W / 2)
		frame := int32((900 - int(remaining)) / 10)
		y := frame / 16
		x := frame % 16
		timerWidget.Dial.Frame.X = x * 64
		timerWidget.Dial.Frame.Y = y * 64
	}
	return nil
}

func (timerWidget *TimerWidget) tick() {
	var delay time.Duration

	remaining := timerWidget.secondsRemaining()
	timerWidget.HideableContainer.Visible = remaining > timerWidget.Blink && remaining < timerWidget.Countdown
	completed := false
	if timerWidget.HideableContainer.Visible {
		// Update every second when visible
		started := time.Now().Local()
		delay = time.Second
		delay -= time.Duration(started.Nanosecond()) * time.Nanosecond
	} else if remaining > 0 {
		// delay count down to 15 minutes before
		delay = time.Duration(remaining-timerWidget.Countdown) * time.Second
		if delay < time.Second {
			delay = time.Second
		}
	} else if timerWidget.Repeat {
		// calculate the delay until the next day
		delay = (24 * time.Hour) + time.Duration(remaining-timerWidget.Countdown)*time.Second
	} else {
		completed = true
	}
	if err := timerWidget.Redraw(); err != nil {
		panic(err)
	}
	if !completed {
		timerWidget.Timeout = engine.SetTimeout(timerWidget.tick, delay)
	}
}

func (timerWidget *TimerWidget) secondsRemaining() int64 {
	now := time.Now().Local()
	target := time.Date(now.Year(), now.Month(), now.Day(), timerWidget.Hour, timerWidget.Minute, timerWidget.Second, 0, time.Local)
	return target.Unix() - now.Unix()
}

var timerWidgetCountdown = 0
var timerWidget *TimerWidget

// TimerWidgetButtonHandler add a timer for 5, 10 or 15 minutes
func TimerWidgetButtonHandler(parent engine.ContainerInterface) {
	if timerWidget != nil {
		timerWidget.Unmount()
		timerWidget = nil
	}
	timerWidgetCountdown += 5
	if timerWidgetCountdown == 20 {
		timerWidgetCountdown = 0
		return
	}
	now := time.Now()
	var err error
	timerWidget, err = NewTimerWidget("timer_background.png", now.Hour(), now.Minute()+timerWidgetCountdown)
	if err != nil {
		panic(err)
	}
	timerWidget.Second = now.Second()
	timerWidget.Blink = -10
	timerWidget.Mount(parent)
}
