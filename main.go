package main

import (
	"flag"
	"fmt"
	"go/build"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"

	"github.com/bfanger/clock/button"
	"github.com/bfanger/clock/clock"
	"github.com/bfanger/clock/display"
	"github.com/bfanger/clock/sprite"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var screenWidth, screenHeight int32 = 240, 320

func main() {

	showFps := flag.Bool("fps", false, "Show FPS counter")
	flag.Parse()

	fmt.Print("Clock")
	defer fmt.Println("bye")

	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		log.Fatalf("Couldn't initialize sdl: %v\n", err)
	}
	defer sdl.Quit()
	if err := ttf.Init(); err != nil {
		log.Fatalf("Couldn't initialize sdl_ttf: %v\n", err)
	}
	defer ttf.Quit()
	w, err := createWindow()
	if err != nil {
		log.Fatalf("Couldn't create window: %v\n", err)
	}
	defer w.Destroy()
	r, err := sdl.CreateRenderer(w, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		log.Fatalf("could not create renderer: %v", err)
	}
	defer r.Destroy()

	scene := display.NewContainer()
	engine := display.NewEngine(r, scene)
	defer engine.Destroy()

	// background := sprite.New("Background", display.NewImage(asset("bedtime.png")), sprite.WithPos(120, 160), sprite.WithAnchor(0.5, 0.5))
	// defer background.Painter.Destroy()
	// scene.AddAt(-1, background)

	c := clock.New(engine, asset("Roboto-Light.ttf"))
	defer c.Destroy()
	scene.Add(c.Layer)
	defer scene.Remove(c.Layer)
	srv := &http.Server{Addr: ":1200"}
	srv.Handler = c.HTTPHandler()

	if runtime.GOOS != "darwin" {
		go handleButtons(c)
	}

	if *showFps {
		fps := display.NewFps(asset("Roboto-Light.ttf"), 14)
		defer fps.Destroy()
		f := sprite.New("FPS-counter", fps, sprite.WithPos(screenWidth-5, 5), sprite.WithAnchor(1, 0))
		engine.Animate(fps)
		defer engine.StopAnimation(fps)
		scene.AddAt(100, f)
		defer scene.Remove(f)
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		_, err := sdl.PushEvent(&sdl.QuitEvent{
			Type:      sdl.QUIT,
			Timestamp: sdl.GetTicks(),
		})
		if err != nil {
			log.Fatalf("could not push quit event: %v", err)
		}
	}()
	graceful := false
	go func() {
		if err := srv.ListenAndServe(); err != nil && graceful == false {
			log.Fatalf("server stopped: %v", err)
		}
	}()
	defer func() {
		graceful = true
		if err := srv.Shutdown(nil); err != nil {
			log.Printf("shutdown failed: %v", err)
		}
	}()

	fmt.Print(" 3")
	if err := engine.Refresh(); err != nil {
		log.Fatal(err)
	}
	engine.Animate(c.Mode(clock.Fullscreen))
	fmt.Println(".0")

	for {
		event := sdl.WaitEvent()
		switch e := event.(type) {
		case *sdl.QuitEvent:
			return
		case *sdl.WindowEvent:
			if e.Event == sdl.WINDOWEVENT_EXPOSED {
				go func() {
					if err := engine.Refresh(); err != nil {
						log.Fatal(err)
					}
				}()
			}
		case *sdl.KeyboardEvent:
			if e.State == sdl.RELEASED {
				// @todo Implement generic gpio 1 - 4
			}
		}
	}

}

// createWindow on the second screen, or in fullscreen mode when the windowsize matches the displaysize
func createWindow() (*sdl.Window, error) {
	n, err := sdl.GetNumVideoDisplays()
	if err != nil {
		return nil, err
	}
	d, err := sdl.GetCurrentDisplayMode(0)
	if err != nil {
		return nil, fmt.Errorf("can't read display mode: %v", err)
	}
	var x, y int32
	var flags uint32
	if n > 1 {
		x, y = sdl.WINDOWPOS_CENTERED_MASK+1, sdl.WINDOWPOS_CENTERED_MASK+1
	} else {
		x, y = sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED
		x = 0
		y = d.H - screenHeight
	}
	if d.W == screenWidth {
		flags += sdl.WINDOW_FULLSCREEN
		if _, err := sdl.ShowCursor(sdl.DISABLE); err != nil {
			return nil, err
		}
	}

	return sdl.CreateWindow("Clock", x, y, screenWidth, screenHeight, flags)
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

func handleButtons(c *clock.Clock) {
	// Change color with the button
	defer log.Println("stopped listening for presses")
	colors := []sdl.Color{clock.Orange, clock.Pink, clock.Green, clock.Blue}
	i := 0
	presses, err := button.Gpio(25) // key: "1", button :4
	if err != nil {
		log.Fatal(err)
	}
	for err := range presses {
		if err != nil {
			log.Fatal(err)
		}
		i++
		if i == len(colors) {
			i = 0
		}

		if err := c.Color(colors[i]); err != nil {
			log.Fatal(err)
		}
	}
}
