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
	server, err := app.NewServer(engine)
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()
	go server.ListenAndServe()

	err = engine.EventLoop(func(event sdl.Event) {
		switch e := event.(type) {
		case *sdl.MouseButtonEvent:
			if e.Type == sdl.MOUSEBUTTONUP {
				server.Toggle()
			}
		case *sdl.KeyboardEvent:
			if e.Type == sdl.KEYDOWN {
				server.Toggle()
			}
		default:
			// log.Printf("%T %v\n", event, event)
		}
	})
	if err != nil {
		log.Fatal(err)
	}

}
