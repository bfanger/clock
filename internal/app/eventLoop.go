package app

import (
	"fmt"
	"os"
	"sync"
	"time"

	"../../internal/engine"
	"github.com/veandco/go-sdl2/sdl"
)

const debug = false

var runningMutex sync.Mutex

// EventLoop handle the events
func EventLoop(world *engine.Container, scene *engine.Container, requestUpdate chan Widget) {
	go renderLoop(world, requestUpdate)
	buttonPressed := make(chan bool)
	if _, err := os.Stat("/sys/class/gpio/"); err == nil {
		go GpioButton(buttonPressed)
	}
	go timerOnClick(buttonPressed, scene, requestUpdate)

	running := true

	for running {
		var event sdl.Event
		// wait here until an event is in the event queue
		sdl.Do(func() {
			event = sdl.PollEvent()
		})
		if event == nil {
			time.Sleep(1 * time.Second / 25)
			continue
		}
		switch t := event.(type) {
		case *sdl.QuitEvent:
			running = false
		// case *sdl.MouseMotionEvent:
		// 	if debug {
		// 		fmt.Printf("[%d ms] MouseMotion\ttype:%d\tid:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d\n",
		// 			t.Timestamp, t.Type, t.Which, t.X, t.Y, t.XRel, t.YRel)
		// 	}
		case *sdl.MouseButtonEvent:
			if t.State == 0 { // Pressed
				buttonPressed <- true
			}
			if debug {
				fmt.Printf("[%d ms] MouseButton\ttype:%d\tid:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\n",
					t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State)
			}
		// case *sdl.MouseWheelEvent:
		// 	if debug {
		// 		fmt.Printf("[%d ms] MouseWheel\ttype:%d\tid:%d\tx:%d\ty:%d\n",
		// 			t.Timestamp, t.Type, t.Which, t.X, t.Y)
		// 	}
		case *sdl.KeyboardEvent:
			if t.Type == sdl.KEYUP && t.Keysym.Sym == sdl.K_ESCAPE {
				running = false
			} else if debug {
				fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
					t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
			}
		case *sdl.WindowEvent:
			if t.Event == sdl.WINDOWEVENT_EXPOSED {
				requestUpdate <- nil
			}
			if debug {
				fmt.Printf("[%d ms] Window\ttype:%d\tevent:%d\n", t.Timestamp, t.Type, t.Event)
			}
		}
	}
}

func renderLoop(world *engine.Container, requestUpdate chan Widget) {
	for {
		sdl.Do(func() {
			if err := world.Render(); err != nil {
				panic(err)
			}
		})

		widget := <-requestUpdate
		if widget != nil {
			sdl.Do(func() {
				if err := widget.Update(); err != nil {
					panic(err)
				}
			})
		}
	}
}

func timerOnClick(buttonPressed chan bool, world *engine.Container, requestUpdate chan Widget) {

	var timer *TimerWidget
	var err error
	countdown := 5
	for {
		if timer != nil {

			select {
			case <-buttonPressed:
				sdl.Do(func() {
					timer.Dispose()
					timer = nil
					if countdown == 20 {
						countdown = 5
						requestUpdate <- nil
					} else {
						_timer, err := createTimer(countdown, world, requestUpdate)
						if err != nil {
							panic(err)
						}
						timer = _timer
						countdown += 5
					}
				})
			case <-timer.Completed:
				sdl.Do(func() {
					timer.Dispose()
					timer = nil
				})
			}
		} else {
			<-buttonPressed
			sdl.Do(func() {
				timer, err = createTimer(countdown, world, requestUpdate)
				if err != nil {
					panic(err)
				}
				countdown += 5
			})
		}
	}
}

func createTimer(countdown int, world *engine.Container, requestUpdate chan Widget) (*TimerWidget, error) {

	now := time.Now()
	timer, err := NewTimerWidget("timer_background.png", now.Hour(), now.Minute()+countdown, world, requestUpdate)

	if err != nil {
		return nil, err
	}

	timer.Second = now.Second()
	timer.Blink = -10
	requestUpdate <- timer
	return timer, nil
}
