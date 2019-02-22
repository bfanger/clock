package main

import (
	"fmt"
	"log"
	"os"

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
	m := &app.Map{
		Zoom:      17,
		Latitude:  52.4900311,
		Longitude: 4.7602125,
		Key:       os.Getenv("MAPTILER_KEY")}

	scene.Append(m)

	err = e.EventLoop(func(e sdl.Event) {

	})
	if err != nil {
		log.Fatalf("eventloop exit: %v", err)
	}
}
