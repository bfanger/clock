package engine

import (
	"errors"
	"fmt"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
)

var world *World

// Init the engine
func Init(window *sdl.Window, renderer *sdl.Renderer) *World {
	windowID, err := window.GetID()
	if err != nil {
		panic(err)
	}
	world = &World{
		Container:      &Container{},
		Window:         window,
		WindowID:       windowID,
		Renderer:       renderer,
		EventQueue:     make(map[int32]func()),
		ButtonHandlers: make(map[int]func())}

	return world
}

// Renderer singleton
func Renderer() *sdl.Renderer {
	if world == nil {
		panic("Call engine.Init() first")
	}
	return world.Renderer
}

// Window singletons
func Window() *sdl.Window {
	if world == nil {
		panic("Call engine.Init() first")
	}
	return world.Window
}

var eventLock sync.Mutex

// PushEvent execute the callback in the UI event thread
func PushEvent(callback func()) error {
	if world == nil {
		return errors.New("Call engine.Init() first")
	}
	eventLock.Lock()
	world.EventAutoIncrement++
	id := world.EventAutoIncrement
	world.EventQueue[id] = callback
	eventLock.Unlock()
	event := sdl.UserEvent{
		Type: sdl.USEREVENT,
		// Timestamp: sdl.GetTicks(),
		WindowID: world.WindowID,
		Code:     id}

	if _, err := sdl.PushEvent(&event); err != nil {
		return err
	}
	return nil
}

// ButtonPressed triggers the button handler
func ButtonPressed(button int) {
	if world == nil {
		panic("Call engine.Init() first")
	}
	if world.ButtonHandlers[button] == nil {
		fmt.Printf("No handler registered for button %d\n", button)
		return
	}
	if err := PushEvent(world.ButtonHandlers[button]); err != nil {
		panic(err)
	}
}
