package app

import (
	"time"

	"github.com/bfanger/clock/pkg/tween"
	"github.com/bfanger/clock/pkg/ui"
)

// Notification for notifications
type Notification struct {
	image  *ui.Image
	sprite *ui.Sprite
	engine *ui.Engine
}

// NewNotification creates a new Notification
func NewNotification(engine *ui.Engine) (*Notification, error) {
	image, err := ui.ImageFromFile(asset("kliko.png"), engine.Renderer)
	if err != nil {
		return nil, err
	}
	sprite := ui.NewSprite(image)
	sprite.X = screenWidth / 2
	sprite.AnchorX = 0.5
	sprite.Y = 80
	engine.Append(sprite)

	return &Notification{
		image:  image,
		sprite: sprite,
		engine: engine}, nil
}

// Close free memory used by the Notification
func (n *Notification) Close() error {
	return n.image.Close()
}

// Show notification
func (n *Notification) Show() {
	go n.engine.Animate(tween.FromToUint8(0, 255, 1500*time.Millisecond, tween.EaseOutQuad, n.sprite.SetAlpha))
}

// Hide notification
func (n *Notification) Hide() {
	go n.engine.Animate(tween.FromToUint8(255, 0, 1000*time.Millisecond, tween.EaseOutQuad, n.sprite.SetAlpha))
}
