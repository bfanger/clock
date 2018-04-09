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
	var e sdl.Event
	var quit, dirty, refresh, render bool
	var refreshQueue int
	for {
		e = sdl.WaitEvent()
		quit, render, refresh = handleEvent(e)
		if quit {
			return nil
		}
		if refresh {
			refreshQueue++
		}
		for {
			e = sdl.PollEvent()
			if e == nil {
				break
			}
			quit, dirty, refresh = handleEvent(e)
			if quit {
				return nil
			}
			if refresh {
				refreshQueue++
			}
			if dirty {
				render = true
			}
		}
		if render {
			r.refresh <- true
		}
		for i := 0; i < refreshQueue; i++ {
			refreshed <- true
		}
		refreshQueue = 0
	}
}

func handleEvent(event sdl.Event) (quit bool, render bool, refresh bool) {
	// log.Printf("%T %+v\n", event, event)
	switch e := event.(type) {
	case *sdl.QuitEvent:
		quit = true
	case *sdl.WindowEvent:
		if e.Event == sdl.WINDOWEVENT_EXPOSED {
			render = true
		}
	case *sdl.UserEvent:
		switch e.Code {
		case quitEvent:
			quit = true
		case refreshEvent:
			render = true
			refresh = true
		}
	}
	return
}
