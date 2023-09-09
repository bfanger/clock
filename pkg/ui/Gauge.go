package ui

import (
	"math"

	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/sdl"
)

// Guage clips a pie piece from an Image.
// @todo Add support for larger pieces than 180deg
type Guage struct {
	imager Imager  // The source image
	start  float64 // in degrees
	end    float64 // in degrees
	image  *Image
}

// NewGuage creates a new guage
func NewGuage(i Imager, start, end float64) *Guage {
	return &Guage{
		start:  start,
		end:    end,
		imager: i}
}

// Close free the texture memory
func (g *Guage) Close() error {
	return g.needsUpdate()
}

// Image creates the texture based of the values
func (g *Guage) Image(r *sdl.Renderer) (*Image, error) {
	if g.image != nil {
		return g.image, nil
	}
	source, err := g.imager.Image(r)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't read image")
	}
	// normalize values
	start := math.Mod(g.start, 360)
	end := math.Mod(g.end, 360)
	if end < start {
		start -= 360
	}
	size := end - start
	diameter := int32(math.Max(float64(source.Frame.W), float64(source.Frame.H)))
	radius := diameter / 2
	g.image = &Image{Frame: sdl.Rect{W: diameter, H: diameter}}
	if g.image.Texture, err = r.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_TARGET, diameter, diameter); err != nil {
		return nil, err
	}
	if err := g.image.Texture.SetBlendMode(sdl.BLENDMODE_BLEND); err != nil {
		return nil, err
	}
	prevTarget := r.GetRenderTarget()
	defer r.SetRenderTarget(prevTarget)
	if size == 0 {
		if err := r.SetRenderTarget(g.image.Texture); err != nil {
			return nil, err
		}
		if err := r.Clear(); err != nil {
			return nil, err
		}
	} else {
		offset, err := r.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_TARGET, radius, diameter)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't create offset texture")
		}
		defer offset.Destroy()
		if err := r.SetRenderTarget(offset); err != nil {
			return nil, err
		}
		if err := r.Clear(); err != nil {
			return nil, err
		}
		pos := source.Frame
		pos.X -= radius
		if err := r.CopyEx(source.Texture, &source.Frame, &pos, -start, nil, sdl.FLIP_NONE); err != nil {
			return nil, err
		}
		limit, err := r.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_TARGET, radius, diameter)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't create limit texture")
		}
		defer limit.Destroy()
		if err := r.SetRenderTarget(limit); err != nil {
			return nil, err
		}
		if err := r.Clear(); err != nil {
			return nil, err
		}
		src := &sdl.Rect{W: radius, H: diameter}
		pivot := &sdl.Point{Y: radius}
		if err := r.CopyEx(offset, src, src, 180-size, pivot, sdl.FLIP_NONE); err != nil {
			return nil, err
		}
		if err := r.SetRenderTarget(g.image.Texture); err != nil {
			return nil, err
		}
		if err := r.Clear(); err != nil {
			return nil, err
		}
		dst := &sdl.Rect{X: radius, W: radius, H: diameter}
		if err := r.CopyEx(limit, src, dst, -180+size+start, pivot, sdl.FLIP_NONE); err != nil {
			return nil, err
		}
	}
	return g.image, nil
}

// SetImager source
func (g *Guage) SetImager(i Imager) error {
	g.imager = i
	return g.needsUpdate()
}

// SetStart angle in degrees
func (g *Guage) SetStart(angle float64) error {
	g.start = angle
	return g.needsUpdate()
}

// SetEnd angle in degrees
func (g *Guage) SetEnd(angle float64) error {
	g.end = angle
	return g.needsUpdate()
}

// needsUpdate destroys the texture so the next call to Image() will generate a new image.
func (g *Guage) needsUpdate() error {
	if g.image != nil {
		if err := g.image.Close(); err != nil {
			return err
		}
		g.image = nil
	}
	return nil
}
