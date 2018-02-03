package main

import (
	"fmt"

	"../../internal/app"
	"../../internal/engine"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func main() {
	fmt.Println("Clock")
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}
	defer sdl.Quit()
	if err := ttf.Init(); err != nil {
		panic(err)
	}

	sdl.Main(run)
}

func run() {

	var err error
	var window *sdl.Window
	sdl.Do(func() {
		window, err = app.CreateWindow()
		if err != nil {
			panic(err)
		}
	})

	defer sdl.Do(func() {
		window.Destroy()
	})
	var renderer *sdl.Renderer
	sdl.Do(func() {
		renderer, err = sdl.CreateRenderer(window, -1, 0)
		if err != nil {
			panic(err)
		}
	})
	defer sdl.Do(func() {
		renderer.Destroy()
	})
	world := engine.NewContainer(renderer)

	requestUpdate := make(chan app.Widget)
	var clock *app.ClockWidget
	sdl.Do(func() {
		clock, err = app.NewClockWidget(world, requestUpdate)
		if err != nil {
			panic(err)
		}
	})
	defer sdl.Do(func() {
		clock.Dispose()
	})

	school, err := app.NewTimerWidget("school_background.png", 8, 15, world, requestUpdate)
	if err != nil {
		panic(err)
	}
	school.Repeat = true
	defer sdl.Do(func() {
		school.Dispose()
	})
	var brightness *app.BrightnessWidget
	sdl.Do(func() {
		brightness, err = app.NewBrightnessWidget(world, requestUpdate)
		if err != nil {
			panic(err)
		}
	})
	defer sdl.Do(func() {
		brightness.Dispose()
	})
	app.EventLoop(world, requestUpdate)

}
