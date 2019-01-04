package app

import (
	"math"
	"time"

	"github.com/bfanger/clock/pkg/tween"
	"github.com/bfanger/clock/pkg/ui"
)

// FeedFishNotification is a animated swimming fish notification
type FeedFishNotification struct {
	*BasicNotification
	animation *tween.Cancelable
}

// NewFeedFishNotification create a Notification
func NewFeedFishNotification(engine *ui.Engine, d time.Duration) (*FeedFishNotification, error) {
	n, err := NewBasicNotification(engine, "vis", d)
	if err != nil {
		return nil, err
	}
	return &FeedFishNotification{BasicNotification: n}, nil
}

// Show the notification and start the swimming animation.
func (n *FeedFishNotification) Show() tween.Tween {
	return tween.Func(func() {
		go n.Swim()
	})
}

// Close the notification and stop the animation
func (n *FeedFishNotification) Close() error {
	n.animation.Cancel()
	return n.BasicNotification.Close()
}

// Swim the fish
func (n *FeedFishNotification) Swim() {
	var y int32 = 400
	tl := &tween.Timeline{}
	tl.Add(tween.Func(func() {
		n.sprite.SetAlpha(255)
	}))
	tl.Add(tween.Repeat(0, tween.FromToInt32(-100, 900, 12*time.Second, tween.Linear, func(x int32) {
		n.sprite.X = x
		n.sprite.Rotation = math.Sin(float64(x)/30) * 7
		n.sprite.Y = y + int32(math.Sin(float64(x)/30-math.Pi/2)*18)
	})))
	n.animation = tween.NewCancelable(tl)
	n.engine.Animate(n.animation)
}
