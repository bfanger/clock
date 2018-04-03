package display

import "github.com/veandco/go-sdl2/sdl"

// Sprite a thing to display on screen
type Sprite struct {
	name    string
	Content Painter
	X, Y    int32
	// @todo Rotation, Pivot, ScaleX, ScaleY, Alpha(prevAlpha)
	AnchorX, AnchorY float32
}

// NewSprite creates a new sprite
func NewSprite(name string, content Painter, x, y int32) *Sprite {
	return &Sprite{name: name, Content: content, X: x, Y: y}
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
	x := s.X - int32(s.AnchorX*float32(t.Frame.W))
	y := s.Y - int32(s.AnchorY*float32(t.Frame.H))
	dst := &sdl.Rect{X: x, Y: y, W: t.Frame.W, H: t.Frame.H}
	r.Copy(t.texture, t.Frame, dst)
	return nil
}
