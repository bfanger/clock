package events

import (
	"fmt"

	"github.com/bfanger/clock/display"
	"github.com/veandco/go-sdl2/sdl"
)

// RefreshEvent triggers a Render
const RefreshEvent int32 = 808

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
	e := sdl.UserEvent{
		Type:      sdl.USEREVENT,
		Timestamp: sdl.GetTicks(),
		WindowID:  0,
		Code:      RefreshEvent,
	}
	sdl.PushEvent(&e)
	<-refreshed
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
			if e.Code == RefreshEvent {
				r.C <- true
				refreshed <- true
			}
		}
	}
}
