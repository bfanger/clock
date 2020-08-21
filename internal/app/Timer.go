package app

import (
	"fmt"
	"time"

	"github.com/bfanger/clock/pkg/ui"
	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/sdl"
)

// Timer is based on a Time Timer
type Timer struct {
	Sprite   *ui.Sprite
	duration time.Duration // The minute the timer stops
	started  time.Time
	scale    time.Duration
	engine   *ui.Engine
	gauge    *ui.Guage
	green    *ui.Image
	orange   *ui.Image
	done     chan bool
}

// NewTimer creates a timer
func NewTimer(e *ui.Engine) (*Timer, error) {
	var err error
	t := &Timer{engine: e, done: make(chan bool)}
	t.green, err = ui.ImageFromFile(Asset("timer/green.png"), e.Renderer)
	if err != nil {
		return nil, err
	}
	t.orange, err = ui.ImageFromFile(Asset("timer/orange.png"), e.Renderer)
	if err != nil {
		return nil, err
	}
	t.gauge = ui.NewGuage(t.green, 0, 0)
	t.Sprite = ui.NewSprite(t.gauge)
	go t.tick()
	return t, nil
}

// Close free resources
func (t *Timer) Close() error {
	close(t.done)
	if err := t.gauge.Close(); err != nil {
		return err
	}
	if err := t.green.Close(); err != nil {
		return err
	}
	return t.orange.Close()
}

// Compose the timer if needed
func (t *Timer) Compose(r *sdl.Renderer) error {
	if t.completed() {
		return nil
	}
	return t.Sprite.Compose(r)
}

// SetDuration of the timer
func (t *Timer) SetDuration(d time.Duration, scale time.Duration) error {
	if d > 30*time.Minute {
		return errors.Errorf("maximun duration of 30 min exceeded, got %v", d)
	}
	if d <= 0 {
		return errors.Errorf("invalid duration: %v", d)
	}
	t.duration = d
	t.scale = scale
	t.started = time.Now()
	return t.update()
}

// func (t *Timer) restart(d time.Duration) error {
// 	t.started = time.Now()
// 	return nil
// }

func (t *Timer) completed() bool {
	if t.duration == 0 {
		return true
	}
	now := time.Now()
	return now.Before(t.started) || now.After(t.started.Add(t.duration))
}

func (t *Timer) update() error {
	now := time.Now()
	if err := t.gauge.SetStart(time2deg(now, t.scale)); err != nil {
		return err
	}
	if err := t.gauge.SetEnd(time2deg(t.started.Add(t.duration), t.scale)); err != nil {
		return err
	}
	imager := t.green
	last10min := now.After(t.started.Add(t.duration - (10 * time.Minute)))
	if t.scale == time.Minute && last10min {
		imager = t.orange
	}
	return t.gauge.SetImager(imager)
}

func (t *Timer) tick() {
	// @todo
	// - only tick when a timer is active
	//   - wait until t.started
	//   - stop when t.duration is reached
	for {
		select {
		case <-t.done:
			return
		case <-time.After(time.Second):
			if t.completed() == false {
				t.engine.Go(t.update)
			}
		}
	}
}

func time2deg(t time.Time, scale time.Duration) float64 {
	switch scale {

	case time.Hour:
		hour := float64(t.Hour() % 12) // 30deg per hour
		minute := float64(t.Minute())  // 5deg per minute
		second := float64(t.Second())  // 1/12deg per sec
		return hour*30 + (minute / 2) + (second / 12)

	case time.Minute:
		minute := float64(t.Minute()) // 6deg per minute
		second := float64(t.Second()) // 0.1deg per sec
		return minute*6 + (second / 10)

	default:
		panic(fmt.Errorf("no time2deg for: %v", scale))
	}
}
