package app

import (
	"fmt"

	"github.com/brian-armstrong/gpio"
)

// GpioButton listen to button 4 on the TFT (GPIO 25 / Pin 22)
func GpioButton(button chan bool) {
	fmt.Println("GPIO Button 4")
	watcher := gpio.NewWatcher()
	watcher.AddPin(25)
	defer watcher.Close()

	for {
		fmt.Println("Starting watcher.")
		_, value := watcher.Watch()
		if value == 0 {
			button <- true
		}
	}
}
