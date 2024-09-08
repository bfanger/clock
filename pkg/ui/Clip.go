package ui

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Clip struct {
	*sdl.Rect
	Composer Composer
}

func (c *Clip) Compose(r *sdl.Renderer) error {
	prev := r.GetClipRect()
	restore := &prev
	if prev.W == 0 {
		restore = nil
	}
	defer r.SetClipRect(restore)

	if err := r.SetClipRect(c.Rect); err != nil {
		return err
	}
	if err := c.Composer.Compose(r); err != nil {
		return err
	}
	return nil

}
