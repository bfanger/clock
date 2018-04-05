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
	fontSize := 95
	orange := sdl.Color{R: 254, G: 110, B: 2, A: 255}
	const y int32 = 80

	hour := display.NewText(font, fontSize, orange, "--")
	s := display.NewSprite("Clock[hour]", hour, 109, y)
	s.AnchorX = 1
	s.AnchorY = 0
	layer.Add(s)

	dot := display.NewText(font, fontSize, orange, ":")
	s = display.NewSprite("Clock[:]", dot, 118, y)
	s.AnchorX = 0.5
	s.AnchorY = 0
	layer.Add(s)

	minute := display.NewText(font, fontSize, orange, "--")
	s = display.NewSprite("Clock[minute]", minute, 130, y)
	s.AnchorX = 0
	s.AnchorY = 0
	layer.Add(s)

	gray := sdl.Color{R: 102, G: 102, B: 102, A: 255}
	date := display.NewText(font, 50, gray, "- ---")
	s = display.NewSprite("Clock[date]", date, 119, y+105)
	s.AnchorX = 0.5
	s.AnchorY = 0
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
		c.hour.Text = "23" //t.Format("15")
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
