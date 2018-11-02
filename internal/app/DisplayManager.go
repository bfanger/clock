package app

import (
	"fmt"
	"time"

	"github.com/bfanger/clock/pkg/tween"
	"github.com/bfanger/clock/pkg/ui"
)

// DisplayManager manages what to show and when.
type DisplayManager struct {
	background    *Background
	clock         *Clock
	notifications []Notification
	engine        *ui.Engine
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
	for _, n := range dm.notifications {
		if err := n.Close(); err != nil {
			return err
		}
	}
	return nil
}

// func (dm *DisplayManager) Notify(n Notification) error {
// notifications
// 	err := s.engine.Do(func() error {
// 		if s.Notification != nil {
// 			if err := s.Notification.Close(); err != nil {
// 				return fmt.Errorf("failed to close notification: %v", err)
// 			}
// 		}
// 		var err error
// 		if icon == "vis" {
// 			s.Notification, err = NewFeedFishNotification(s.engine)
// 		} else {
// 			s.Notification, err = NewBasicNotification(s.engine, icon)
// 		}
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// 	s.ShowNotification(s.Notification)
// } else {
// 	if err := s.HideNotification(); err != nil {
// 		panic(err)
// 	}
// return nil
// }

// ShowNotification display a new notification
func (dm *DisplayManager) Notify(n Notification) {
	// @todo lock access to notifications slice?
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
	time.Sleep(n.Duration())
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

// HideNotification hides the active notification
func (dm *DisplayManager) HideNotification() error {
	// if s.Notification == nil {
	// 	return errors.New("no notification active")
	// }
	// n := s.Notification
	// d.Notification = nil

	// if err := s.engine.Do(n.Close); err != nil {
	// 	return fmt.Errorf("failed to close notification: %v", err)
	// }
	return nil
}
