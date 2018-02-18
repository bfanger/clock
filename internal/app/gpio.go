package app

import (
	"fmt"
	"os"

	"../../internal/engine"
	"github.com/brian-armstrong/gpio"
)

// HandleGpioButtons listen to button 4 on the TFT (GPIO 25 / Pin 22)
func HandleGpioButtons() {
	if _, err := os.Stat("/sys/class/gpio/"); err != nil {
		return
	}
	fmt.Println("GPIO Button 4")
	watcher := gpio.NewWatcher()
	watcher.AddPin(25)
	defer watcher.Close()
	for {
		fmt.Println("Starting watcher.")
		_, value := watcher.Watch()
		if value == 0 {
			engine.ButtonPressed(4)
		}
	}
}
