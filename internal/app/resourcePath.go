package app

import (
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

// ResourcePath to a absolute path for a file in the assets folder
func ResourcePath(filename string) string {
	assetPath := sdl.GetBasePath() + "assets/"
	if _, err := os.Stat(assetPath); err != nil {
		assetPath = "./assets/"
	}
	return assetPath + filename

}
