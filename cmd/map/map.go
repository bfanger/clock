package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bfanger/clock/pkg/ui"
	"github.com/joho/godotenv"

	"github.com/bfanger/clock/internal/app"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
	d, err := app.NewDisplay()
	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()
	scene := &ui.Container{}
	e := ui.NewEngine(scene, d.Renderer)
	m := app.NewMap(os.Getenv("MAPTILER_KEY"), e)
	m.Zoom = 17
	m.Latitude = 52.4900311
	m.Longitude = 4.7602125

	scene.Append(m)
	go e.Animate(m.PanTo(m.Latitude, m.Longitude+0.0002, time.Second))

	err = e.EventLoop(func(e sdl.Event) {

	})
	if err != nil {
		log.Fatalf("eventloop exit: %v", err)
	}
}
