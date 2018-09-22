package main

import (
	"log"
	"runtime"
	"sync"

	"github.com/bfanger/clock/internal/app"
	"github.com/bfanger/clock/pkg/ui"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	runtime.LockOSThread()

	display, err := app.NewDisplay()
	if err != nil {
		log.Fatal(err)
	}
	defer display.Close()

	engine := ui.NewEngine(display.Renderer)
	time, err := app.NewTime(engine)
	if err != nil {
		log.Fatal(err)
	}
	defer time.Close()
	ready := sync.Once{}
	intro := func() {
		if err := time.Intro(); err != nil {
			log.Fatal(err)
		}
	}

	windowEvents := map[uint8]string{
		sdl.WINDOWEVENT_NONE:         "NONE",
		sdl.WINDOWEVENT_SHOWN:        "SHOWN",
		sdl.WINDOWEVENT_HIDDEN:       "HIDDEN",
		sdl.WINDOWEVENT_EXPOSED:      "EXPOSED",
		sdl.WINDOWEVENT_MOVED:        "MOVED",
		sdl.WINDOWEVENT_RESIZED:      "RESIZED",
		sdl.WINDOWEVENT_SIZE_CHANGED: "SIZE_CHANGED",
		sdl.WINDOWEVENT_MINIMIZED:    "MINIMIZED",
		sdl.WINDOWEVENT_MAXIMIZED:    "MAXIMIZED",
		sdl.WINDOWEVENT_RESTORED:     "RESTORED",
		sdl.WINDOWEVENT_ENTER:        "ENTER",
		sdl.WINDOWEVENT_LEAVE:        "LEAVE",
		sdl.WINDOWEVENT_FOCUS_GAINED: "FOCUS_GAINED",
		sdl.WINDOWEVENT_FOCUS_LOST:   "FOCUS_LOST",
		sdl.WINDOWEVENT_CLOSE:        "CLOSE",
		sdl.WINDOWEVENT_TAKE_FOCUS:   "TAKE_FOCUS",
		sdl.WINDOWEVENT_HIT_TEST:     "HIT_TEST"}

	err = engine.EventLoop(func(event sdl.Event) {
		switch e := event.(type) {
		case *sdl.MouseButtonEvent:
			if e.Type == sdl.MOUSEBUTTONUP {
				// tap
			}
		case *sdl.WindowEvent:
			ready.Do(intro)
			log.Printf("WindowEvent: %s\n", windowEvents[e.Event])
		default:
			// log.Printf("%T %v\n", event, event)
		}
	})
	if err != nil {
		log.Fatal(err)
	}

}
