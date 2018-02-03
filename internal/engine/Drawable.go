package engine

// Drawable interface
type Drawable interface {
	Draw() error
	Dispose() error
}
