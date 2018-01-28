package main

import (
	"fmt"
	"os"
	"time"

	"../../internal/app"
	"../../internal/engine"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var sprites = make([]engine.Sprite, 0)

func main() {
	fmt.Println("Clock")
	assetPath := sdl.GetBasePath() + "assets/"
	if _, err := os.Stat(assetPath); err != nil {
		assetPath = "./assets/"
	}

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
	backgroundSprite, err := engine.TextureSpriteFromImage(renderer, assetPath+"background.png")
	if err != nil {
		panic(err)
	}
	defer backgroundSprite.Destroy()
	sprites = append(sprites, backgroundSprite)

	// Time
	font, err := ttf.OpenFont(assetPath+"Teko-Light.ttf", 135)
	if err != nil {
		panic(err)
	}
	font.SetHinting(ttf.HINTING_NORMAL)
	date := time.Now()
	timeSprite, err := engine.NewTextSprite(
		font,
		engine.White(),
		date.Format("15:04"),
		renderer)
	if err != nil {
		panic(err)
	}
	timeSprite.TextureSprite.Destination.X = 95
	timeSprite.TextureSprite.Destination.Y = 80

	defer timeSprite.Destroy()

	sprites = append(sprites, timeSprite)

	// Brightness
	var displayMode sdl.DisplayMode
	if err := sdl.GetCurrentDisplayMode(0, &displayMode); err != nil {
		panic(err)
	}
	if displayMode.W <= 320 {
		brightnessSprite, err := engine.NewBrightnessSprite(renderer, 128)
		if err != nil {
			panic(err)
		}
		defer brightnessSprite.Destroy()
		sprites = append(sprites, brightnessSprite)
	}

	quit := make(chan bool)
	go ticker(timeSprite, renderer, quit)

	// Main loop
	render(renderer)
	app.EventLoop()
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
			// fmt.Println(date.Format("15:04"))
			timeSprite.Text = date.Format("15:04")
			if err := timeSprite.Update(); err != nil {
				panic(err)
			}
			render(renderer)
		}
	}
}
