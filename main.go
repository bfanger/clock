package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/bfanger/clock/display"
	"github.com/bfanger/clock/events"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var screenWidth, screenHeight int32 = 320, 240

func main() {
	defer fmt.Println("bye")
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
	image := display.NewImage("/Sites/clock3/src/github.com/bfanger/clock/assets/image.jpg")
	defer image.Destroy()
	r.Add(0, display.NewSprite("Background", image, 0, 0))

	white := sdl.Color{R: 255, G: 255, B: 255, A: 255}
	text := display.NewText("/Sites/clock3/src/github.com/bfanger/clock/assets/Roboto-Bold.ttf", 90, white, "23:99")
	defer text.Destroy()
	r.Add(2, display.NewSprite("Time", text, 50, 70))
	var m sync.Mutex
	render := make(chan bool)
	defer close(render)
	go func() {
		for range render {
			m.Lock()
			fmt.Print(".")
			if err = r.Render(); err != nil {
				panic(err)
			}
			m.Unlock()
		}
	}()

	events.Init()
	defer events.Quit()
	fmt.Println(" 3.0")

	go func() {
		count := 0
		for {
			count++
			m.Lock()
			text.Text = strconv.Itoa(count)
			m.Unlock()
			time.Sleep(time.Second)
			events.Refresh()
		}

	}()

	if err := events.EventLoop(render); err != nil {
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
	}

	sdl.ShowCursor(sdl.DISABLE)

	return sdl.CreateWindow("Clock", x, y, screenWidth, screenHeight, flags)

}

// isRaspberryPi checks if the display size is 320x240
func isRaspberryPi() bool {
	d, err := sdl.GetCurrentDisplayMode(0)
	if err != nil {
		log.Fatalf("can't read display mode: %v\n", err)
	}
	return d.W == screenWidth && d.H == screenHeight
}
