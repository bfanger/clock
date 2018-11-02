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
		log.Fatalf("failed to create display: %v", err)
	}
	defer display.Close()

	engine := ui.NewEngine(display.Renderer)
	displayManager, err := app.NewDisplayManager(engine)
	if err != nil {
		log.Fatal(err)
	}
	defer displayManager.Close()

	server := app.NewServer(displayManager, engine)
	go server.ListenAndServe()

	err = engine.EventLoop(func(event sdl.Event) {
		switch e := event.(type) {
		case *sdl.MouseButtonEvent:
			if e.Type == sdl.MOUSEBUTTONUP {

			}
		case *sdl.KeyboardEvent:
			if e.Type == sdl.KEYDOWN {
			}
		default:
			// log.Printf("%T %v\n", event, event)
		}
	})
	if err != nil {
		log.Fatalf("eventloop exit: %v", err)
	}

}
