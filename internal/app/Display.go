package app

import (
	"fmt"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const screenWidth, screenHeight int32 = 320, 240

// Display encapsulate setting up and cleaning up a SDL renderer
type Display struct {
	Renderer   *sdl.Renderer
	window     *sdl.Window
	Fullscreen bool
}

// NewDisplay initializes SDL and creates a window
func NewDisplay() (*Display, error) {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return nil, fmt.Errorf("couldn't initialize sdl: %v", err)
	}
	if err := ttf.Init(); err != nil {
		return nil, fmt.Errorf("couldn't initialize sdl_ttf: %v", err)
	}
	if err := img.Init(img.INIT_PNG); err != img.INIT_PNG {
		return nil, fmt.Errorf("couldn't initialize sdl_img: %v", err)
	}
	w, err := createWindow(screenWidth, screenHeight)
	if err != nil {
		return nil, fmt.Errorf("couldn't create window: %v", err)
	}
	r, err := sdl.CreateRenderer(w, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		return nil, fmt.Errorf("could not create renderer: %v", err)
	}
	return &Display{
		Renderer: r,
		window:   w}, nil
}

// Close open resources
func (d *Display) Close() error {
	if err := d.Renderer.Destroy(); err != nil {
		return err
	}
	if err := d.window.Destroy(); err != nil {
		return err
	}
	ttf.Quit()
	img.Quit()
	sdl.Quit()
	return nil
}

// createWindow on the second screen, or in fullscreen mode when the windowsize matches the displaysize
func createWindow(width, height int32) (*sdl.Window, error) {
	n, err := sdl.GetNumVideoDisplays()
	if err != nil {
		return nil, err
	}
	d, err := sdl.GetCurrentDisplayMode(0)
	if err != nil {
		return nil, fmt.Errorf("can't read display mode: %v", err)
	}
	var x, y int32
	var flags uint32
	if n > 1 {
		x, y = sdl.WINDOWPOS_CENTERED_MASK+1, sdl.WINDOWPOS_CENTERED_MASK+1
	} else {
		x, y = sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED
		x = 0
		y = d.H - height
	}
	if d.W == width {
		flags += sdl.WINDOW_FULLSCREEN
		if _, err := sdl.ShowCursor(sdl.DISABLE); err != nil {
			return nil, err
		}
	}

	return sdl.CreateWindow("Clock", x, y, width, height, flags)
}
