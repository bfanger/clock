package app

import (
	"github.com/bfanger/clock/pkg/ui"
)

// Marker on the map
type Marker struct {
	Latitude  float64
	Longitude float64
	Sprite    *ui.Sprite
}
