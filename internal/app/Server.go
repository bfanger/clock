package app

import (
	"net/http"

	"github.com/bfanger/clock/pkg/ui"
)

type Server struct {
	Background   *Background
	Clock        *Time
	Notification *Notification
	maximized    bool
}

func NewServer(engine *ui.Engine) (*Server, error) {
	s := &Server{maximized: true}
	var err error
	s.Background, err = NewBackground(engine)
	if err != nil {
		return nil, err
	}
	s.Notification, err = NewNotification(engine)
	if err != nil {
		return nil, err
	}
	s.Clock, err = NewTime(engine)
	if err != nil {
		return nil, err
	}
	return s, nil
}
func (s *Server) ListenAndServe() {
	http.HandleFunc("/", s.handleToggle)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func (s *Server) Toggle() {
	if s.maximized {
		s.Clock.Minimize()
		s.Background.Maximize()
		s.Notification.Show()
	} else {
		s.Clock.Maximize()
		s.Background.Minimize()
		s.Notification.Hide()
	}
	s.maximized = !s.maximized
}

func (s *Server) Close() {
	defer s.Background.Close()
	defer s.Clock.Close()
	defer s.Notification.Close()
}
func (s *Server) handleToggle(w http.ResponseWriter, r *http.Request) {
	// message := r.URL.Path
	// message = strings.TrimPrefix(message, "/")
	// message = "Hello " + message
	s.Toggle()
	if s.maximized {
		w.Write([]byte("maximized"))
	} else {
		w.Write([]byte("minimized"))
	}
}
