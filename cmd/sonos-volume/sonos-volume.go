package main

import (
	"fmt"

	"github.com/bfanger/clock/internal/sonos"
)

var room = "Woonkamer"

func main() {
	speaker, err := sonos.FindRoom(room)
	fmt.Printf("Found: %s\n", speaker.Name)
	if err != nil {
		panic(err)
	}
	volume, err := speaker.GetVolume()
	if err != nil {
		panic(err)
	}
	sendVolume(volume)

}

func sendVolume(volume int) {
	fmt.Printf("Sonos volume: %d\n", volume)

}
