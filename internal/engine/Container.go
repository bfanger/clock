package engine

import (
	"errors"

	"github.com/veandco/go-sdl2/sdl"
)

// Container with multiple Drawables
type Container struct {
	Renderer *sdl.Renderer
	Items    []Drawable
}

// NewContainer creates a ready tot use container
func NewContainer(renderer *sdl.Renderer) *Container {
	return &Container{Renderer: renderer}
}

// Draw all items
func (container *Container) Draw() error {
	for _, item := range container.Items {
		if err := item.Draw(); err != nil {
			return err
		}
	}
	return nil
}

// Dispose nothing
func (container *Container) Dispose() error {
	return nil
}

// DisposeItems disposes all items
func (container *Container) DisposeItems() error {
	for _, item := range container.Items {
		if err := item.Dispose(); err != nil {
			return err
		}
	}
	container.Items = nil
	return nil
}

// Render the frame
func (container *Container) Render() error {
	if err := container.Renderer.Clear(); err != nil {
		panic(err)
	}

	if err := container.Draw(); err != nil {
		panic(err)
	}
	container.Renderer.Present()
	return nil
}

// Add item
func (container *Container) Add(item Drawable) {
	container.Items = append(container.Items, item)
}

// Remove item
func (container *Container) Remove(item Drawable) error {

	for index, _item := range container.Items {
		if _item == item {
			container.Items = append(container.Items[:index], container.Items[index+1:]...)
			return nil
		}
	}
	return errors.New("Item not found")
}
