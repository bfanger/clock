package app

import (
	"html/template"
	"net/http"
	"sync"
	"time"

	"github.com/bfanger/clock/pkg/tween"
	"github.com/bfanger/clock/pkg/ui"
)

type Server struct {
	Background   *Background
	Clock        *Time
	Notification *Notification
	engine       *ui.Engine
	serialize    sync.Mutex
}

// NewServer creates a new webserver and creates the widgets controlled by the endpoints
func NewServer(engine *ui.Engine) *Server {
	return &Server{engine: engine}
}

// ListenAndServe start listening to requests and serving responses
func (s *Server) ListenAndServe() {
	http.HandleFunc("/", s.handleToggle)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

type toggleRequest struct {
	Show bool
}

func (s *Server) handleToggle(w http.ResponseWriter, r *http.Request) {
	s.serialize.Lock()
	defer s.serialize.Unlock()
	data := toggleRequest{}
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			panic(err)
		}
		data.Show = r.PostForm.Get("action") == "Show"
		if data.Show {
			s.showNotification()
		} else {
			s.hideNotification()
		}
	}

	t, err := template.ParseFiles(asset("form.html"))
	if err != nil {
		panic(err)
	}
	w.Header().Add("Content-Type", "text/html")
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

func (s *Server) showNotification() {
	tl := &tween.Timeline{}
	tl.Add(s.Clock.Minimize())
	tl.AddAt(200*time.Millisecond, s.Background.Maximize())
	tl.AddAt(800*time.Millisecond, s.Notification.Show())
	s.engine.Animate(tl)
}

func (s *Server) hideNotification() {
	tl := &tween.Timeline{}
	tl.Add(s.Notification.Hide())
	tl.AddAt(100*time.Millisecond, s.Clock.Maximize())
	tl.AddAt(100*time.Millisecond, s.Background.Minimize())
	s.engine.Animate(tl)
}
