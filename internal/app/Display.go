package app

import (
	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const screenWidth, screenHeight int32 = 800, 480

// Display encapsulate setting up and cleaning up a SDL renderer
type Display struct {
	Renderer   *sdl.Renderer
	Window     *sdl.Window
	Fullscreen bool
}

// NewDisplay initializes SDL and creates a window
func NewDisplay() (*Display, error) {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return nil, errors.Wrap(err, "couldn't initialize sdl")
	}
	if err := ttf.Init(); err != nil {
		return nil, errors.Wrap(err, "couldn't initialize sdl_ttf")
	}
	if err := img.Init(img.INIT_PNG); err != nil {
		return nil, errors.Errorf("couldn't initialize sdl_img: %d", err)
	}
	n, err := sdl.GetNumVideoDisplays()
	if err != nil {
		return nil, err
	}
	m, err := sdl.GetCurrentDisplayMode(0)
	if err != nil {
		return nil, errors.Wrap(err, "can't read display mode")
	}
	var x, y int32 = sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED
	var rendererFlags uint32 = sdl.RENDERER_SOFTWARE
	var windowFlags uint32
	width, height := screenWidth, screenHeight
	if m.W == screenWidth {
		// fullscreen mode when the windowsize matches the displaysize
		windowFlags |= sdl.WINDOW_FULLSCREEN
		if _, err := sdl.ShowCursor(sdl.DISABLE); err != nil {
			return nil, err
		}
	} else {
		windowFlags |= sdl.WINDOW_ALLOW_HIGHDPI
		rendererFlags = sdl.RENDERER_ACCELERATED | sdl.RENDERER_PRESENTVSYNC
		width /= 2
		height /= 2
		if n == 1 {
			// Single monitor setup, show the clock bottom left.
			x, y = 0, m.H-height
		} else {
			// In a multi monitor setup, show the clock the second screen.
			x, y = sdl.WINDOWPOS_CENTERED_MASK+1, sdl.WINDOWPOS_CENTERED_MASK+1
		}
	}
	d := &Display{}
	d.Window, err = sdl.CreateWindow("Clock", x, y, width, height, windowFlags)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create window")
	}
	d.Renderer, err = sdl.CreateRenderer(d.Window, -1, rendererFlags)
	if err != nil {
		return nil, errors.Wrap(err, "could not create renderer")
	}
	return d, nil
}

// Close open resources
func (d *Display) Close() error {
	if err := d.Renderer.Destroy(); err != nil {
		return err
	}
	if err := d.Window.Destroy(); err != nil {
		return err
	}
	ttf.Quit()
	img.Quit()
	sdl.Quit()
	return nil
}
