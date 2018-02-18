package engine

// Widget interface
type Widget interface {
	Mount(container *ContainerInterface) error
	Unmount() error
}
