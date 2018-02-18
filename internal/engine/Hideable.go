package engine

// Hideable wraps a Drawable and adds the abillity to toggle visibility
type Hideable struct {
	Drawable Drawable
	Visible  bool
}

// NewHideable creates a ready tot use hideable
func NewHideable(Drawable Drawable) *Hideable {
	return &Hideable{
		Drawable: Drawable,
		Visible:  true}
}

// Draw the items when it's visible
func (hideable *Hideable) Draw() error {
	if hideable.Visible {
		return hideable.Drawable.Draw()
	}
	return nil
}

// Dispose hideable
func (hideable *Hideable) Dispose() error {
	return nil
}
