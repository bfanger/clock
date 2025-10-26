package app

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/bfanger/clock/pkg/ui"
)

// Server handles API call
type Server struct {
	wm     *WidgetManager
	engine *ui.Engine
}

// NewServer creates a new webserver and creates the widgets controlled by the endpoints
func NewServer(wm *WidgetManager, e *ui.Engine) *Server {
	return &Server{wm: wm, engine: e}
}

// ListenAndServe start listening to requests and serving responses
func (s *Server) ListenAndServe() {
	http.HandleFunc("/", s.notify)
	http.HandleFunc("/notify", s.notify)
	http.HandleFunc("/button", s.button)
	http.HandleFunc("/volume", s.volumeHandler)
	http.HandleFunc("/rainfall", s.rainfallHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

type formViewModel struct {
	Show bool
}

func (s *Server) notify(w http.ResponseWriter, r *http.Request) {
	vm := formViewModel{}
	defer r.Body.Close()
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			panic(err)
		}
		icon := r.PostForm.Get("icon")
		err := s.engine.Do(func() error {
			var n Notification
			d, err := strconv.Atoi(r.PostForm.Get("duration"))
			if err != nil {
				return err
			}
			duration := time.Duration(d) * time.Second
			if icon == "ice" || icon == "plastic" || icon == "papier" || icon == "gft" {
				n, err = NewTrayNotification(icon, s.engine, duration)
			} else if icon == "vis" {
				n, err = NewFeedFishNotification(s.engine, duration)
			} else if icon == "gps" {
				lat, err := strconv.ParseFloat(r.PostForm.Get("latitude"), 64)
				if err != nil {
					return err
				}
				lng, err := strconv.ParseFloat(r.PostForm.Get("longitude"), 64)
				if err != nil {
					return err
				}
				n, err = NewGPSNotification(lat, lng, s.engine, s.wm.background)
				if err != nil {
					return err
				}
			} else {
				n, err = NewBasicNotification(s.engine, icon, duration)
			}
			if err != nil {
				return err
			}
			timer, err := strconv.Atoi(r.PostForm.Get("timer"))
			if err == nil && timer > 0 {
				if err := s.wm.timer.SetDuration(time.Duration(timer)*time.Minute, time.Minute); err != nil {
					return err
				}
			}
			go s.wm.Notify(n)
			return nil
		})
		if err != nil {
			panic(err)
		}
	}

	t, err := template.ParseFiles(Asset("notify.html"))
	if err != nil {
		panic(err)
	}
	w.Header().Add("Content-Type", "text/html")
	if err := t.Execute(w, vm); err != nil {
		panic(err)
	}
}

func (s *Server) button(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	t, err := template.ParseFiles(Asset("button.html"))
	if err != nil {
		panic(err)
	}
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			panic(err)
		}
		go s.wm.ButtonPressed()
	}
	w.Header().Add("Content-Type", "text/html")
	if err := t.Execute(w, struct{}{}); err != nil {
		panic(err)
	}
}

func (s *Server) volumeHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	t, err := template.ParseFiles(Asset("volume.html"))
	if err != nil {
		panic(err)
	}
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			panic(err)
		}
		value, err := strconv.Atoi(r.PostForm.Get("volume"))
		if err != nil {
			panic(err)
		}
		go s.wm.VolumeChanged(value)
	}
	w.Header().Add("Content-Type", "text/html")
	if err := t.Execute(w, struct{}{}); err != nil {
		panic(err)
	}
}

func (s *Server) rainfallHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method == "POST" {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		data := []struct {
			PercentageValue float64 `json:"percentageValue"`
			UTCDateTime     string  `json:"utcDateTime"`
		}{}
		if err := json.Unmarshal(body, &data); err != nil {
			panic(err)
		}
		forecasts := make([]RainfallForecast, len(data))
		for i, v := range data {
			t, err := time.Parse("2006-01-02T15:04:05", v.UTCDateTime)
			if err != nil {
				panic(err)
			}
			forecasts[i] = RainfallForecast{
				Timestamp:  t,
				Percentage: v.PercentageValue,
			}
		}

		go s.wm.rainfall.SetForecasts(forecasts)
	}
	w.Write([]byte("null"))

}
