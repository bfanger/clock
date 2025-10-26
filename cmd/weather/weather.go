package main

import (
	"encoding/json"
	"fmt"
	"io"
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
	weathermapParams := os.Getenv("OPENWEATHERMAP_PARAMS")
	buienradarParams := os.Getenv("BUIENRADAR_PARAMS")

	if weathermapParams == "" {
		app.Fatal(errors.New("Missing OPENWEATHERMAP_PARAMS"))
	}

	if buienradarParams != "" {
		go buienradar(buienradarParams)
	}
	if weathermapParams != "" {
		openweathermap(weathermapParams)
	}
}

func buienradar(params string) {
	res, err := http.Get("https://graphdata.buienradar.nl/3.0/forecast/geo/RainHistoryForecast?" + params)
	if err != nil {
		app.Fatal(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		app.Fatal(errors.Wrap(err, "invalid response"))
	}
	var data struct {
		Forecasts []BuienradarForecast `json:"forecasts"`
	}
	fmt.Println(string(body))
	if err := json.Unmarshal(body, &data); err != nil {
		app.Fatal(errors.Wrap(err, "invalid json"))
	}
	for _, forecast := range data.Forecasts {
		fmt.Printf("%s: %.1fmm/h ( %.0f %% )\n", forecast.Timestamp.Format("2-15:04"), forecast.Value, forecast.Percentage)
	}
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

type BuienradarForecast struct {
	Timestamp  Timestamp `json:"utcDateTime"`
	Value      float64   `json:"dataValue"`
	Percentage float64   `json:"percentageValue"`
}

type Timestamp struct {
	time.Time
}

func (p *Timestamp) UnmarshalJSON(bytes []byte) error {
	var raw string
	err := json.Unmarshal(bytes, &raw)

	if err != nil {
		return errors.Wrap(err, "error decoding timestamp")
	}
	timestamp, err := time.Parse("2006-01-02T15:04:05", raw)
	if err != nil {
		return errors.Wrap(err, "error parsing timestamp")
	}
	fmt.Println(timestamp)
	p.Time = timestamp
	return nil
}
