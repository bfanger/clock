package engine

import (
	"github.com/veandco/go-sdl2/sdl"
)

// World ...
type World struct {
	*Container
	Window             *sdl.Window
	WindowID           uint32
	Renderer           *sdl.Renderer
	EventAutoIncrement int32
	EventQueue         map[int32]func()
	ButtonHandlers     map[int]func()
}

// Render the frame
func (world *World) Render() error {
	if err := world.Renderer.Clear(); err != nil {
		return err
	}

	if err := world.Draw(); err != nil {
		return err
	}
	world.Renderer.Present()
	return nil
}

// const debug = true

// HandleEvents handle the events
func (world *World) HandleEvents() {
	for {
		redraw := true
		event := sdl.WaitEvent()

		switch t := event.(type) {
		case *sdl.QuitEvent:
			return

		case *sdl.MouseButtonEvent:
			redraw = false
			if t.State == 0 { // Pressed
				// @todo button based on mouse position?
				ButtonPressed(4)
			}
		case *sdl.KeyboardEvent:
			redraw = false
			if t.Type == sdl.KEYUP {
				if t.Keysym.Sym == sdl.K_ESCAPE {
					return
				}
				if t.Keysym.Sym == sdl.K_4 {
					ButtonPressed(4)
				}
			}
		case *sdl.WindowEvent:
			if t.Event != sdl.WINDOWEVENT_EXPOSED {
				redraw = false
			}
		case *sdl.UserEvent:
			world.EventQueue[t.Code]()
			delete(world.EventQueue, t.Code)

		default:
			redraw = false

		}
		if redraw {
			world.Render()
		}
	}
}
