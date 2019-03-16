package app

import (
	"os"

	"github.com/pkg/errors"

	"github.com/bfanger/clock/pkg/ui"
	"github.com/veandco/go-sdl2/sdl"
)

// Location manages displaying the map
type Location struct {
	m *Map
}

// NewLocation creates a new location
func NewLocation(e *ui.Engine) (*Location, error) {
	key := os.Getenv("MAPTILER_KEY")
	if key == "" {
		return nil, errors.New("Invalid MAPTILER_KEY")
	}
	m := NewMap(key, e)
	m.Alpha = 180
	m.Zoom = 16
	m.Latitude = 52.4900311
	m.Longitude = 4.7602125
	icon, err := ui.ImageFromFile(Asset("position.png"), e.Renderer)
	if err != nil {
		return nil, err
	}
	m.CenterOffsetX = 240
	marker := &Marker{
		Latitude:  m.Latitude,
		Longitude: m.Longitude,
		Sprite:    ui.NewSprite(icon)}
	marker.Sprite.AnchorX = 0.5
	marker.Sprite.AnchorY = 0.5
	m.Markers = append(m.Markers, marker)

	return &Location{m: m}, nil
}

// Close the location viewer
func (l *Location) Close() error {
	return l.m.Close()
}

// Compose the location viewer
func (l *Location) Compose(r *sdl.Renderer) error {
	return nil
	// return l.m.Compose(r)
}
