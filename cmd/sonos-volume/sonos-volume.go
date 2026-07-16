package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/bfanger/clock/internal/app"
	"github.com/bfanger/clock/internal/sonos"
)

var room = "Woonkamer"

func main() {
	speaker, err := sonos.FindRoom(room)
	if err != nil {
		app.Fatal(err)
	}
	fmt.Printf("Found: \"%s\" (%s) in \"%s\"\n", speaker.Name, speaker.IP.String(), speaker.Room)

	err = speaker.HandleVolumeEvents(sendVolume)
	if err != nil {
		app.Fatal(err)
	}
}

func sendVolume(volume int) {
	fmt.Printf("Volume: %d\n", volume)
	data := url.Values{}
	data.Set("volume", fmt.Sprintf("%d", volume))
	r, err := http.PostForm(app.Endpoint("/volume"), data)
	if err != nil {
		app.Fatal(err)
	}
	defer r.Body.Close()
}
