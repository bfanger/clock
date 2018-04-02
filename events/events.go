package events

import (
	"fmt"
	"sync"

	"github.com/bfanger/clock/display"
	"github.com/veandco/go-sdl2/sdl"
)

var events chan sdl.Event
var eventMutex sync.RWMutex
var refresh bool
var refreshMutex sync.Mutex

var (
	windowListeners []chan *sdl.WindowEvent
)

const (
	refreshEvent int32 = 0 //int itoa
)

// Init creates the event channel and starts multiplexing the events
func Init() {
	events = make(chan sdl.Event)

	go func() {
		events := OnWindowEvents()
		for e := range events {
			if e.Event == sdl.WINDOWEVENT_EXPOSED {
				Refresh()
			}
		}
	}()

	go func() {
		for e := range events {
			switch t := e.(type) {
			case *sdl.QuitEvent:
			case *sdl.WindowEvent:
				eventMutex.RLock()
				for _, c := range windowListeners {
					select {
					case c <- t:
					default:
						fmt.Println("blocked on WindowEvent")
						c <- t
						fmt.Println("resumed")
					}
				}
				eventMutex.RUnlock()
			default:
				// fmt.Printf("%T%+v\n", e, e)
			}
		}
	}()
}

// Quit closes the event stream
func Quit() {
	// @todo Remove eventListeners
	close(events)
}

// Refresh triggers a screen update
func Refresh() {
	refreshMutex.Lock()
	defer refreshMutex.Unlock()
	if refresh == false {
		e := &sdl.UserEvent{
			Type:      sdl.USEREVENT,
			Timestamp: sdl.GetTicks(),
			WindowID:  0,
			Code:      refreshEvent}

		sdl.PushEvent(e)
	}
	refresh = true
}

// OnWindowEvents start listening to window events
func OnWindowEvents() chan *sdl.WindowEvent {
	eventMutex.Lock()
	defer eventMutex.Unlock()
	c := make(chan *sdl.WindowEvent)
	windowListeners = append(windowListeners, c)
	return c
}

// EventLoop start the main event loop and keep running until a quit event
func EventLoop(r *display.Renderer) error {
	if events == nil {
		return fmt.Errorf("events not initialized")
	}
	var e sdl.Event
	for {
		refreshMutex.Lock()
		if refresh {
			r.C <- true
			refresh = false
		}
		refreshMutex.Unlock()
		e = sdl.WaitEvent()
		switch e.(type) {
		case *sdl.QuitEvent:
			return nil
		}
		events <- e
	}
}
