package app

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/bfanger/clock/pkg/tween"
	"github.com/bfanger/clock/pkg/ui"
	"github.com/veandco/go-sdl2/sdl"
)

const tileSize = 512

// Map showing latlon
// https://www.maptiler.com
type Map struct {
	X, Y          int32
	W, H          int32
	Zoom          int
	Latitude      float64
	Longitude     float64
	CenterOffsetX int32 // Offset the Latitude/Longitude from the center in pixels
	CenterOffsetY int32
	key           string
	engine        *ui.Engine
	tiles         map[int]map[int]*ui.Image
	downloads     map[string]bool
}

// NewMap create a new map
func NewMap(key string, e *ui.Engine) *Map {
	return &Map{
		Zoom:      2,
		W:         screenWidth,
		H:         screenHeight,
		key:       key,
		engine:    e,
		tiles:     make(map[int]map[int]*ui.Image),
		downloads: make(map[string]bool)}
}

// Compose the map
func (m *Map) Compose(r *sdl.Renderer) error {
	lat, lon := coords(m.Latitude, m.Longitude, m.Zoom)
	x := int(lat)
	y := int(lon)
	offsetX := int32((lat-math.Floor(lat))*tileSize) - m.CenterOffsetX
	offsetY := int32((lon-math.Floor(lon))*tileSize) - m.CenterOffsetY
	minX, maxX := minmax(m.W, offsetX)
	minY, maxY := minmax(m.H, offsetY)
	prev := r.GetClipRect()
	defer r.SetClipRect(&prev)
	r.SetClipRect(&sdl.Rect{X: m.X, Y: m.Y, W: m.W, H: m.H})
	for dx := minX; dx <= maxX; dx++ {
		for dy := minY; dy <= maxY; dy++ {
			image := m.getTile(x+dx, y+dy)
			if image != nil {
				src, dst := m.tileRects(dx, dy, offsetX, offsetY)
				if err := r.Copy(image.Texture, src, dst); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// PanTo a position on the map
func (m *Map) PanTo(latitude, longitude float64, d time.Duration) *tween.Timeline {
	tl := &tween.Timeline{}
	tl.Add(tween.FromToFloat64(m.Latitude, latitude, d, tween.EaseInOutQuad, func(v float64) {
		m.Latitude = v
	}))
	tl.AddAt(0, tween.FromToFloat64(m.Longitude, longitude, d, tween.EaseInOutQuad, func(v float64) {
		m.Longitude = v
	}))
	return tl
}
func (m *Map) getTile(x, y int) *ui.Image {
	if m.tiles == nil {

	}
	if m.tiles[x] == nil {
		m.tiles[x] = make(map[int]*ui.Image)
	}
	if m.tiles[x][y] == nil {
		url := fmt.Sprintf("http://maps.tilehosting.com/styles/hybrid/%d/%d/%d@2x.jpg?key=%s", m.Zoom, x, y, m.key)
		if m.downloads[url] == false {
			m.downloads[url] = true
			go func() {
				body, err := download(url)
				if err != nil {
					log.Fatal(err)
				}
				err = m.engine.Do(func() error {
					image, err := ui.ImageFromBytes(body, m.engine.Renderer)
					if err != nil {
						return err
					}
					m.tiles[x][y] = image
					return nil
				})
				if err != nil {
					log.Fatal(err)
				}

			}()
		}
	}
	return m.tiles[x][y]
}
func (m *Map) tileRects(dx, dy int, offsetX, offsetY int32) (*sdl.Rect, *sdl.Rect) {
	src := &sdl.Rect{W: tileSize, H: tileSize}
	dst := &sdl.Rect{W: tileSize, H: tileSize}
	dst.X = m.X + (m.W / 2) - offsetX + (int32(dx) * tileSize)
	dst.Y = m.Y + (m.H / 2) - offsetY + (int32(dy) * tileSize)
	return src, dst
}

func download(url string) ([]byte, error) {
	// @todo Implement cache?
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed: %s", response.Status)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func minmax(size, offset int32) (min, max int) {
	before := size/2 - offset
	if before > 0 {
		min = 0 - int(math.Ceil(float64(before)/float64(tileSize)))
	}
	after := size/2 - (tileSize - offset)
	if after > 0 {
		max = int(math.Ceil(float64(after) / float64(tileSize)))
	}
	return
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
