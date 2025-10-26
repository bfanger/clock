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
	weerliveParams := os.Getenv("WEERLIVE_PARAMS")

	if weathermapParams == "" && weerliveParams == "" {
		app.Fatal(errors.New("Missing OPENWEATHERMAP_PARAMS or WEERLIVE_PARAMS"))
	}
	if weerliveParams == "" {
		openweathermap(weathermapParams)
	} else {
		go weerlive(weerliveParams)
		openweathermap(weathermapParams)
	}
}

func weerlive(apiparams string) {
	res, err := http.Get("https://weerlive.nl/api/weerlive_api_v2.php?" + apiparams)
	if err != nil {
		app.Fatal(err)
	}
	defer res.Body.Close()
	var data struct {
		Hours []WeerliveHour `json:"uur_verw"`
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		app.Fatal(errors.Wrap(err, "invalid response"))
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		app.Fatal(errors.Wrap(err, "invalid format"))
	}
	for _, hour := range data.Hours {
		fmt.Printf("%s: %.2f°C, rainfall: %.2fmm, %s\n", hour.Timestamp.Format("15:00"), hour.Temperature, hour.Rainfall, hour.Keyword)
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

type WeerliveHour struct {
	Hour      string    `json:"uur"`
	Timestamp Timestamp `json:"timestamp"`
	Keyword   string    `json:"image"`
	// Expected temperature in °C
	Temperature float64 `json:"temp"`
	// Cumulative rainfall in mm
	Rainfall float64 `json:"neersl"`
}

type Timestamp struct {
	time.Time
}

func (p *Timestamp) UnmarshalJSON(bytes []byte) error {
	var raw int64
	err := json.Unmarshal(bytes, &raw)

	if err != nil {
		return errors.Wrap(err, "error decoding timestamp")
	}
	p.Time = time.Unix(raw, 0)
	return nil
}
