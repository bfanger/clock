package display

import (
	"fmt"
	"sort"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
)

// Container is a layer that contains other layers.
type Container struct {
	layers map[int][]Layer
	depths []int
	mu     sync.RWMutex
}

// NewContainer creates a Container
func NewContainer() *Container {
	return &Container{
		layers: make(map[int][]Layer),
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
	c.mu.RLock()
	defer c.mu.RUnlock()
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
	c.AddAt(l, 0)
}

// AddAt a layer at specific depth
func (c *Container) AddAt(l Layer, depth int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	new := c.layers[depth] == nil
	c.layers[depth] = append(c.layers[depth], l)
	if new {
		c.depths = append(c.depths, depth)
		sort.Ints(c.depths)
	}
}

// Remove a layer
func (c *Container) Remove(l Layer) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, z := range c.depths {
		for i, layer := range c.layers[z] {
			if layer == l {
				c.layers[z] = append(c.layers[z][:i], c.layers[z][i+1:]...)
				return true
			}
		}
	}
	return false
}

// RemoveAt removes a layer at a specific depth
func (c *Container) RemoveAt(l Layer, depth int) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i, layer := range c.layers[depth] {
		if layer == l {
			c.layers[depth] = append(c.layers[depth][:i], c.layers[depth][i+1:]...)
			return true
		}
	}
	return false
}

// Move the contents
func (c *Container) Move(dx, dy int32) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, z := range c.depths {
		for _, layer := range c.layers[z] {
			layer.Move(dx, dy)
		}
	}
}
