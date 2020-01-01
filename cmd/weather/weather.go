package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/bfanger/clock/internal/app"
	"github.com/bfanger/clock/internal/schedule"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err)
	}
	trigger := schedule.RepeatedAppointment{
		Notification: "ice",
		Hour:         6,
		Minute:       30,
		Duration:     2 * time.Hour,
		Repeat:       schedule.RepeatDays{Monday: true, Tuesday: true, Wednesday: true, Thursday: true, Friday: true},
	}
	appid := os.Getenv("OPENWEATHERMAP_APPID")
	if appid == "" {
		app.Fatal(errors.New("Missing OPENWEATHERMAP_APPID"))
	}

	for {
		appointment, err := trigger.Planned()
		if err != nil {
			app.Fatal(err)
		}
		appointment.Wait()
		fmt.Println("Getting temp from openweathermap")
		temp, err := getTemp(appid)
		if err != nil {
			app.Fatal(err)
		}
		fmt.Printf("temp: %.2f°C\n", temp)
		if temp < 3 {
			app.ShowAppointment(appointment)
		}
	}
}

func getTemp(appid string) (float64, error) {
	r, err := http.Get("https://api.openweathermap.org/data/2.5/weather?units=metric&lat=52.49&lon=4.76&appid=" + appid)
	if err != nil {
		return 0, err
	}
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return 0, errors.Wrap(err, "invalid response")
	}
	var obj struct {
		Message string
		Main    struct {
			Temp float64
		}
	}
	if err = json.Unmarshal(data, &obj); err != nil {
		return 0, errors.Wrap(err, "invalid json")
	}
	if r.StatusCode != http.StatusOK {
		return 0, errors.Errorf("status: %d: %s", r.StatusCode, obj.Message)
	}
	return obj.Main.Temp, nil
}
