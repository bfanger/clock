package app

import (
	"fmt"
	"sync"
	"time"

	"github.com/bfanger/clock/pkg/tween"
	"github.com/bfanger/clock/pkg/ui"
)

// DisplayManager manages what to show and when.
type DisplayManager struct {
	background       *Background
	clock            *Clock
	notifications    []Notification
	notificationLock sync.Mutex
	engine           *ui.Engine
}

// NewDisplayManager create a new DisplayManager
func NewDisplayManager(e *ui.Engine) (*DisplayManager, error) {
	dm := &DisplayManager{engine: e}
	var err error
	dm.background, err = NewBackground(e)
	if err != nil {
		return nil, fmt.Errorf("failed to create background: %v", err)
	}
	dm.clock, err = NewClock(e)
	if err != nil {
		return nil, fmt.Errorf("failed to create clock: %v", err)
	}
	return dm, nil

}

// Close free memory used by the display elements
func (dm *DisplayManager) Close() error {
	if err := dm.background.Close(); err != nil {
		return err
	}
	if err := dm.clock.Close(); err != nil {
		return err
	}
	// @todo use notificationLock?
	for _, n := range dm.notifications {
		if err := n.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Notify display a new notification
func (dm *DisplayManager) Notify(n Notification) {
	dm.notificationLock.Lock()
	dm.notifications = append(dm.notifications, n)
	if len(dm.notifications) == 1 {
		tl := &tween.Timeline{}
		tl.Add(dm.clock.Minimize())
		tl.AddAt(200*time.Millisecond, dm.background.Maximize())
		tl.AddAt(800*time.Millisecond, n.Show())
		dm.engine.Animate(tl)
	} else {
		dm.engine.Animate(n.Show())
	}
	dm.notificationLock.Unlock()
	time.Sleep(n.Duration())
	dm.notificationLock.Lock()
	defer dm.notificationLock.Unlock()
	if len(dm.notifications) == 1 {
		tl := &tween.Timeline{}
		tl.Add(n.Hide())
		tl.AddAt(100*time.Millisecond, dm.clock.Maximize())
		tl.AddAt(100*time.Millisecond, dm.background.Minimize())
		dm.engine.Animate(tl)
	} else {
		dm.engine.Animate(n.Hide())
	}
	for i, x := range dm.notifications {
		if n == x {
			dm.notifications = append(dm.notifications[:i], dm.notifications[i+1:]...)
			break
		}
	}
	dm.engine.Go(n.Close)
}
