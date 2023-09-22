package app

import (
	"errors"
	"os"
	"time"

	"github.com/bfanger/clock/pkg/tween"
	"github.com/bfanger/clock/pkg/ui"
	"github.com/veandco/go-sdl2/sdl"
)

var gps = struct {
	Map        *Map
	Marker     *Marker
	Active     *GPSNotification
	Timeout    time.Duration
	LastUpdate time.Time
}{
	Timeout: 5 * time.Minute,
}

// GPSNotification manages displaying the map
type GPSNotification struct {
	Latitude  float64
	Longitude float64
	container *ui.Container
	e         *ui.Engine
}

// NewGPSNotification creates a new GPSNotification
func NewGPSNotification(latitude, longitude float64, e *ui.Engine, c *ui.Container) (*GPSNotification, error) {
	if gps.Map == nil {
		key := os.Getenv("MAPTILER_KEY")
		if key == "" {
			return nil, errors.New("invalid MAPTILER_KEY")
		}
		gps.Map = NewMap(key, e)
		gps.Map.Alpha = 210
		gps.Map.Zoom = 16
		gps.Map.Latitude = latitude
		gps.Map.Longitude = longitude
		home, err := ui.ImageFromFile(Asset("map/home.png"), e.Renderer)
		if err != nil {
			return nil, err
		}
		gps.Map.CenterOffsetX = 240
		gps.Marker = &Marker{
			Latitude:  52.4900311,
			Longitude: 4.7602125,
			Sprite:    ui.NewSprite(home)}
		gps.Marker.Sprite.AnchorX = 0.5
		gps.Marker.Sprite.AnchorY = 0.5
		gps.Map.Markers = append(gps.Map.Markers, gps.Marker)

		charlie, err := ui.ImageFromFile(Asset("map/charlie.png"), e.Renderer)
		if err != nil {
			return nil, err
		}
		gps.Map.CenterOffsetX = 240
		gps.Marker = &Marker{
			Latitude:  latitude,
			Longitude: longitude,
			Sprite:    ui.NewSprite(charlie)}
		gps.Marker.Sprite.AnchorX = 0.5
		gps.Marker.Sprite.AnchorY = 0.5
		gps.Map.Markers = append(gps.Map.Markers, gps.Marker)
	}

	return &GPSNotification{
		Latitude:  latitude,
		Longitude: longitude,
		e:         e,
		container: c,
	}, nil
}

// Wait until no there are no new notifications for x minutes
func (l *GPSNotification) Wait() {
	if gps.Active == l {
		for {
			d := gps.Timeout - time.Since(gps.LastUpdate)
			if d < 0 {
				break
			}
			time.Sleep(d)
		}
	}
}

// Close the location viewer
func (l *GPSNotification) Close() error {
	// return l.gps.Map.Close()
	return nil
}

// Compose to match Notification interface
func (l *GPSNotification) Compose(_ *sdl.Renderer) error {
	return nil
}

// Show the map
func (l *GPSNotification) Show() tween.Tween {
	gps.LastUpdate = time.Now()
	if gps.Active != nil {
		tl := tween.Timeline{}
		tl.AddAt(0, tween.FromTo(gps.Marker.Latitude, l.Latitude, 3*time.Second, tween.EaseOutQuad, func(v float64) {
			gps.Marker.Latitude = v
		}))
		tl.AddAt(0, tween.FromTo(gps.Marker.Longitude, l.Longitude, 3*time.Second, tween.EaseOutQuad, func(v float64) {
			gps.Marker.Longitude = v
		}))
		max := 0.0012
		if l.Latitude-max > gps.Map.Latitude || l.Latitude+max < gps.Map.Latitude {
			tl.AddAt(0, tween.FromTo(gps.Map.Latitude, l.Latitude, 3*time.Second, tween.EaseOutQuad, func(v float64) {
				gps.Map.Latitude = v
			}))
		}
		if l.Longitude+max < gps.Map.Longitude || l.Longitude-max > gps.Map.Longitude {
			tl.AddAt(0, tween.FromTo(gps.Map.Longitude, l.Longitude, 3*time.Second, tween.EaseOutQuad, func(v float64) {
				gps.Map.Longitude = v
			}))
		}
		return &tl

	}
	gps.Active = l
	l.container.Append(gps.Map)
	return tween.Empty()
}

// Hide the map
func (l *GPSNotification) Hide() tween.Tween {
	if gps.Active == l {
		l.container.Remove(gps.Map)
		gps.Active = nil
	}
	return tween.Empty()
}
