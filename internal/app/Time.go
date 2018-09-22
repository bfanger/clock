package app

import (
	"fmt"
	"time"

	"github.com/bfanger/clock/pkg/tween"
	"github.com/bfanger/clock/pkg/ui"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var white = sdl.Color{R: 255, G: 255, B: 255}
var orange = sdl.Color{R: 254, G: 110, B: 2, A: 255}

// Time display the current time
type Time struct {
	engine *ui.Engine
	font   *ttf.Font
	text   *ui.Text
	done   chan bool
}

// NewTime creats a new time widget
func NewTime(engine *ui.Engine) (*Time, error) {
	font, err := ttf.OpenFont(asset("Roboto-Light.ttf"), 110)
	if err != nil {
		return nil, fmt.Errorf("unable to open font: %v", err)
	}
	text := ui.NewText("", font, orange)
	text.Y = screenHeight
	t := &Time{
		engine: engine,
		text:   text,
		font:   font,
		done:   make(chan bool)}

	if err := t.updateTime(); err != nil {
		return nil, err
	}
	engine.Append(text)

	go t.tick()

	return t, nil
}

// Close free resources
func (t *Time) Close() error {
	t.engine.Remove(t.text)
	if err := t.text.Close(); err != nil {
		return err
	}
	close(t.done)
	t.font.Close()
	return nil
}

func (t *Time) updateTime() error {
	now := time.Now()
	time := fmt.Sprintf("%d%s", now.Hour(), now.Format(":04"))
	if err := t.text.SetText(time); err != nil {
		return err
	}
	image, err := t.text.Image(t.engine.Renderer)
	if err != nil {
		return err
	}
	if image != nil {
		t.text.X = (320 / 2) - (image.Frame.W / 2)
	}
	return nil
}

func (t *Time) Intro() error {
	height, err := t.text.Height(t.engine.Renderer)
	if err != nil {
		return err
	}
	tween := tween.FromToInt32(screenHeight, screenHeight/2-(height/2), 1*time.Second, tween.EaseInOutQuad, func(y int32) {
		t.text.Y = y
	})
	go t.engine.Animate(tween)
	return nil
}

func (t *Time) tick() {
	for {
		select {
		case <-t.done:
			return
		case <-time.After(time.Until(next(time.Minute, time.Now()))):
			t.engine.Go(t.updateTime)
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
