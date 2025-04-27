package main

import (
	"time"

	"github.com/bfanger/clock/internal/app"
	"github.com/bfanger/clock/pkg/ui"
	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	display, err := app.NewDisplay()
	if err != nil {
		app.Fatal(errors.New("failed to create display"))
	}
	defer display.Close()
	img, err := ui.ImageFromFile(app.Asset("poc.png"), display.Renderer)
	if err != nil {
		app.Fatal(errors.New("failed to load image"))
	}
	engine := ui.NewEngine(img, display.Renderer)
	engine.Wait = time.Second / 120 // Limit framerate (VSYNC doesn't work on macOS Mohave)

	err = engine.EventLoop(func(event sdl.Event) {
	})
	if err != nil {
		app.Fatal(errors.Wrap(err, "event loop"))
	}
}
