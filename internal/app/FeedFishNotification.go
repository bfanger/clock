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
}

// NewFeedFishNotification create a Notification
func NewFeedFishNotification(engine *ui.Engine) (*FeedFishNotification, error) {
	n, err := NewBasicNotification(engine, "vis")
	if err != nil {
		return nil, err
	}
	return &FeedFishNotification{BasicNotification: n}, nil
}

// Show the fish
func (n *FeedFishNotification) Show() tween.Tween {
	tl := &tween.Timeline{}
	y := n.sprite.Y
	tl.Add(tween.Func(func() {
		n.sprite.SetAlpha(255)
	}))
	tl.Add(tween.Repeat(0, tween.FromToInt32(-90, 365, 10*time.Second, tween.Linear, func(x int32) {
		n.sprite.X = x
		n.sprite.Rotation = math.Sin(float64(x)/20) * 7
		n.sprite.Y = y + int32(math.Sin(float64(x)/20-math.Pi/2)*10)
	})))
	return tl
}
