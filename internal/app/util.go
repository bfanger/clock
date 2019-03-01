package app

import (
	"fmt"
	"go/build"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/sdl"
)

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
