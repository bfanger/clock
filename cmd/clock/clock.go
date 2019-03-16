package main

import (
	"flag"
	"fmt"
	"runtime"
	"time"

	"github.com/bfanger/clock/internal/app"
	"github.com/bfanger/clock/pkg/ui"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func main() {
	runtime.LockOSThread()
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
	fpsVisible := flag.Bool("fps", false, "Show FPS counter")
	flag.Parse()

	display, err := app.NewDisplay()
	if err != nil {
		app.Fatal(errors.Wrap(err, "failed to create display"))
	}
	defer display.Close()
	scene := &ui.Container{}
	engine := ui.NewEngine(scene, display.Renderer)
	engine.Wait = time.Second / 120 // Limit framerate (VSYNC doesn't work on macOS Mohave)
	wm, err := app.NewWidgetManager(scene, engine)
	if err != nil {
		app.Fatal(err)
	}
	defer wm.Close()

	if *fpsVisible {
		font, err := ttf.OpenFont(app.Asset("Roboto-Light.ttf"), 24)
		if err != nil {
			app.Fatal(errors.Wrap(err, "unable to open font"))
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
				// a := app.Alarm{Notification: "vis", Duration: 10 * time.Second}
				// go a.Activate()
				go wm.ButtonPressed()

			}
		case *sdl.KeyboardEvent:
			if e.Type == sdl.KEYDOWN {
			}
		case *sdl.WindowEvent:
			{
				if e.Event == sdl.WINDOWEVENT_RESIZED {
					if err := display.Resized(); err != nil {
						app.Fatal(errors.Wrap(err, "failed to respond to resize event"))
					}
					go engine.Go(func() error { return nil })
				}
			}
		default:
			// log.Printf("%T %v\n", event, event)
		}
	})
	if err != nil {
		app.Fatal(errors.Wrap(err, "eventloop exit"))
	}

}
