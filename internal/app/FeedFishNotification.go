package app

import (
	"time"

	"github.com/bfanger/clock/pkg/tween"
	"github.com/bfanger/clock/pkg/ui"
)

// import (
// 	"time"

// 	"github.com/bfanger/clock/pkg/tween"
// )

type FeedFishNotification struct {
	*BasicNotification
}

func NewFeedFishNotification(engine *ui.Engine) (*FeedFishNotification, error) {
	n, err := NewBasicNotification(engine, "vis")
	if err != nil {
		return nil, err
	}
	n.sprite.SetAlpha(255)
	return &FeedFishNotification{BasicNotification: n}, nil

}
func (n *FeedFishNotification) Show() tween.Tween {
	tl := &tween.Timeline{}
	for i := 0; i < 20; i++ {
		tl.Add(tween.FromToInt32(-90, 365, 8000*time.Millisecond, tween.Linear, func(x int32) {
			n.sprite.X = x
		}))
	}
	return tl
}
