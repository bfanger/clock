package main

import (
	"fmt"
	"runtime"

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

	scene := &engine.Container{}
	world.Add(scene)

	clock := app.ClockWidget{}
	if err = clock.Mount(scene); err != nil {
		panic(err)
	}
	defer clock.Unmount()

	school, err := app.NewTimerWidget("school_background.png", 8, 15)
	if err != nil {
		panic(err)
	}
	school.Repeat = true
	if err = school.Mount(scene); err != nil {
		panic(err)
	}
	defer school.Unmount()

	world.ButtonHandlers[4] = func() {
		app.TimerWidgetButtonHandler(scene)
	}

	go app.HandleGpioButtons()

	brightness := app.BrightnessWidget{}
	if err := brightness.Mount(world); err != nil {
		panic(err)
	}
	defer brightness.Unmount()

	world.HandleEvents()
}
