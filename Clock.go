package main

import (
	"strconv"
	"time"

	"github.com/bfanger/clock/display"
	"github.com/bfanger/clock/events"
	"github.com/veandco/go-sdl2/sdl"
)

// Clock displays the current time
type Clock struct {
	Layer display.Layer
	text  *display.Text
}

// NewClock create a new clock and updates every minute
func NewClock(r *display.Renderer) *Clock {
	orange := sdl.Color{R: 251, G: 140, B: 63, A: 255}
	t := display.NewText(asset("Audiowide-Regular.ttf"), 80, orange, "--:--")
	s := display.NewSprite("Clock", t, 160, 120)
	s.AnchorX = 0.5
	s.AnchorY = 0.5
	c := &Clock{
		Layer: s,
		text:  t,
	}
	go c.eventLoop(r)
	return c
}

// Destroy the clock
func (c *Clock) Destroy() error {
	err := c.text.Destroy()
	c.text = nil
	return err
}

func (c *Clock) eventLoop(r *display.Renderer) {
	for {
		r.Mutex.Lock()
		t := time.Now()
		if c.text == nil {
			return
		}
		c.text.Text = strconv.Itoa(t.Hour()) + t.Format(":04")
		events.Refresh()
		r.Mutex.Unlock()
		delay := time.Duration(time.Minute + (time.Second / 100))
		delay -= time.Duration(t.Second()) * time.Second
		delay -= time.Duration(t.Nanosecond()) * time.Nanosecond

		time.Sleep(delay)
	}
}
