package app

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const debug = false

// EventLoop handle the events
func EventLoop() {

	running := true
	for running {
		event := sdl.WaitEvent() // wait here until an event is in the event queue
		switch t := event.(type) {
		case *sdl.QuitEvent:
			running = false
		case *sdl.MouseMotionEvent:
			if debug {
				fmt.Printf("[%d ms] MouseMotion\ttype:%d\tid:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d\n",
					t.Timestamp, t.Type, t.Which, t.X, t.Y, t.XRel, t.YRel)
			}
		case *sdl.MouseButtonEvent:
			if debug {
				fmt.Printf("[%d ms] MouseButton\ttype:%d\tid:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\n",
					t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State)
			}
		case *sdl.MouseWheelEvent:
			if debug {
				fmt.Printf("[%d ms] MouseWheel\ttype:%d\tid:%d\tx:%d\ty:%d\n",
					t.Timestamp, t.Type, t.Which, t.X, t.Y)
			}
		case *sdl.KeyboardEvent:
			if t.Type == sdl.KEYUP && t.Keysym.Sym == sdl.K_ESCAPE {
				running = false
			} else if debug {
				fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
					t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
			}
		case *sdl.WindowEvent:
			if debug {
				fmt.Printf("[%d ms] Window\ttype:%d\tevent:%d\n", t.Timestamp, t.Type, t.Event)
			}
		}
	}
}
