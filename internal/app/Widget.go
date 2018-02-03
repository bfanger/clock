package app

// Widget interface, Update must be call from the ui thread (with sdl.Do)
type Widget interface {
	Update() error
}
