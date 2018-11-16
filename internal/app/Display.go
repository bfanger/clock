package app

import (
	"fmt"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const screenWidth, screenHeight int32 = 800, 480

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
	width, height := screenWidth, screenHeight
	if d.W == screenWidth {
		// fullscreen mode when the windowsize matches the displaysize
		flags += sdl.WINDOW_FULLSCREEN
		if _, err := sdl.ShowCursor(sdl.DISABLE); err != nil {
			return nil, err
		}
	} else {
		flags += sdl.WINDOW_ALLOW_HIGHDPI
		width /= 2
		height /= 2
	}
	if n == 1 {
		// Single monitor setup, show the clock bottom left.
		x, y = 0, d.H-height
	} else {
		// In a multi monitor setup, show the clock the second screen.
		x, y = sdl.WINDOWPOS_CENTERED_MASK+1, sdl.WINDOWPOS_CENTERED_MASK+1
	}
	w, err := sdl.CreateWindow("Clock", x, y, width, height, flags)
	if err != nil {
		return nil, fmt.Errorf("couldn't create window: %v", err)
	}
	r, err := sdl.CreateRenderer(w, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		return nil, fmt.Errorf("could not create renderer: %v", err)
	}
	drawWidth, _ := w.GLGetDrawableSize()
	if drawWidth != screenWidth {
		// scale the renderer on non-retina screens
		scale := float32(drawWidth) / float32(screenWidth)
		r.SetScale(scale, scale)
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
