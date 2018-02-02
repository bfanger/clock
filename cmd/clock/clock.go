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

	window, err := app.CreateWindow()
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, 0)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()
	world := engine.NewContainer(renderer)

	requestUpdate := make(chan app.Widget)

	clock, err := app.NewClockWidget(world, requestUpdate)
	if err != nil {
		panic(err)
	}
	defer clock.Dispose()

	school, err := app.NewTimerWidget("school_background.png", 8, 15, world, requestUpdate)
	if err != nil {
		panic(err)
	}
	defer school.Dispose()

	brightness, err := app.NewBrightnessWidget(world, requestUpdate)
	if err != nil {
		panic(err)
	}
	defer brightness.Dispose()

	// Main loop
	go renderLoop(world, requestUpdate)

	app.EventLoop()
}

func renderLoop(world *engine.Container, requestUpdate chan app.Widget) {
	for {
		if err := world.Renderer.Clear(); err != nil {
			panic(err)
		}

		if err := world.Render(); err != nil {
			panic(err)
		}
		world.Renderer.Present()

		widget := <-requestUpdate
		if err := widget.Update(); err != nil {
			panic(err)
		}
	}
}
