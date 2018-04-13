package main

import (
	"flag"
	"fmt"
	"go/build"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bfanger/clock/clock"
	"github.com/bfanger/clock/display"
	"github.com/bfanger/clock/sprite"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var screenWidth, screenHeight int32 = 240, 320

func main() {
	defer fmt.Println("bye")

	showFps := flag.Bool("fps", false, "Show FPS counter")
	flag.Parse()

	fmt.Print("Clock")

	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		log.Fatalf("Couldn't initialize sdl: %v\n", err)
	}
	defer sdl.Quit()
	if err := ttf.Init(); err != nil {
		log.Fatalf("Couldn't initialize sdl_ttf: %v\n", err)
	}
	defer ttf.Quit()
	window, err := createWindow()
	if err != nil {
		log.Fatalf("Couldn't create window: %v\n", err)
	}
	defer window.Destroy()
	r, err := display.NewRenderer(window)
	if err != nil {
		log.Fatalf("Couldn't create renderer: %v\n", err)
	}
	defer r.Destroy()

	// image := display.NewImage(asset("image.jpg"))
	// defer image.Destroy()
	// r.AddAt(sprite.New("Background", image), -1)

	display.Init(r)
	defer display.Quit()
	fmt.Println(" 3.0")

	c := clock.New(&r.Mutex, asset("Roboto-Light.ttf"))
	defer c.Destroy()
	r.Add(c.Layer)
	defer r.Remove(c.Layer)
	r.Animate(c.Mode(clock.Fullscreen))
	http.HandleFunc("/top", func(w http.ResponseWriter, req *http.Request) {
		r.Animate(c.Mode(clock.Top))
		w.Write([]byte("<a href=\"fullscreen\">goto fullscreen</a>"))
	})
	http.HandleFunc("/fullscreen", func(w http.ResponseWriter, req *http.Request) {
		r.Animate(c.Mode(clock.Fullscreen))
		w.Write([]byte("<a href=\"top\">goto top</a><br><a href=\"hidden\">hide</a>"))
	})
	http.HandleFunc("/hidden", func(w http.ResponseWriter, req *http.Request) {
		r.Animate(c.Mode(clock.Hidden))
		w.Write([]byte("<a href=\"fullscreen\">goto fullscreen</a>"))
	})

	if *showFps {
		fps := display.NewFps(r, asset("Roboto-Light.ttf"), 14)
		f := sprite.New("FPS-counter", fps, sprite.WithPos(screenWidth-5, 5), sprite.WithAnchor(1, 0))
		r.Animate(fps)
		r.AddAt(100, f)
		defer func() {
			r.StopAnimation(fps)
			r.Remove(f)
			r.Mutex.Lock()
			defer r.Mutex.Unlock()
			fps.Destroy()
		}()
	}

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		display.Shutdown()
	}()

	go func() {
		http.ListenAndServe(":8000", nil)
	}()

	if err := display.EventLoop(r); err != nil {
		log.Fatalf("eventLoop: %v\n", err)
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
		sdl.ShowCursor(sdl.DISABLE)
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

// isRaspberryPi checks if the display size is 240x320
func isRaspberryPi() bool {
	d, err := sdl.GetCurrentDisplayMode(0)
	if err != nil {
		log.Fatalf("can't read display mode: %v\n", err)
	}
	return d.W == screenWidth && d.H == screenHeight
}
