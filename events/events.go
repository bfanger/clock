package events

import (
	"fmt"

	"github.com/bfanger/clock/display"
	"github.com/veandco/go-sdl2/sdl"
)

// RefreshEvent triggers a Render
const (
	RefreshEvent int32 = iota + 808000
	QuitEvent
)

var refreshed chan bool

// Init creates the event channel and starts multiplexing the events
func Init(r *display.Renderer) {
	refreshed = make(chan bool)
}

// Quit closes the event stream
func Quit() {
	close(refreshed)
}

// Refresh triggers a screen update
func Refresh() {
	// @todo Merge
	e := sdl.UserEvent{
		Type:      sdl.USEREVENT,
		Timestamp: sdl.GetTicks(),
		WindowID:  0,
		Code:      RefreshEvent,
	}
	sdl.PushEvent(&e)
	<-refreshed
}

// Shutdown the event handler
func Shutdown() {
	e := sdl.UserEvent{
		Type:      sdl.USEREVENT,
		Timestamp: sdl.GetTicks(),
		WindowID:  0,
		Code:      QuitEvent,
	}
	sdl.PushEvent(&e)
}

// EventLoop start the main event loop and keep running until a quit event
func EventLoop(r *display.Renderer) error {
	if refreshed == nil {
		return fmt.Errorf("events.Init not called")
	}
	r.C <- true
	for {
		event := sdl.WaitEvent()
		switch e := event.(type) {
		case *sdl.QuitEvent:
			return nil
		case *sdl.WindowEvent:
			if e.Event == sdl.WINDOWEVENT_EXPOSED {
				r.C <- true
			}
		case *sdl.UserEvent:
			switch e.Code {
			case QuitEvent:
				return nil
			case RefreshEvent:
				r.C <- true
				refreshed <- true
			}
		}
	}
}
