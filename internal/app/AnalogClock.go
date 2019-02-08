package app

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"time"

	"github.com/bfanger/clock/pkg/ui"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// AnalogClock displays the current time
type AnalogClock struct {
	engine    *ui.Engine
	container *ui.Container
	face      struct {
		image  *ui.Image
		sprite *ui.Sprite
	}
	timer *Timer
	font  *ttf.Font
	hours [12]struct {
		text   *ui.Text
		sprite *ui.Sprite
	}
	hourHand struct {
		image  *ui.Image
		sprite *ui.Sprite
	}
	minuteHand struct {
		image  *ui.Image
		sprite *ui.Sprite
	}
	x, y int32
	done chan bool
}

const radius = 184.0
const fontSize = 70

var color = sdl.Color{R: 103, G: 103, B: 109, A: 255}
var activeColor = sdl.Color{R: 203, G: 222, B: 198, A: 255}

// NewAnalogClock creats a new time widget
func NewAnalogClock(engine *ui.Engine) (*AnalogClock, error) {
	c := &AnalogClock{
		engine:    engine,
		container: &ui.Container{},
		done:      make(chan bool)}
	i, err := ui.ImageFromFile(Asset("analog-clock/face.png"), engine.Renderer)
	if err != nil {
		return nil, err
	}
	c.face.image = i
	c.face.sprite = ui.NewSprite(c.face.image)
	c.face.sprite.AnchorX = 0.5
	c.face.sprite.AnchorY = 0.5
	c.container.Append(c.face.sprite)

	c.timer, err = NewTimer(engine)
	if err != nil {
		return nil, fmt.Errorf("couldn't create timer: %v", err)
	}
	c.timer.Sprite.AnchorX = 0.5
	c.timer.Sprite.AnchorY = 0.5
	c.container.Append(c.timer)

	c.font, err = ttf.OpenFont(Asset("Roboto-Regular.ttf"), fontSize)
	if err != nil {
		return nil, err
	}
	// c.font.SetStyle()

	for i := 0; i < 12; i++ {
		text := strconv.Itoa(i)
		if text == "0" {
			text = "12"
		}
		c.hours[i].text = ui.NewText(text, c.font, color)
		hour := ui.NewSprite(c.hours[i].text)
		hour.AnchorX = 0.5
		hour.AnchorY = 0.5
		c.container.Append(hour)
		c.hours[i].sprite = hour
	}
	// minute hand
	i, err = ui.ImageFromFile(Asset("analog-clock/minute-hand.png"), engine.Renderer)
	if err != nil {
		return nil, err
	}
	c.minuteHand.image = i
	c.minuteHand.sprite = ui.NewSprite(c.minuteHand.image)
	c.minuteHand.sprite.AnchorX = 0.5
	c.minuteHand.sprite.AnchorY = 0.5
	c.minuteHand.sprite.SetAlpha(180)
	c.container.Append(c.minuteHand.sprite)

	// hour hand
	i, err = ui.ImageFromFile(Asset("analog-clock/hour-hand.png"), engine.Renderer)
	if err != nil {
		return nil, err
	}
	c.hourHand.image = i
	c.hourHand.sprite = ui.NewSprite(c.hourHand.image)
	c.hourHand.sprite.AnchorX = 0.5
	c.hourHand.sprite.AnchorY = 0.5
	c.hourHand.sprite.SetAlpha(160)
	c.container.Append(c.hourHand.sprite)
	c.MoveTo(screenWidth/2, screenHeight/2)
	go c.tick()

	return c, nil
}

// MoveTo positions the clock
func (c *AnalogClock) MoveTo(x, y int32) {
	c.face.sprite.X = x
	c.face.sprite.Y = y
	c.timer.Sprite.X = x
	c.timer.Sprite.Y = y
	c.hourHand.sprite.X = x
	c.hourHand.sprite.Y = y
	c.minuteHand.sprite.X = x
	c.minuteHand.sprite.Y = y

	for i := 0; i < 12; i++ {
		angle := math.Pi * (float64(i) / 6)
		c.hours[i].sprite.X = x + int32(math.Sin(angle)*radius)
		c.hours[i].sprite.Y = y + int32(math.Cos(angle)*-radius)
	}
	c.x = x
	c.y = y
	c.updateTime()
}

// Close frees related resources
func (c *AnalogClock) Close() error {
	close(c.done)
	c.font.Close()
	closers := []io.Closer{
		c.face.image,
		c.timer,
		c.minuteHand.image,
		c.hourHand.image}
	for i := 0; i < 12; i++ {
		closers = append(closers, c.hours[i].text)
	}
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			return err
		}
	}

	return nil
}

// Compose the clock
func (c *AnalogClock) Compose(r *sdl.Renderer) error {
	return c.container.Compose(r)
}

// Update the clock
func (c *AnalogClock) updateTime() error {
	now := time.Now()
	hour := now.Hour() % 12
	minute := now.Minute()
	second := now.Second()
	// hour
	c.hours[hour].text.SetColor(activeColor)
	previous := hour - 1
	if previous == -1 {
		previous = 11
	}
	c.hours[previous].text.SetColor(color)
	c.minuteHand.sprite.Rotation = (float64(minute) * 6) + (float64(second) * 0.1)
	c.hourHand.sprite.Rotation = (360 * (float64(hour) / 12)) + (float64(minute) * 0.5)
	// minute
	index := int((float32(minute) + 2.5) / 5)
	previous = index - 1
	if previous == -1 {
		previous = 11
	}
	if index == 12 {
		index = 0
	}

	return nil
}

func (c *AnalogClock) tick() {
	for {
		select {
		case <-c.done:
			return
		case <-time.After(time.Until(next(time.Second, time.Now()))):
			c.engine.Go(c.updateTime)
		}
	}
}
