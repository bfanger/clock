package main

import (
	"fmt"
	"runtime"
	"time"

	"../../internal/app"
	"../../internal/engine"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func main() {
	fmt.Println("Clock")
	runtime.LockOSThread()
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
	world := engine.Init(window, renderer)

	engine.SetTimeout(func() {
		app.Boot(world)
	}, time.Second)

	world.HandleEvents()
}
