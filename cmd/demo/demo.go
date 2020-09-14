package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

const SCREEN_WIDTH = 800
const SCREEN_HEIGHT = 480

func main() {
	//Initialize SDL
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}

	// var window, err := sdl.CreateWindow("Clock", x, y, width, height, flags)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "couldn't create window")
	// }
	// d.Renderer, err = sdl.CreateRenderer(d.window, -1, 0) //sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC
	// if err != nil {
	// 	return nil, errors.Wrap(err, "could not create renderer")
	// }
	//The surface contained by the window
	// SDL_Surface* screenSurface = NULL;

	//Create window
	window, err := sdl.CreateWindow("SDL Tutorial", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, SCREEN_WIDTH, SCREEN_HEIGHT, sdl.WINDOW_FULLSCREEN|sdl.WINDOW_SHOWN|sdl.WINDOW_OPENGL)
	if err != nil {
		panic(err)
	}

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	if err := surface.FillRect(nil, sdl.MapRGB(surface.Format, 255, 255, 255)); err != nil {
		panic(err)
	}
	if err := window.UpdateSurface(); err != nil {
		panic(err)
	}
	sdl.Delay(2000)

	if err := window.Destroy(); err != nil {
		panic(err)
	}
	sdl.Quit()
}
