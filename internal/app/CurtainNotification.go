package app

import (
	"time"

	"github.com/bfanger/clock/pkg/tween"
	"github.com/bfanger/clock/pkg/ui"
	"github.com/veandco/go-sdl2/sdl"
)

// CurtainNotification is a animated swimming fish notification
type CurtainNotification struct {
	*BasicNotification
	right *ui.Sprite
	left  *ui.Sprite
}

// NewCurtainNotification create a Notification
func NewCurtainNotification(engine *ui.Engine, d time.Duration) (*CurtainNotification, error) {
	n, err := NewBasicNotification(engine, "curtain", d)
	if err != nil {
		return nil, err
	}
	left := ui.NewSprite(n.image)
	left.X = -168

	right := ui.NewSprite(n.image)
	right.X = screenWidth
	right.ScaleX = -1

	return &CurtainNotification{BasicNotification: n, left: left, right: right}, nil
}

// Show the curtains
func (n *CurtainNotification) Show() tween.Tween {
	return n.Animation()
}

// Hide the curtains
func (n *CurtainNotification) Hide() tween.Tween {
	return tween.Reverse(n.Animation())
}

// Compose the CurtainNotification
func (n *CurtainNotification) Compose(r *sdl.Renderer) error {
	if err := n.left.Compose(r); err != nil {
		return err
	}
	return n.right.Compose(r)
}

func (c *CurtainNotification) Animation() tween.Tween {
	tl := &tween.Timeline{}
	duration := 1200 * time.Millisecond

	tl.Add(tween.FromTo(0, 190, duration, tween.EaseOutQuad, func(a int32) {
		c.left.X = a - 190
		c.right.X = screenWidth + 22 - a
	}))
	tl.AddAt(0, tween.FromTo(-50, 0, duration, tween.EaseOutQuad, func(a int32) {
		c.left.Y = a
		c.right.Y = a
	}))
	tl.AddAt(0, tween.FromTo(10, 0, duration, tween.EaseOutQuad, func(a float64) {
		c.left.Rotation = a
		c.right.Rotation = -a
	}))
	return tl
}
