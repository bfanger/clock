package main

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"runtime"

	"github.com/bfanger/clock/internal/app"
	"github.com/bfanger/clock/pkg/ui"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func init() {
	runtime.LockOSThread()
}
func main() {
	fmt.Println("Clock 4.0")

	d, err := app.NewDisplay(320, 240)
	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()

	engine := ui.NewEngine(d.Renderer)

	font, err := ttf.OpenFont(asset("Roboto-Light.ttf"), 110)
	if err != nil {
		log.Fatalf("unable to open font: %v", err)
	}
	defer font.Close()

	go app.Time(engine, font)

	err = engine.EventLoop(func(event sdl.Event) {
		// log.Printf("%T %v\n", event, event)
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Bye")
}

// asset returns the absolute path for a file in the assets folder
func asset(filename string) string {
	binPath := sdl.GetBasePath() + "assets/"
	_, err := os.Stat(binPath)
	if err == nil {
		return binPath + filename
	}
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	packagePath := gopath + "/src/github.com/bfanger/clock/assets/"
	_, err = os.Stat(packagePath)
	if err == nil {
		return packagePath + filename
	}
	return "./assets/" + filename
}
