package main

import (
	"time"

	"github.com/bfanger/clock/display"
	"github.com/bfanger/clock/events"
	"github.com/veandco/go-sdl2/sdl"
)

// Clock displays the current time
type Clock struct {
	Layer  display.Layer
	hour   *display.Text
	dot    *display.Text
	minute *display.Text
	date   *display.Text
}

// NewClock create a new clock and updates every minute
func NewClock(r *display.Renderer) *Clock {
	layer := display.NewContainer()
	font := asset("Roboto-Light.ttf")
	fontSize := 115
	gray := sdl.Color{R: 127, G: 126, B: 126, A: 255}
	orange := sdl.Color{R: 254, G: 110, B: 2, A: 255}

	hour := display.NewText(font, fontSize, gray, "00")
	s := display.NewSprite("Clock[hour]", hour, 148, 170)
	s.AnchorX = 1
	s.AnchorY = 1
	layer.Add(s)

	dot := display.NewText(font, fontSize, gray, ":")
	s = display.NewSprite("Clock[:]", dot, 159, 170)
	s.AnchorX = 0.5
	s.AnchorY = 1
	layer.Add(s)

	minute := display.NewText(font, fontSize, orange, "00")
	s = display.NewSprite("Clock[minute]", minute, 172, 170)
	s.AnchorX = 0
	s.AnchorY = 1
	layer.Add(s)

	darkGray := sdl.Color{R: 102, G: 102, B: 102, A: 255}
	date := display.NewText(font, 55, darkGray, "- ---")
	s = display.NewSprite("Clock[date]", date, 161, 195)
	s.AnchorX = 0.5
	s.AnchorY = 0.5
	layer.Add(s)

	c := &Clock{
		Layer:  layer,
		hour:   hour,
		dot:    dot,
		minute: minute,
		date:   date,
	}
	go c.eventLoop(r)
	return c
}

// Destroy the clock
func (c *Clock) Destroy() error {
	err := c.hour.Destroy()
	c.hour = nil
	return err
}

func (c *Clock) eventLoop(r *display.Renderer) {
	for {
		r.Mutex.Lock()
		t := time.Now()
		c.hour.Text = t.Format("15")
		c.minute.Text = t.Format("04")
		c.date.Text = t.Format("02 Jan")
		r.Mutex.Unlock()
		events.Refresh()
		SleepUntilNext(time.Minute, t)
	}
}

// Next creates the next rounded date based on the step
func Next(d time.Duration, since time.Time) time.Time {
	t := time.Date(since.Year(), since.Month(), since.Day(), since.Hour(), 0, 0, 10000000, since.Location())
	if d >= 60*time.Minute {
		panic("large durations not implemented")
	} else if d >= time.Minute {
		t = t.Add(time.Duration(since.Minute())*time.Minute + d)
	} else {
		t = t.Add(time.Duration(since.Minute())*time.Minute + time.Duration(since.Second())*time.Second + d)
	}
	return t
}

// SleepUntilNext whole minute or second.
func SleepUntilNext(d time.Duration, since time.Time) {
	time.Sleep(time.Until(Next(d, since)))
}
