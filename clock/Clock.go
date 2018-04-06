package clock

import (
	"fmt"
	"time"

	"github.com/bfanger/clock/display"
	"github.com/bfanger/clock/tween"
	"github.com/veandco/go-sdl2/sdl"
)

// Clock displays the current time
type Clock struct {
	Layer  display.Layer
	hour   *display.Text
	dot    *display.Text
	minute *display.Text
	date   *display.Text
	quit   chan bool
}

// New create a new clock and updates every minute
func New(r *display.Renderer, font string) *Clock {
	layer := display.NewContainer()
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
		quit:   make(chan bool),
	}
	go c.eventLoop(r, c.quit)
	return c
}

// Destroy the clock
func (c *Clock) Destroy() error {
	c.quit <- true
	close(c.quit)
	if err := c.hour.Destroy(); err != nil {
		return fmt.Errorf("could not destroy hour: %v", err)
	}
	if err := c.dot.Destroy(); err != nil {
		return fmt.Errorf("could not destroy dot: %v", err)
	}
	if err := c.minute.Destroy(); err != nil {
		return fmt.Errorf("could not destroy minute: %v", err)
	}
	if err := c.date.Destroy(); err != nil {
		return fmt.Errorf("could not destroy date: %v", err)
	}
	return nil
}

func (c *Clock) eventLoop(r *display.Renderer, quit <-chan bool) {
	for {
		r.Mutex.Lock()
		t := time.Now()
		c.hour.Text = t.Format("15")
		c.minute.Text = t.Format("04")
		c.date.Text = t.Format("02 Jan")
		r.Mutex.Unlock()
		display.Refresh()
		select {
		case <-quit:
			return
		case <-time.After(time.Until(Next(time.Minute, t))):
		}
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

// Show the clock
func (c *Clock) Show(r *display.Renderer, animated bool) {
	r.Add(c.Layer)
	if animated == false {
		return
	}
	c.Layer.Move(0, 320)
	var prev int32
	r.Animate(tween.FromToInt32(0, -320, 2*time.Second, func(v int32) {
		d := v - prev
		prev = v
		c.Layer.Move(0, d)
	}))
}
