package main

import (
	"fmt"

	"github.com/bfanger/clock/display"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func main() {
	fmt.Print("Clock")

	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}
	defer sdl.Quit()
	if err := ttf.Init(); err != nil {
		panic(err)
	}
	defer ttf.Quit()
	window, err := sdl.CreateWindow("Clock", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		320, 240, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	r, err := display.NewRenderer(window)
	if err != nil {
		panic(err)
	}
	image := display.NewImage("/Sites/clock3/src/github.com/bfanger/clock/assets/image.jpg")
	defer image.Destroy()
	r.Add(0, display.NewSprite("Background", image, 0, 0))

	fmt.Print(" 3")
	white := sdl.Color{R: 255, G: 255, B: 255, A: 255}
	text := display.NewText("/Sites/clock3/src/github.com/bfanger/clock/assets/Roboto-Bold.ttf", 90, white, "23:99")
	defer text.Destroy()
	r.Add(2, display.NewSprite("Time", text, 50, 70))

	first := true
	quit := false
	repaint := false
	var event sdl.Event
	for {
		event = sdl.WaitEvent()
		repaint, quit = handleEvent(event)

		if quit {
			return
		}
		if repaint || first {
			if err = r.Render(); err != nil {
				panic(err)
			}
		}
		if first {
			fmt.Println(".0")
			first = false
		}
	}
}

func handleEvent(event sdl.Event) (repaint bool, quit bool) {
	switch e := event.(type) {
	case *sdl.QuitEvent:
		quit = true
	case *sdl.TouchFingerEvent:
	case *sdl.MouseMotionEvent:
	case *sdl.KeyboardEvent:
	case *sdl.WindowEvent:
		if e.Event == sdl.WINDOWEVENT_EXPOSED {
			repaint = true
		}
	default:
		// fmt.Printf("%T%+v\n", event, event)
		repaint = false
	}
	return
}
