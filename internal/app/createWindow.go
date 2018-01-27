package app

import "github.com/veandco/go-sdl2/sdl"

// CreateWindow is documented
func CreateWindow() (*sdl.Window, error) {

	var displays, err = sdl.GetNumVideoDisplays()
	if err != nil {
		panic(err)
	}
	var mode sdl.DisplayMode
	if err := sdl.GetCurrentDisplayMode(0, &mode); err != nil {
		panic(err)
	}
	var x, y int32
	var flags uint32 = sdl.WINDOW_BORDERLESS
	if displays > 1 {
		x, y = sdl.WINDOWPOS_CENTERED_MASK+1, sdl.WINDOWPOS_CENTERED_MASK+1
	} else {
		x, y = sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED
		x = 0
		y = mode.H - 240
	}
	if mode.W == 320 {
		flags += sdl.WINDOW_FULLSCREEN
	}

	sdl.ShowCursor(sdl.DISABLE)
	window, err := sdl.CreateWindow("Klok", x, y, 320, 240, flags)

	return window, err

}
