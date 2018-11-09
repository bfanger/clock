package app

import (
	"time"

	"github.com/bfanger/clock/pkg/tween"
)

// Notification widget
type Notification interface {
	Show() tween.Tween
	Hide() tween.Tween
	Duration() time.Duration
	Close() error
}
