package engine

// Renderable interface
type Renderable interface {
	Render() error
	Dispose() error
}
