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

	needsRedraw := make(chan bool)

	time, err := app.NewTimeWidget(world, needsRedraw)
	if err != nil {
		panic(err)
	}
	defer time.Dispose()

	brightness, err := app.NewBrightnessWidget(world, needsRedraw)
	if err != nil {
		panic(err)
	}
	defer brightness.Dispose()

	// Main loop
	go renderLoop(world, needsRedraw)

	app.EventLoop()
}

func renderLoop(world *engine.Container, needsRedraw chan bool) {
	for {
		world.Render()
		world.Renderer.Present()
		<-needsRedraw
	}
}
