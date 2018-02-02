package engine

import (
	"errors"

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

// Dispose nothing
func (container *Container) Dispose() error {
	// for _, item := range container.Items {
	// 	if err := item.Dispose(); err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

// Add item
func (container *Container) Add(item Renderable) {
	container.Items = append(container.Items, item)
}

// Remove item
func (container *Container) Remove(item Renderable) error {

	for index, _item := range container.Items {
		if _item == item {
			container.Items = append(container.Items[:index], container.Items[index+1:]...)
			return nil
		}
	}
	return errors.New("Item not found")
}
