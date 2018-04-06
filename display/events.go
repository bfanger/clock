package display

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	// refreshEvent triggers a Render
	refreshEvent int32 = iota + 808000
	// quitEvent stops the eventLoop
	quitEvent
)

var refreshed chan bool

// Init creates the event channel and starts multiplexing the events
func Init(r *Renderer) {
	refreshed = make(chan bool)
}

// Quit closes the event stream
func Quit() {
	close(refreshed)
}

// Refresh triggers a screen update
func Refresh() {
	// @todo prevent builing a queue of refresh events.
	e := sdl.UserEvent{
		Type:      sdl.USEREVENT,
		Timestamp: sdl.GetTicks(),
		WindowID:  0,
		Code:      refreshEvent,
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
		Code:      quitEvent,
	}
	sdl.PushEvent(&e)
}

// EventLoop start the main event loop and keep running until a quit event
func EventLoop(r *Renderer) error {
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
			case quitEvent:
				return nil
			case refreshEvent:
				r.C <- true
				refreshed <- true
			}
		}
	}
}
