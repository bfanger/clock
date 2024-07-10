package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/bfanger/clock/internal/sonos"
)

var room = "Woonkamer"

func main() {
	speaker, err := sonos.FindRoom(room)
	fmt.Printf("Found: %s\n", speaker.Name)
	if err != nil {
		panic(err)
	}
	err = speaker.HandleVolumeEvents(func(volume int) { sendVolume(volume) })
	if err != nil {
		panic(err)
	}
}

func sendVolume(volume int) error {
	fmt.Printf("Sonos volume: %d\n", volume)
	data := url.Values{}
	data.Set("volume", fmt.Sprintf("%d", volume))
	r, err := http.PostForm("http://localhost:8080/volume", data)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return nil
}
