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
	engine := ui.NewEngine(scene, display.Renderer)
	engine.Wait = time.Second / 120 // Limit framerate (VSYNC doesn't work on macOS Mohave)
	wm, err := app.NewWidgetManager(scene, engine)
	if err != nil {
		app.Fatal(err)
	}
	defer wm.Close()

	addGauge(scene, engine)
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

func addGauge(scene *ui.Container, engine *ui.Engine) {
	conic, err := ui.ImageFromFile(app.Asset("conic.png"), engine.Renderer)
	if err != nil {
		app.Fatal(err)
	}
	gauge1 := ui.NewGuage(conic, -45, 45)
	sprite1 := ui.NewSprite(gauge1)
	scene.Append(sprite1)

	gauge2 := ui.NewGuage(conic, 0, 200)
	sprite2 := ui.NewSprite(gauge2)
	sprite2.X = 300
	scene.Append(sprite2)

	// go engine.Animate(tween.FromTo(120, 270, 5*time.Second, tween.Linear, func(x float64) {
	// 	gauge.SetStart(45)
	// 	gauge.SetEnd(45 + x)
	// }))

	// sprite1.X = 180
	// sprite1.Y = 30
	// sprite1.SetScale(1.5)

}
