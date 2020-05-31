package app

import (
	"html/template"
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
	if r.Method == "POST" {
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
