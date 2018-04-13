package clock

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/bfanger/clock/display"
	"github.com/bfanger/clock/sprite"
	"github.com/bfanger/clock/tween"
	"github.com/veandco/go-sdl2/sdl"
)

type mode int

const (
	Hidden mode = iota
	Fullscreen
	Top
)

// Clock displays the current time
type Clock struct {
	Layer      *display.Container
	hour       *sprite.Sprite
	dot        *sprite.Sprite
	minute     *sprite.Sprite
	time       *display.Container
	date       *sprite.Sprite
	busy       *sync.Mutex
	quit       chan bool
	mode       mode
	transition display.Animater
}

// New creates a clock and updates every minute
func New(busy *sync.Mutex, font string) *Clock {
	fontSize := 95
	orange := sdl.Color{R: 254, G: 110, B: 2, A: 255}
	time := display.NewContainer()

	hour := sprite.New("Clock[hour]", display.NewText(font, fontSize, orange, "99"))
	hour.AnchorX = 1
	hour.AnchorY = 0

	dot := sprite.New("Clock[:]", display.NewText(font, fontSize, orange, ":"))
	dot.AnchorX = 0.5
	dot.AnchorY = 0

	minute := sprite.New("Clock[minute]", display.NewText(font, fontSize, orange, "99"))
	minute.AnchorX = 0
	minute.AnchorY = 0
	time.Add(hour, dot, minute)

	gray := sdl.Color{R: 127, G: 126, B: 126, A: 255}
	date := sprite.New("Clock[date]", display.NewText(font, 50, gray, "- ---"))
	date.AnchorX = 0.5
	date.AnchorY = 0

	c := &Clock{
		Layer:  display.NewContainer(),
		hour:   hour,
		dot:    dot,
		minute: minute,
		time:   time,
		date:   date,
		mode:   Hidden,
		busy:   busy,
		quit:   make(chan bool),
	}
	go c.eventLoop()
	return c
}

// Destroy the clock
func (c *Clock) Destroy() error {
	c.busy.Lock()
	defer c.busy.Unlock()
	c.quit <- true
	close(c.quit)
	if err := c.hour.Painter.Destroy(); err != nil {
		return fmt.Errorf("could not destroy hour: %v", err)
	}
	if err := c.dot.Painter.Destroy(); err != nil {
		return fmt.Errorf("could not destroy dot: %v", err)
	}
	if err := c.minute.Painter.Destroy(); err != nil {
		return fmt.Errorf("could not destroy minute: %v", err)
	}
	if err := c.date.Painter.Destroy(); err != nil {
		return fmt.Errorf("could not destroy date: %v", err)
	}
	return nil
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

// Mode set the mode
func (c *Clock) Mode(m mode) display.Animater {
	prev := c.mode
	c.mode = m
	c.Layer.Remove(c.time)
	c.Layer.Remove(c.date)
	hour := c.hour.Painter.(*display.Text)
	dot := c.dot.Painter.(*display.Text)
	minute := c.minute.Painter.(*display.Text)
	var yFullscreen int32 = 80
	sizeFullscreen := 95
	var yTop int32 = 5
	sizeTop := 60

	switch m {
	case Hidden:
	case Fullscreen:
		hour.Size = sizeFullscreen
		dot.Size = sizeFullscreen
		minute.Size = sizeFullscreen
		c.hour.X, c.hour.Y = 109, yFullscreen
		c.dot.X, c.dot.Y = 118, yFullscreen
		c.minute.X, c.minute.Y = 130, yFullscreen
		c.date.X, c.date.Y = 119, yFullscreen+105
		c.Layer.Add(c.time, c.date)
		if prev == Hidden {
			var once sync.Once
			setup := func() {
				c.time.Move(0, 230)
				c.date.Move(0, 130)
			}
			tl := &tween.Timeline{}
			tl.Add(tween.FromToInt32Delta(230, 0, time.Second, tween.EaseInOutQuad, func(d int32) {
				once.Do(setup)
				c.time.Move(0, d)
			}))
			tl.AddAt(400*time.Millisecond, tween.FromToInt32Delta(130, 0, 700*time.Millisecond, tween.EaseInOutQuad, func(d int32) {
				c.date.Move(0, d)
			}))
			return tl
		}
		if prev == Top {
			var once sync.Once
			scale := float32(sizeTop) / float32(sizeFullscreen)
			distance := yTop - yFullscreen
			setup := func() {
				c.time.Move(0, distance)
			}
			return tween.New(700*time.Millisecond, tween.EaseInOutQuad, func(d float32) {
				once.Do(setup)
				s := 1 + (scale-1)*(1-d)
				c.hour.SetScale(s)
				c.dot.SetScale(s)
				c.minute.SetScale(s)
				y := int32(float32(yFullscreen) + (float32(distance) * (1 - d)))
				c.hour.Y, c.dot.Y, c.minute.Y = y, y, y
			})
		}
	case Top:
		hour.Size = sizeTop
		dot.Size = sizeTop
		minute.Size = sizeTop
		c.hour.X, c.hour.Y = 109, yTop
		c.dot.X, c.dot.Y = 118, yTop
		c.minute.X, c.minute.Y = 130, yTop
		c.date.X, c.date.Y = 119, yTop+105
		c.Layer.Add(c.time)
		if prev == Fullscreen {
			var once sync.Once
			scale := float32(sizeFullscreen) / float32(sizeTop)
			distance := yFullscreen - yTop
			setup := func() {
				c.time.Move(0, distance)
			}
			return tween.New(700*time.Millisecond, tween.EaseInOutQuad, func(d float32) {
				once.Do(setup)
				s := 1 + (scale-1)*(1-d)
				c.hour.SetScale(s)
				c.dot.SetScale(s)
				c.minute.SetScale(s)
				y := int32(float32(yTop) + (float32(distance) * (1 - d)))
				c.hour.Y, c.dot.Y, c.minute.Y = y, y, y
			})
		}

	default:
		log.Panicf("Could not set clock mode to: %v", m)
	}
	return tween.Empty()
}

// Show the clock
func (c *Clock) Show() display.Animater {
	return c.Mode(Fullscreen)
}

// Hide the clock
func (c *Clock) Hide() display.Animater {
	return c.Mode(Hidden)
}

func (c *Clock) eventLoop() {
	for {
		c.busy.Lock()
		t := time.Now()
		c.hour.Painter.(*display.Text).Text = t.Format("15")
		c.minute.Painter.(*display.Text).Text = t.Format("04")
		c.date.Painter.(*display.Text).Text = t.Format("02 Jan")
		c.busy.Unlock()
		display.Refresh()
		select {
		case <-c.quit:
			return
		case <-time.After(time.Until(Next(time.Minute, t))):
		}
	}
}
