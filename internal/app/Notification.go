package app

import (
	"github.com/bfanger/clock/pkg/tween"
	"github.com/veandco/go-sdl2/sdl"
)

// Notification widget
type Notification interface {
	Show() tween.Tween
	Hide() tween.Tween
	Wait()
	Close() error
	Compose(*sdl.Renderer) error
}
