package engine

import (
	"errors"
)

// ContainerInterface is used by Mount & Unmount
type ContainerInterface interface {
	Add(Drawable)
	Remove(Drawable) error
}

// Container with multiple Drawables
type Container struct {
	Items []Drawable
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
