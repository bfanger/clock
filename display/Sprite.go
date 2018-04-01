package display

import "github.com/veandco/go-sdl2/sdl"

// Sprite a thing to display on screen
type Sprite struct {
	name    string
	Content Painter
	X, Y    int32
	// @todo Rotation, Pivot, ScaleX, ScaleY, Alpha(prevAlpha) AnchorX, AnchorY
}

func NewSprite(name string, content Painter, x, y int32) *Sprite {
	return &Sprite{name: name, Content: content, X: x, Y: y}
}

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
	dst := &sdl.Rect{X: s.X, Y: s.Y, W: t.Frame.W, H: t.Frame.H}
	r.Copy(t.texture, t.Frame, dst)
	return nil
}
