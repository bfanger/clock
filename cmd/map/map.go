package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bfanger/clock/pkg/tween"

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
	m.Zoom = 16
	m.Latitude = 52.4900311
	m.Longitude = 4.7602125
	icon, err := ui.ImageFromFile(app.Asset("position.png"), d.Renderer)
	if err != nil {
		log.Fatal(err)
	}
	m.CenterOffsetX = 240
	marker := &app.Marker{
		Latitude:  m.Latitude,
		Longitude: m.Longitude,
		Sprite:    ui.NewSprite(icon)}
	marker.Sprite.AnchorX = 0.5
	marker.Sprite.AnchorY = 0.5
	m.Markers = append(m.Markers, marker)

	scene.Append(m)
	clock, err := app.NewAnalogClock(e)
	if err != nil {
		log.Fatal(err)
	}
	clock.MoveTo(240, 240)
	scene.Append(clock)
	go func() {
		time.Sleep(1 * time.Second)
		e.Animate(tween.FromToFloat64(marker.Longitude, marker.Longitude+0.0005, time.Second, tween.EaseInOutQuad, func(lon float64) {
			marker.Longitude = lon
		}))
		time.Sleep(2 * time.Second)
		e.Animate(m.PanTo(m.Latitude, m.Longitude+0.0003, 2*time.Second))
	}()

	err = e.EventLoop(func(e sdl.Event) {

	})
	if err != nil {
		log.Fatalf("eventloop exit: %v", err)
	}
}
