
setup:
	go install github.com/bokwoon95/wgo@latest

setup-macos: setup
	brew install pkg-config sdl2 sdl2_image sdl2_ttf
	
setup-linux: setup
	go install github.com/bokwoon95/wgo@latest
	sudo apt update
	sudo apt install libsdl2-dev libsdl2-image-dev libsdl2-ttf-dev

install:
	go install ./cmd/clock
	go install ./cmd/garbage-truck
	go install ./cmd/school-schedule
	go install ./cmd/reminders
	go install ./cmd/weather

dev:
	wgo run ./cmd/clock

test:
	go test ./pkg/tween ./internal/app ./internal/schedule

