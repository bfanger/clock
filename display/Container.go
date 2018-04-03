package display

import (
	"fmt"
	"sort"

	"github.com/veandco/go-sdl2/sdl"
)

// Container is a layer that contains other layers.
type Container struct {
	layers map[int][]Layer
	depths []int
}

// NewContainer creates a Container
func NewContainer() *Container {
	return &Container{
		layers: make(map[int][]Layer),
		depths: []int{0},
	}
}

// Name of the container shows number of layers
func (c *Container) Name() string {
	n := 0
	for _, l := range c.layers {
		n += len(l)
	}
	return fmt.Sprintf("Container[%d]", n)
}

// Render all layers
func (c *Container) Render(r *sdl.Renderer) error {
	for _, z := range c.depths {
		for _, layer := range c.layers[z] {
			if err := layer.Render(r); err != nil {
				return fmt.Errorf("%s failed to render: %v", layer.Name(), err)
			}
		}
	}
	return nil
}

// Add a layer
func (c *Container) Add(l Layer) {
	c.layers[0] = append(c.layers[0], l)
}

// AddAt a layer at specific depth
func (c *Container) AddAt(l Layer, depth int) {
	new := c.layers[depth] == nil
	c.layers[depth] = append(c.layers[depth], l)
	if new {
		c.depths = append(c.depths, depth)
		sort.Ints(c.depths)
	}
}
