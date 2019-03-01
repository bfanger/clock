package app

import (
	"fmt"
	"time"

	"github.com/bfanger/clock/pkg/tween"
	"github.com/bfanger/clock/pkg/ui"
	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var white = sdl.Color{R: 255, G: 255, B: 255}
var orange = sdl.Color{R: 203, G: 87, B: 0, A: 255}

// DigitalClock displays the current time
type DigitalClock struct {
	engine *ui.Engine
	font   *ttf.Font
	text   *ui.Text
	sprite *ui.Sprite
	done   chan bool
}

// NewDigitalClock creats a new time widget
func NewDigitalClock(engine *ui.Engine) (*DigitalClock, error) {
	font, err := ttf.OpenFont(Asset("Roboto-Light.ttf"), 180)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open font")
	}
	text := ui.NewText("", font, white)
	sprite := ui.NewSprite(text)
	sprite.X = screenWidth / 2
	sprite.Y = screenHeight / 2
	sprite.AnchorX = 0.5
	sprite.AnchorY = 0.5
	sprite.SetAlpha(160)

	// sprite.SetScale(0.2)

	c := &DigitalClock{
		engine: engine,
		text:   text,
		font:   font,
		sprite: sprite,
		done:   make(chan bool)}

	if err := c.updateTime(); err != nil {
		return nil, err
	}
	go c.tick()

	return c, nil
}

// Close free resources
func (c *DigitalClock) Close() error {
	if err := c.text.Close(); err != nil {
		return err
	}
	close(c.done)
	c.font.Close()
	return nil
}

// MoveTo positions the clock
func (c *DigitalClock) MoveTo(x, y int32) {
	c.sprite.X = x
	c.sprite.Y = y
}

// Minimize time to make room for notifications
func (c *DigitalClock) Minimize() tween.Tween {
	tl := &tween.Timeline{}
	tl.Add(tween.FromToFloat32(1, 0.78, 1*time.Second, tween.EaseInOutQuad, c.sprite.SetScale))
	tl.AddAt(150*time.Millisecond, tween.FromToInt32(screenHeight/2, 110, 850*time.Millisecond, tween.EaseInOutQuad, func(v int32) {
		c.sprite.Y = v
	}))
	return tl
}

// Maximize time
func (c *DigitalClock) Maximize() tween.Tween {
	tl := &tween.Timeline{}
	tl.Add(tween.FromToFloat32(0.78, 1, 1*time.Second, tween.EaseInOutQuad, c.sprite.SetScale))
	tl.AddAt(150*time.Millisecond, tween.FromToInt32(110, screenHeight/2, 850*time.Millisecond, tween.EaseInOutQuad, func(v int32) {
		c.sprite.Y = v
	}))
	return tl
}

func (c *DigitalClock) updateTime() error {
	now := time.Now()
	time := fmt.Sprintf("%d%s", now.Hour(), now.Format(":04"))
	if err := c.text.SetText(time); err != nil {
		return err
	}
	return nil
}

func (c *DigitalClock) tick() {
	for {
		select {
		case <-c.done:
			return
		case <-time.After(time.Until(next(time.Minute, time.Now()))):
			c.engine.Go(c.updateTime)
		}
	}
}

// next calculates the time at which on the next d (Minute/Second) starts
func next(d time.Duration, since time.Time) time.Time {
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
