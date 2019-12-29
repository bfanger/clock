package app

import (
	"time"

	"github.com/bfanger/clock/pkg/tween"
	"github.com/bfanger/clock/pkg/ui"
)

// TrayNotification is a animated swimming fish notification
type TrayNotification struct {
	*BasicNotification
}

// NewTrayNotification create a Notification
func NewTrayNotification(icon string, engine *ui.Engine, d time.Duration) (*TrayNotification, error) {
	n, err := NewBasicNotification(engine, icon, d)
	n.sprite.X = screenWidth - 5
	n.sprite.AnchorX = 1
	n.sprite.AnchorY = 0
	n.sprite.Y = 5
	if err != nil {
		return nil, err
	}
	return &TrayNotification{BasicNotification: n}, nil
}

// Show the notification and start the swimming animation.
func (n *TrayNotification) Show() tween.Tween {
	return tween.FromToUint8(0, 255, 1000*time.Millisecond, tween.EaseOutQuad, n.sprite.SetAlpha)
}

// Close the notification and stop the animation
func (n *TrayNotification) Close() error {
	return n.BasicNotification.Close()
}
