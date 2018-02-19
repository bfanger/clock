package app

import (
	"../../internal/engine"
	"github.com/veandco/go-sdl2/sdl"
)

// Boot the interface
func Boot(world *engine.World) {
	scene := &engine.Container{}
	world.Add(scene)
	if IsRaspberryPi() {
		brightness := BrightnessWidget{}
		if err := brightness.Mount(world); err != nil {
			panic(err)
		}
		// defer brightness.Unmount()
	}

	clock := ClockWidget{}
	if err := clock.Mount(scene); err != nil {
		panic(err)
	}
	// defer clock.Unmount()

	school, err := NewTimerWidget("school_background.png", 8, 15)
	if err != nil {
		panic(err)
	}
	school.Repeat = true
	if err = school.Mount(scene); err != nil {
		panic(err)
	}
	// defer school.Unmount()

	world.ButtonHandlers[4] = func() {
		TimerWidgetButtonHandler(scene)
	}

	go HandleGpioButtons()
}

// WindowWidth 320
const WindowWidth = 320

// WindowHeight 240
const WindowHeight = 240

// IsRaspberryPi checks if the display size is 320x240
func IsRaspberryPi() bool {
	var displayMode sdl.DisplayMode
	if err := sdl.GetCurrentDisplayMode(0, &displayMode); err != nil {
		panic(err)
	}
	return displayMode.W == WindowWidth && displayMode.H == WindowHeight
}

// CreateWindow on the second screen, or in fullscreen mode when the windowsize matches the displaysize
func CreateWindow() (*sdl.Window, error) {
	displays, err := sdl.GetNumVideoDisplays()
	if err != nil {
		return nil, err
	}
	var mode sdl.DisplayMode
	if err := sdl.GetCurrentDisplayMode(0, &mode); err != nil {
		return nil, err
	}
	var x, y int32
	var flags uint32
	if displays > 1 {
		x, y = sdl.WINDOWPOS_CENTERED_MASK+1, sdl.WINDOWPOS_CENTERED_MASK+1
	} else {
		x, y = sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED
		x = 0
		y = mode.H - WindowHeight
	}
	if mode.W == WindowWidth {
		flags += sdl.WINDOW_FULLSCREEN
	}

	sdl.ShowCursor(sdl.DISABLE)

	return sdl.CreateWindow("Klok", x, y, WindowWidth, WindowHeight, flags)

}
