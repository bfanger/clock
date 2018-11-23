package main

import (
	"flag"
	"log"
	"runtime"
	"time"

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
	engine.Wait = time.Second / 120 // Limit framerate (VSYNC doesn't work on macOS Mohave)
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

	err = engine.EventLoop(func(event sdl.Event) {
		switch e := event.(type) {
		case *sdl.MouseButtonEvent:
			if e.Type == sdl.MOUSEBUTTONUP {
				a := app.Alarm{Notification: "vis", Duration: 10 * time.Second}
				go a.Activate()

			}
		case *sdl.KeyboardEvent:
			if e.Type == sdl.KEYDOWN {
			}
		case *sdl.WindowEvent:
			{
				if e.Event == sdl.WINDOWEVENT_RESIZED {
					if err := display.Resized(); err != nil {
						log.Fatalf("failed to respond to resize event: %v", err)
					}
				}
			}
		default:
			// log.Printf("%T %v\n", event, event)
		}
	})
	if err != nil {
		log.Fatalf("eventloop exit: %v", err)
	}

}
