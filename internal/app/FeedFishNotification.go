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
	step := 1
	// for i := 0; i < 20; i++ {
	tl.Add(tween.FromToInt32(-90, 365, 8000*time.Millisecond, tween.Linear, func(x int32) {
		if step == 1 {
			n.sprite.X = x
		}
		if x == 365 {
			step = 2
		}
	}))
	tl.Add(tween.FromToInt32(-90, 366, 20000*time.Millisecond, tween.Linear, func(x int32) {
		if step == 2 {
			n.sprite.X = x
			if x == 366 {
				step = 3
			}
		}
	}))
	tl.Add(tween.FromToInt32(0, 1, 4000*time.Millisecond, tween.Linear, func(x int32) {
		if step == 3 {
			n.sprite.X = 160
		}
	}))
	// }
	return tl
}
