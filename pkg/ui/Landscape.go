package ui

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Landscape struct {
	Composer    Composer
	passthrough bool
	texture     *sdl.Texture
	src         sdl.Rect
	dest        sdl.Rect
	pivot       sdl.Point
}

func NewLandscape(composer Composer) *Landscape {
	return &Landscape{
		Composer: composer,
	}
}

func (l *Landscape) Compose(r *sdl.Renderer) error {
	if l.passthrough {
		return l.Composer.Compose(r)
	}
	if l.texture == nil {
		viewport := r.GetViewport()
		if viewport.W > viewport.H {
			l.passthrough = true
			return l.Composer.Compose(r)
		}
		l.src = sdl.Rect{W: viewport.H, H: viewport.W}
		l.dest = sdl.Rect{X: 0, Y: -viewport.W, W: viewport.H, H: viewport.W}
		l.pivot = sdl.Point{X: 0, Y: viewport.W}

		texture, err := r.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_TARGET, l.src.W, l.src.H)
		if err != nil {
			return err
		}
		l.texture = texture
	}
	prevTarget := r.GetRenderTarget()
	defer r.SetRenderTarget(prevTarget)
	if err := r.SetRenderTarget(l.texture); err != nil {
		return err
	}
	if err := r.Clear(); err != nil {
		return err
	}
	if err := l.Composer.Compose(r); err != nil {
		return err
	}
	if err := r.SetRenderTarget(prevTarget); err != nil {
		return err
	}
	return r.CopyEx(l.texture, &l.src, &l.dest, 90, &l.pivot, sdl.FLIP_NONE)
}

func (l *Landscape) Close() error {
	if l.texture == nil {
		return nil
	}
	return l.texture.Destroy()
}
