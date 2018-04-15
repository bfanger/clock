package button

import (
	"runtime"
)

func Combi(char rune, gpio int) (<-chan error, error) {
	var presses chan error
	keyboard, err := Keyboard(char)
	if err != nil {
		return nil, err
	}
	if runtime.GOOS == "darwin" {
		presses = keyboard
	} else {
		presses = make(chan error)
		button, err := Gpio(gpio)
		if err != nil {
			return nil, err
		}
		go func() {
			defer close(presses)
			for {
				select {
				case e, ok := <-keyboard:
					if !ok {
						keyboard = nil
					} else {
						err = e
					}
				case e, ok := <-button:
					if !ok {
						button = nil
					} else {
						err = e
					}
				}
				if button == nil && keyboard == nil {
					return
				}
				presses <- err
				if err != nil {
					return
				}
			}
		}()
	}
	return presses, nil

}
