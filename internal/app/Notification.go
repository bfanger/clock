package app

import (
	"net/http"
	"net/url"

	"github.com/bfanger/clock/pkg/tween"
)

// Notification widget
type Notification interface {
	Show() tween.Tween
	Hide() tween.Tween
	Close() error
}

const endpoint = "http://localhost:8080/"

// ShowNotification is a helper for clock related services
func ShowNotification(icon string) error {
	data := url.Values{}
	data.Set("action", "show")
	data.Set("icon", icon)
	if _, err := http.PostForm(endpoint, data); err != nil {
		return err
	}
	return nil
}

// HideNotification is a helper for clock related services
func HideNotification() error {
	data := url.Values{}
	data.Set("action", "hide")
	if _, err := http.PostForm(endpoint, data); err != nil {
		return err
	}
	return nil
}
