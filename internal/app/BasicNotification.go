package app

import (
	"time"

	"github.com/bfanger/clock/pkg/tween"
	"github.com/bfanger/clock/pkg/ui"
	"github.com/veandco/go-sdl2/sdl"
)

// BasicNotification for notifications
type BasicNotification struct {
	image    *ui.Image
	sprite   *ui.Sprite
	duration time.Duration
	engine   *ui.Engine
}

// NewBasicNotification creates a new Notification
func NewBasicNotification(engine *ui.Engine, icon string, d time.Duration) (*BasicNotification, error) {
	image, err := ui.ImageFromFile(Asset("notifications/"+icon+".png"), engine.Renderer)
	if err != nil {
		return nil, err
	}
	sprite := ui.NewSprite(image)
	sprite.X = screenWidth - ((screenWidth - 480) / 2)
	sprite.AnchorX = 0.5
	sprite.AnchorY = 0.5
	sprite.Y = 240
	sprite.SetAlpha(0)

	return &BasicNotification{
		image:    image,
		sprite:   sprite,
		duration: d,
		engine:   engine}, nil
}

// Close free memory used by the Notification
func (n *BasicNotification) Close() error {
	return n.image.Close()
}

// Compose the notification
func (n *BasicNotification) Compose(r *sdl.Renderer) error {
	return n.sprite.Compose(r)
}

// Show notification
func (n *BasicNotification) Show() tween.Tween {
	return tween.FromToUint8(0, 255, 1000*time.Millisecond, tween.EaseOutQuad, n.sprite.SetAlpha)
}

// Hide notification
func (n *BasicNotification) Hide() tween.Tween {
	return tween.FromToUint8(255, 0, 500*time.Millisecond, tween.EaseOutQuad, n.sprite.SetAlpha)
}

// Wait the configured duration
func (n *BasicNotification) Wait() {
	time.Sleep(n.duration)
}
