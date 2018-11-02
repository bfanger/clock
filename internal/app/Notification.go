package app

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/bfanger/clock/pkg/tween"
)

// Notification widget
type Notification interface {
	Show() tween.Tween
	Hide() tween.Tween
	Duration() time.Duration
	Close() error
}

const endpoint = "http://localhost:8080/"

// ShowNotification is a helper for clock related services
func ShowNotification(icon string, d time.Duration) error {
	data := url.Values{}
	data.Set("action", "notify")
	data.Set("icon", icon)
	data.Set("duration", strconv.Itoa(int(d.Seconds())))
	if _, err := http.PostForm(endpoint, data); err != nil {
		return err
	}
	return nil
}
