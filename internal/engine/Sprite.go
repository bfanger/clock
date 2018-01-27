package engine

// Sprite interface
type Sprite interface {
	Render() error
	Destroy() error
}
