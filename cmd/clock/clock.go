package main

import (
	"log"
	"runtime"

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
	maximized := true
	toggle := func() {
		if maximized {
			time.Minimize()
		} else {
			time.Maximize()
		}
		maximized = !maximized
	}

	err = engine.EventLoop(func(event sdl.Event) {
		switch e := event.(type) {
		case *sdl.MouseButtonEvent:
			if e.Type == sdl.MOUSEBUTTONUP {
				toggle()
			}
		case *sdl.KeyboardEvent:
			if e.Type == sdl.KEYDOWN {
				toggle()
			}
		default:
			// log.Printf("%T %v\n", event, event)
		}
	})
	if err != nil {
		log.Fatal(err)
	}

}
