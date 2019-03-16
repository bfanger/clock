package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/bfanger/clock/internal/app"
	"github.com/bfanger/clock/internal/pubsub"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}
	ttn, err := pubsub.NewTheThingsNetwork()
	if err != nil {
		app.Fatal(err)
	}
	mqtt := os.Getenv("MQTT_URL")
	if mqtt == "" {
		app.Fatal(errors.New("Missing MQTT_URL"))
	}
	mapbox := os.Getenv("MAPBOX_TOKEN")
	if mapbox == "" {
		app.Fatal(errors.New("Missing MAPBOX_TOKEN"))
	}
	c, err := pubsub.NewConnection(mqtt)
	if err != nil {
		app.Fatal(err)
	}
	defer c.Close()

	log.Println("Connected to mqtt")

	go c.HandleRPC("mapbox_token", func(_ []byte) []byte {
		return []byte(mapbox)
	})

	go c.HandleRPC("history/gps/charlie", func(payload []byte) []byte {
		latlngs, err := ttn.History(string(payload))
		if err != nil {
			panic(err)
		}
		response, err := json.Marshal(latlngs)
		if err != nil {
			panic(err)
		}
		return response
	})

	for l := range ttn.Updates() {
		log.Printf("%+v\n", l)
		update, err := json.Marshal(l)
		if err != nil {
			app.Fatal(err)
		}
		if err := c.Publish("sensors/gps/charlie", update); err != nil {
			app.Fatal(err)
		}
	}
	app.Fatal(errors.New("quit?"))
}
