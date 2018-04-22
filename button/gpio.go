package button

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
)

func init() {
	// Load all the drivers:
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}
}

// Gpio creates a signal when a button is pressed
func Gpio(bcm int) (chan error, error) {
	p := gpioreg.ByName(strconv.Itoa(bcm))
	if p == nil {
		return nil, errors.New("could not register pin")
	}
	if err := p.In(gpio.PullUp, gpio.RisingEdge); err != nil {
		return nil, fmt.Errorf("could not setup pin: %v", err)
	}
	c := make(chan error)
	go func() {
		for {
			p.WaitForEdge(-1)
			c <- nil
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return c, nil
}
