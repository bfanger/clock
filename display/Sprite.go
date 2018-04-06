package display

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

// Sprite a thing to display on screen
type Sprite struct {
	name    string
	Content Painter
	X, Y    int32
	// @todo Rotation, Pivot, , Alpha(prevAlpha)
	AnchorX, AnchorY float32
	ScaleX, ScaleY   float32
}

// NewSprite creates a new sprite
func NewSprite(name string, content Painter, x, y int32) *Sprite {
	return &Sprite{name: name, Content: content, X: x, Y: y, ScaleX: 1, ScaleY: 1}
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
	t, err := s.Content.Paint(r)
	if err != nil {
		return err
	}
	if t == nil {
		return fmt.Errorf("paint result was nil. %T %+v", s.Content, s.Content)
	}
	x := s.X - int32(s.AnchorX*float32(t.Frame.W))
	y := s.Y - int32(s.AnchorY*float32(t.Frame.H))
	w := int32(s.ScaleX * float32(t.Frame.W))
	h := int32(s.ScaleY * float32(t.Frame.H))
	dst := &sdl.Rect{X: x, Y: y, W: w, H: h}
	r.Copy(t.texture, t.Frame, dst)
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
