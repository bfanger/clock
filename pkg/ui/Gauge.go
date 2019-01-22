package ui

import (
	"fmt"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

// Guage clips a pie piece from an Image.
// @todo Add support for larger pieces than 180deg
type Guage struct {
	Sprite   *Sprite // The output
	Imager   Imager  // The input
	start    float64 // in degrees
	end      float64 // in degrees
	canvas   *Image
	renderer *sdl.Renderer
}

// NewGuage creates a new guage
func NewGuage(i Imager, start, end float64, r *sdl.Renderer) (*Guage, error) {
	canvas := &Image{}
	g := &Guage{
		Sprite:   NewSprite(canvas),
		Imager:   i,
		canvas:   canvas,
		renderer: r}
	if err := g.Set(start, end); err != nil {
		return nil, fmt.Errorf("couldn't create guage: %v", err)
	}
	return g, nil
}

// Close free resources
func (g *Guage) Close() error {
	return g.canvas.Close()
}

// Compose the gauge
func (g *Guage) Compose(r *sdl.Renderer) error {
	return g.Sprite.Compose(r)
}

// Set values and update the guage
func (g *Guage) Set(start, end float64) error {
	// normalize values
	start = math.Mod(start, 360)
	end = math.Mod(end, 360)
	if end < start {
		start -= 360
	}
	g.start = start
	g.end = end
	return g.update()
}

func (g *Guage) update() error {
	size := g.end - g.start
	// @todo optimize "size == 0"?
	r := g.renderer
	image, err := g.Imager.Image(r)
	if err != nil {
		return fmt.Errorf("couldn't read image: %v", err)
	}
	diameter := int32(math.Max(float64(image.Frame.W), float64(image.Frame.H)))
	radius := diameter / 2
	offset, err := r.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_TARGET, radius, diameter)
	if err != nil {
		return fmt.Errorf("couldn't create offset texture: %v", err)
	}
	defer offset.Destroy()
	if err := offset.SetBlendMode(sdl.BLENDMODE_BLEND); err != nil {
		return err
	}
	prevTarget := r.GetRenderTarget()
	if err := r.SetRenderTarget(offset); err != nil {
		return err
	}
	defer r.SetRenderTarget(prevTarget)
	if err := r.Clear(); err != nil {
		return err
	}
	pos := image.Frame
	pos.X -= radius
	if err := r.CopyEx(image.Texture, &image.Frame, &pos, -g.start, nil, sdl.FLIP_NONE); err != nil {
		return err
	}
	limit, err := r.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_TARGET, radius, diameter)
	if err != nil {
		return fmt.Errorf("couldn't create limit texture: %v", err)
	}
	defer limit.Destroy()
	if err := limit.SetBlendMode(sdl.BLENDMODE_BLEND); err != nil {
		return err
	}
	if err := r.SetRenderTarget(limit); err != nil {
		return err
	}
	if err := r.Clear(); err != nil {
		return err
	}
	src := &sdl.Rect{W: radius, H: diameter}
	pivot := &sdl.Point{Y: radius}
	if err := r.CopyEx(offset, src, src, 180-size, pivot, sdl.FLIP_NONE); err != nil {
		return err
	}
	if g.canvas.Texture != nil {
		if err := g.canvas.Texture.Destroy(); err != nil {
			return err
		}
	}
	if g.canvas.Texture, err = r.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_TARGET, diameter, diameter); err != nil {
		return err
	}
	if err := g.canvas.Texture.SetBlendMode(sdl.BLENDMODE_BLEND); err != nil {
		return err
	}
	g.canvas.Frame = sdl.Rect{W: diameter, H: diameter}
	if err := r.SetRenderTarget(g.canvas.Texture); err != nil {
		return err
	}
	if err := r.Clear(); err != nil {
		return err
	}
	dst := &sdl.Rect{X: radius, W: radius, H: diameter}
	if err := r.CopyEx(limit, src, dst, -180+size+g.start, pivot, sdl.FLIP_NONE); err != nil {
		return err
	}
	return nil
}
