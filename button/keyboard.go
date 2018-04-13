package button

import (
	termbox "github.com/nsf/termbox-go"
)

func Keyboard(key rune) (chan error, error) {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	c := make(chan error)
	go func() {
		defer termbox.Close()
		defer close(c)
		for {
			event := termbox.PollEvent()
			if event.Key == termbox.KeyCtrlC {
				break
			}
			if event.Key == termbox.KeyEsc {
				break
			}

			if key == event.Ch {
				c <- nil
			}
		}
	}()

	return c, nil
}
