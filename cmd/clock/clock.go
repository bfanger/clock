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
	server := app.NewServer(engine)
	server.Background, err = app.NewBackground(engine)
	defer server.Background.Close()

	if err != nil {
		log.Fatal(err)
	}
	server.Notification, err = app.NewNotification(engine)
	if err != nil {
		log.Fatal(err)
	}
	defer server.Notification.Close()
	server.Clock, err = app.NewTime(engine)
	if err != nil {
		log.Fatal(err)
	}
	defer server.Clock.Close()

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
		log.Fatal(err)
	}

}
