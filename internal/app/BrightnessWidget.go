package app

import (
	"time"

	"../../internal/engine"
)

// BrightnessWidget adapts the brighness to time of day
type BrightnessWidget struct {
	Parent     engine.ContainerInterface
	Brightness *engine.Brightness
	Timeout    *engine.Timeout
}

// Mount creates an active BrightnessWidget
func (*BrightnessWidget) Mount(parent engine.ContainerInterface) error {

	brightnessWidget := &BrightnessWidget{}

	brightness, err := engine.NewBrightness(brightnessAlphaForTime(time.Now().Local()))
	if err != nil {
		return err
	}
	brightnessWidget.Brightness = brightness
	parent.Add(brightnessWidget.Brightness)
	brightnessWidget.Parent = parent
	brightnessWidget.tick()
	return nil
}

// Unmount and dispose resources
func (brightnessWidget *BrightnessWidget) Unmount() error {
	if brightnessWidget.Parent == nil {
		return nil // Software brighness was disabled
	}
	brightnessWidget.Timeout.Cancel()
	if err := brightnessWidget.Parent.Remove(brightnessWidget.Brightness); err != nil {
		return err
	}
	return brightnessWidget.Brightness.Dispose()
}

// Update the brightness texture based on the current time
func (brightnessWidget *BrightnessWidget) tick() {
	now := time.Now().Local()
	brightnessWidget.Brightness.Alpha = brightnessAlphaForTime(now)
	if err := brightnessWidget.Brightness.Update(); err != nil {
		panic(err)
	}
	// Calculate the delay to the begin of the next hour
	delay := (time.Duration(1) * time.Hour)
	delay -= (time.Duration(now.Minute()) * time.Minute)
	delay -= (time.Duration(now.Second()) * time.Second)

	brightnessWidget.Timeout = engine.SetTimeout(brightnessWidget.tick, delay)
}

func brightnessAlphaForTime(time time.Time) uint8 {
	levels := map[int]uint8{
		0:  180,
		1:  180,
		2:  180,
		3:  180,
		4:  160,
		5:  140,
		6:  120,
		7:  100,
		8:  90,
		9:  80,
		10: 50,
		11: 30,
		12: 10,
		13: 10,
		14: 20,
		15: 0,
		16: 80,
		17: 100,
		18: 110,
		19: 120,
		20: 130,
		21: 150,
		22: 160,
		23: 180,
	}
	return levels[time.Hour()]
}
