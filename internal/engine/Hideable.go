package engine

// Hideable wraps a renderable and adds the abillity to toggle visibility
type Hideable struct {
	Renderable Renderable
	Visible    bool
}

// NewHideable creates a ready tot use hideable
func NewHideable(renderable Renderable) *Hideable {
	return &Hideable{
		Renderable: renderable,
		Visible:    true}
}

// Render all items
func (hideable *Hideable) Render() error {
	if hideable.Visible {
		return hideable.Renderable.Render()
	}
	return nil
}

// Dispose hideable
func (hideable *Hideable) Dispose() error {
	return nil
}
