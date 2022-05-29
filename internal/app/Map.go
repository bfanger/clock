package app

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"time"

	"github.com/bfanger/clock/pkg/tween"
	"github.com/bfanger/clock/pkg/ui"
	"github.com/pkg/errors"
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
	Markers       []*Marker
	Alpha         uint8
	key           string
	engine        *ui.Engine
	tiles         map[int]map[int]*ui.Image
	downloads     map[string]bool
	err           error
}

// NewMap create a new map
func NewMap(key string, e *ui.Engine) *Map {
	return &Map{
		Zoom:      2,
		W:         screenWidth,
		H:         screenHeight,
		Alpha:     255,
		key:       key,
		engine:    e,
		tiles:     make(map[int]map[int]*ui.Image),
		downloads: make(map[string]bool)}
}

// Close map and free memory used by the tiles
func (m *Map) Close() error {
	for y := range m.tiles {
		for _, img := range m.tiles[y] {
			if err := img.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

// Compose the map
func (m *Map) Compose(r *sdl.Renderer) error {
	if m.err != nil {
		return m.err
	}
	lat, lon := coords(m.Latitude, m.Longitude, m.Zoom)
	x := int(lat)
	y := int(lon)
	offsetX := int32((lat-math.Floor(lat))*tileSize) - m.CenterOffsetX
	offsetY := int32((lon-math.Floor(lon))*tileSize) - m.CenterOffsetY
	minX, maxX := minmax(m.W, offsetX)
	minY, maxY := minmax(m.H, offsetY)
	prev := r.GetClipRect()
	restore := &prev
	if prev.W == 0 {
		restore = nil
	}
	defer r.SetClipRect(restore)
	r.SetClipRect(&sdl.Rect{X: m.X, Y: m.Y, W: m.W, H: m.H})
	for dx := minX; dx <= maxX; dx++ {
		for dy := minY; dy <= maxY; dy++ {
			image := m.getTile(x+dx, y+dy)
			if image != nil {
				if m.Alpha == 255 {
					image.Texture.SetBlendMode(sdl.BLENDMODE_NONE)
				} else {
					image.Texture.SetBlendMode(sdl.BLENDMODE_BLEND)
					if err := image.Texture.SetAlphaMod(m.Alpha); err != nil {
						return err
					}
				}
				src, dst := m.tileRects(dx, dy, offsetX, offsetY)
				if err := r.Copy(image.Texture, src, dst); err != nil {
					return err
				}
			}
		}
	}
	for _, marker := range m.Markers {
		marker.Sprite.X, marker.Sprite.Y = m.XY(marker.Latitude, marker.Longitude)
		if err := marker.Sprite.Compose(r); err != nil {
			return err
		}
	}
	return nil
}

// PanTo a position on the map
func (m *Map) PanTo(latitude, longitude float64, d time.Duration) *tween.Timeline {
	tl := &tween.Timeline{}
	tl.Add(tween.FromTo(m.Latitude, latitude, d, tween.EaseInOutQuad, func(v float64) {
		m.Latitude = v
	}))
	tl.AddAt(0, tween.FromTo(m.Longitude, longitude, d, tween.EaseInOutQuad, func(v float64) {
		m.Longitude = v
	}))
	return tl
}

// XY get the screen location in pixels of the given latitude and longitude
func (m *Map) XY(latitude, longitude float64) (int32, int32) {
	mapX, mapY := coords(m.Latitude, m.Longitude, m.Zoom)
	tileX, tileY := coords(latitude, longitude, m.Zoom)
	x := m.X + (m.W / 2) + m.CenterOffsetX + int32((tileX-mapX)*tileSize)
	y := m.Y + (m.H / 2) + m.CenterOffsetY + int32((tileY-mapY)*tileSize)
	return x, y
}

func (m *Map) getTile(x, y int) *ui.Image {
	if m.tiles[x] == nil {
		m.tiles[x] = make(map[int]*ui.Image)
	}
	if m.tiles[x][y] == nil {
		url := fmt.Sprintf("https://api.maptiler.com/maps/basic/256/%d/%d/%d.png?key=%s", m.Zoom, x, y, m.key)
		if m.downloads[url] == false {
			m.downloads[url] = true
			go func() {
				body, err := download(url)
				if err != nil {
					m.err = err
					return
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
					m.err = err
					return
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
		return nil, errors.WithStack(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, errors.Errorf("request \"%s\" failed: %s", url, response.Status)
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
