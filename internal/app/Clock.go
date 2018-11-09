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
var orange = sdl.Color{R: 203, G: 87, B: 0, A: 255}

// Clock displays the current time
type Clock struct {
	engine *ui.Engine
	font   *ttf.Font
	text   *ui.Text
	sprite *ui.Sprite
	done   chan bool
}

// NewClock creats a new time widget
func NewClock(engine *ui.Engine) (*Clock, error) {
	font, err := ttf.OpenFont(Asset("Roboto-Light.ttf"), 220)
	if err != nil {
		return nil, fmt.Errorf("unable to open font: %v", err)
	}
	text := ui.NewText("", font, orange)
	sprite := ui.NewSprite(text)
	sprite.X = screenWidth / 2
	sprite.Y = screenHeight / 2
	sprite.AnchorX = 0.5
	sprite.AnchorY = 0.5

	// sprite.SetScale(0.2)

	t := &Clock{
		engine: engine,
		text:   text,
		font:   font,
		sprite: sprite,
		done:   make(chan bool)}

	if err := t.updateTime(); err != nil {
		return nil, err
	}
	engine.Append(t.sprite)
	go t.tick()

	return t, nil
}

// Close free resources
func (t *Clock) Close() error {
	t.engine.Remove(t.text)
	if err := t.text.Close(); err != nil {
		return err
	}
	close(t.done)
	t.font.Close()
	return nil
}

// Minimize time to make room for notifications
func (t *Clock) Minimize() tween.Tween {
	tl := &tween.Timeline{}
	tl.Add(tween.FromToFloat32(1, 0.72, 1*time.Second, tween.EaseInOutQuad, t.sprite.SetScale))
	tl.AddAt(150*time.Millisecond, tween.FromToInt32(screenHeight/2, 90, 850*time.Millisecond, tween.EaseInOutQuad, func(v int32) {
		t.sprite.Y = v
	}))
	return tl
}

// Maximize time
func (t *Clock) Maximize() tween.Tween {
	tl := &tween.Timeline{}
	tl.Add(tween.FromToFloat32(0.72, 1, 1*time.Second, tween.EaseInOutQuad, t.sprite.SetScale))
	tl.AddAt(150*time.Millisecond, tween.FromToInt32(90, screenHeight/2, 850*time.Millisecond, tween.EaseInOutQuad, func(v int32) {
		t.sprite.Y = v
	}))
	return tl
}

func (t *Clock) updateTime() error {
	now := time.Now()
	time := fmt.Sprintf("%d%s", now.Hour(), now.Format(":04"))
	if err := t.text.SetText(time); err != nil {
		return err
	}
	return nil
}

func (t *Clock) tick() {
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
