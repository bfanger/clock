package app

import (
	"sync"
	"time"

	"github.com/bfanger/clock/pkg/tween"
	"github.com/bfanger/clock/pkg/ui"
	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/sdl"
)

// WidgetManager manages what to show and when.
type WidgetManager struct {
	Scene *ui.Container
	clock interface {
		Close() error
		MoveTo(x, y int32)
		Compose(*sdl.Renderer) error
	}
	timer            *Timer
	splash           *Splash
	background       *ui.Container
	notifications    []Notification
	notificationLock sync.Mutex
	engine           *ui.Engine
}

// NewWidgetManager create a new WidgetManager
func NewWidgetManager(scene *ui.Container, e *ui.Engine) (*WidgetManager, error) {
	wm := &WidgetManager{engine: e, Scene: scene}
	var err error
	wm.background = &ui.Container{}
	if err != nil {
		return nil, err
	}
	wm.Scene.Append(wm.background)

	clock, err := NewAnalogClock(e)

	// clock, err := NewDigitalClock(e)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create clock")
	}
	wm.clock = clock
	wm.timer = clock.timer
	wm.Scene.Append(wm.clock)
	if wm.splash, err = NewSplash(e.Renderer); err != nil {
		return nil, errors.Wrap(err, "failed to create splash")
	}
	return wm, nil
}

// Close free memory used by the display elements
func (wm *WidgetManager) Close() error {
	wm.Scene.Remove(wm.clock)
	if err := wm.clock.Close(); err != nil {
		return err
	}
	wm.Scene.Remove(wm.splash)
	if err := wm.splash.Close(); err != nil {
		return err
	}
	// @todo use notificationLock?
	for _, n := range wm.notifications {
		if err := n.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Notify display a new notification
func (wm *WidgetManager) Notify(n Notification) {
	wm.notificationLock.Lock()
	wm.notifications = append(wm.notifications, n)
	wm.Scene.Append(n)
	if len(wm.notifications) == 1 {
		tl := &tween.Timeline{}
		tl.Add(tween.FromTo(screenWidth/2, 240, 700*time.Millisecond, tween.EaseInOutQuad, func(x int32) {
			wm.clock.MoveTo(x, screenHeight/2)
		}))
		tl.AddAt(800*time.Millisecond, n.Show())
		wm.engine.Animate(tl)
	} else {
		wm.engine.Animate(n.Show())
	}
	wm.notificationLock.Unlock()
	n.Wait()
	wm.notificationLock.Lock()
	defer wm.notificationLock.Unlock()
	if len(wm.notifications) == 1 {
		tl := &tween.Timeline{}
		tl.Add(n.Hide())
		tl.Add(tween.FromTo(240, screenWidth/2, 100*time.Millisecond, tween.EaseInOutQuad, func(x int32) {
			wm.clock.MoveTo(x, screenHeight/2)
		}))
		wm.engine.Animate(tl)
	} else {
		wm.engine.Animate(n.Hide())
	}
	for i, x := range wm.notifications {
		if n == x {
			wm.notifications = append(wm.notifications[:i], wm.notifications[i+1:]...)
			break
		}
	}
	wm.engine.Go(func() error {
		wm.Scene.Remove(n)
		return n.Close()
	})
}

// ButtonPressed show the splash image for a second
func (wm *WidgetManager) ButtonPressed() {
	wm.Scene.Append(wm.splash)
	wm.engine.Animate(wm.splash.Splash())
	wm.Scene.Remove(wm.splash)
}
