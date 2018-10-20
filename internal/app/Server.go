package app

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"sync"
	"time"

	"github.com/bfanger/clock/pkg/tween"
	"github.com/bfanger/clock/pkg/ui"
)

// Server handles API call
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

type formViewModel struct {
	Show bool
}

func (s *Server) handleToggle(w http.ResponseWriter, r *http.Request) {
	s.serialize.Lock()
	defer s.serialize.Unlock()
	vm := formViewModel{}
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			panic(err)
		}
		vm.Show = r.PostForm.Get("action") == "show"
		icon := r.PostForm.Get("icon")
		if vm.Show {
			err := s.engine.Do(func() error {
				if s.Notification != nil {
					if err := s.Notification.Close(); err != nil {
						return fmt.Errorf("failed to close notification: %v", err)
					}
				}
				n, err := NewNotification(s.engine, icon)
				if err != nil {
					return err
				}
				s.Notification = n
				return nil
			})
			if err != nil {
				panic(err)
			}
			s.ShowNotification(s.Notification)
		} else {
			if err := s.HideNotification(); err != nil {
				panic(err)
			}
		}
	} else {
		vm.Show = s.Notification != nil
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

// ShowNotification display a new notification
func (s *Server) ShowNotification(n *Notification) {
	tl := &tween.Timeline{}
	tl.Add(s.Clock.Minimize())
	tl.AddAt(200*time.Millisecond, s.Background.Maximize())
	tl.AddAt(800*time.Millisecond, n.Show())
	s.engine.Animate(tl)
}

// HideNotification hides the active notification
func (s *Server) HideNotification() error {
	if s.Notification == nil {
		return errors.New("no notification active")
	}
	tl := &tween.Timeline{}
	n := s.Notification
	s.Notification = nil
	tl.Add(n.Hide())
	tl.AddAt(100*time.Millisecond, s.Clock.Maximize())
	tl.AddAt(100*time.Millisecond, s.Background.Minimize())
	s.engine.Animate(tl)
	if err := s.engine.Do(n.Close); err != nil {
		return fmt.Errorf("failed to close notification: %v", err)
	}
	return nil
}
