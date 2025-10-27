package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
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
	weathermapParams := os.Getenv("OPENWEATHERMAP_PARAMS")
	buienradarParams := os.Getenv("BUIENRADAR_PARAMS")

	if weathermapParams == "" {
		app.Fatal(errors.New("Missing OPENWEATHERMAP_PARAMS"))
	}
	wg := &sync.WaitGroup{}
	if buienradarParams != "" {
		wg.Go(func() { buienradar(buienradarParams) })
	}
	if weathermapParams != "" {
		wg.Go(func() { openweathermap(weathermapParams) })
	}
	wg.Wait()
}

func buienradar(params string) {
	time.Sleep(2 * time.Second)
	fmt.Println("buienradar enabled")
	for {
		if err := sendForecast(params); err != nil {
			app.Fatal(err)
		}
		time.Sleep(10 * time.Minute)
	}
}

func sendForecast(params string) error {
	res, err := http.Get("https://graphdata.buienradar.nl/3.0/forecast/geo/RainHistoryForecast?" + params)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "invalid response")
	}
	var data struct {
		Forecasts []struct {
			UTCDateTime     string  `json:"utcDateTime"`
			PercentageValue float64 `json:"percentageValue"`
		} `json:"forecasts"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		app.Fatal(errors.Wrap(err, "invalid json"))
	}
	buffer, err := json.Marshal(data.Forecasts)
	if err != nil {
		app.Fatal(err)
	}
	r, err := http.Post("http://localhost:8080/rainfall", "application/json", bytes.NewBuffer(buffer))
	if err != nil {
		app.Fatal(err)
	}
	defer r.Body.Close()
	return nil
}

func openweathermap(params string) {
	trigger := schedule.RepeatedAppointment{
		Notification: "ice",
		Hour:         6,
		Minute:       30,
		Duration:     2 * time.Hour,
		Repeat:       schedule.RepeatDays{Monday: true, Tuesday: true, Wednesday: true, Thursday: true, Friday: true},
	}
	for {
		appointment, err := trigger.Planned()
		if err != nil {
			app.Fatal(err)
		}
		appointment.Wait()
		fmt.Println("Getting temp from openweathermap")
		temp, err := getTemperature(params)
		if err != nil {
			app.Fatal(err)
		}
		fmt.Printf("temp: %.2f°C\n", temp)
		if temp < 3 {
			app.ShowAppointment(appointment)
		}
	}
}

func getTemperature(params string) (float64, error) {
	r, err := http.Get("https://api.openweathermap.org/data/2.5/weather?" + params)
	if err != nil {
		return 0, err
	}
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
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
