package main

import (
	"flag"
	"log"
	"runtime"

	"github.com/bfanger/clock/internal/app"
	"github.com/bfanger/clock/pkg/ui"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func main() {
	runtime.LockOSThread()
	fpsVisible := flag.Bool("fps", false, "Show FPS counter")
	flag.Parse()

	display, err := app.NewDisplay()
	if err != nil {
		log.Fatalf("failed to create display: %v", err)
	}
	defer display.Close()

	engine := ui.NewEngine(display.Renderer)
	wm, err := app.NewWidgetManager(engine)
	if err != nil {
		log.Fatal(err)
	}
	defer wm.Close()

	if *fpsVisible {
		font, err := ttf.OpenFont(app.Asset("Roboto-Light.ttf"), 24)
		if err != nil {
			log.Fatalf("unable to open font: %v", err)
		}
		fps := ui.NewFps(engine, font)
		defer fps.Close()
	}

	server := app.NewServer(wm, engine)
	go server.ListenAndServe()

	// a := app.Alarm{Notification: "vis", Duration: time.Minute}
	// go a.Activate()

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
