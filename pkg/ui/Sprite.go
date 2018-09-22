package ui

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

// Sprite a thing to display on screen
type Sprite struct {
	Imager Imager
	X, Y   int32
	// @todo Rotation, Pivot
	AnchorX, AnchorY float32
	ScaleX, ScaleY   float32
	alpha            uint8
	image            *Image
}

// NewSprite creates a new sprite
func NewSprite(imager Imager) *Sprite {
	return &Sprite{Imager: imager, ScaleX: 1, ScaleY: 1, alpha: 255}
}

// Compose the sprite
func (s *Sprite) Compose(r *sdl.Renderer) error {
	img, err := s.Imager.Image(r)
	if err != nil {
		return err
	}
	if img == nil {
		return fmt.Errorf("Image() returned nil. %T %+v", s.Imager, s.Imager)
	}
	if s.image != img {
		err = img.Texture.SetAlphaMod(s.alpha)
		if err != nil && s.alpha != 255 {
			return err
		}
		s.image = img
	}
	w := int32(s.ScaleX * float32(img.Frame.W))
	h := int32(s.ScaleY * float32(img.Frame.H))
	x := s.X - int32(s.AnchorX*float32(w))
	y := s.Y - int32(s.AnchorY*float32(h))
	dst := &sdl.Rect{X: x, Y: y, W: w, H: h}
	return r.Copy(img.Texture, &img.Frame, dst)
}

// SetScale in both X & Y direction
func (s *Sprite) SetScale(scale float32) {
	s.ScaleX = scale
	s.ScaleY = scale
}

// SetAlpha sets the alpha
func (s *Sprite) SetAlpha(a uint8) {
	if s.alpha != a {
		s.alpha = a
		s.image = nil
	}
}

// Move the sprite
// func (s *Sprite) Move(dx, dy int32) {
// 	s.X += dx
// 	s.Y += dy
// }

// // Option of sprite.New
// type Option func(*Sprite)

// // WithPos sets the postion of the sprite
// func WithPos(x, y int32) Option {
// 	return func(s *Sprite) {
// 		s.X = x
// 		s.Y = y
// 	}
// }

// // WithAnchor sets the anchor of the sprite
// func WithAnchor(x, y float32) Option {
// 	return func(s *Sprite) {
// 		s.AnchorX = x
// 		s.AnchorY = y
// 	}
// }

// // WithScale sets the scale of the sprite
// func WithScale(x, y float32) Option {
// 	return func(s *Sprite) {
// 		s.ScaleX = x
// 		s.ScaleY = y
// 	}
// }

// // WithAlpha sets the opacity of the sprite
// func WithAlpha(a uint8) Option {
// 	return func(s *Sprite) {
// 		s.Alpha = a
// 	}
// }
