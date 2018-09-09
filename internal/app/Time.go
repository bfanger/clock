package app

import (
	"sync"
	"time"

	"github.com/bfanger/clock/pkg/tween"
	"github.com/bfanger/clock/pkg/ui"
	"github.com/bfanger/clock/pkg/ui/text"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var white = sdl.Color{R: 255, G: 255, B: 255}
var orange = sdl.Color{R: 254, G: 110, B: 2, A: 255}

// Time display the current time
func Time(engine *ui.Engine, font *ttf.Font) {
	text := text.New(time.Now().Format("15:04"), font, text.WithColor(orange))
	var textHeight int32
	wg := sync.WaitGroup{}
	wg.Add(1)
	engine.Go(func() error {
		image, err := text.Image(engine.Renderer)
		if err != nil {
			return err
		}
		if image != nil {
			text.X = (320 / 2) - (image.Frame.W / 2)
			text.Y = 240
			textHeight = image.Frame.H
		}
		engine.Append(text)
		wg.Done()
		return nil
	})
	wg.Wait()
	time.Sleep(200 * time.Millisecond)
	t := tween.FromToInt32(240, (240/2)-(textHeight/2), 1*time.Second, tween.EaseInOutQuad, func(y int32) {
		text.Y = y
	})
	engine.Animate(t)

	for {
		engine.Go(func() error {
			if err := text.SetText(time.Now().Format("15:04")); err != nil {
				return err
			}
			image, err := text.Image(engine.Renderer)
			if err != nil {
				return err
			}
			if image != nil {
				text.X = (320 / 2) - (image.Frame.W / 2)
			}
			return nil
		})
		time.Sleep(time.Until(next(time.Minute, time.Now())))
	}
}

// Next creates the next rounded date based on the step
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
