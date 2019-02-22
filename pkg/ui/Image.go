package ui

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"unsafe"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

// Image is the result of a paint and the base ingredient for the Compose()
type Image struct {
	Texture *sdl.Texture
	Frame   sdl.Rect
}

// ImageFromTexture creates an Image from a texture
func ImageFromTexture(t *sdl.Texture, frame sdl.Rect) *Image {
	return &Image{Texture: t, Frame: frame}
}

// Close frees the texture memory used by the image
func (i *Image) Close() error {
	return i.Texture.Destroy()
}

// Image also implements the Imager interface
func (i *Image) Image(r *sdl.Renderer) (*Image, error) {
	return i, nil
}

// Compose the image
func (i *Image) Compose(r *sdl.Renderer) error {
	return r.Copy(i.Texture, &i.Frame, &i.Frame)
}

// ImageFromSurface creates an Image from a surface
func ImageFromSurface(s *sdl.Surface, r *sdl.Renderer) (*Image, error) {
	t, err := r.CreateTextureFromSurface(s)
	if err != nil {
		return nil, err
	}
	return &Image{Texture: t, Frame: sdl.Rect{W: s.W, H: s.H}}, nil
}

// ImageFromFile loads an image from disk
func ImageFromFile(filename string, r *sdl.Renderer) (*Image, error) {
	s, err := img.Load(filename)
	if err != nil {
		return nil, fmt.Errorf("could not load image from %s: %v", filename, err)
	}
	defer s.Free()
	s.SetBlendMode(sdl.BLENDMODE_BLEND)
	return ImageFromSurface(s, r)
}

// ImageFromURL downloads the image from the web
func ImageFromURL(url string, r *sdl.Renderer) (*Image, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed: %s", response.Status)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	buffer := sdl.RWFromMem(unsafe.Pointer(&body[0]), len(body))
	surface, err := img.LoadRW(buffer, true)
	if err != nil {
		return nil, err
	}
	defer surface.Free()
	return ImageFromSurface(surface, r)
}
