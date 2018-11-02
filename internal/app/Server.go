package app

import (
	"html/template"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/bfanger/clock/pkg/ui"
)

// Server handles API call
type Server struct {
	displayManager *DisplayManager
	engine         *ui.Engine
	serialize      sync.Mutex
}

// NewServer creates a new webserver and creates the widgets controlled by the endpoints
func NewServer(d *DisplayManager, e *ui.Engine) *Server {
	return &Server{displayManager: d, engine: e}
}

// ListenAndServe start listening to requests and serving responses
func (s *Server) ListenAndServe() {
	http.HandleFunc("/", s.notify)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

type formViewModel struct {
	Show bool
}

func (s *Server) notify(w http.ResponseWriter, r *http.Request) {
	s.serialize.Lock()
	defer s.serialize.Unlock()
	vm := formViewModel{}
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			panic(err)
		}
		icon := r.PostForm.Get("icon")
		err := s.engine.Do(func() error {
			var n Notification
			duration, err := strconv.Atoi(r.PostForm.Get("duration"))
			if err != nil {
				return err
			}
			if icon == "vis" {
				n, err = NewFeedFishNotification(s.engine, time.Duration(duration)*time.Second)
			} else {
				n, err = NewBasicNotification(s.engine, icon, time.Duration(duration)*time.Second)
			}
			if err != nil {
				return err
			}
			go s.displayManager.Notify(n)
			return nil
		})
		if err != nil {
			panic(err)
		}
	}

	t, err := template.ParseFiles(asset("form.html"))
	if err != nil {
		panic(err)
	}
	w.Header().Add("Content-Type", "text/html")
	if err := t.Execute(w, vm); err != nil {
		panic(err)
	}
}
