package app

import (
	"fmt"
	"go/build"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/bfanger/clock/internal/schedule"
	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/sdl"
)

const endpoint = "http://localhost:8080/"

// Asset returns the absolute path for a file in the assets folder
func Asset(filename string) string {
	binPath := sdl.GetBasePath() + "assets/"
	_, err := os.Stat(binPath)
	if err == nil {
		return binPath + filename
	}
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	packagePath := gopath + "/src/github.com/bfanger/clock/assets/"
	_, err = os.Stat(packagePath)
	if err == nil {
		return packagePath + filename
	}
	return "./assets/" + filename
}

// Fatal exit with a formatted error.
func Fatal(err error) {
	RED := "\033[1;31m"
	GRAY := "\033[1;30m"
	NC := "\033[0m"
	log.Println(RED + err.Error() + NC)
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}
	if err, ok := err.(stackTracer); ok {
		for _, f := range err.StackTrace() {
			fmt.Printf(GRAY+"%+s:%d\n"+NC, f, f)
		}
	}
	os.Exit(1)
}

// NotificationOption for dynamic activation options
type NotificationOption struct {
	Key   string
	Value string
}

// ShowNotification to clock
func ShowNotification(notification string, d time.Duration, opts ...NotificationOption) error {
	fmt.Printf("Sending notification %s\n", notification)
	data := url.Values{}
	data.Set("action", "notify")
	data.Set("icon", notification)
	data.Set("duration", strconv.Itoa(int(d.Seconds())))
	for _, o := range opts {
		data.Set(o.Key, o.Value)
	}
	r, err := http.PostForm(endpoint+"notify", data)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return nil
}

// ShowAppointment on the clock
func ShowAppointment(a *schedule.Appointment) error {
	var opts []NotificationOption
	if a.Timer != 0 {
		opts = append(opts, NotificationOption{Key: "timer", Value: strconv.Itoa(int(a.Timer.Minutes()))})
	}
	return ShowNotification(a.Notification, a.Duration, opts...)

}
