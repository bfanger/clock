package app

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/bfanger/clock/pkg/tween"
	"github.com/bfanger/clock/pkg/ui"
	"github.com/veandco/go-sdl2/sdl"
)

// Map showing latlon
// https://www.maptiler.com
type Map struct {
	Key       string
	Zoom      int
	Latitude  float64
	Longitude float64
	tiles     xtile
}
type xtile map[int]ytile
type ytile map[int]*ui.Image

// Compose the map
func (m *Map) Compose(r *sdl.Renderer) error {

	lat, lon := coords(m.Latitude, m.Longitude, m.Zoom)
	x := int(lat)
	y := int(lon)

	offsetX := int32((lat - math.Floor(lat)) * 512)
	offsetY := int32((lon - math.Floor(lon)) * 512)

	for dx := -1; dx <= 1; dx++ {
		for dy := 0; dy <= 1; dy++ {
			image := m.getTile(x+dx, y+dy, r)
			dst := &sdl.Rect{W: 512, H: 512, X: 400 - offsetX + (int32(dx) * 512), Y: 240 - offsetY + (int32(dy) * 512)}
			if err := r.Copy(image.Texture, &image.Frame, dst); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *Map) PanTo(latitude, longitude float64) *tween.Timeline {
	tl := &tween.Timeline{}
	tl.Add(tween.FromToFloat64(m.Latitude, latitude, time.Second, tween.EaseInOutQuad, func(lat float64) {
		m.Latitude = lat
	}))
	return tl
}
func (m *Map) getTile(x, y int, r *sdl.Renderer) *ui.Image {
	if m.tiles == nil {
		m.tiles = make(map[int]ytile)
	}
	if m.tiles[x] == nil {
		m.tiles[x] = make(map[int]*ui.Image)
	}
	if m.tiles[x][y] == nil {
		url := fmt.Sprintf("http://maps.tilehosting.com/styles/hybrid/%d/%d/%d@2x.jpg?key=%s", m.Zoom, x, y, m.Key)
		image, err := ui.ImageFromURL(url, r)
		if err != nil {
			log.Fatal(err)
		}
		m.tiles[x][y] = image
	}
	return m.tiles[x][y]
}

// Convert latitude and longitude into x and y of the tile
func coords(latitude, longitude float64, zoom int) (float64, float64) {
	n := math.Pow(2.0, float64(zoom))
	x := (longitude + 180.0) / 360.0 * n
	l := deg2rad(latitude)
	y := (1.0 - math.Log(math.Tan(l)+(1/math.Cos(l)))/math.Pi) / 2.0 * n
	return x, y
}

func deg2rad(angle float64) float64 {
	return (math.Pi / 180) * angle
}
