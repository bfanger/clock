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
	date := time.Now()
	fmt.Println("\nClock", date.Format("15:04"))
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

	background, err := engine.TextureSpriteFromImage("./assets/background.png", renderer)
	if err != nil {
		panic(err)
	}
	defer background.Destroy()

	time, err := engine.TextSpriteFromText(date.Format("15:04"), renderer)
	if err != nil {
		panic(err)
	}
	defer time.Destroy()

	sprites = append(sprites, background)
	sprites = append(sprites, time)
	// ticker := time.NewTicker(1 * time.Second)

	// for {
	// 	select {
	// 	case <-ticker.C:
	// 		date := time.Now()
	// 		fmt.Println("\nClock", date)
	// 	}
	// }

	app.EventLoop(render, renderer)
}

func render(renderer *sdl.Renderer) {
	fmt.Println("Draw")
	count := len(sprites)
	for i := 0; i < count; i++ {
		sprites[i].Render()
	}
	renderer.Present()
}
