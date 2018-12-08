package app

import (
	"fmt"
	"sync"
	"time"

	"github.com/bfanger/clock/pkg/tween"
	"github.com/bfanger/clock/pkg/ui"
)

// WidgetManager manages what to show and when.
type WidgetManager struct {
	background       *Background
	clock            *Clock
	splash           *Splash
	notifications    []Notification
	notificationLock sync.Mutex
	engine           *ui.Engine
}

// NewWidgetManager create a new WidgetManager
func NewWidgetManager(e *ui.Engine) (*WidgetManager, error) {
	wm := &WidgetManager{engine: e}
	var err error
	wm.background, err = NewBackground(e)
	if err != nil {
		return nil, fmt.Errorf("failed to create background: %v", err)
	}
	wm.clock, err = NewClock(e)
	if err != nil {
		return nil, fmt.Errorf("failed to create clock: %v", err)
	}
	wm.splash, err = NewSplash(e)
	if err != nil {
		return nil, fmt.Errorf("failed to create splash: %v", err)
	}
	return wm, nil

}

// Close free memory used by the display elements
func (wm *WidgetManager) Close() error {
	if err := wm.background.Close(); err != nil {
		return err
	}
	if err := wm.clock.Close(); err != nil {
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
	if len(wm.notifications) == 1 {
		tl := &tween.Timeline{}
		tl.Add(wm.clock.Minimize())
		tl.AddAt(200*time.Millisecond, wm.background.Maximize())
		tl.AddAt(800*time.Millisecond, n.Show())
		wm.engine.Animate(tl)
	} else {
		wm.engine.Animate(n.Show())
	}
	wm.notificationLock.Unlock()
	time.Sleep(n.Duration())
	wm.notificationLock.Lock()
	defer wm.notificationLock.Unlock()
	if len(wm.notifications) == 1 {
		tl := &tween.Timeline{}
		tl.Add(n.Hide())
		tl.AddAt(100*time.Millisecond, wm.clock.Maximize())
		tl.AddAt(100*time.Millisecond, wm.background.Minimize())
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
	wm.engine.Go(n.Close)
}

func (wm *WidgetManager) ButtonPressed() {
	wm.engine.Animate(wm.splash.Splash())
}
