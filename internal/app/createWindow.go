package app

import "github.com/veandco/go-sdl2/sdl"

// WindowWidth 320
const WindowWidth = 320

// WindowHeight 240
const WindowHeight = 240

// CreateWindow is documented
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
