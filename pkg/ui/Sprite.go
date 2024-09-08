package ui

import (
	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/sdl"
)

// Sprite a thing to display on screen
type Sprite struct {
	Imager           Imager
	X, Y             int32
	AnchorX, AnchorY float32
	ScaleX, ScaleY   float32
	Rotation         float64
	image            *Image
	alpha            uint8
	dst              *sdl.Rect
}

// NewSprite creates a new sprite
func NewSprite(imager Imager) *Sprite {
	return &Sprite{
		Imager: imager,
		ScaleX: 1,
		ScaleY: 1,
		alpha:  255,
		dst:    &sdl.Rect{}}
}

// Compose the sprite
func (s *Sprite) Compose(r *sdl.Renderer) error {
	if s.alpha == 0 {
		return nil
	}
	if s.Imager == nil {
		return errors.New("Imager is required")
	}
	img, err := s.Imager.Image(r)
	if err != nil {
		return err
	}
	if img == nil {
		return errors.Errorf("Image() returned nil. %T %+v", s.Imager, s.Imager)
	}
	if s.image != img {
		err = img.Texture.SetAlphaMod(s.alpha)
		if err != nil && s.alpha != 255 {
			return err
		}
		s.image = img
	}
	flip := sdl.FLIP_NONE
	scaleX, scaleY := s.ScaleX, s.ScaleY
	if scaleX < 0 {
		scaleX *= -1
		flip = sdl.FLIP_HORIZONTAL
	}
	if scaleY < 0 {
		scaleY *= -1
		if flip == sdl.FLIP_HORIZONTAL {
			flip |= sdl.FLIP_VERTICAL
		} else {
			flip = sdl.FLIP_VERTICAL
		}
	}

	s.dst.W = int32(scaleX * float32(img.Frame.W))
	s.dst.H = int32(scaleY * float32(img.Frame.H))
	s.dst.X = s.X - int32(s.AnchorX*float32(s.dst.W))
	s.dst.Y = s.Y - int32(s.AnchorY*float32(s.dst.H))
	if flip == sdl.FLIP_NONE && s.Rotation == 0 {
		return r.Copy(img.Texture, &img.Frame, s.dst)
	}
	// @todo Pivot
	return r.CopyEx(img.Texture, &img.Frame, s.dst, s.Rotation, nil, flip)
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

func (s *Sprite) GetAlpha() uint8 {
	return s.alpha
}
