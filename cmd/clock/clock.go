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
	if err := godotenv.Load(); err != nil {
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
	rotate := ui.NewLandscape(scene)
	engine := ui.NewEngine(rotate, display.Renderer)
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
			if e.Type == sdl.MOUSEBUTTONDOWN {
				go wm.ButtonPressed()
				// go func() {
				// 	if err := app.ShowNotification("vis", 10*time.Second); err != nil {
				// 		app.Fatal(err)
				// 	}
				// }()
			}
		default:
			// fmt.Printf("%T %v\n", event, event)
		}
	})
	if err != nil {
		app.Fatal(errors.Wrap(err, "eventloop exit"))
	}
}
