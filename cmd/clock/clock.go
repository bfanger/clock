package main

import (
	"fmt"
	"time"

	"../../internal/app"
	"../../internal/engine"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var sprites = make([]engine.Sprite, 0)

func main() {

	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}
	defer sdl.Quit()
	if err := ttf.Init(); err != nil {
		panic(err)
	}

	window, err := app.CreateWindow()
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, 0)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	// Background
	backgroundSprite, err := engine.TextureSpriteFromImage(renderer, "./assets/background.png")
	if err != nil {
		panic(err)
	}
	defer backgroundSprite.Destroy()
	sprites = append(sprites, backgroundSprite)

	// Time
	font, err := ttf.OpenFont("./assets/Teko-Light.ttf", 130)
	if err != nil {
		panic(err)
	}
	date := time.Now()
	timeSprite := engine.TextSprite{
		Font:     font,
		Color:    engine.White(),
		Text:     date.Format("15:04"),
		Renderer: renderer}

	defer timeSprite.Destroy()
	timeSprite.Update()

	sprites = append(sprites, &timeSprite)

	quit := make(chan bool)
	go ticker(&timeSprite, renderer, quit)

	// Main loop
	app.EventLoop(render, renderer)
	quit <- true
}

// render all sprite layers
func render(renderer *sdl.Renderer) {
	count := len(sprites)
	for i := 0; i < count; i++ {
		err := sprites[i].Render()
		if err != nil {
			fmt.Println(err)
		}
	}
	renderer.Present()
}

// ticker updates the time every 15 seconds
func ticker(timeSprite *engine.TextSprite, renderer *sdl.Renderer, quit chan bool) {
	ticker := time.NewTicker(15 * time.Second)

	for {
		select {
		case <-quit:
			ticker.Stop()
			return
		case <-ticker.C:
			date := time.Now()
			fmt.Println("\nClock", date.Format("15:04"))
			timeSprite.Text = date.Format("15:04")
			timeSprite.Update()
			render(renderer)
		}
	}
}
