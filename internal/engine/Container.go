package engine

import (
	"github.com/veandco/go-sdl2/sdl"
)

// Container with multiple renderables
type Container struct {
	Renderer *sdl.Renderer
	Items    []Renderable
}

// NewContainer creates a ready tot use container
func NewContainer(renderer *sdl.Renderer) *Container {
	return &Container{
		Renderer: renderer,
		Items:    make([]Renderable, 0)}
}

// Render all items
func (container *Container) Render() error {
	for _, item := range container.Items {
		if err := item.Render(); err != nil {
			return err
		}
	}
	return nil
}

// Dispose all items
func (container *Container) Dispose() error {
	return nil
}

// Add item
func (container *Container) Add(item Renderable) {
	container.Items = append(container.Items, item)
}
