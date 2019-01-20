package ui

import "github.com/veandco/go-sdl2/sdl"

// Container a collection of Composers
type Container struct {
	Layers []Composer
}

// Compose the layers
func (c *Container) Compose(r *sdl.Renderer) error {
	for _, layer := range c.Layers {
		if err := layer.Compose(r); err != nil {
			return err
		}
	}
	return nil
}

// Append a composer
func (c *Container) Append(layer Composer) {
	c.Layers = append(c.Layers, layer)
}

// Remove a composer
func (c *Container) Remove(layer Composer) {
	for i, l := range c.Layers {
		if l == layer {
			c.Layers = append(c.Layers[:i], c.Layers[i+1:]...)
			return
		}
	}
}
