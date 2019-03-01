package main

import (
	"log"
	"time"

	"github.com/bfanger/clock/internal/app"
	"github.com/bfanger/clock/pkg/ui"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	display, err := app.NewDisplay()
	if err != nil {
		log.Fatalf("failed to create display")
	}
	defer display.Close()
	engine := ui.NewEngine(display.Renderer)
	engine.Wait = time.Second / 120 // Limit framerate (VSYNC doesn't work on macOS Mohave)
	img, err := ui.ImageFromFile(app.Asset("poc.png"), display.Renderer)
	if err != nil {
		log.Fatalf("failed to load image")
	}
	engine.Scene.Append(img)
	err = engine.EventLoop(func(event sdl.Event) {

	})
	if err != nil {
		log.Fatalf("eventloop exit")
	}
}
