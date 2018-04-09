package sprite

import (
	"fmt"

	"github.com/bfanger/clock/display"
	"github.com/veandco/go-sdl2/sdl"
)

// Sprite a thing to display on screen
type Sprite struct {
	name    string
	Painter display.Painter
	X, Y    int32
	// @todo Rotation, Pivot, , Alpha(prevAlpha)
	AnchorX, AnchorY float32
	ScaleX, ScaleY   float32
}

// New creates a new sprite
func New(name string, painter display.Painter, opts ...Option) *Sprite {
	s := &Sprite{name: name, Painter: painter, ScaleX: 1, ScaleY: 1}
	for _, opt := range opts {
		opt(s)
	}
	// , X: x, Y: y
	return s
}

// Name of the sprite
func (s *Sprite) Name() string {
	if s.name == "" {
		return "Sprite"
	}
	return s.name
}

// Render the sprite
func (s *Sprite) Render(r *sdl.Renderer) error {
	t, err := s.Painter.Paint(r)
	if err != nil {
		return err
	}
	if t == nil {
		return fmt.Errorf("Paint() returned nil. %T %+v", s.Painter, s.Painter)
	}
	x := s.X - int32(s.AnchorX*float32(t.Frame.W))
	y := s.Y - int32(s.AnchorY*float32(t.Frame.H))
	w := int32(s.ScaleX * float32(t.Frame.W))
	h := int32(s.ScaleY * float32(t.Frame.H))
	dst := &sdl.Rect{X: x, Y: y, W: w, H: h}
	r.Copy(t.Texture, t.Frame, dst)
	return nil
}

// SetScale in both X & Y direction
func (s *Sprite) SetScale(scale float32) {
	s.ScaleX = scale
	s.ScaleY = scale
}

// Move the sprite
func (s *Sprite) Move(dx, dy int32) {
	s.X += dx
	s.Y += dy
}

// Option of sprite.New
type Option func(*Sprite)

// WithPos sets the postion of the sprite
func WithPos(x, y int32) Option {
	return func(s *Sprite) {
		s.X = x
		s.Y = y
	}
}

// WithAnchor sets the anchor of the sprite
func WithAnchor(x, y float32) Option {
	return func(s *Sprite) {
		s.AnchorX = x
		s.AnchorY = y
	}
}
